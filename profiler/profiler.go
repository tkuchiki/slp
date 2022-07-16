package profiler

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/percona/go-mysql/log"
	"github.com/percona/go-mysql/log/slow"
	"github.com/tkuchiki/slp/abstractor"
	"github.com/tkuchiki/slp/helper"
	"github.com/tkuchiki/slp/options"
	"github.com/tkuchiki/slp/parser/sqlparser"
	"github.com/tkuchiki/slp/stats"
)

type Profiler struct {
	sqlp      *sqlparser.SQLParser
	slowp     *slow.SlowLogParser
	abst      *abstractor.SQLAbstractor
	opts      *options.Options
	readBytes uint64
}

func NewProfiler(file *os.File, logopt log.Options, opts *options.Options) *Profiler {
	return &Profiler{
		sqlp:  sqlparser.NewSQLParser(),
		slowp: slow.NewSlowLogParser(file, logopt),
		abst:  abstractor.NewSQLAbstractor(opts.BundleWhereIn, opts.BundleValues),
		opts:  opts,
	}
}

func (p *Profiler) OpenPosFile(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
}

func (p *Profiler) ReadPosFile(f *os.File) (int64, error) {
	reader := bufio.NewReader(f)
	pos, _, err := reader.ReadLine()
	if err != nil {
		return 0, err
	}

	return helper.StringToInt64(string(pos))
}

func (p *Profiler) ReadBytes() int64 {
	return int64(p.readBytes)
}

func (p *Profiler) SetReadBytes(n int64) {
	p.readBytes = uint64(n)
}

func (p *Profiler) Seek(file *os.File, n int64) error {
	_, err := file.Seek(n, io.SeekCurrent)
	return err
}

func (p *Profiler) Profile(qstats *stats.QueryStats) error {
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

		matched, err := qstats.DoFilter(e)
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
