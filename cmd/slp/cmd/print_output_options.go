package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tkuchiki/slp/helper"
	"github.com/tkuchiki/slp/stats"
)

func NewPrintOutputOptionsCmd() *cobra.Command {
	var printOutputOptionsCmd = &cobra.Command{
		Use:   "print-output-options",
		Short: "Print --output options",
		Long:  `Print --output options`,
		RunE: func(cmd *cobra.Command, args []string) error {
			sep, err := cmd.PersistentFlags().GetString("sep")
			if err != nil {
				return err
			}

			ps, err := cmd.PersistentFlags().GetString("percentiles")
			if err != nil {
				return err
			}

			var percentiles []int
			if ps != "" {
				percentiles, err = helper.SplitCSVIntoInts(ps)
				if err != nil {
					return err
				}

				if err = helper.ValidatePercentiles(percentiles); err != nil {
					return err
				}
			}

			fmt.Println(strings.Join(stats.OutputKeywords(percentiles), sep))

			return nil
		},
	}

	printOutputOptionsCmd.PersistentFlags().StringP("sep", "", "\n", "Separator")
	printOutputOptionsCmd.PersistentFlags().StringP("percentiles", "", "", "Specifies the percentiles separated by commas")

	return printOutputOptionsCmd
}
