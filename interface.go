package levenshtein

type FuzzyDiff interface {
	Compare(str1, str2 string) presenter.Output
	CompareMult(input []input.Strings) []presenter.Output
	Opts() []Option
}
