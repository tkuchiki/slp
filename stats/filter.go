package stats

import (
	mlog "github.com/percona/go-mysql/log"
	"github.com/tkuchiki/slp/errors"
	"github.com/tkuchiki/slp/options"
)

type Filter struct {
	options *options.Options
	expeval *ExpEval
}

func NewFilter(options *options.Options) *Filter {
	return &Filter{
		options: options,
	}
}

func (f *Filter) Init() error {
	var err error

	if err != nil {
		return err
	}

	if f.options.Filters != "" {
		var ee *ExpEval
		ee, err = NewExpEval(f.options.Filters)
		if err != nil {
			return err
		}

		f.expeval = ee
	}

	return nil
}

func (f *Filter) isEnable() bool {
	if f.expeval != nil {
		return true
	}

	return false
}

func (f *Filter) Do(stat *mlog.Event) error {
	if !f.isEnable() {
		return nil
	}

	if f.expeval != nil {
		matched, err := f.expeval.Run(stat)
		if err != nil {
			return err
		}

		if !matched {
			return errors.SkipReadLineErr
		}
	}

	return nil
}
