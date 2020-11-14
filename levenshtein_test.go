package levenshtein_test

import (
	"fmt"
	"testing"

	"github.com/gnames/levenshtein"
	"github.com/gnames/levenshtein/presenter"
	"github.com/stretchr/testify/assert"
)

func TestMax(t *testing.T) {
	testData := []struct {
		str1     string
		str2     string
		editDist int
		aborted  bool
		tags1    string
		tags2    string
	}{
		{"Pomatomus", "Bomatomus", 1, false,
			"<subst>P</subst>omatomus", "<subst>B</subst>omatomus"},
		{"Poma tomus", "Pomatomos", 2, false,
			"Poma<ins> </ins>tom<subst>u</subst>s",
			"Poma<del> </del>tom<subst>o</subst>s"},
		{"Boston", "Chigago", 2, true, "", ""},
	}

	var fd levenshtein.Levenshtein
	opts := []levenshtein.Option{
		levenshtein.OptWithDiff(true),
		levenshtein.OptMaxEditDist(2),
	}
	fd = levenshtein.NewLevenshtein(opts...)
	for _, v := range testData {
		msg := fmt.Sprintf("'%s' vs '%s'", v.str1, v.str2)
		out := fd.Compare(v.str1, v.str2)
		assert.Equal(t, out.EditDist, v.editDist, msg)
		assert.Equal(t, out.Aborted, v.aborted, msg)
		assert.Equal(t, out.Tags1, v.tags1, msg)
		assert.Equal(t, out.Tags2, v.tags2, msg)
	}
}

func TestNoMax(t *testing.T) {
	testData := []struct {
		str1     string
		str2     string
		editDist int
		aborted  bool
		tags1    string
		tags2    string
	}{
		{"Pomatomus", "Bomatomus", 1, false,
			"<subst>P</subst>omatomus", "<subst>B</subst>omatomus"},
		{"Poma tomus", "Pomatomos", 2, false,
			"Poma<ins> </ins>tom<subst>u</subst>s",
			"Poma<del> </del>tom<subst>o</subst>s"},
		{"Boston", "Chicago", 7, false,
			"<del>C</del><subst>Boston</subst>",
			"<ins>C</ins><subst>hicago</subst>"},
	}

	var fd levenshtein.Levenshtein
	opts := []levenshtein.Option{
		levenshtein.OptWithDiff(true),
	}
	fd = levenshtein.NewLevenshtein(opts...)
	for _, v := range testData {
		msg := fmt.Sprintf("'%s' vs '%s'", v.str1, v.str2)
		out := fd.Compare(v.str1, v.str2)
		assert.Equal(t, out.EditDist, v.editDist, msg)
		assert.Equal(t, out.Aborted, v.aborted, msg)
		assert.Equal(t, out.Tags1, v.tags1, msg)
		assert.Equal(t, out.Tags2, v.tags2, msg)
	}
}

func TestNoDiff(t *testing.T) {
	testData := []struct {
		str1     string
		str2     string
		editDist int
		aborted  bool
		tags1    string
		tags2    string
	}{
		{"Pomatomus", "Bomatomus", 1, false, "", ""},
		// {"Poma tomus", "Pomatomos", 2, false, "", ""},
		// {"Boston", "Chicago", 7, false, "", ""},
	}

	var fd levenshtein.Levenshtein
	fd = levenshtein.NewLevenshtein()
	for _, v := range testData {
		msg := fmt.Sprintf("'%s' vs '%s'", v.str1, v.str2)
		out := fd.Compare(v.str1, v.str2)
		assert.Equal(t, out.EditDist, v.editDist, msg)
		assert.Equal(t, out.Aborted, v.aborted, msg)
		assert.Equal(t, out.Tags1, v.tags1, msg)
		assert.Equal(t, out.Tags2, v.tags2, msg)
	}
}

func TestMult(t *testing.T) {
	testData := []struct {
		str1     string
		str2     string
		editDist int
		aborted  bool
		tags1    string
		tags2    string
	}{
		{"Puma", "Poma", 1, false,
			"P<subst>u</subst>ma", "P<subst>o</subst>ma"},
		{"Puma concolor", "Pomacancolor", 3, false,
			"", ""},
	}
	str := make([]levenshtein.Strings, len(testData))
	for i, v := range testData {
		str[i] = levenshtein.Strings{String1: v.str1, String2: v.str2}
	}
	var fd levenshtein.Levenshtein
	fd = levenshtein.NewLevenshtein()
	out := fd.CompareMult(str)
	for i, v := range out {
		assert.Equal(t, v.EditDist, testData[i].editDist)
	}
}

// BenchmarkCompare checks the speed of fuzzy matching. Run it with:
// `go test -bench=. -benchmem -count=10 -run=XXX > bench.txt && benchstat bench.txt`
func BenchmarkCompare(b *testing.B) {
	d := levenshtein.NewLevenshtein()
	ops := []levenshtein.Option{levenshtein.OptWithDiff(true)}
	dDiff := levenshtein.NewLevenshtein(ops...)
	var out presenter.Output
	b.Run("CompareOnce", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			out = d.Compare("Pomatomus solatror", "Pomatomus saltator")
		}
		_ = fmt.Sprintf("%d\n", out.EditDist)
	})
	b.Run("CompareOnceDiff", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			out = dDiff.Compare("Pomatomus solatror", "Pomatomus saltator")
		}
		_ = fmt.Sprintf("%d\n", out.EditDist)
	})
}
