package internal

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/formatters"
)

type ChecksFlags struct {
	lang   string
	format string
}

func getChecksCmd() *cobra.Command {
	flags := ChecksFlags{}

	cmd := cobra.Command{
		Use:               "checks [flags]",
		Short:             "Display optional checks",
		Long:              `Display optional checks that can be passed to 'breaking' and 'changelog' with the 'include-checks' flag`,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions, // see https://github.com/spf13/cobra/issues/1969
		RunE: func(cmd *cobra.Command, args []string) error {
			// by now flags have been parsed successfully, so we don't need to show usage on any errors
			cmd.Root().SilenceUsage = true

			if err := runChecks(cmd.OutOrStdout(), flags); err != nil {
				setReturnValue(cmd, err.Code)
				return err
			}

			return nil
		},
	}

	enumWithOptions(&cmd, newEnumValue([]string{LangEn, LangRu}, LangDefault, &flags.lang), "lang", "l", "language for localized output")
	enumWithOptions(&cmd, newEnumValue(formatters.SupportedFormatsByContentType("checks"), string(formatters.FormatText), &flags.format), "format", "f", "output format")

	return &cmd
}

func runChecks(stdout io.Writer, flags ChecksFlags) *ReturnError {
	rules := checker.GetAllRules()

	if err := outputChecks(stdout, flags.lang, flags.format, rules); err != nil {
		return err
	}

	return nil
}

func outputChecks(stdout io.Writer, lang string, format string, rules []checker.BackwardCompatibilityRule) *ReturnError {
	// formatter lookup
	formatter, err := formatters.Lookup(format, formatters.FormatterOpts{
		Language: lang,
	})
	if err != nil {
		return getErrUnsupportedChecksFormat(format)
	}

	// render
	bytes, err := formatter.RenderChecks(rules, formatters.RenderOpts{})
	if err != nil {
		return getErrFailedPrint("checks "+format, err)
	}

	// print output
	_, _ = fmt.Fprintf(stdout, "%s\n", bytes)

	return nil
}
