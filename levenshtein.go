package levenshtein

import (
	"sync"

	"github.com/gnames/gnlib/gnuuid"
	"github.com/gnames/levenshtein/presenter"
	u "github.com/google/uuid"
)

var Batch = 10_000
var jobs = 16
var Version = "v.0.0.0"

type Option func(*levenshtein)

func OptWithDiff(b bool) Option {
	return func(l *levenshtein) {
		l.withDiff = b
	}
}

type levenshtein struct {
	withDiff bool
}

func NewLevenshtein(opts ...Option) levenshtein {
	return levenshtein{}
}

func (fd fuzzyDiff) Compare(str1, str2 string) output.Output {
	outRaw := fd.Differ.Compare([]rune(str1), []rune(str2))
	return output.NewOutput(outRaw)
}

func (fd fuzzyDiff) Opts() []differ.Option {
	return fd.Differ.Opts()
}

func (fd fuzzyDiff) CompareMult(inp []input.Strings) []output.Output {
	outMap := make(map[u.UUID]output.Output)
	chIn := make(chan input.Strings)
	chOut := make(chan output.Output)
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
		go fd.compareWorker(chIn, chOut, &wgWorker)
	}
	go func() {
		defer wgFin.Done()
		for out := range chOut {
			outMap[getUUID(out.Str1, out.Str2)] = out
		}
	}()

	wgWorker.Wait()
	close(chOut)

	wgFin.Wait()

	res := make([]output.Output, len(inp))
	for i, v := range inp {
		res[i] = outMap[getUUID(v.String1, v.String2)]
	}
	return res
}

func (fd fuzzyDiff) compareWorker(
	chIn <-chan input.Strings,
	chOut chan<- output.Output,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	fd = NewFuzzyDiff(fd.Opts()...)
	for v := range chIn {
		chOut <- fd.Compare(v.String1, v.String2)
	}
}

func getUUID(str1, str2 string) u.UUID {
	return gnuuid.New(str1 + "\v|\v" + str2)
}
