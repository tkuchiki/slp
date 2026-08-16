package profiler

import (
	"fmt"
	"os"

	"github.com/tkuchiki/slp/stats"

	"github.com/percona/go-mysql/log"
	"github.com/percona/go-mysql/log/slow"
	sqlv1 "github.com/tkuchiki/logschema/sql/v1"
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

		lockDuration, err := optionalDuration(e.TimeMetrics, "Lock_time")
		if err != nil {
			continue
		}

		record, err := newQueryRecord("mysql", e.Query, !p.opts.NoAbstract, e.TimeMetrics["Query_time"], sqlv1.QueryData{
			LockDurationNano: lockDuration,
			RowsSent:         optionalDecimal(e.NumberMetrics, "Rows_sent"),
			RowsExamined:     optionalDecimal(e.NumberMetrics, "Rows_examined"),
			RowsAffected:     optionalDecimal(e.NumberMetrics, "Rows_affected"),
			BytesSent:        optionalDecimal(e.NumberMetrics, "Bytes_sent"),
		})
		if err != nil {
			continue
		}

		matched, err := qstats.DoFilter(record)
		if err != nil {
			return err
		}

		if !matched {
			continue
		}

		qstats.Observe(record)

		if qstats.CountQueries() > p.opts.Limit {
			return fmt.Errorf("Too many Queries (%d or less)", p.opts.Limit)
		}

		p.readBytes = e.OffsetEnd
	}

	return nil
}

func (p *MySQLProfiler) ReadBytesInt64() int64 {
	return int64(p.readBytes)
}
