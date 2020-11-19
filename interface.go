package levenshtein

import (
	"fmt"

	"github.com/gnames/levenshtein/presenter"
)

// Levenshtein is the main API interface for the library. It provides methods
// for calculating edit distance between two strings, or between a slice of
// two strings.
type Levenshtein interface {
	// Compare calculates edit distance between two strings.
	Compare(str1, str2 string) presenter.Output

	// CompareMult, calculates edit distance between many name-strings.
	// The job is parallelized, and then reassembled in the same order as
	// the input.
	CompareMult(input []Strings) []presenter.Output

	// Option returns back options applied to the Levenshtein implementation.
	Opts() []Option
}

func Example() {
	opts := []Option{
		OptWithDiff(true),
		OptMaxEditDist(1),
	}
	l := NewLevenshtein(opts...)
	out := l.Compare("Something", "smoething")
	fmt.Printf("%+v\n", out)
}
