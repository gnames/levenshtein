// Package editdist includes a Levenshtein automaton as well as
// a traditional implementation to calculate Levenshtein Distance.
// The code is based on an excellent levenshtein implementation at
// https://github.com/agnivade/levenshtein
package editdist

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

type eventType uint8

const (
	none eventType = iota
	same
	subst
	ins
	del
)

func (e eventType) String() string {
	switch e {
	case subst:
		return "subst"
	case ins:
		return "ins"
	case del:
		return "del"
	default:
		return ""
	}
}

// ComputeDistanceMax computes the levenshtein distance between the two
// strings passed as an argument. It stops execution if edit distance grows
// a certain max value. It returns edit distance and a boolean. The boolean is
// true when calculation was aborted by the `max` value.
func ComputeDistanceMax(a, b string, max int) (int, bool) {
	if len(a) == 0 {
		dist := utf8.RuneCountInString(b)
		if max > 0 && dist > max {
			return max, true
		}
		return dist, false
	}

	if len(b) == 0 {
		dist := utf8.RuneCountInString(a)
		if max > 0 && dist > max {
			return max, true
		}
		return dist, false
	}

	if a == b {
		return 0, false
	}

	// We need to convert to []rune if the strings are non-ASCII.
	// This could be avoided by using utf8.RuneCountInString
	// and then doing some juggling with rune indices,
	// but leads to far more bounds checks. It is a reasonable trade-off.
	s1 := []rune(a)
	s2 := []rune(b)

	// swap to save some memory O(min(a,b)) instead of O(a)
	if len(s1) > len(s2) {
		s1, s2 = s2, s1
	}
	lenS1 := len(s1)
	lenS2 := len(s2)

	// init the row
	x := make([]uint8, lenS1+1)
	// we start from 1 because index 0 is already 0.
	for i := 1; i < len(x); i++ {
		x[i] = uint8(i)
	}

	// make a dummy bounds check to prevent the 2 bounds check down below.
	// The one inside the loop is particularly costly.
	_ = x[lenS1]
	// fill in the rest
	var rowDist uint8
	for i := 1; i <= lenS2; i++ {
		prev := uint8(i)
		rowDist = 255
		for j := 1; j <= lenS1; j++ {
			current := x[j-1] // match
			if s2[i-1] != s1[j-1] {
				current =
					min(
						min(x[j-1]+1, // substitution
							prev+1), // insertion
						x[j]+1) // deletion
			}
			if current < rowDist || rowDist == 255 {
				rowDist = current
			}
			x[j-1] = prev
			prev = current
		}

		if max > 0 && rowDist > uint8(max) {
			return max, true
		}
		x[lenS1] = prev
	}
	return int(x[lenS1]), false
}

// ComputeDistance computes the levenshtein distance between the two
// strings passed as arguments. The third argument is a flag that
// would trigger creation of tagged strings that show how exactly
// the two strings differ. If the diff artument is true, the tagged
// strings will be provided in the output.
func ComputeDistance(a, b string, diff bool) (int, string, string) {
	if a == b {
		return 0, a, b
	}

	if len(a) == 0 {
		return utf8.RuneCountInString(b),
			"<del>" + b + "</del>",
			"<ins>" + b + "</ins>"
	}

	if len(b) == 0 {
		return utf8.RuneCountInString(a),
			"<ins>" + a + "</ins>",
			"<del>" + a + "</del>"
	}

	// We need to convert to []rune if the strings are non-ASCII.
	// This could be avoided by using utf8.RuneCountInString
	// and then doing some juggling with rune indices,
	// but leads to far more bounds checks. It is a reasonable trade-off.
	s1 := []rune(a)
	s2 := []rune(b)

	lenS1 := len(s1)
	lenS2 := len(s2)

	rl := lenS1 + 1
	cl := lenS2 + 1

	var m []uint8

	if diff {
		m = make([]uint8, 0, cl*rl)
	}

	// init the row
	x := make([]uint8, lenS1+1)
	// we start from 1 because index 0 is already 0.
	for i := 1; i < len(x); i++ {
		x[i] = uint8(i)
	}
	if diff {
		m = append(m, x...)
	}

	// make a dummy bounds check to prevent the 2 bounds check down below.
	// The one inside the loop is particularly costly.
	_ = x[lenS1]
	// fill in the rest
	for i := 1; i <= lenS2; i++ {
		prev := uint8(i)
		for j := 1; j <= lenS1; j++ {
			current := x[j-1] // match
			if s2[i-1] != s1[j-1] {
				current =
					min(
						x[j-1]+1, // substitution (go left)
						min(prev+1, // insertion (go diag)
							x[j]+1), // deletion (go up)
					)
			}
			x[j-1] = prev
			prev = current
		}

		x[lenS1] = prev
		if diff {
			m = append(m, x...)
		}
	}
	var d1, d2 string
	if diff {
		d1, d2 = traceBack(s1, s2, m)
	}
	return int(x[lenS1]), d1, d2
}

// ComputeDistanceTerm comutes edit distance between two strings and
// teturns back edit distance and strings that show the difference
// between compared strings using terminal colors.
func ComputeDistanceTerm(a, b string) (int, string, string) {
	ed, aDiff, bDiff := ComputeDistance(a, b, true)
	aDiff = termColors(aDiff)
	bDiff = termColors(bDiff)
	return ed, aDiff, bDiff
}

func termColors(s string) string {
	rgDel := regexp.MustCompile(`<del>([^<]+)</del>`)
	s = rgDel.ReplaceAllStringFunc(s, func(m string) string {
		txt := m[5 : len(m)-6]
		var res []rune
		for _, v := range txt {
			res = append(res, v)
			res = append(res, '̶')
		}
		return "<del>" + string(res) + "</del>"
	})
	var (
		red    = "\033[1;31m"
		green  = "\033[1;32m"
		yellow = "\033[1;33m"
		end    = "\033[0m"
	)
	s = strings.ReplaceAll(s, "<ins>", green)
	s = strings.ReplaceAll(s, "</ins>", end)
	s = strings.ReplaceAll(s, "<del>", red)
	s = strings.ReplaceAll(s, "</del>", end)
	s = strings.ReplaceAll(s, "<subst>", yellow)
	s = strings.ReplaceAll(s, "</subst>", end)
	return s
}

func min(a, b uint8) uint8 {
	if b < a {
		return b
	}
	return a
}

func traceBack(s1, s2 []rune, m []uint8) (string, string) {
	var e eventType
	var dist, prevDist int
	var iDel, jDel, iIns, jIns, iSubst, jSubst int
	lenS1 := len(s1)
	lenS2 := len(s2)
	rl := lenS1 + 1
	events := make([]eventType, 0, lenS1+lenS2)
	i := lenS2
	j := lenS1
	prevDist = int(m[rl*i+j])
	for !(i == 0 && j == 0) {
		e = same
		iDel, jDel = i-1, j
		iIns, jIns = i, j-1
		iSubst, jSubst = i-1, j-1
		i, j = iSubst, jSubst
		distSubst, distIns, distDel := -1, -1, -1
		if iSubst >= 0 && jSubst >= 0 {
			distSubst = int(m[rl*iSubst+jSubst])
		}
		if jIns >= 0 {
			distIns = int(m[rl*iIns+jIns])
		}
		if iDel >= 0 {
			distDel = int(m[rl*iDel+jDel])
		}
		dist = prevDist
		if distSubst >= 0 && distSubst < dist {
			e = subst
			dist = distSubst
		}
		if distIns >= 0 && distIns < dist {
			e = ins
			i, j = iIns, jIns
			dist = distIns
		}
		if distDel >= 0 && distDel < dist {
			e = del
			i, j = iDel, jDel
			dist = distDel
		}
		prevDist = dist
		events = append(events, e)
	}
	return diffs(s1, s2, events)
}

func diffs(s1, s2 []rune, events []eventType) (string, string) {
	var prev, event eventType
	var deletes, inserts int
	lenS1 := len(s1)
	lenS2 := len(s2)
	d1 := make([]rune, 0, (lenS1+lenS2)*2)
	d2 := make([]rune, 0, (lenS1+lenS2)*2)

	i := 0
	for j := len(events) - 1; j >= 0; j-- {
		event = events[j]
		// init prev event
		if prev == none {
			if event != same {
				d1 = append(d1, []rune("<"+event.String()+">")...)
				d2 = append(d2, []rune("<"+invert(event).String()+">")...)
			}
		} else if event != prev {
			if prev != same {
				d1 = append(d1, []rune("</"+prev.String()+">")...)
				d2 = append(d2, []rune("</"+invert(prev).String()+">")...)
			}
			if event != same {

				d1 = append(d1, []rune("<"+event.String()+">")...)

				d2 = append(d2, []rune("<"+invert(event).String()+">")...)
			}
		}
		switch event {
		case del:
			c1 := s2[i-inserts]
			c2 := s2[i-inserts]
			d1 = append(d1, c1)
			d2 = append(d2, c2)
			deletes++
		case ins:
			c1 := s1[i-deletes]
			c2 := s1[i-deletes]
			d1 = append(d1, c1)
			d2 = append(d2, c2)
			inserts++
		default:
			c1 := s1[i-deletes]
			c2 := s2[i-inserts]
			d1 = append(d1, c1)
			d2 = append(d2, c2)
		}
		prev = event
		i++
	}
	if event != same {
		d1 = append(d1, []rune("</"+event.String()+">")...)
		d2 = append(d2, []rune("</"+invert(event).String()+">")...)
	}
	return string(d1), string(d2)
}

func invert(e eventType) eventType {
	if e == ins {
		return del
	}
	if e == del {
		return ins
	}
	return e
}
