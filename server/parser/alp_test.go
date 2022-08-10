package parser

import "testing"

func TestParseAlpRow(t *testing.T) {
	line := `|     1 |   0 |    1 |   0 |   0 |   0 | GET    | /api/organizer/competitions                  | 0.000 |  0.000 |    0.000 | 0.000 | 0.000 | 0.000 |  0.000 |  0.000 |   177.000 |    177.000 |      177.000 |   177.000 |`
	row, ok := ParseAlpRow(line)
	if !ok {
		t.Error("Expected to parse row")
	}
	if row.Count != 1 {
		t.Error("Expected count to be 1")
	}
	if row.Count100 != 0 {
		t.Error("Expected count100 to be 0")
	}
	if row.Count200 != 1 {
		t.Error("Expected count200 to be 1")
	}
	if row.Count300 != 0 {
		t.Error("Expected count300 to be 0")
	}
	if row.Count400 != 0 {
		t.Error("Expected count400 to be 0")
	}
	if row.Count500 != 0 {
		t.Error("Expected count500 to be 0")
	}
	if row.Method != "GET" {
		t.Error("Expected method to be GET")
	}
	if row.URI != "/api/organizer/competitions" {
		t.Error("Expected uri to be /api/organizer/competitions")
	}
	if row.Min != 0.000 {
		t.Error("Expected min to be 0.000")
	}
	if row.Max != 0.000 {
		t.Error("Expected max to be 0.000")
	}
	if row.Sum != 0.000 {
		t.Error("Expected sum to be 0.000")
	}
	if row.Avg != 0.000 {
		t.Error("Expected avg to be 0.000")
	}
	if row.P90 != 0.000 {
		t.Error("Expected p90 to be 0.000")
	}
	if row.P95 != 0.000 {
		t.Error("Expected p95 to be 0.000")
	}
	if row.P99 != 0.000 {
		t.Error("Expected p99 to be 0.000")
	}
	if row.Stddev != 0.000 {
		t.Error("Expected stddev to be 0.000")
	}
	if row.MinBody != 177.000 {
		t.Error("Expected minbody to be 177.000")
	}
	if row.MaxBody != 177.000 {
		t.Error("Expected maxbody to be 177.000")
	}
	if row.SumBody != 177.000 {
		t.Error("Expected sumbody to be 177.000")
	}
	if row.AvgBody != 177.000 {
		t.Error("Expected avgbody to be 177.000")
	}
}
