package profiler

import (
	"math"
	"testing"

	sqlv1 "github.com/tkuchiki/logschema/sql/v1"
)

func TestNewQueryRecord(t *testing.T) {
	record, err := newQueryRecord("postgresql", "SELECT N", true, 0.25, sqlv1.QueryData{})
	if err != nil {
		t.Fatal(err)
	}

	if err := record.Validate(); err != nil {
		t.Fatalf("record is invalid: %v", err)
	}
	if got, want := record.Data.Fingerprint.Value, "SELECT N"; got != want {
		t.Fatalf("fingerprint = %q, want %q", got, want)
	}
	if got, want := uint64(record.DurationNano), uint64(250_000_000); got != want {
		t.Fatalf("duration = %d, want %d", got, want)
	}
}

func TestNewIdentityQueryRecordKeepsSensitiveTextExplicit(t *testing.T) {
	record, err := newQueryRecord("mysql", "SELECT * FROM users WHERE id = 42", false, 0.25, sqlv1.QueryData{})
	if err != nil {
		t.Fatal(err)
	}

	if record.Data.QueryText == nil || *record.Data.QueryText != "SELECT * FROM users WHERE id = 42" {
		t.Fatalf("query text = %v", record.Data.QueryText)
	}
	if record.Data.QuerySummary != nil {
		t.Fatalf("identity query has unexpected summary %q", *record.Data.QuerySummary)
	}
}

func TestOptionalMetricsPreserveMissingValues(t *testing.T) {
	if value := optionalDecimal(map[string]uint64{}, "Rows_sent"); value != nil {
		t.Fatalf("missing metric = %d, want nil", *value)
	}
	if value := optionalDecimal(map[string]uint64{"Rows_sent": 0}, "Rows_sent"); value == nil || *value != 0 {
		t.Fatalf("measured zero = %v, want pointer to zero", value)
	}
}

func TestSecondsToNanosecondsRejectsInvalidValues(t *testing.T) {
	for _, value := range []float64{-1, math.Inf(1), math.NaN()} {
		if _, err := secondsToNanoseconds(value); err == nil {
			t.Fatalf("invalid duration %v was accepted", value)
		}
	}
}
