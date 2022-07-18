package cmd

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/percona/go-mysql/log"
	"github.com/spf13/cobra"
	"github.com/tkuchiki/slp/options"
	"github.com/tkuchiki/slp/profiler"
	"github.com/tkuchiki/slp/stats"
)

func NewRootCmd(version string) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:     "slp",
		Short:   "Slow Log Profiler",
		Long:    `slp is a (MySQL) SlowLog Profiler`,
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			sortOptions := stats.NewSortOptions()
			opts, err := createOptions(cmd, sortOptions)
			if err != nil {
				return err
			}

			var slowlogFile *os.File
			if opts.File != "" {
				slowlogFile, err = os.Open(opts.File)
				if err != nil {
					return err
				}
			} else {
				slowlogFile = os.Stdin
			}
			defer slowlogFile.Close()

			printOptions := stats.NewPrintOptions(opts.NoHeaders, opts.ShowFooters, opts.PaginationLimit)
			printer := stats.NewPrinter(os.Stdout, opts.Output, opts.Format, opts.Percentiles, printOptions)
			if err = printer.Validate(); err != nil {
				return err
			}

			dump, err := cmd.PersistentFlags().GetString("dump")
			if err != nil {
				return err
			}

			load, err := cmd.PersistentFlags().GetString("load")
			if err != nil {
				return err
			}

			var sts *stats.QueryStats

			if len(opts.Percentiles) == 0 {
				sts = stats.NewQueryStats(false, false, false, false, false, false)
			} else {
				sts = stats.NewQueryStats(true, true, true, true, true, true)
			}

			sts.SetOptions(opts)
			sts.SetSortOptions(sortOptions)

			if load != "" {
				lf, err := os.Open(load)
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

			if len(opts.MatchingGroups) > 0 {
				err = sts.SetQueryMatchingGroups(opts.MatchingGroups)
				if err != nil {
					return err
				}
			}

			opt := log.Options{
				DefaultLocation: time.Local,
			}
			p := profiler.NewProfiler(slowlogFile, opt, opts)

			var posfile *os.File
			var pos int64
			if opts.PosFile != "" {
				posfile, err = p.OpenPosFile(opts.PosFile)
				if err != nil {
					return err
				}
				defer posfile.Close()

				pos, err = p.ReadPosFile(posfile)
				if err != nil && err != io.EOF {
					return err
				}

				err = p.Seek(slowlogFile, pos)
				if err != nil {
					return err
				}

				p.SetReadBytes(pos)
			}

			err = sts.InitFilter(opts)
			if err != nil {
				return err
			}

			err = p.Profile(sts)
			if err != nil {
				return err
			}

			if dump != "" {
				df, err := os.OpenFile(dump, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
				err = sts.DumpStats(df)
				if err != nil {
					return err
				}
				defer df.Close()
			}

			if !opts.NoSavePos && opts.PosFile != "" {
				posfile.Seek(0, 0)
				if pos > 0 {
					pos += p.ReadBytes()
				} else {
					pos = p.ReadBytes()
				}

				_, err = posfile.Write([]byte(fmt.Sprint(pos)))
				if err != nil {
					return err
				}
			}

			sts.SortWithOptions()

			printer.Print(sts, nil)

			return nil
		},
	}

	rootCmd.PersistentFlags().StringP("config", "", "", "The configuration file")
	rootCmd.PersistentFlags().StringP("file", "", "", "The slowlog file")
	rootCmd.PersistentFlags().StringP("dump", "", "", "Dump profiled data as YAML")
	rootCmd.PersistentFlags().StringP("load", "", "", "Load the profiled YAML data")
	rootCmd.PersistentFlags().StringP("format", "", options.DefaultFormatOption, "The output format (table, markdown, tsv, csv and html)")
	rootCmd.PersistentFlags().StringP("sort", "", options.DefaultSortOption, "Output the results in sorted order")
	rootCmd.PersistentFlags().BoolP("reverse", "r", false, "Sort results in reverse order")
	rootCmd.PersistentFlags().BoolP("noheaders", "", false, "Output no header line at all (only --format=tsv, csv)")
	rootCmd.PersistentFlags().BoolP("show-footers", "", false, "Output footer line at all (only --format=table, markdown)")
	rootCmd.PersistentFlags().IntP("limit", "", options.DefaultLimitOption, "The maximum number of results to display")
	rootCmd.PersistentFlags().StringP("output", "o", options.DefaultOutputOption, "Specifies the results to display, separated by commas")
	rootCmd.PersistentFlags().StringP("matching-groups", "m", "", "Specifies Query matching groups separated by commas")
	rootCmd.PersistentFlags().StringP("filters", "f", "", "Only the logs are profiled that match the conditions")
	rootCmd.PersistentFlags().StringP("pos", "", "", "The position file")
	rootCmd.PersistentFlags().BoolP("nosave-pos", "", false, "Do not save position file")
	rootCmd.PersistentFlags().StringP("percentiles", "", "", "Specifies the percentiles separated by commas")
	rootCmd.PersistentFlags().BoolP("assemble-where-in", "", false, "Assemble WHERE IN statement")
	rootCmd.PersistentFlags().BoolP("assemble-values", "", false, "Assemble VALUES of INSERT statement")
	rootCmd.PersistentFlags().BoolP("noabstract", "a", false, "Do not abstract all numbers to N and strings to 'S'")
	rootCmd.PersistentFlags().IntP("page", "", options.DefaultPaginationLimit, "Number of pages of pagination")

	rootCmd.AddCommand(NewPrintOutputOptionsCmd())
	rootCmd.AddCommand(NewDiffCmd(rootCmd))
	rootCmd.SetVersionTemplate(fmt.Sprintln(version))

	return rootCmd
}

func Execute(version string) error {
	rootCmd := NewRootCmd(version)
	return rootCmd.Execute()
}
