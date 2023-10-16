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

func (p *Profiler) ReadPositionFile(posFilePath string, slowlogFile *os.File) (*os.File, int64, error) {
	if posFilePath == "" {
		return nil, 0, nil
	}

	posFile, err := p.OpenPosFile(posFilePath)
	if err != nil {
		return nil, 0, err
	}

	position, err := p.ReadPosFile(posFile)
	if err != nil && err != io.EOF {
		return nil, 0, err
	}

	err = p.Seek(slowlogFile, position)
	if err != nil {
		return nil, 0, err
	}

	p.SetReadBytes(position)

	return posFile, position, nil
}

func (p *Profiler) WritePositionFile(posFile *os.File, noSavePos bool, position int64) error {
	if !noSavePos && posFile != nil {
		posFile.Seek(0, 0)
		if position > 0 {
			position += p.ReadBytes()
		} else {
			position = p.ReadBytes()
		}

		_, err := posFile.Write([]byte(fmt.Sprint(position)))
		if err != nil {
			return err
		}
	}

	return nil
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

type ProfileHelper struct {
	inReader *os.File
}

func NewProfileHelper() *ProfileHelper {
	return &ProfileHelper{
		inReader: os.Stdin,
	}
}

func (ph *ProfileHelper) SetInReader(f *os.File) {
	ph.inReader = f
}

func (ph *ProfileHelper) Open(filename string) (*os.File, error) {
	var f *os.File
	var err error

	if filename != "" {
		f, err = os.Open(filename)
	} else {
		f = ph.inReader
	}

	return f, err
}

func (ph *ProfileHelper) LoadAndPrint(dumpFile string, sts *stats.QueryStats, printer *stats.Printer) error {
	lf, err := os.Open(dumpFile)
	if err != nil {
		return err
	}
	err = sts.LoadStats(lf)
	if err != nil {
		return err
	}
	defer lf.Close()

	sts.SortWithOptions()
	printer.Print(sts, nil)
	return nil
}

func (ph *ProfileHelper) Dump(dumpFile string, sts *stats.QueryStats) error {
	if dumpFile == "" {
		return nil
	}

	df, err := os.OpenFile(dumpFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	err = sts.DumpStats(df)
	if err != nil {
		return err
	}
	defer df.Close()

	return nil
}
