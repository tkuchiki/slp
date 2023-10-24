package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tkuchiki/slp/profiler"
	"github.com/tkuchiki/slp/stats"
)

func newPostgreSQLCmd(flags *flags) *cobra.Command {
	var pgCmd = &cobra.Command{
		Use:           "pg",
		Short:         "Profile the slowlogs for PostgreSQL",
		Long:          `Profile the slowlogs for PostgreSQL`,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := flags.createPGOptions(cmd)
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
				sts = stats.NewQueryStats(true, false, false, false, false, false)
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

			p, err := profiler.NewPGProfiler(slowlogFile, opts)
			if err != nil {
				return err
			}

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

	flags.defineProfileOptions(pgCmd)
	flags.definePGOptions(pgCmd)

	pgCmd.Flags().SortFlags = false
	pgCmd.PersistentFlags().SortFlags = false
	pgCmd.InheritedFlags().SortFlags = false

	return pgCmd
}

func newPostgreSQLPrintOutputOptionsCmd(flags *flags) *cobra.Command {
	postgresqlPrintOutputOptionsCmd := newPrintOutputOptionsSubCmd()
	postgresqlPrintOutputOptionsCmd.RunE = func(cmd *cobra.Command, args []string) error {
		opts, err := flags.createPrintOutputOptionsOptions(cmd)
		if err != nil {
			return err
		}

		fmt.Println(strings.Join(stats.PGOutputKeywords(opts.Percentiles), opts.Separator))

		return nil
	}

	postgresqlPrintOutputOptionsCmd.Flags().SortFlags = false
	postgresqlPrintOutputOptionsCmd.PersistentFlags().SortFlags = false
	postgresqlPrintOutputOptionsCmd.InheritedFlags().SortFlags = false

	return postgresqlPrintOutputOptionsCmd
}
