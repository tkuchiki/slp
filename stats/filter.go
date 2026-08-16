package stats

import (
	sqlv1 "github.com/tkuchiki/logschema/sql/v1"
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
	if f.options.Filters != "" {
		ee, err := NewExpEval(f.options.Filters)
		if err != nil {
			return err
		}

		f.expeval = ee
	}

	return nil
}

func (f *Filter) isEnable() bool {
	return f.expeval != nil
}

func (f *Filter) Do(record *sqlv1.Query) error {
	if !f.isEnable() {
		return nil
	}

	if f.expeval != nil {
		matched, err := f.expeval.Run(record)
		if err != nil {
			return err
		}

		if !matched {
			return errors.SkipReadLineErr
		}
	}

	return nil
}
