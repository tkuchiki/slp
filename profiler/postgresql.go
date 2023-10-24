package profiler

import (
	"fmt"
	"os"
	"time"

	slowlogparser "github.com/tkuchiki/go-pg-slowlog/parser"
	"github.com/tkuchiki/slp/abstractor"
	"github.com/tkuchiki/slp/options"
	"github.com/tkuchiki/slp/stats"
)

type PGProfiler struct {
	slowp     *slowlogparser.PGSlowLogParser
	abst      *abstractor.PGAbstractor
	opts      *options.Options
	readBytes uint64
}

func NewPGProfiler(file *os.File, opts *options.Options) (*PGProfiler, error) {
	slowlogParser, err := slowlogparser.NewPGSlowLogParser(file, opts.LogLinePrefix)
	if err != nil {
		return nil, err
	}

	return &PGProfiler{
		slowp: slowlogParser,
		abst:  abstractor.NewPGAbstractor(opts.NoAbstract, opts.BundleWhereIn, opts.BundleValues),
		opts:  opts,
	}, nil
}

func (p *PGProfiler) Profile(qstats *stats.QueryStats) error {
	defer p.slowp.Stop()
	go p.slowp.Start()

	for logEntry := range p.slowp.LogEntryChan() {
		var query string
		if !p.opts.NoAbstract {
			q, err := p.abst.Abstract(logEntry.Statement)
			if err != nil {
				continue
			}

			query = q
		}

		metrics := &stats.QueryMetrics{
			Query:     query,
			QueryTime: float64(logEntry.Duration) / float64(time.Second),
		}

		matched, err := qstats.DoFilter(metrics)
		if err != nil {
			return err
		}

		if !matched {
			continue
		}

		qstats.Set(query, metrics.QueryTime, 0, 0, 0, 0, 0)

		if qstats.CountUris() > p.opts.Limit {
			return fmt.Errorf("Too many Queries (%d or less)", p.opts.Limit)
		}
	}

	p.readBytes = uint64(p.slowp.ReadBytes())

	return nil
}

func (p *PGProfiler) ReadBytesInt64() int64 {
	return int64(p.readBytes)
}
