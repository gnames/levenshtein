# levenshtein

The package calculates edit distance between two strings according to
Levenshtein algorithm. It also can show differences beteen two strings. This
package extends an [excellent levenshtein
implementation][agnivade levenshtsin](https://github.com/agnivade/levenshtein) by @agnivade.

## Installation

### Installation of a binary file

You can download binary `fzdiff` from the [latest release] and install it
somewhere in your PATH.

### Installation with Go

If you have Go installed on your system, run:

```bash
go get github.com/gnames/levenshtein/fzdiff
```

## Usage

### Command line interface

- Run `fzdiff` on two strings.

    ```bash
    fzdiff "Something" "smoething"
    # output:
    String1,String2,Tags1,Tags2,EditDistance,Aborted
    Something,smoething,,,3,false
    ```

- Change output.

    ```bash
    fzdiff "Something" "smoething" -f compact
    {"string1":"Something","string2":"smoething","editDistance":3}


    fzdiff "Something" "smoething" -f pretty
    {
      "string1": "Something",
      "string2": "smoething",
      "editDistance": 3
    }
    ```

- Run `fzdiff` with max edit distance constraint:

    ```bash
    fzdiff "Something" "smoething" -m 1
    String1,String2,Tags1,Tags2,EditDistance,Aborted
    Something,smoething,,,1,true
    ```

- Run `fzdiff` with tags output:

    ```bash
    fzdiff "Something" "smoething" -t
    String1,String2,Tags1,Tags2,EditDistance,Aborted
    Something,smoething,<subst>Som</subst>ething,<subst>smo</subst>ething,3,false
    ```

- Run `fzdiff` on a CSV file to compare the first 2 fields.

    ```bash
    fzdiff strings.csv -t
    ```

- Run `fzdiff` using pipes with STDIN, STDOUT

    ```bash
    echo "Something,smoething" | fzdiff -t
    Id,Verbatim,Cardinality,CanonicalFull,CanonicalSimple,CanonicalStem,Authorship,Year,Quality
    Something,smoething,<subst>Som</subst>ething,<subst>smo</subst>ething,3,false

    # or

    cat strings.csv | fzdiff -t > diffs.csv
    ```

### Usage as a library

```go
package main

import (
 "fmt"

 "github.com/gnames/levenshtein"
)

func main() {
 opts := []levenshtein.Option{
  levenshtein.OptWithDiff(true),
  levenshtein.OptMaxEditDist(1),
 }
 l := levenshtein.NewLevenshtein(opts...)
 out := l.Compare("Something", "smoething")
 fmt.Printf("%+v\n", out)

 strs := []levenshtein.Strings{
  {String1: "Something", String2: "smoething"},
  {String1: "one", String2: "two"},
 }

 outs := l.CompareMult(strs)
 for _, out := range outs {
  fmt.Printf("%+v\n", out)
 }
}
```

## Testing

From the `root` of the project:

```bash
go test ./...
```

### Benchmarking

You need to install [benchstat] for more readable restuls.

```bash
cd entity/editdist
go test -bench=. -benchmem -count=10 -run=XXX > bench.txt
benchstat bench.txt
```

An example of the benchmarking:

```bash
cd entity/editdist
benchstat bench.txt

name                        time/op
Dist/CompareOnceMaxOff-8     740ns ± 1%
Dist/CompareOnceMax-8        580ns ± 1%
Dist/CompareDiffOffEqual-8  3.93ns ± 1%
Dist/CompareDiffOnEqual-8   3.85ns ± 5%
Dist/CompareDiffOff-8        488ns ± 1%
Dist/CompareDiffOn-8        2.37µs ± 6%

name                        alloc/op
Dist/CompareOnceMaxOff-8     32.0B ± 0%
Dist/CompareOnceMax-8        32.0B ± 0%
Dist/CompareDiffOffEqual-8   0.00B
Dist/CompareDiffOnEqual-8    0.00B
Dist/CompareDiffOff-8        32.0B ± 0%
Dist/CompareDiffOn-8        1.17kB ± 0%

name                        allocs/op
Dist/CompareOnceMaxOff-8      1.00 ± 0%
Dist/CompareOnceMax-8         1.00 ± 0%
Dist/CompareDiffOffEqual-8    0.00
Dist/CompareDiffOnEqual-8     0.00
Dist/CompareDiffOff-8         1.00 ± 0%
Dist/CompareDiffOn-8          7.00 ± 0%
```

## License

Released under [MIT license]

[agnivade levenshtsin]: https://github.com/agnivade/levenshtein
[latest release]: https://github.com/gnames/levenshtein/releases/latest
[benchstat]: https://github.com/golang/perf/tree/master/cmd/benchstat
[MIT license]: https://raw.githubusercontent.com/gnames/levenshtein/master/LICENSE

