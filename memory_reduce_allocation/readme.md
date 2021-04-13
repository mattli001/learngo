

```
matt@matt-macbook /s/g/t/l/memory_reduce_allocation> go test -v -benchmem -bench=BenchmarkMake
goos: darwin
goarch: amd64
pkg: github.com/mattli001/learngo/memory_reduce_allocation
BenchmarkMake
BenchmarkMake-4         1000000000               0.000000 ns/op        0 B/op          0 allocs/op
PASS
ok      github.com/mattli001/learngo/memory_reduce_allocation   0.014s
```

# Directory "convert"
Fork from https://github.com/appleboy/com/tree/master/convert



matt@matt-macbook /s/g/t/l/memory_reduce_allocation> go test -v -benchmem -run='^$' -bench=^Benchmark ./convert/...
goos: darwin
goarch: amd64
pkg: github.com/mattli001/learngo/memory_reduce_allocation/convert
BenchmarkCountParamsOld
BenchmarkCountParamsOld-4                 575650              2320 ns/op               0 B/op          0 allocs/op
BenchmarkCountParamsNew
BenchmarkCountParamsNew-4               12708914                89.7 ns/op             0 B/op          0 allocs/op
BenchmarkBytesToStrOld01
BenchmarkBytesToStrOld01-4               8028020               147 ns/op        6968.95 MB/s        1024 B/op          1 allocs/op
BenchmarkBytesToStrOld2
BenchmarkBytesToStrOld2-4               777117932                1.56 ns/op     654482.91 MB/s         0 B/op          0 allocs/op
BenchmarkBytesToStrNew
BenchmarkBytesToStrNew-4                1000000000               0.606 ns/op    1690216.00 MB/s        0 B/op          0 allocs/op
BenchmarkStr2BytesOld01
BenchmarkStr2BytesOld01-4                7094371               153 ns/op        6690.42 MB/s        1024 B/op          1 allocs/op
BenchmarkStr2BytesOld02
BenchmarkStr2BytesOld02-4               392720553                3.08 ns/op     332554.21 MB/s         0 B/op          0 allocs/op
BenchmarkStr2BytesNew
BenchmarkStr2BytesNew-4                 749864353                1.56 ns/op     657279.18 MB/s         0 B/op          0 allocs/op
BenchmarkConvertOld
BenchmarkConvertOld-4                    3727018               298 ns/op            2048 B/op          2 allocs/op
BenchmarkConvertNew
BenchmarkConvertNew-4                   973641290                1.22 ns/op            0 B/op          0 allocs/op
BenchmarkSnakeCasedNameRegex
BenchmarkSnakeCasedNameRegex-4             35252             37670 ns/op            4812 B/op         80 allocs/op
BenchmarkSnakeCasedNameOld
BenchmarkSnakeCasedNameOld-4              386313              2926 ns/op            2328 B/op          9 allocs/op
BenchmarkSnakeCasedNameNew
BenchmarkSnakeCasedNameNew-4             3201142               379 ns/op             624 B/op          2 allocs/op
BenchmarkTitleCasedName
BenchmarkTitleCasedName-4                1651602               785 ns/op             224 B/op          1 allocs/op
