package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tkuchiki/slp/helper"
	"github.com/tkuchiki/slp/options"
	"github.com/tkuchiki/slp/stats"
)

func createOptions(rootCmd *cobra.Command, sortOptions *stats.SortOptions) (*options.Options, error) {
	config, err := rootCmd.PersistentFlags().GetString("config")
	if err != nil {
		return nil, err
	}

	file, err := rootCmd.PersistentFlags().GetString("file")
	if err != nil {
		return nil, err
	}

	format, err := rootCmd.PersistentFlags().GetString("format")
	if err != nil {
		return nil, err
	}

	sort, err := rootCmd.PersistentFlags().GetString("sort")
	if err != nil {
		return nil, err
	}

	err = sortOptions.SetAndValidate(sort)
	if err != nil {
		return nil, err
	}

	reverse, err := rootCmd.PersistentFlags().GetBool("reverse")
	if err != nil {
		return nil, err
	}

	noHeaders, err := rootCmd.PersistentFlags().GetBool("noheaders")
	if err != nil {
		return nil, err
	}

	showFooters, err := rootCmd.PersistentFlags().GetBool("show-footers")
	if err != nil {
		return nil, err
	}

	limit, err := rootCmd.PersistentFlags().GetInt("limit")
	if err != nil {
		return nil, err
	}

	output, err := rootCmd.PersistentFlags().GetString("output")
	if err != nil {
		return nil, err
	}

	matchingGroups, err := rootCmd.PersistentFlags().GetString("matching-groups")
	if err != nil {
		return nil, err
	}

	filters, err := rootCmd.PersistentFlags().GetString("filters")
	if err != nil {
		return nil, err
	}

	pos, err := rootCmd.PersistentFlags().GetString("pos")
	if err != nil {
		return nil, err
	}

	noSavePos, err := rootCmd.PersistentFlags().GetBool("nosave-pos")
	if err != nil {
		return nil, err
	}

	ps, err := rootCmd.PersistentFlags().GetString("percentiles")
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

	bundleWhereIn, err := rootCmd.PersistentFlags().GetBool("bundle-where-in")
	if err != nil {
		return nil, err
	}

	bundleValues, err := rootCmd.PersistentFlags().GetBool("bundle-values")
	if err != nil {
		return nil, err
	}

	noAbstract, err := rootCmd.PersistentFlags().GetBool("noabstract")
	if err != nil {
		return nil, err
	}

	paginationLimit, err := rootCmd.PersistentFlags().GetInt("page")
	if err != nil {
		return nil, err
	}

	var opts *options.Options
	if config != "" {
		cf, err := os.Open(config)
		if err != nil {
			return nil, err
		}
		defer cf.Close()

		opts, err = options.LoadOptionsFromReader(cf)
		if err != nil {
			return nil, err
		}

		err = sortOptions.SetAndValidate(opts.Sort)
		if err != nil {
			return nil, err
		}

		percentiles = opts.Percentiles
	} else {
		opts = options.NewOptions()
	}

	return options.SetOptions(opts,
		options.File(file),
		options.Sort(sortOptions.SortType()),
		options.Reverse(reverse),
		options.Format(format),
		options.Limit(limit),
		options.Output(output),
		options.NoHeaders(noHeaders),
		options.ShowFooters(showFooters),
		options.CSVGroups(matchingGroups),
		options.Filters(filters),
		options.PosFile(pos),
		options.NoSavePos(noSavePos),
		options.Percentiles(percentiles),
		options.BundleWhereIn(bundleWhereIn),
		options.BundleValues(bundleValues),
		options.NoAbstract(noAbstract),
		options.PaginationLimit(paginationLimit),
	), nil
}
