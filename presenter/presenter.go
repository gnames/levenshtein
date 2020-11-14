package presenter

import (
	"strconv"

	gncsv "github.com/gnames/gnlib/csv"
	"github.com/gnames/gnlib/encode"
	"github.com/gnames/gnlib/format"
)

type Output struct {
	String1  string `json:"string1"`
	String2  string `json:"string2"`
	Tags1    string `json:"tags1,omitempty"`
	Tags2    string `json:"tags2,omitempty"`
	EditDist int    `json:"editDistance"`
	Aborted  bool   `json:"aborted,omitempty"`
}

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
