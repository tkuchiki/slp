package stats

type QueryMetrics struct {
	Query        string
	QueryTime    float64
	LockTime     float64
	RowsSent     uint64
	RowsExamined uint64
	RowsAffected uint64
	BytesSent    uint64
}
