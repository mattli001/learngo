go test -bench=. -benchmem



The benchmark results on a Macbook Pro:

🍎 strfmt → go test -bench=. -benchmem
goos: darwin
goarch: amd64
pkg: github.com/sjwhitworth/perfblog/strfmt
BenchmarkStrconv-8      30000000            39.5 ns/op        32 B/op          1 allocs/op
BenchmarkFmt-8          10000000           143 ns/op          72 B/op          3 allocs/op
We can see that the strconv version is 3.5x faster, results in 1/3rd the number of

```
matt@matt-macbook /s/g/t/l/strconv_instead_fmt> go test -bench=. -benchmem
goos: darwin
goarch: amd64
pkg: github.com/mattli001/learngo/strconv_instead_fmt
BenchmarkStrconv-4      23853715                51.3 ns/op            32 B/op          1 allocs/op
BenchmarkFmt-4           7814708               158 ns/op              64 B/op          2 allocs/op
PASS
ok      github.com/mattli001/learngo/strconv_instead_fmt        2.683s
```