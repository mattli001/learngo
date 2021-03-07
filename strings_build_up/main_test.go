package main

import (
	"testing"
)

var str string

func BenchmarkStringBuildNaive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str = buildStrNaive()
	}
}
func BenchmarkStringBuildBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str = buildStrBuilder()
	}
}
