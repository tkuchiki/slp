package stats

import (
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
	if f.expeval != nil {
		return true
	}

	return false
}

func (f *Filter) Do(metrics *QueryMetrics) error {
	if !f.isEnable() {
		return nil
	}

	if f.expeval != nil {
		matched, err := f.expeval.Run(metrics)
		if err != nil {
			return err
		}

		if !matched {
			return errors.SkipReadLineErr
		}
	}

	return nil
}
