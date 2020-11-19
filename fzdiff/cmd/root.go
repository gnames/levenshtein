/*
Copyright © 2020 Dmitry Mozzherin <dmozzherin@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gnames/gnlib/format"
	"github.com/gnames/gnlib/sys"
	"github.com/gnames/levenshtein"
	"github.com/gnames/levenshtein/presenter"
	"github.com/spf13/cobra"
	"gitlab.com/gogna/gnparser/output"
)

var opts []levenshtein.Option

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fzdiff",
	Short: "Finds Levenshtein edit distance and tags differences between strings.",
	Long: `fzdiff takes two strings and returns back edit distance between them
  according to Levenshtein algorithm. It has an option to abort progress after
  reaching a certain edit distance. Also it may insrt tags into strings showing
  where did edit events happen for either input.`,
	Run: func(cmd *cobra.Command, args []string) {
		if showVersionFlag(cmd) {
			os.Exit(0)
		}

		formatString, _ := cmd.Flags().GetString("format")
		frmt, _ := format.NewFormat(formatString)
		if frmt == format.FormatNone {
			log.Printf("Cannot set format from '%s', setting format to csv.",
				formatString)
			frmt = format.CSV
		}

		withTags, _ := cmd.Flags().GetBool("tags")
		opts = append(opts, levenshtein.OptWithDiff(withTags))

		maxEditDist, _ := cmd.Flags().GetInt("max_edit_distance")
		opts = append(opts, levenshtein.OptMaxEditDist(maxEditDist))

		l := levenshtein.NewLevenshtein(opts...)

		if len(args) == 0 {
			processStdin(cmd, l, frmt)
			os.Exit(0)
		}

		data := getInput(cmd, args)
		compare(l, data, frmt)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("version", "V", false, "Prints version information")
	rootCmd.Flags().BoolP("tags", "t", false, "Adds diff tags into strings.")
	rootCmd.Flags().IntP("max_edit_distance", "m", 0, "Max threshold for edit distance.")
	rootCmd.Flags().StringP("format", "f", "csv", `Format of the output: "compact", "pretty", "csv".
  compact: compact JSON,
  pretty: pretty JSON,
  csv: CSV (DEFAULT)`)
}

// showVersionFlag provides version and the build timestamp. If it returns
// true, it means that version flag was given.
func showVersionFlag(cmd *cobra.Command) bool {
	hasVersionFlag, _ := cmd.Flags().GetBool("version")

	if hasVersionFlag {
		fmt.Printf("\nversion: %s\n\n", levenshtein.Version)
	}
	return hasVersionFlag
}

func processStdin(cmd *cobra.Command, l levenshtein.Levenshtein,
	frmt format.Format) {
	if !checkStdin() {
		_ = cmd.Help()
		os.Exit(0)
	}
	compareFile(l, os.Stdin, frmt)
}

func checkStdin() bool {
	stdInFile := os.Stdin
	stat, err := stdInFile.Stat()
	if err != nil {
		log.Panic(err)
	}
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func getInput(cmd *cobra.Command, args []string) []string {
	var data []string
	switch len(args) {
	case 1:
		data = []string{args[0]}
	case 2:
		data = []string{args[0], args[1]}
	default:
		_ = cmd.Help()
		os.Exit(0)
	}
	return data
}

func compare(l levenshtein.Levenshtein, data []string,
	frmt format.Format) {
	if len(data) == 1 {
		path := string(data[0])
		if sys.FileExists(path) {
			f, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			compareFile(l, f, frmt)
			f.Close()
			os.Exit(0)
		}
	}

	compareStrings(l, data, frmt)
}

func compareFile(l levenshtein.Levenshtein, f io.Reader, frmt format.Format) {
	batch := make([]levenshtein.Strings, 0, levenshtein.Batch)
	if frmt == format.CSV {
		fmt.Println(output.CSVHeader())
	}
	r := csv.NewReader(f)

	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Cannot read CSV file: %s", err)
		}
		if len(row) < 2 {
			log.Fatalf("There is less than 2 strings in a row %+v: %s", row, err)
		}
		if len(batch) < levenshtein.Batch {
			batch = append(batch, levenshtein.Strings{String1: row[0], String2: row[1]})
		} else {
			out := l.CompareMult(batch)
			for _, v := range out {
				res, err := v.Encode(frmt)
				if err != nil {
					log.Fatalf("cannot encode %s: %s", frmt.String(), err)
				}
				fmt.Println(res)
			}
			batch = batch[:0]
		}
	}
	out := l.CompareMult(batch)
	for _, v := range out {
		res, err := v.Encode(frmt)
		if err != nil {
			log.Fatalf("cannot encode %s: %s", frmt.String(), err)
		}
		fmt.Println(res)
	}
}

func compareStrings(l levenshtein.Levenshtein, data []string,
	frmt format.Format) {
	out := l.Compare(data[0], data[1])
	res, err := out.Encode(frmt)
	if err != nil {
		log.Fatal(err)
	}
	if frmt == format.CSV {
		fmt.Println(presenter.CSVHeader())
	}
	fmt.Println(res)
}
