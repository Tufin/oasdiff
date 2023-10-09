package internal

import (
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/tufin/oasdiff/load"
)

func getParseArgs(flags Flags) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("please specify base and revision arguments as a path to a file, a glob (in composed mode), a URL, or '-' to read standard input")
		}
		if len(args) > 2 {
			return fmt.Errorf("extra arguments after base and revision")
		}
		if flags.getComposed() {
			if args[0] == "-" || args[1] == "-" {
				return errors.New("can't read from stdin in composed mode")
			}
		}

		return nil
	}
}

type runner func(flags Flags, stdout io.Writer) (bool, *ReturnError)

func getRun(flags Flags, runner runner) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {

		flags.setBase(load.GetSource(args[0]))
		flags.setRevision(load.GetSource(args[1]))

		// by now flags have been parsed successfully so we don't need to show usage on any errors
		cmd.Root().SilenceUsage = true

		failEmpty, err := runner(flags, cmd.OutOrStdout())
		if err != nil {
			setReturnValue(cmd, err.Code)
			return err
		}

		if failEmpty {
			setReturnValue(cmd, 1)
		}

		return nil
	}
}
