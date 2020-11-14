package levenshtein

import "github.com/gnames/levenshtein/presenter"

type Levenshtein interface {
	Compare(str1, str2 string) presenter.Output
	CompareMult(input []Strings) []presenter.Output
	Opts() []Option
}
