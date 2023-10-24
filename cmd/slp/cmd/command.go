package cmd

import (
	"github.com/spf13/cobra"
)

type Command struct {
	// slp
	rootCmd *cobra.Command

	// slp my
	mysqlCmd *cobra.Command
	// slp my print-output-options
	mysqlPrintOutputOptionsCmd *cobra.Command

	// slp pg
	postgresqlCmd *cobra.Command
	// slp pg print-output-options
	postgresqlPrintOutputOptionsCmd *cobra.Command

	// slp diff
	diffCmd *cobra.Command

	flags *flags
}

func NewCommand(version string) *Command {
	command := &Command{}
	command.flags = newFlags()

	command.rootCmd = newRootCmd(version, command.flags)

	command.flags.defineGlobalOptions(command.rootCmd)

	// slp my
	command.mysqlCmd = newMySQLCmd(command.flags)
	command.rootCmd.AddCommand(command.mysqlCmd)

	// slp my print-output-options
	command.mysqlPrintOutputOptionsCmd = newMySQLPrintOutputOptionsCmd(command.flags)
	command.mysqlCmd.AddCommand(command.mysqlPrintOutputOptionsCmd)

	// slp pg
	command.postgresqlCmd = newPostgreSQLCmd(command.flags)
	command.rootCmd.AddCommand(command.postgresqlCmd)
	// slp pg print-output-options
	command.postgresqlPrintOutputOptionsCmd = newPostgreSQLPrintOutputOptionsCmd(command.flags)
	command.postgresqlCmd.AddCommand(command.postgresqlPrintOutputOptionsCmd)

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
