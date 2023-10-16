package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/percona/go-mysql/log"
	"github.com/spf13/cobra"
	"github.com/tkuchiki/slp/profiler"
	"github.com/tkuchiki/slp/stats"
)

func newRootCmd(version string, flags *flags) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:           "slp",
		Short:         "Slow Log Profiler",
		Long:          `slp is a (MySQL) SlowLog Profiler`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       version,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := flags.createOptions(cmd)
			if err != nil {
				return err
			}

			phelper := profiler.NewProfileHelper()

			slowlogFile, err := phelper.Open(opts.File)
			if err != nil {
				return err
			}
			defer slowlogFile.Close()

			printOptions := stats.NewPrintOptions(opts.NoHeaders, opts.ShowFooters, opts.PaginationLimit)
			printer := stats.NewPrinter(os.Stdout, opts.Output, opts.Format, opts.Percentiles, printOptions)
			if err = printer.Validate(); err != nil {
				return err
			}

			var sts *stats.QueryStats

			if len(opts.Percentiles) == 0 {
				sts = stats.NewQueryStats(false, false, false, false, false, false)
			} else {
				sts = stats.NewQueryStats(true, true, true, true, true, true)
			}

			sts.SetOptions(opts)
			sts.SetSortOptions(flags.sortOptions)

			if opts.Load != "" {
				return phelper.LoadAndPrint(opts.Load, sts, printer)
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

			posFile, pos, err := p.ReadPositionFile(opts.PosFile, slowlogFile)
			if err != nil {
				return err
			}
			defer posFile.Close()

			err = sts.InitFilter(opts)
			if err != nil {
				return err
			}

			err = p.Profile(sts)
			if err != nil {
				return err
			}

			err = phelper.Dump(opts.Dump, sts)
			if err != nil {
				return err
			}

			err = p.WritePositionFile(posFile, opts.NoSavePos, pos)
			if err != nil {
				return err
			}

			sts.SortWithOptions()

			printer.Print(sts, nil)

			return nil
		},
	}

	flags.defineProfileOptions(rootCmd)
	rootCmd.SetVersionTemplate(fmt.Sprintln(version))

	rootCmd.Flags().SortFlags = false
	rootCmd.PersistentFlags().SortFlags = false
	rootCmd.InheritedFlags().SortFlags = false

	return rootCmd
}
