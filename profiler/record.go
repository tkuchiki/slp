package profiler

import (
	"fmt"
	"math"
	"time"

	corev1 "github.com/tkuchiki/logschema/core/v1"
	sqlv1 "github.com/tkuchiki/logschema/sql/v1"
)

const fingerprintVersion = "1"

func newQueryRecord(dbSystem, query string, abstracted bool, durationSeconds float64, data sqlv1.QueryData) (*sqlv1.Query, error) {
	durationNano, err := secondsToNanoseconds(durationSeconds)
	if err != nil {
		return nil, err
	}

	data.DBSystem = dbSystem
	fingerprintAlgorithm := "identity"
	if abstracted {
		data.QuerySummary = stringPointer(query)
		fingerprintAlgorithm = "slp." + dbSystem + ".abstract"
	} else {
		data.QueryText = stringPointer(query)
	}
	data.Fingerprint = sqlv1.Fingerprint{
		Value:     query,
		Algorithm: fingerprintAlgorithm,
		Version:   fingerprintVersion,
	}

	record, err := sqlv1.NewQuery(durationNano, corev1.Source{Kind: corev1.SourceOther}, data)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

func secondsToNanoseconds(seconds float64) (corev1.DecimalUint64, error) {
	if math.IsNaN(seconds) || math.IsInf(seconds, 0) || seconds < 0 {
		return 0, fmt.Errorf("query duration must be a finite non-negative number")
	}

	nanoseconds := math.Round(seconds * float64(time.Second))
	if nanoseconds >= math.Exp2(64) {
		return 0, fmt.Errorf("query duration exceeds the LogSchema duration range")
	}

	return corev1.DecimalUint64(uint64(nanoseconds)), nil
}

func durationPointer(seconds float64) (*corev1.DecimalUint64, error) {
	value, err := secondsToNanoseconds(seconds)
	if err != nil {
		return nil, err
	}

	return &value, nil
}

func decimalPointer(value uint64) *corev1.DecimalUint64 {
	converted := corev1.DecimalUint64(value)

	return &converted
}

func optionalDuration(metrics map[string]float64, key string) (*corev1.DecimalUint64, error) {
	value, ok := metrics[key]
	if !ok {
		return nil, nil
	}

	return durationPointer(value)
}

func optionalDecimal(metrics map[string]uint64, key string) *corev1.DecimalUint64 {
	value, ok := metrics[key]
	if !ok {
		return nil
	}

	return decimalPointer(value)
}

func stringPointer(value string) *string {
	return &value
}
