package malloc_test

/*
#cgo LDFLAGS: /usr/local/lib/libjemalloc.a -L/usr/local/lib -Wl,-rpath,/usr/local/lib -ljemalloc -lm -lstdc++ -pthread -ldl
#include <stdlib.h>
#include <jemalloc/jemalloc.h>
*/
import "C"
import (
	"unsafe"
)

const MaxArrayLen = 1<<50 - 1

//go:linkname throw runtime.throw
func throw(s string)

func Calloc(n int) []byte {
	if n == 0 {
		return make([]byte, 0)
	}
	// We need to be conscious of the Cgo pointer passing rules:
	//
	//   https://golang.org/cmd/cgo/#hdr-Passing_pointers
	//
	//   ...
	//   Note: the current implementation has a bug. While Go code is permitted
	//   to write nil or a C pointer (but not a Go pointer) to C memory, the
	//   current implementation may sometimes cause a runtime error if the
	//   contents of the C memory appear to be a Go pointer. Therefore, avoid
	//   passing uninitialized C memory to Go code if the Go code is going to
	//   store pointer values in it. Zero out the memory in C before passing it
	//   to Go.

	ptr := C.je_calloc(C.size_t(n), 1)
	if ptr == nil {
		// NB: throw is like panic, except it guarantees the process will be
		// terminated. The call below is exactly what the Go runtime invokes when
		// it cannot allocate memory.
		throw("out of memory")
	}
	uptr := unsafe.Pointer(ptr)

	// Interpret the C pointer as a pointer to a Go array, then slice.
	return (*[MaxArrayLen]byte)(uptr)[:n:n]
}

// Free frees the specified slice.
func Free(b []byte) {
	if sz := cap(b); sz != 0 {
		b = b[:cap(b)]
		ptr := unsafe.Pointer(&b[0])
		C.je_free(ptr)
	}
}
