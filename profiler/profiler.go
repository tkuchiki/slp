package profiler

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/tkuchiki/slp/helper"
	"github.com/tkuchiki/slp/stats"
)

type Profiler interface {
	Profile(qstats *stats.QueryStats) error
}

type ProfileHelper struct {
	inReader  *os.File
	readBytes uint64
}

func NewProfileHelper() *ProfileHelper {
	return &ProfileHelper{
		inReader: os.Stdin,
	}
}

func (ph *ProfileHelper) OpenPosFile(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
}

func (ph *ProfileHelper) ReadPosFile(f *os.File) (int64, error) {
	reader := bufio.NewReader(f)
	pos, _, err := reader.ReadLine()
	if err != nil {
		return 0, err
	}

	return helper.StringToInt64(string(pos))
}

func (ph *ProfileHelper) ReadBytes() int64 {
	return int64(ph.readBytes)
}

func (ph *ProfileHelper) SetReadBytes(n int64) {
	ph.readBytes = uint64(n)
}

func (ph *ProfileHelper) Seek(file *os.File, n int64) error {
	_, err := file.Seek(n, io.SeekCurrent)
	return err
}

func (ph *ProfileHelper) ReadPositionFile(posFilePath string, slowlogFile *os.File) (*os.File, int64, error) {
	if posFilePath == "" {
		return nil, 0, nil
	}

	posFile, err := ph.OpenPosFile(posFilePath)
	if err != nil {
		return nil, 0, err
	}

	position, err := ph.ReadPosFile(posFile)
	if err != nil && err != io.EOF {
		return nil, 0, err
	}

	err = ph.Seek(slowlogFile, position)
	if err != nil {
		return nil, 0, err
	}

	ph.SetReadBytes(position)

	return posFile, position, nil
}

func (ph *ProfileHelper) WritePositionFile(posFile *os.File, noSavePos bool, position int64) error {
	if !noSavePos && posFile != nil {
		posFile.Seek(0, 0)
		if position > 0 {
			position += ph.ReadBytes()
		} else {
			position = ph.ReadBytes()
		}

		_, err := posFile.Write([]byte(fmt.Sprint(position)))
		if err != nil {
			return err
		}
	}

	return nil
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
