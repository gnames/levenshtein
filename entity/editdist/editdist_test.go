package editdist_test

import (
	"fmt"
	"testing"

	"github.com/gnames/levenshtein/entity/editdist"
	"github.com/stretchr/testify/assert"
)

func TestDist(t *testing.T) {
	testData := []struct {
		str1, str2 string
		dist       int
	}{
		{"Hello", "He1lo", 1},
		{"Pomatomus", "Pom-tomus", 1},
		{"Pomatomus", "Poma  tomus", 2},
		{"Pomatomus", "Pomщtomus", 1},
		{"sitting", "kitten", 3},
		{"Boston", "Chicago", 7},
		{"Chicago", "Boston", 7},
	}

	for _, v := range testData {
		msg := fmt.Sprintf("'%s' vs '%s'", v.str1, v.str2)
		dist, _, _ := editdist.ComputeDistance(v.str1, v.str2, false)
		assert.Equal(t, dist, v.dist, msg)
	}
}

func TestMax(t *testing.T) {
	testData := []struct {
		str1, str2 string
		dist       int
		abort      bool
	}{
		{"Hello", "Hello", 0, false},
		{"Pomatomus", "Pom-tomus", 1, false},
		{"Pomatomus", "Poma  tomus", 2, false},
		{"Pomatomus", "Pomщtomus", 1, false},
		{"pOMatomus", "Pomatomus", 2, true},
		{"Boston", "Chicago", 2, true},
		{"Chicago", "Boston", 2, true},
	}

	for _, v := range testData {
		msg := fmt.Sprintf("'%s' vs '%s'", v.str1, v.str2)
		dist, ab := editdist.ComputeDistanceMax(v.str1, v.str2, 2)
		assert.Equal(t, dist, v.dist, msg)
		assert.Equal(t, ab, v.abort, msg)
	}
}

func TestDiff(t *testing.T) {
	testData := []struct {
		str1, str2 string
		dist       int
		d1, d2     string
	}{
		{"Hello", "He1lo", 1, "He<subst>l</subst>lo", "He<subst>1</subst>lo"},
		{"Pomatomus", "Poma  tomus", 2, "Poma<del>  </del>tomus", "Poma<ins>  </ins>tomus"},
		{"Poma  tomus", "Pomatomus", 2, "Poma<ins>  </ins>tomus", "Poma<del>  </del>tomus"},
		{"Boston", "Chicago", 7, "<del>C</del><subst>Boston</subst>", "<ins>C</ins><subst>hicago</subst>"},
		{"Chicago", "Boston", 7, "<ins>C</ins><subst>hicago</subst>", "<del>C</del><subst>Boston</subst>"},
		{"ebas", "bac", 2, "<ins>e</ins>ba<subst>s</subst>", "<del>e</del>ba<subst>c</subst>"},
		{"rebase", "basic", 4, "<ins>re</ins>bas<del>i</del><subst>e</subst>", "<del>re</del>bas<ins>i</ins><subst>c</subst>"},
		{"test1", "", 5, "<ins>test1</ins>", "<del>test1</del>"},
		{"", "test2", 5, "<del>test2</del>", "<ins>test2</ins>"},
		{"", "", 0, "", ""},
	}

	for _, v := range testData {
		msg := fmt.Sprintf("'%s' vs '%s'", v.str1, v.str2)
		dist, d1, d2 := editdist.ComputeDistance(v.str1, v.str2, true)
		assert.Equal(t, dist, v.dist, msg)
		assert.Equal(t, d1, v.d1, msg)
		assert.Equal(t, d2, v.d2, msg)
	}
}

// BenchmarkDist checks the speed of editdist matching. Run it with:
// `go test -bench=. -benchmem -count=10 -run=XXX > bench.txt && benchstat bench.txt`

func BenchmarkDist(b *testing.B) {
	var out int
	b.Run("CompareOnceMaxOff", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			out, _ = editdist.ComputeDistanceMax("Pomatomus solatror", "Pomatomus saltator", 0)
		}
		_ = fmt.Sprintf("%d\n", out)
	})
	b.Run("CompareOnceMax", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			out, _ = editdist.ComputeDistanceMax("Pomatomus solatror", "Pomatomus saltator", 1)
		}
		_ = fmt.Sprintf("%d\n", out)
	})
	b.Run("CompareDiffOffEqual", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			out, _, _ = editdist.ComputeDistance("Pomatomus saltator", "Pomatomus saltator", false)
		}
		_ = fmt.Sprintf("%d\n", out)
	})
	b.Run("CompareDiffOnEqual", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			out, _, _ = editdist.ComputeDistance("Pomatomus saltator", "Pomatomus saltator", true)
		}
		_ = fmt.Sprintf("%d\n", out)
	})
	b.Run("CompareDiffOff", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			out, _, _ = editdist.ComputeDistance("Pomatomus solatror", "Pomatomus saltator", false)
		}
		_ = fmt.Sprintf("%d\n", out)
	})
	b.Run("CompareDiffOn", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			out, _, _ = editdist.ComputeDistance("Pomatomus solatror", "Pomatomus saltator", true)
		}
		_ = fmt.Sprintf("%d\n", out)
	})
}
