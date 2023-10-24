package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/percona/go-mysql/log"
	"github.com/spf13/cobra"
	"github.com/tkuchiki/slp/profiler"
	"github.com/tkuchiki/slp/stats"
)

func newMySQLCmd(flags *flags) *cobra.Command {
	var mysqlCmd = &cobra.Command{
		Use:           "my",
		Short:         "Profile the slowlogs for MySQL",
		Long:          `Profile the slowlogs for MySQL`,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := flags.createMySQLOptions(cmd)
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
			if err = printer.Validate(opts.Sort); err != nil {
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
			p := profiler.NewMySQLProfiler(slowlogFile, opt, opts)

			posFile, pos, err := phelper.ReadPositionFile(opts.PosFile, slowlogFile)
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

			err = phelper.WritePositionFile(posFile, opts.NoSavePos, pos+p.ReadBytesInt64())
			if err != nil {
				return err
			}

			sts.SortWithOptions()

			printer.Print(sts, nil)

			return nil
		},
	}

	flags.defineProfileOptions(mysqlCmd)

	mysqlCmd.Flags().SortFlags = false
	mysqlCmd.PersistentFlags().SortFlags = false
	mysqlCmd.InheritedFlags().SortFlags = false

	return mysqlCmd
}

func newMySQLPrintOutputOptionsCmd(flags *flags) *cobra.Command {
	mysqlPrintOutputOptionsCmd := newPrintOutputOptionsSubCmd()
	mysqlPrintOutputOptionsCmd.RunE = func(cmd *cobra.Command, args []string) error {
		opts, err := flags.createPrintOutputOptionsOptions(cmd)
		if err != nil {
			return err
		}

		fmt.Println(strings.Join(stats.MySQLOutputKeywords(opts.Percentiles), opts.Separator))

		return nil
	}

	mysqlPrintOutputOptionsCmd.Flags().SortFlags = false
	mysqlPrintOutputOptionsCmd.PersistentFlags().SortFlags = false
	mysqlPrintOutputOptionsCmd.InheritedFlags().SortFlags = false

	return mysqlPrintOutputOptionsCmd
}
