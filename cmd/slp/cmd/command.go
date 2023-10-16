package cmd

import (
	"github.com/spf13/cobra"
)

type Command struct {
	// slp
	rootCmd *cobra.Command

	// slp diff
	diffCmd *cobra.Command

	flags *flags
}

func NewCommand(version string) *Command {
	command := &Command{}
	command.flags = newFlags()

	command.rootCmd = newRootCmd(version, command.flags)

	command.flags.defineGlobalOptions(command.rootCmd)

	// slp diff
	command.diffCmd = newDiffCmd(command.flags)
	command.rootCmd.AddCommand(command.diffCmd)

	return command
}

func (c *Command) Execute() error {
	return c.rootCmd.Execute()
}

func (c *Command) setArgs(args []string) {
	c.rootCmd.SetArgs(args)
}
