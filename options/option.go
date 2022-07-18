package options

import (
	"io"

	"github.com/tkuchiki/slp/helper"
	"gopkg.in/yaml.v2"
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
	File            string   `yaml:"file"`
	Sort            string   `yaml:"sort"`
	Reverse         bool     `yaml:"reverse"`
	Format          string   `yaml:"format"`
	NoHeaders       bool     `yaml:"noheaders"`
	ShowFooters     bool     `yaml:"show_footers"`
	Limit           int      `yaml:"limit"`
	MatchingGroups  []string `yaml:"matching_groups"`
	Filters         string   `yaml:"filters"`
	PosFile         string   `yaml:"pos_file"`
	NoSavePos       bool     `yaml:"nosave_pos"`
	Output          string   `yaml:"output"`
	Percentiles     []int    `yaml:"percentiles"`
	AssembleWhereIn bool     `yaml:"assemble_where_in"`
	AssembleValues  bool     `yaml:"assemble_values"`
	NoAbstract      bool     `yaml:"noabstract"`
	PaginationLimit int      `yaml:"pagination_limit"`
}

type Option func(*Options)

func File(s string) Option {
	return func(opts *Options) {
		if s != "" {
			opts.File = s
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

func AssembleValues(b bool) Option {
	return func(opts *Options) {
		if b {
			opts.AssembleValues = b
		}
	}
}

func AssembleWhereIn(b bool) Option {
	return func(opts *Options) {
		if b {
			opts.AssembleWhereIn = b
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

func LoadOptionsFromReader(r io.Reader) (*Options, error) {
	opts := NewOptions()
	buf, err := io.ReadAll(r)
	if err != nil {
		return opts, err
	}

	configs := NewOptions()
	err = yaml.Unmarshal(buf, configs)

	opts = SetOptions(opts,
		Sort(configs.Sort),
		Limit(configs.Limit),
		Output(configs.Output),
		Reverse(configs.Reverse),
		File(configs.File),
		Format(configs.Format),
		NoHeaders(configs.NoHeaders),
		ShowFooters(configs.ShowFooters),
		PosFile(configs.PosFile),
		NoSavePos(configs.NoSavePos),
		MatchingGroups(configs.MatchingGroups),
		Filters(configs.Filters),
		Percentiles(configs.Percentiles),
		AssembleWhereIn(configs.AssembleWhereIn),
		AssembleValues(configs.AssembleValues),
		NoAbstract(configs.NoAbstract),
		PaginationLimit(configs.PaginationLimit),
	)

	return opts, err
}
