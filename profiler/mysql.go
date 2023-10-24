package profiler

import (
	"fmt"
	"os"

	"github.com/tkuchiki/slp/stats"

	"github.com/percona/go-mysql/log"
	"github.com/percona/go-mysql/log/slow"
	"github.com/tkuchiki/slp/abstractor"
	"github.com/tkuchiki/slp/mysql/parser/sqlparser"
	"github.com/tkuchiki/slp/options"
)

type MySQLProfiler struct {
	sqlp      *sqlparser.MySQLSQLParser
	slowp     *slow.SlowLogParser
	abst      *abstractor.MySQLAbstractor
	opts      *options.Options
	readBytes uint64
}

func NewMySQLProfiler(file *os.File, logopt log.Options, opts *options.Options) *MySQLProfiler {
	return &MySQLProfiler{
		sqlp:  sqlparser.NewMySQLSQLParser(),
		slowp: slow.NewSlowLogParser(file, logopt),
		abst:  abstractor.NewMySQLAbstractor(opts.BundleWhereIn, opts.BundleValues),
		opts:  opts,
	}
}

func (p *MySQLProfiler) Profile(qstats *stats.QueryStats) error {
	defer p.slowp.Stop()
	go p.slowp.Start()

	for e := range p.slowp.EventChan() {
		if !p.opts.NoAbstract {
			stmt, err := p.sqlp.Parse(e.Query)
			if err != nil {
				continue
			}

			query, err := p.abst.Abstract(stmt)
			if err != nil {
				continue
			}

			e.Query = query
		}

		metrics := &stats.QueryMetrics{
			Query:        e.Query,
			QueryTime:    e.TimeMetrics["Query_time"],
			LockTime:     e.TimeMetrics["Lock_time"],
			RowsSent:     e.NumberMetrics["Rows_sent"],
			RowsExamined: e.NumberMetrics["Rows_examined"],
			RowsAffected: e.NumberMetrics["Rows_affected"],
			BytesSent:    e.NumberMetrics["Bytes_sent"],
		}

		matched, err := qstats.DoFilter(metrics)
		if err != nil {
			return err
		}

		if !matched {
			continue
		}

		qstats.Set(e.Query, e.TimeMetrics["Query_time"], e.TimeMetrics["Lock_time"], e.NumberMetrics["Rows_sent"], e.NumberMetrics["Rows_examined"], e.NumberMetrics["Rows_affected"], e.NumberMetrics["Bytes_sent"])

		if qstats.CountUris() > p.opts.Limit {
			return fmt.Errorf("Too many Queries (%d or less)", p.opts.Limit)
		}

		p.readBytes = e.OffsetEnd
	}

	return nil
}

func (p *MySQLProfiler) ReadBytesInt64() int64 {
	return int64(p.readBytes)
}
