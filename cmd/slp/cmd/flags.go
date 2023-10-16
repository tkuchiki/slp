package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tkuchiki/slp/helper"
	"github.com/tkuchiki/slp/options"
	"github.com/tkuchiki/slp/stats"
)

const (
	flagConfig             = "config"
	flagFile               = "file"
	flagDump               = "dump"
	flagLoad               = "load"
	flagFormat             = "format"
	flagSort               = "sort"
	flagReverse            = "reverse"
	flagNoHeaders          = "noheaders"
	flagShowFooters        = "show-footers"
	flagLimit              = "limit"
	flagOutput             = "output"
	flagMatchingGroups     = "matching-groups"
	flagFilters            = "filters"
	flagPositionFile       = "pos"
	flagNoSavePositionFile = "nosave-pos"
	flagPercentiles        = "percentiles"
	flagBundleWhereIn      = "bundle-where-in"
	flagBundleValues       = "bundle-values"
	flagNoAbstract         = "noabstract"
	flagPage               = "page"
)

type flags struct {
	config      string
	sortOptions *stats.SortOptions
}

func newFlags() *flags {
	return &flags{
		sortOptions: stats.NewSortOptions(),
	}
}

func (f *flags) defineConfig(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&f.config, flagConfig, "", "The configuration file")
}

func (f *flags) defineFile(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagFile, "", "", "The access log file")
}

func (f *flags) defineDump(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagDump, "", "", "Dump profiled data as YAML")
}

func (f *flags) defineLoad(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagLoad, "", "", "Load the profiled YAML data")
}

func (f *flags) defineFormat(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagFormat, "", options.DefaultFormatOption, "The output format (table, markdown, tsv, csv, html)")
}

func (f *flags) defineSort(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagSort, "", options.DefaultSortOption, "Output the results in sorted order")
}

func (f *flags) defineReverse(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolP(flagReverse, "r", false, "Sort results in reverse order")
}

func (f *flags) defineNoHeaders(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolP(flagNoHeaders, "", false, "Output no header line at all (only --format=tsv, csv)")
}

func (f *flags) defineShowFooters(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolP(flagShowFooters, "", false, "Output footer line at all (only --format=table, markdown)")
}

func (f *flags) defineLimit(cmd *cobra.Command) {
	cmd.PersistentFlags().IntP(flagLimit, "", options.DefaultLimitOption, "The maximum number of results to display")
}

func (f *flags) defineOutput(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagOutput, "o", options.DefaultOutputOption, "Specifies the results to display, separated by commas")
}

func (f *flags) defineMatchingGroups(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagMatchingGroups, "m", "", "Specifies Query matching groups separated by commas")
}

func (f *flags) defineFilters(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagFilters, "f", "", "Only the logs are profiled that match the conditions")
}

func (f *flags) definePositionFile(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagPositionFile, "", "", "The position file")
}

func (f *flags) defineNoSavePositionFile(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolP(flagNoSavePositionFile, "", false, "Do not save position file")
}

func (f *flags) definePercentiles(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagPercentiles, "", "", "Specifies the percentiles separated by commas")
}

func (f *flags) definePage(cmd *cobra.Command) {
	cmd.PersistentFlags().IntP(flagPage, "", options.DefaultPaginationLimit, "Number of pages of pagination")
}

func (f *flags) defineBundleWhereIN(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolP(flagBundleWhereIn, "", false, "Bundle WHERE IN conditions")
}

func (f *flags) defineBundleValues(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolP(flagBundleValues, "", false, "Bundle VALUES of INSERT statement")
}

func (f *flags) defineNoAbstract(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolP(flagNoAbstract, "a", false, "Do not abstract all numbers to N and strings to 'S'")
}

func (f *flags) defineGlobalOptions(cmd *cobra.Command) {
	f.defineConfig(cmd)
}

func (f *flags) defineProfileOptions(cmd *cobra.Command) {
	f.defineFile(cmd)
	f.defineDump(cmd)
	f.defineLoad(cmd)
	f.defineFormat(cmd)
	f.defineSort(cmd)
	f.defineReverse(cmd)
	f.defineNoHeaders(cmd)
	f.defineShowFooters(cmd)
	f.defineLimit(cmd)
	f.defineOutput(cmd)
	f.defineMatchingGroups(cmd)
	f.defineFilters(cmd)
	f.definePositionFile(cmd)
	f.defineNoSavePositionFile(cmd)
	f.definePercentiles(cmd)
	f.definePage(cmd)
	f.defineBundleWhereIN(cmd)
	f.defineBundleValues(cmd)
	f.defineNoAbstract(cmd)
}

func (f *flags) defineDiffOptions(cmd *cobra.Command) {
	f.defineFormat(cmd)
	f.defineSort(cmd)
	f.defineReverse(cmd)
	f.defineNoHeaders(cmd)
	f.defineShowFooters(cmd)
	f.defineLimit(cmd)
	f.defineOutput(cmd)
	f.defineMatchingGroups(cmd)
	f.defineFilters(cmd)
	f.definePercentiles(cmd)
	f.definePage(cmd)
	f.defineBundleWhereIN(cmd)
	f.defineBundleValues(cmd)
	f.defineNoAbstract(cmd)
}

func (f *flags) bindFlags(cmd *cobra.Command) {
	viper.BindPFlag("file", cmd.PersistentFlags().Lookup(flagFile))
	viper.BindPFlag("dump", cmd.PersistentFlags().Lookup(flagDump))
	viper.BindPFlag("load", cmd.PersistentFlags().Lookup(flagLoad))
	viper.BindPFlag("format", cmd.PersistentFlags().Lookup(flagFormat))
	viper.BindPFlag("noheaders", cmd.PersistentFlags().Lookup(flagNoHeaders))
	viper.BindPFlag("show_footers", cmd.PersistentFlags().Lookup(flagShowFooters))
	viper.BindPFlag("limit", cmd.PersistentFlags().Lookup(flagLimit))
	viper.BindPFlag("matching_groups", cmd.PersistentFlags().Lookup(flagMatchingGroups))
	viper.BindPFlag("filters", cmd.PersistentFlags().Lookup(flagFilters))
	viper.BindPFlag("pos_file", cmd.PersistentFlags().Lookup(flagPositionFile))
	viper.BindPFlag("nosave_pos", cmd.PersistentFlags().Lookup(flagNoSavePositionFile))
	viper.BindPFlag("output", cmd.PersistentFlags().Lookup(flagOutput))
	viper.BindPFlag("pagination_limit", cmd.PersistentFlags().Lookup(flagPage))
	viper.BindPFlag("bundle_where_in", cmd.PersistentFlags().Lookup(flagBundleWhereIn))
	viper.BindPFlag("bundle_values", cmd.PersistentFlags().Lookup(flagBundleValues))
	viper.BindPFlag("noabstract", cmd.PersistentFlags().Lookup(flagNoAbstract))
}

func (f *flags) createOptionsFromConfig(cmd *cobra.Command) (*options.Options, error) {
	opts := options.NewOptions()
	viper.SetConfigFile(f.config)
	viper.SetConfigType("yaml")

	// Start workaround
	// viper seems to merge slices, so we'll set empty slice and overwrite it manually.
	opts.Percentiles = []int{}

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(opts); err != nil {
		return nil, err
	}

	if len(opts.Percentiles) == 0 {
		opts.Percentiles = options.DefaultPercentilesOption
	}
	// End workaround

	percentilesFlag := cmd.PersistentFlags().Lookup(flagPercentiles)
	if percentilesFlag != nil && percentilesFlag.Changed {
		ps := cmd.PersistentFlags().Lookup(flagPercentiles).Value.String()
		var percentiles []int
		var err error
		if ps != "" {
			percentiles, err = helper.SplitCSVIntoInts(ps)
			if err != nil {
				return nil, err
			}

			if err = helper.ValidatePercentiles(percentiles); err != nil {
				return nil, err
			}
		}
		opts.Percentiles = percentiles
	}

	if err := f.sortOptions.SetAndValidate(opts.Sort); err != nil {
		return nil, err
	}
	opts.Sort = f.sortOptions.SortType()

	return opts, nil
}

func (f *flags) setOptions(cmd *cobra.Command, opts *options.Options, flags []string) (*options.Options, error) {
	for _, flag := range flags {
		switch flag {
		case flagFile:
			file, err := cmd.PersistentFlags().GetString(flagFile)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.File(file))
		case flagDump:
			dump, err := cmd.PersistentFlags().GetString(flagDump)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.Dump(dump))
		case flagLoad:
			load, err := cmd.PersistentFlags().GetString(flagLoad)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.Load(load))
		case flagFormat:
			format, err := cmd.PersistentFlags().GetString(flagFormat)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.Format(format))
		case flagSort:
			sort, err := cmd.PersistentFlags().GetString(flagSort)
			if err != nil {
				return nil, err
			}

			err = f.sortOptions.SetAndValidate(sort)
			if err != nil {
				return nil, err
			}

			opts = options.SetOptions(opts, options.Sort(sort))
		case flagReverse:
			reverse, err := cmd.PersistentFlags().GetBool(flagReverse)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.Reverse(reverse))
		case flagNoHeaders:
			noHeaders, err := cmd.PersistentFlags().GetBool(flagNoHeaders)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.NoHeaders(noHeaders))
		case flagShowFooters:
			showFooters, err := cmd.PersistentFlags().GetBool(flagShowFooters)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.ShowFooters(showFooters))
		case flagLimit:
			limit, err := cmd.PersistentFlags().GetInt(flagLimit)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.Limit(limit))
		case flagOutput:
			output, err := cmd.PersistentFlags().GetString(flagOutput)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.Output(output))
		case flagMatchingGroups:
			matchingGroups, err := cmd.PersistentFlags().GetString(flagMatchingGroups)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.CSVGroups(matchingGroups))
		case flagFilters:
			filters, err := cmd.PersistentFlags().GetString(flagFilters)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.Filters(filters))
		case flagPositionFile:
			pos, err := cmd.PersistentFlags().GetString(flagPositionFile)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.PosFile(pos))
		case flagNoSavePositionFile:
			noSavePos, err := cmd.PersistentFlags().GetBool(flagNoSavePositionFile)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.NoSavePos(noSavePos))
		case flagPercentiles:
			ps, err := cmd.PersistentFlags().GetString(flagPercentiles)
			if err != nil {
				return nil, err
			}

			var percentiles []int
			if ps != "" {
				percentiles, err = helper.SplitCSVIntoInts(ps)
				if err != nil {
					return nil, err
				}

				if err = helper.ValidatePercentiles(percentiles); err != nil {
					return nil, err
				}
			}
			opts = options.SetOptions(opts, options.Percentiles(percentiles))
		case flagPage:
			paginationLimit, err := cmd.PersistentFlags().GetInt(flagPage)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.PaginationLimit(paginationLimit))
		case flagBundleWhereIn:
			bundleWhereIn, err := cmd.PersistentFlags().GetBool(flagBundleWhereIn)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.BundleWhereIn(bundleWhereIn))
		case flagBundleValues:
			bundleValues, err := cmd.PersistentFlags().GetBool(flagBundleValues)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.BundleValues(bundleValues))
		case flagNoAbstract:
			noAbstract, err := cmd.PersistentFlags().GetBool(flagNoAbstract)
			if err != nil {
				return nil, err
			}
			opts = options.SetOptions(opts, options.NoAbstract(noAbstract))
		}
	}

	return opts, nil
}

func (f *flags) setProfileOptions(cmd *cobra.Command, opts *options.Options) (*options.Options, error) {
	_flags := []string{
		flagFile,
		flagBundleWhereIn,
		flagBundleValues,
		flagNoAbstract,
		flagSort,
		flagReverse,
		flagOutput,
		flagMatchingGroups,
		flagFilters,
		flagPercentiles,
		flagPositionFile,
		flagNoSavePositionFile,
		flagFormat,
		flagDump,
		flagLoad,
		flagNoHeaders,
		flagShowFooters,
		flagLimit,
		flagPage,
	}

	return f.setOptions(cmd, opts, _flags)
}

func (f *flags) setDiffOptions(cmd *cobra.Command, opts *options.Options) (*options.Options, error) {
	_flags := []string{
		flagSort,
		flagReverse,
		flagOutput,
		flagFilters,
		flagPercentiles,
		flagFormat,
		flagNoHeaders,
		flagShowFooters,
		flagPage,
	}

	return f.setOptions(cmd, opts, _flags)
}

// slp
func (f *flags) createOptions(cmd *cobra.Command) (*options.Options, error) {
	if f.config != "" {
		f.bindFlags(cmd)
		return f.createOptionsFromConfig(cmd)
	}

	opts, err := f.setProfileOptions(cmd, options.NewOptions())
	if err != nil {
		return nil, err
	}

	return f.setProfileOptions(cmd, opts)
}

// slp diff
func (f *flags) createDiffOptions(cmd *cobra.Command) (*options.Options, error) {
	if f.config != "" {
		f.bindFlags(cmd)
		return f.createOptionsFromConfig(cmd)
	}

	return f.setDiffOptions(cmd, options.NewOptions())
}
