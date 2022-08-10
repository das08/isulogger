package parser

import (
	"errors"
	"log"
	"reflect"
	"strconv"
	"strings"
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
		"Count", "Count100", "Count200", "Count300", "Count400", "Count500", "Method", "URI", "Min", "Max", "Sum", "Avg", "P90", "P95", "P99", "Stddev", "MinBody", "MaxBody", "SumBody", "AvgBody",
	}
}

func ParseAlpData(alpstring string) (AlpData, error) {
	lines := strings.Split(alpstring, "\n")
	alpdata := make(AlpData, 0)
	for _, line := range lines[3:] {
		row, ok := ParseAlpRow(line)
		if !ok {
			errors.New("failed to parse")
		}
		alpdata = append(alpdata, row)
	}
	return alpdata, nil
}

func ParseAlpRow(alpline string) (AlpRow, bool) {
	var row AlpRow
	split := strings.Split(alpline, "|")

	prrow := reflect.ValueOf(&row)
	rrow := prrow.Elem()

	for i, columnName := range alpColumns {
		if i+1 >= len(split) {
			return AlpRow{}, false
		}
		columnStr := split[i+1]
		columnStr = strings.TrimSpace(columnStr)
		field := rrow.FieldByName(columnName)
		switch field.Kind() {
		case reflect.Int:
			if val, err := strconv.Atoi(columnStr); err == nil {
				field.SetInt(int64(val))
			} else {
				log.Printf("Error parsing %s: %s", columnName, err)
				return AlpRow{}, false
			}
		case reflect.Float64:
			if val, err := strconv.ParseFloat(columnStr, 64); err == nil {
				field.SetFloat(val)
			} else {
				log.Printf("Error parsing %s: %s", columnName, err)
				return AlpRow{}, false
			}
		case reflect.String:
			field.SetString(columnStr)
		default:
			ty := field.Kind().String()
			panic("unreachable: " + ty + " (" + columnName + ")")
		}
	}
	return row, true
}
