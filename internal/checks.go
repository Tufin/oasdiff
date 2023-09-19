package internal

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/tufin/oasdiff/checker"
)

func getChecksCmd() *cobra.Command {

	cmd := cobra.Command{
		Use:   "checks [flags]",
		Short: "Display optional checks",
		Long:  `Display optional checks`,
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {

			// by now flags have been parsed successfully so we don't need to show usage on any errors
			cmd.Root().SilenceUsage = true

			if err := runChecks(cmd.OutOrStdout()); err != nil {
				setReturnValue(cmd, err.Code)
				return err
			}

			return nil
		},
	}

	return &cmd
}

func runChecks(stdout io.Writer) *ReturnError {
	if err := printYAML(stdout, checker.GetOptionalChecks()); err != nil {
		return getErrFailedPrint("optional breaking changes checks", err)
	}
	return nil
}
