package main

import "strings"

var strs = []string{
	"here's",
	"a",
	"some",
	"long",
	"list",
	"of",
	"strings",
	"for",
	"you",
}

func buildStrNaive() string {
	var s string

	for _, v := range strs {
		s += v
	}

	return s
}

func buildStrBuilder() string {
	b := strings.Builder{}

	// Grow the buffer to a decent length, so we don't have to continually
	// re-allocate.
	b.Grow(60)

	for _, v := range strs {
		b.WriteString(v)
	}

	return b.String()
}
