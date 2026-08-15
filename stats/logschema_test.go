package stats

import (
	"testing"
	"time"

	corev1 "github.com/tkuchiki/logschema/core/v1"
	sqlv1 "github.com/tkuchiki/logschema/sql/v1"
)

func TestQueryStatsObserveLogSchemaRecord(t *testing.T) {
	record := testQueryRecord(t)
	stats := NewQueryStats(true, true, true, true, true, true)

	stats.Observe(&record)

	if got, want := len(stats.Stats()), 1; got != want {
		t.Fatalf("stats length = %d, want %d", got, want)
	}
	stat := stats.Stats()[0]
	if got, want := stat.Query, "SELECT * FROM users WHERE id = N"; got != want {
		t.Fatalf("query = %q, want %q", got, want)
	}
	if got, want := stat.MaxQueryTime(), 0.125; got != want {
		t.Fatalf("query time = %v, want %v", got, want)
	}
	if got, want := stat.MaxLockTime(), 0.005; got != want {
		t.Fatalf("lock time = %v, want %v", got, want)
	}
	if got, want := stat.MaxRowsExamined(), uint64(42); got != want {
		t.Fatalf("rows examined = %d, want %d", got, want)
	}
}

func TestQueryStatsIgnoresMissingOptionalMetrics(t *testing.T) {
	withMetrics := testQueryRecord(t)
	withoutMetrics, err := sqlv1.NewQuery(250_000_000, corev1.Source{Kind: corev1.SourceOther}, sqlv1.QueryData{
		DBSystem:     "mysql",
		QuerySummary: stringPointer("SELECT * FROM users WHERE id = N"),
		Fingerprint:  sqlv1.Fingerprint{Value: "SELECT * FROM users WHERE id = N", Algorithm: "slp.mysql.abstract", Version: "1"},
	})
	if err != nil {
		t.Fatal(err)
	}

	stats := NewQueryStats(true, true, true, true, true, true)
	stats.Observe(&withMetrics)
	stats.Observe(&withoutMetrics)
	stat := stats.Stats()[0]

	if got, want := stat.AvgLockTime(), 0.005; got != want {
		t.Fatalf("average lock time = %v, want %v", got, want)
	}
	if got, want := stat.PNRowsExamined(100), uint64(42); got != want {
		t.Fatalf("rows examined p100 = %d, want %d", got, want)
	}
	if got := stat.AvgRowsSent(); got != 0 {
		t.Fatalf("average rows sent = %v, want 0", got)
	}
	if got := stat.PNRowsSent(99); got != 0 {
		t.Fatalf("rows sent p99 = %d, want 0", got)
	}
}

func TestNumberStatsPreservesMeasuredZeroAsMinimum(t *testing.T) {
	stats := newNumberStats(false)
	stats.Set(0)
	stats.Set(10)

	if got := stats.Min; got != 0 {
		t.Fatalf("minimum = %d, want 0", got)
	}
}

func TestExpEvalUsesLogSchemaRecord(t *testing.T) {
	record := testQueryRecord(t)
	evaluator, err := NewExpEval(`Query matches "SELECT" && QueryTime == 0.125 && LockTime == 0.005 && RowsExamined == 42`)
	if err != nil {
		t.Fatal(err)
	}

	matched, err := evaluator.Run(&record)
	if err != nil {
		t.Fatal(err)
	}
	if !matched {
		t.Fatal("LogSchema query did not match the filter")
	}
}

func testQueryRecord(t *testing.T) sqlv1.Query {
	t.Helper()

	lockDuration := corev1.DecimalUint64(5 * time.Millisecond)
	rowsExamined := corev1.DecimalUint64(42)
	record, err := sqlv1.NewQuery(125_000_000, corev1.Source{Kind: corev1.SourceOther}, sqlv1.QueryData{
		DBSystem:         "mysql",
		QuerySummary:     stringPointer("SELECT * FROM users WHERE id = N"),
		Fingerprint:      sqlv1.Fingerprint{Value: "SELECT * FROM users WHERE id = N", Algorithm: "slp.mysql.abstract", Version: "1"},
		LockDurationNano: &lockDuration,
		RowsExamined:     &rowsExamined,
	})
	if err != nil {
		t.Fatal(err)
	}

	return record
}

func stringPointer(value string) *string {
	return &value
}
