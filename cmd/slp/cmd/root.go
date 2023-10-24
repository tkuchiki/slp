package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newRootCmd(version string, flags *flags) *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:           "slp",
		Short:         "Slow Log Profiler",
		Long:          `slp is a MySQL/PostgreSQL SlowLog Profiler`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       version,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				cmd.Help()
			}

			return nil
		},
	}

	rootCmd.SetVersionTemplate(fmt.Sprintln(version))

	rootCmd.Flags().SortFlags = false
	rootCmd.PersistentFlags().SortFlags = false
	rootCmd.InheritedFlags().SortFlags = false

	return rootCmd
}
