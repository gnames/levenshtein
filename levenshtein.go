package levenshtein

import (
	"sync"

	"github.com/gnames/gnlib/gnuuid"
	"github.com/gnames/levenshtein/entity/editdist"
	"github.com/gnames/levenshtein/presenter"
	u "github.com/google/uuid"
)

// Version of the package.
var Version = "v.0.1.0"

// Batch can be used to break a very large input into chunks.
var Batch = 10_000

// job is a hardcoded number of jobs used to parallelize CompareMult
// jobs. 16 seem to be a reasonable number, that would not slow things
// down even on a 1 processor machine, and more than 16 go-routines
// seem to show diminishing returns.
var jobs = 16

// Option is an 'interface' for creating Options for Levensthsein,
// that modify its behavior.
type Option func(*levenshtein)

// OptWithDiff if set to true will return restuls where difference between
// strings will be marked by "ins", "del", "subst" tags. The option will
// slowdown calculations 4-5 times.
func OptWithDiff(b bool) Option {
	return func(l *levenshtein) {
		l.withDiff = b
	}
}

// OptMaxEditDist sets the maximum edit distance after which the progression
// of calculations will be aborted. Such constraint does not seem to be
// always beneficial.
func OptMaxEditDist(i int) Option {
	return func(l *levenshtein) {
		l.maxEditDist = i
	}
}

// levenshtein is an implementation of Levenshtein interface.
type levenshtein struct {
	withDiff    bool
	maxEditDist int
}

// NewLevenshtein returns an object that implements Levenshtein
// interface.
func NewLevenshtein(opts ...Option) levenshtein {
	l := levenshtein{}
	for _, opt := range opts {
		opt(&l)
	}
	return l
}

// Strings structure is used to feed CompareMult function with data.
type Strings struct {
	String1, String2 string
}

// Compare is an implementation of Levenshtein interface.
func (l levenshtein) Compare(str1, str2 string) presenter.Output {
	var ed int
	var t1, t2 string
	var aborted bool
	if l.maxEditDist > 0 {
		ed, aborted = editdist.ComputeDistanceMax(str1, str2, l.maxEditDist)
	}

	if !aborted {
		ed, t1, t2 = editdist.ComputeDistance(str1, str2, l.withDiff)
	}

	return presenter.Output{
		String1:  str1,
		String2:  str2,
		Tags1:    t1,
		Tags2:    t2,
		EditDist: ed,
		Aborted:  aborted,
	}
}

// Opts is an implementation of Levenshtein interface.
func (l levenshtein) Opts() []Option {
	return []Option{OptWithDiff(l.withDiff)}
}

// CompareMult is an implementation of Levenshtein interface.
func (l levenshtein) CompareMult(inp []Strings) []presenter.Output {
	outMap := make(map[u.UUID]presenter.Output)
	chIn := make(chan Strings)
	chOut := make(chan presenter.Output)
	var wgWorker sync.WaitGroup
	wgWorker.Add(jobs)
	var wgFin sync.WaitGroup
	wgFin.Add(1)

	go func() {
		for _, v := range inp {
			chIn <- v
		}
		close(chIn)
	}()
	for i := 0; i < jobs; i++ {
		go l.compareWorker(chIn, chOut, &wgWorker)
	}
	go func() {
		defer wgFin.Done()
		for out := range chOut {
			outMap[getUUID(out.String1, out.String2)] = out
		}
	}()

	wgWorker.Wait()
	close(chOut)

	wgFin.Wait()

	res := make([]presenter.Output, len(inp))
	for i, v := range inp {
		res[i] = outMap[getUUID(v.String1, v.String2)]
	}
	return res
}

func (l levenshtein) compareWorker(
	chIn <-chan Strings,
	chOut chan<- presenter.Output,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	l = NewLevenshtein(l.Opts()...)
	for v := range chIn {
		chOut <- l.Compare(v.String1, v.String2)
	}
}

func getUUID(str1, str2 string) u.UUID {
	return gnuuid.New(str1 + "\v|\v" + str2)
}
