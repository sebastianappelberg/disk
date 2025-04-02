package cmd

import (
	"github.com/spf13/cobra"
)

var version string

func NewCmdRoot() *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "disk",
		Short:   "disk is a tool that helps you identify files that you can remove.",
		Long:    "disk is a tool that helps you identify files that you can remove.",
		Version: version,
	}

	// Subcommands
	cmd.AddCommand(NewCmdClean())
	cmd.AddCommand(NewCmdTree())
	cmd.AddCommand(NewCmdUsage())

	return cmd
}
