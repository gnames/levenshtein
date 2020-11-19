package presenter

import (
	"strconv"

	gncsv "github.com/gnames/gnlib/csv"
	"github.com/gnames/gnlib/encode"
	"github.com/gnames/gnlib/format"
)

// Output is a representation of edit distance calculation results.
type Output struct {
	// String1 is the first input string.
	String1 string `json:"string1"`
	// String2 is the second input string.
	String2 string `json:"string2"`
	// Tags1, would contain tagged version of the first string
	// that will show where it differs from the second string.
	Tags1 string `json:"tags1,omitempty"`
	// Tags2, would contain tagged version of the second string
	// that will show where it differs from the first string.
	Tags2 string `json:"tags2,omitempty"`
	// EditDist is the calculated edit distance according to
	// Levenshtein algorithm. In case if Maximum Edit Distance is
	// set it might not show the actual edit distance between two
	// strings.
	EditDist int `json:"editDistance"`
	// Aborted is true if Maximum Edit Distance is provided, and
	// it was exceeded during calculations.
	Aborted bool `json:"aborted,omitempty"`
}

// Encode method produces representation of Output for consumption
// either by a CLI user, or a WEB client. It supports 3 possible formats:
// CSV, JSON (pretty), JSON compact.
func (o Output) Encode(f format.Format) (string, error) {
	switch f {
	case format.CSV:
		return o.encodeCSV()
	case format.PrettyJSON:
		return o.encodeJSON(true)
	case format.CompactJSON:
		return o.encodeJSON(false)
	default:
		return o.encodeCSV()
	}
}

// CSVHeader produces a CSV header compatible with CSV output
// of Output's Encode method.
func CSVHeader() string {
	return "String1,String2,Tags1,Tags2,EditDistance,Aborted"
}

func (o Output) encodeCSV() (string, error) {
	row := []string{o.String1, o.String2, o.Tags1, o.Tags2,
		strconv.Itoa(o.EditDist), strconv.FormatBool(o.Aborted),
	}
	return gncsv.ToCSV(row), nil
}

func (o Output) encodeJSON(pretty bool) (string, error) {
	enc := encode.GNjson{Pretty: pretty}
	res, err := enc.Encode(o)
	return string(res), err
}
