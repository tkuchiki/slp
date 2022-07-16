package stats

import (
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	mlog "github.com/percona/go-mysql/log"
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

func (ee *ExpEval) Run(stat *mlog.Event) (bool, error) {
	env := &ExpEvalEnv{
		Query:        stat.Query,
		QueryTime:    stat.TimeMetrics["Query_time"],
		LockTime:     stat.TimeMetrics["Lock_time"],
		RowsSent:     stat.NumberMetrics["Rows_sent"],
		RowsExamined: stat.NumberMetrics["Rows_examined"],
		RowsAffected: stat.NumberMetrics["Rows_affected"],
		BytesSent:    stat.NumberMetrics["Bytes_sent"],
	}

	output, err := expr.Run(ee.program, env)
	if err != nil {
		return false, err
	}

	return output.(bool), nil
}
