package cmd

import (
	"github.com/tkuchiki/slp/options"

	"github.com/spf13/cobra"
)

func newPrintOutputOptionsSubCmd() *cobra.Command {
	var printOutputOptionsCmd = &cobra.Command{
		Use:   "print-output-options",
		Short: "Print --output/--sort options",
		Long:  `Print --output/--sort options`,
	}

	printOutputOptionsCmd.PersistentFlags().StringP("sep", "", "\n", "Separator")
	printOutputOptionsCmd.PersistentFlags().StringP("percentiles", "", "", "Specifies the percentiles separated by commas")

	return printOutputOptionsCmd
}

func runPrintOutputOptions(opts *options.Options) {

}
