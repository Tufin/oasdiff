package internal

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/tufin/oasdiff/checker"
)

func getBreakingChangesCmd() *cobra.Command {

	flags := ChangelogFlags{}

	cmd := cobra.Command{
		Use:   "breaking base revision [flags]",
		Short: "Display breaking changes",
		Long:  "Display breaking changes between base and revision specs." + specHelp,
		Args:  getParseArgs(&flags),
		RunE:  getRun(&flags, runBreakingChanges),
	}

	addCommonDiffFlags(&cmd, &flags)
	addCommonBreakingFlags(&cmd, &flags)
	enumWithOptions(&cmd, newEnumValue([]string{LevelErr, LevelWarn}, "", &flags.failOn), "fail-on", "o", "exit with return code 1 when output includes errors with this level or higher")

	return &cmd
}

func runBreakingChanges(flags Flags, stdout io.Writer) (bool, *ReturnError) {
	return getChangelog(flags, stdout, checker.WARN)
}
