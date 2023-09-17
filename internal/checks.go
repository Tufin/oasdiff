package internal

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/tufin/oasdiff/checker"
)

func getChecksCmd() *cobra.Command {

	cmd := cobra.Command{
		Use:   "checks [flags]",
		Short: "Display optional breaking changes checks",
		Long:  `Display optional breaking changes checks`,
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			// by now flags have been parsed successfully so we don't need to show usage on any errors
			cmd.Root().SilenceUsage = true

			runChecks(cmd.OutOrStdout())

			return nil
		},
	}

	return &cmd
}

func runChecks(stdout io.Writer) {
	printYAML(stdout, checker.GetOptionalChecks())
}
