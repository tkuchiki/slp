package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tkuchiki/slp/stats"
)

func NewDiffCmd(rootCmd *cobra.Command) *cobra.Command {
	var diffCmd = &cobra.Command{
		Use:   "diff <from> <to>",
		Args:  cobra.ExactArgs(2),
		Short: "Show the difference between the two profile results",
		Long:  `Show the difference between the two profile results`,
		RunE: func(cmd *cobra.Command, args []string) error {
			sortOptions := stats.NewSortOptions()
			opts, err := createOptions(rootCmd, sortOptions)
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

			fromPath := args[0]
			toPath := args[1]

			from, err := os.Open(fromPath)
			if err != nil {
				return err
			}
			err = sts.LoadStats(from)
			if err != nil {
				return err
			}
			defer from.Close()

			sts.SortWithOptions()

			var toSts *stats.QueryStats

			if len(opts.Percentiles) == 0 {
				toSts = stats.NewQueryStats(false, false, false, false, false, false)
			} else {
				toSts = stats.NewQueryStats(true, true, true, true, true, true)
			}

			err = toSts.InitFilter(opts)
			if err != nil {
				return err
			}

			toSts.SetOptions(opts)
			toSts.SetSortOptions(sortOptions)

			to, err := os.Open(toPath)
			if err != nil {
				return err
			}
			err = toSts.LoadStats(to)
			if err != nil {
				return err
			}
			defer to.Close()

			toSts.SortWithOptions()

			printOptions := stats.NewPrintOptions(opts.NoHeaders, opts.ShowFooters, opts.PaginationLimit)
			printer := stats.NewPrinter(os.Stdout, opts.Output, opts.Format, opts.Percentiles, printOptions)
			if err = printer.Validate(); err != nil {
				return err
			}

			printer.Print(sts, toSts)

			return nil
		},
	}

	return diffCmd
}
