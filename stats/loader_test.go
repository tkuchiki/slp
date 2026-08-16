package stats

import (
	"strings"
	"testing"
)

func TestLoadStatsRestoresLegacySampleCounts(t *testing.T) {
	const dump = `- query: SELECT N
  count: 2
  query_time:
    max: 0.2
    min: 0.1
    sum: 0.3
    usepercentile: false
    percentiles: []
  lock_time:
    max: 0.01
    min: 0
    sum: 0.01
    usepercentile: false
    percentiles: []
  rows_sent:
    max: 0
    min: 0
    sum: 0
    usepercentile: false
    percentiles: []
  rows_examined:
    max: 10
    min: 0
    sum: 10
    usepercentile: false
    percentiles: []
  rows_affected:
    max: 0
    min: 0
    sum: 0
    usepercentile: false
    percentiles: []
  bytes_sent:
    max: 0
    min: 0
    sum: 0
    usepercentile: false
    percentiles: []
`

	stats := NewQueryStats(false, false, false, false, false, false)
	if err := stats.LoadStats(strings.NewReader(dump)); err != nil {
		t.Fatal(err)
	}

	stat := stats.Stats()[0]
	if got, want := stat.AvgLockTime(), 0.005; got != want {
		t.Fatalf("legacy lock time average = %v, want %v", got, want)
	}
	if got, want := stat.AvgRowsExamined(), float64(5); got != want {
		t.Fatalf("legacy rows examined average = %v, want %v", got, want)
	}
}
