package parser

import (
	"log"
	"reflect"
	"strconv"
)

type AlpData []AlpRow

type AlpRow struct {
	Count    int     `json:"count"`
	Count100 int     `json:"count100"`
	Count200 int     `json:"count200"`
	Count300 int     `json:"count300"`
	Count400 int     `json:"count400"`
	Count500 int     `json:"count500"`
	Method   string  `json:"method"`
	URI      string  `json:"uri"`
	Min      float64 `json:"min"`
	Max      float64 `json:"max"`
	Sum      float64 `json:"sum"`
	Avg      float64 `json:"avg"`
	P90      float64 `json:"p90"`
	P95      float64 `json:"p95"`
	P99      float64 `json:"p99"`
	Stddev   float64 `json:"stddev"`
	MinBody  float64 `json:"minbody"`
	MaxBody  float64 `json:"maxbody"`
	SumBody  float64 `json:"sumbody"`
	AvgBody  float64 `json:"avgbody"`
}

var (
	alpColumns []string
)

func init() {
	alpColumns = []string{
		"count", "count100", "count200", "count300", "count400", "count500", "method", "uri", "min", "max", "sum", "avg", "p90", "p95", "p99", "stddev", "minbody", "maxbody", "sumbody", "avgbody",
	}
}

func ParseAlpData(alpstring string) (AlpData, error) {

}

func ParseAlpRow(alpline string) (AlpRow, bool) {
	i := 0
	k := -1
	substrStart := 0
	substrEnd := 0
	var alpRow AlpRow
	for {
		columnEnd := false
		if string(alpline[i]) == " " {
			i += 1
			columnEnd = true
		}
		if string(alpline[i]) == "|" {
			i += 1
			k += 1
			columnEnd = true
		}

		if columnEnd {
			substrEnd = i
			substr := alpline[substrStart:substrEnd]
			columnName := alpColumns[k]
			field := reflect.ValueOf(alpRow).FieldByName(columnName)
			switch field.Kind() {
			case reflect.Int:
				i, err := strconv.ParseInt(substr, 10, 64)
				if err != nil {
					log.Println(err)
					return alpRow, false
				}
				field.SetInt(i)
			case reflect.String:
				field.SetString(substr)
			case reflect.Float64:
				f, err := strconv.ParseFloat(substr, 64)
				if err != nil {
					log.Println(err)
					return alpRow, false
				}
				field.SetFloat(f)
			default:
				panic("unreachable")
			}
			continue
		}

		if k < 0 {
			return AlpRow{}, false
		}
		if k >= len(alpColumns) {
			break
		}

		if substrStart == substrEnd {
			substrStart = i
			substrEnd = i + 1
		}

		i += 1
	}
}
