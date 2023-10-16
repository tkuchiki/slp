package options

import (
	"github.com/tkuchiki/slp/helper"
)

const (
	DefaultSortOption      = "count"
	DefaultFormatOption    = "table"
	DefaultLimitOption     = 5000
	DefaultOutputOption    = "simple"
	DefaultPaginationLimit = 100
)

var DefaultPercentilesOption = []int{}

type Options struct {
	File            string   `mapstructure:"file"`
	Dump            string   `mapstructure:"dump"`
	Load            string   `mapstructure:"load"`
	Sort            string   `mapstructure:"sort"`
	Reverse         bool     `mapstructure:"reverse"`
	Format          string   `mapstructure:"format"`
	NoHeaders       bool     `mapstructure:"noheaders"`
	ShowFooters     bool     `mapstructure:"show_footers"`
	Limit           int      `mapstructure:"limit"`
	MatchingGroups  []string `mapstructure:"matching_groups"`
	Filters         string   `mapstructure:"filters"`
	PosFile         string   `mapstructure:"pos_file"`
	NoSavePos       bool     `mapstructure:"nosave_pos"`
	Output          string   `mapstructure:"output"`
	Percentiles     []int    `mapstructure:"percentiles"`
	BundleWhereIn   bool     `mapstructure:"bundle_where_in"`
	BundleValues    bool     `mapstructure:"bundle_values"`
	NoAbstract      bool     `mapstructure:"noabstract"`
	PaginationLimit int      `mapstructure:"pagination_limit"`
}

type Option func(*Options)

func File(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.File = s
		}
	}
}

func Dump(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.Dump = s
		}
	}
}

func Load(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.Load = s
		}
	}
}

func Sort(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.Sort = s
		}
	}
}

func Reverse(b bool) Option {
	return func(opts *Options) {
		if b {
			opts.Reverse = b
		}
	}
}

func Format(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.Format = s
		}
	}
}

func NoHeaders(b bool) Option {
	return func(opts *Options) {
		if b {
			opts.NoHeaders = b
		}
	}
}

func ShowFooters(b bool) Option {
	return func(opts *Options) {
		if b {
			opts.ShowFooters = b
		}
	}
}

func Limit(i int) Option {
	return func(opts *Options) {
		if i > 0 {
			opts.Limit = i
		}
	}
}

func MatchingGroups(values []string) Option {
	return func(opts *Options) {
		if len(values) > 0 {
			opts.MatchingGroups = values
		}
	}
}

func CSVGroups(csv string) Option {
	return func(opts *Options) {
		a := helper.SplitCSV(csv)
		if len(a) > 0 {
			opts.MatchingGroups = a
		}
	}
}

func Filters(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.Filters = s
		}
	}
}

func Output(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.Output = s
		}
	}
}

func PosFile(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.PosFile = s
		}
	}
}

func NoSavePos(b bool) Option {
	return func(opts *Options) {
		if b {
			opts.NoSavePos = b
		}
	}
}

func Percentiles(i []int) Option {
	return func(opts *Options) {
		if len(i) > 0 {
			opts.Percentiles = i
		}
	}
}

func BundleValues(b bool) Option {
	return func(opts *Options) {
		if b {
			opts.BundleValues = b
		}
	}
}

func BundleWhereIn(b bool) Option {
	return func(opts *Options) {
		if b {
			opts.BundleWhereIn = b
		}
	}
}

func NoAbstract(b bool) Option {
	return func(opts *Options) {
		if b {
			opts.NoAbstract = b
		}
	}
}

func PaginationLimit(i int) Option {
	return func(opts *Options) {
		if i > 0 {
			opts.PaginationLimit = i
		}
	}
}

func NewOptions(opt ...Option) *Options {
	options := &Options{
		Sort:            DefaultSortOption,
		Format:          DefaultFormatOption,
		Limit:           DefaultLimitOption,
		Output:          DefaultOutputOption,
		Percentiles:     DefaultPercentilesOption,
		PaginationLimit: DefaultPaginationLimit,
	}

	for _, o := range opt {
		o(options)
	}

	return options
}

func SetOptions(options *Options, opt ...Option) *Options {
	for _, o := range opt {
		o(options)
	}

	return options
}
