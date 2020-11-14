package levenshtein

import (
	"sync"

	"github.com/gnames/gnlib/gnuuid"
	"github.com/gnames/levenshtein/entity/editdist"
	"github.com/gnames/levenshtein/presenter"
	u "github.com/google/uuid"
)

var Version = "v.0.0.0"
var Batch = 10_000
var jobs = 16

type Option func(*levenshtein)

func OptWithDiff(b bool) Option {
	return func(l *levenshtein) {
		l.withDiff = b
	}
}

func OptMaxEditDist(i int) Option {
	return func(l *levenshtein) {
		l.maxEditDist = i
	}
}

type levenshtein struct {
	withDiff    bool
	maxEditDist int
}

func NewLevenshtein(opts ...Option) levenshtein {
	l := levenshtein{}
	for _, opt := range opts {
		opt(&l)
	}
	return l
}

type Strings struct {
	String1, String2 string
}

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

func (l levenshtein) Opts() []Option {
	return []Option{OptWithDiff(l.withDiff)}
}

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
