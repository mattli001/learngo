

# [原始英文位置連結](https://dgraph.io/blog/post/manual-memory-management-golang-jemalloc/)
# [翻譯版結連](https://blog.csdn.net/Rong_Toa/article/details/110095720)

# 使用jemalloc在Go中進行手動內存管理

自2015年成立以來，Dgraph Labs一直是Go語言的用戶。五年之後，Go代碼達到20萬行，我們很高興地報告，我們仍然堅信Go是並且仍然是正確的選擇。我們對Go的興奮不僅限於構建系統，還使我們甚至可以使用Go編寫腳本，這些腳本通常是用Bash或Python編寫的。我們發現使用Go可以幫助我們構建乾淨，可讀，可維護並且最重要的是高效並發的代碼庫。

但是，自早期以來，我們一直關注的一個領域是： 內存管理。我們沒有反對Go垃圾收集器的方法，但是儘管它為開發人員提供了便利，但它具有其他內存垃圾收集器所面臨的相同問題：*`它根本無法與手動內存管理的效率競爭`*。


當您手動管理內存時，內存使用率較低，可預測，並且允許突發的內存分配不會引起內存使用率的瘋狂飆升。對於使用Go內存的Dgraph，所有這些都是問題1。實際上，Dgraph遇到內存不足，是我們從用戶那裡聽到的非常普遍的抱怨。

諸如Rust之類的語言之所以得到發展，部分原因是它允許安全的手動內存管理。我們可以完全理解。

根據我們的經驗，與嘗試使用具有垃圾回收的語言優化內存使用2相比，進行手動內存分配和解決潛在的內存洩漏花費的精力更少。在構建具有幾乎無限的可伸縮性的數據庫系統時，手動內存管理非常值得解決。

我們對Go的熱愛和避免使用Go GC的需求，使我們找到了在Go中進行手動內存管理的新穎方法。當然，*`大多數Go用戶將永遠不需要手動進行內存管理。除非您需要，否則我們建議您不要這樣做。`* 當您確實需要它時，您就會知道。

在這篇文章中，我將分享我們在Dgraph Labs中從對手動內存管理的探索中學到的知識，並說明我們如何在Go中手動管理內存。

## 通過Cgo創建內存

靈感來自 Cgo Wiki 的 [將C數組轉換為Go切片](https://github.com/golang/go/wiki/cgo#turning-c-arrays-into-go-slices)的文章。我們可以使用 mallocC 分配內存，然後將 unsafe 其傳遞給 Go，而不會受到 Go GC 的干擾。

```
import "C"
import "unsafe"
...
        var theCArray *C.YourType = C.getTheArray()
        length := C.getTheArrayLength()
        slice := (*[1 << 28]C.YourType)(unsafe.Pointer(theCArray))[:length:length]
```

但是，如 [golang.org/cmd/cgo](https://golang.org/cmd/cgo/#hdr-Passing_pointers) 中所述，以上內容帶有警告。


>注意：當前的實作存在一個錯誤。儘管允許Go代碼向C存儲器寫入nil或C指針（而不是Go指針），但是如果C存儲器的內容似乎是Go指針，則當前實作有時可能會導致運行時錯誤。因此，如果Go代碼要在其中存儲指針值，請避免將未初始化的C內存傳遞給Go代碼。將C中的內存清零，然後再傳遞給Go。

因此，malloc 我們不使用它，而是使用它的代價稍高的兄弟 calloc。calloc與的工作方式相同malloc，除了在將內存返回給調用者之前將內存清零。


*`我們僅通過實現基本Calloc和Free函數開始，這些函數通過Cgo為Go分配和取消分配字節片`*。為了測試這些功能，我們開發並運行了[連續內存使用的測試](https://github.com/dgraph-io/ristretto/tree/master/contrib)。該測試無休止地重複了一個分配/取消分配的循環，在該循環中，它首先分配了各種隨機大小的內存塊，直到分配了16GB的內存，然後釋放了這些塊，直到僅分配了1GB的內存。

此程序的C等效行為符合預期。我們會看到RSS內存在 `htop` 增加到16GB，然後下降到1GB，再增加回 16GB，依此類推。但是，Go程序在每個週期後使用 `Calloc` 和 `Free` 會逐漸使用更多的內存（請參見下表）。

我們將這種行為歸因於內存碎片，因為默認C.calloc調用中缺乏線程意識。在Go #dark-artsSlack頻道提供了一些幫助 （特別感謝Kale Blankenship）之後，我們決定jemalloc嘗試一下。

## jemalloc

>jemalloc 是通用的 malloc(3) 實現，它強調避免碎片和可擴展的並發支持。jemalloc於2005年首次用作FreeBSD libc分配器，此後，它便被許多應用程式使用，以依賴其可預測的行為。— http://jemalloc.net

我們的API，改用 jemalloc [3](https://dgraph.io/blog/post/manual-memory-management-golang-jemalloc/#fn:5) 用於 `calloc` 和 `free` 調用。它執行得很漂亮：jemalloc 本地支持幾乎沒有內存碎片的線程。來自我們的內存使用情況監視測試的分配-解除分配週期在預期的限制之間循環，而忽略了運行測試所需的少量開銷。

為了確保我們正在使用jemalloc並避免名稱衝突，我們在使用的過程中添加了 `je_` 前綴，因此我們的 API 現在正在調用 `je_calloc` 與 `je_free`，取代 `calloc ` 與 `free`。

![[go-mem-allocation.png]]

在上圖中，通過 C.calloc 分配 Go 內存導致主要的內存碎片，導致該程序在第11個週期佔用了20GB的內存。而 jemalloc 等同的代碼沒有明顯的碎片，每個週期下降到接近1GB。

在程序結束時（最右邊的小凹處），釋放了所有分配的內存之後，C.calloc 程序仍然佔用了約 20GB 的內存，而 jemalloc 則顯示只用了 400MB 的內存。

要安裝jemalloc，請從[此處下載](https://github.com/jemalloc/jemalloc/releases)它，然後運行以下命令：

```
./configure --with-jemalloc-prefix='je_' --with-malloc-conf='background_thread:true,metadata_thp:auto'
make
sudo make install
```

golang 的 [Calloc 代碼](https://github.com/dgraph-io/ristretto/blob/e2057c125fc2c91db8342a0b27f709acf2fd136f/z/calloc_jemalloc.go#L51) 大致如下：
```
	ptr := C.je_calloc(C.size_t(n), 1)
	if ptr == nil {
		// NB: throw is like panic, except it guarantees the process will be
		// terminated. The call below is exactly what the Go runtime invokes when
		// it cannot allocate memory.
		throw("out of memory")
	}
	uptr := unsafe.Pointer(ptr)

	atomic.AddInt64(&numBytes, int64(n))
	// Interpret the C pointer as a pointer to a Go array, then slice.
	return (*[MaxArrayLen]byte)(uptr)[:n:n]
```


我們將此代碼作為 [`Ristretto's z`](https://github.com/dgraph-io/ristretto/tree/master/z) 軟件包的一部分，因此 Dgraph 和 Badger 都可以使用它。為了使我們的代碼切換到使用 jemalloc 分配字節片，我們添加了一個 build 標籤 jemalloc。為了進一步簡化我們的部署，我們 jemalloc 通過設置正確的LDFLAGS ，使庫在任何生成的Go二進制文件中靜態鏈接。

## Laying out Go structs on byte slices

現在我們有了分配和釋放字節片的方法，下一步是使用它來佈局Go結構。我們可以從一個基本的[（完整代碼）](https://github.com/dgraph-io/ristretto/blob/master/contrib/demo) 開始。

```
type node struct {
    val  int
    next *node
}
 
var nodeSz = int(unsafe.Sizeof(node{}))
 
func newNode(val int) *node {
    b := z.Calloc(nodeSz)
    n := (*node)(unsafe.Pointer(&b[0]))
    n.val = val
    return n
}
 
func freeNode(n *node) {
    buf := (*[z.MaxArrayLen]byte)(unsafe.Pointer(n))[:nodeSz:nodeSz]
    z.Free(buf)
}
```

在上面的代碼中，我們使用佈局了C分配的內存中的Go結構newNode。我們創建了一個相應的freeNode函數，一旦完成了該結構，便可以釋放內存。Go結構具有基本數據類型int和指向下一個節點結構的指針，所有這些都已在程序中設置和訪問。我們分配了2M個節點物件，並從其中創建了一個鏈接列表，以演示jemalloc的正常功能。

使用默認的Go內存管理機制，我們看到為鏈結列表分配了 31 MiB的 Heap 內存，包含2M個物件，但沒有通過jemalloc分配。
```

$ go run .
Allocated memory: 0 Objects: 2000001
node: 0
...
node: 2000000
After freeing. Allocated memory: 0
HeapAlloc: 31 MiB
```

使用jemalloc構建標記，我們看到通過 jemalloc 分配了30 MiB的內存，在釋放鏈接列表後，該內存下降為零。Go Heap分配僅為399 KiB，這可能來自運行程序的開銷。
```
$ go run -tags=jemalloc .
Allocated memory: 30 MiB Objects: 2000001
node: 0
...
node: 2000000
After freeing. Allocated memory: 0

HeapAlloc: 399 KiB
```


