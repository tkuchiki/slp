package stats

import (
	"time"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	sqlv1 "github.com/tkuchiki/logschema/sql/v1"
)

type ExpEval struct {
	program *vm.Program
}

type ExpEvalEnv struct {
	Query        string
	QueryTime    float64
	LockTime     float64
	RowsSent     uint64
	RowsExamined uint64
	RowsAffected uint64
	BytesSent    uint64
}

func NewExpEval(input string) (*ExpEval, error) {
	program, err := expr.Compile(input, expr.Env(&ExpEvalEnv{}), expr.AsBool())
	if err != nil {
		return nil, err
	}

	return &ExpEval{
		program: program,
	}, nil
}

func (ee *ExpEval) Run(record *sqlv1.Query) (bool, error) {
	env := &ExpEvalEnv{
		Query:        record.Data.Fingerprint.Value,
		QueryTime:    float64(record.DurationNano) / float64(time.Second),
		LockTime:     durationSeconds(record.Data.LockDurationNano),
		RowsSent:     decimalUint64(record.Data.RowsSent),
		RowsExamined: decimalUint64(record.Data.RowsExamined),
		RowsAffected: decimalUint64(record.Data.RowsAffected),
		BytesSent:    decimalUint64(record.Data.BytesSent),
	}

	output, err := expr.Run(ee.program, env)
	if err != nil {
		return false, err
	}

	return output.(bool), nil
}
