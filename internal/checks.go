package internal

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/formatters"
	"golang.org/x/exp/slices"
)

type ChecksFlags struct {
	lang     string
	format   string
	severity []string
	tags     []string
	required string
}

func getChecksCmd() *cobra.Command {
	flags := ChecksFlags{}

	cmd := cobra.Command{
		Use:               "checks [flags]",
		Short:             "Display checks",
		Long:              `Display a list of all checks that oasdiff supports.`,
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
	enumWithOptions(&cmd, newEnumSliceValue([]string{"info", "warn", "error"}, nil, &flags.severity), "severity", "s", "list of severities to include (experimental)")
	cmd.PersistentFlags().StringSliceVarP(&flags.tags, "tags", "t", []string{}, "list of tags to include, eg. parameter, request (experimental)")
	enumWithOptions(&cmd, newEnumValue([]string{"true", "false", "all"}, "all", &flags.required), "required", "r", "filter by required / optional")

	return &cmd
}

func runChecks(stdout io.Writer, flags ChecksFlags) *ReturnError {
	rules := checker.GetAllRules()

	if err := outputChecks(stdout, flags, rules); err != nil {
		return err
	}

	return nil
}

func outputChecks(stdout io.Writer, flags ChecksFlags, rules []checker.BackwardCompatibilityRule) *ReturnError {
	// formatter lookup
	formatter, err := formatters.Lookup(flags.format, formatters.FormatterOpts{
		Language: flags.lang,
	})
	if err != nil {
		return getErrUnsupportedChecksFormat(flags.format)
	}

	// filter rules
	filteredRules := make([]checker.BackwardCompatibilityRule, 0, len(rules))
	for _, rule := range rules {
		// severity
		if len(flags.severity) > 0 {
			if rule.Level == checker.ERR && !slices.Contains(flags.severity, "error") {
				continue
			}
			if rule.Level == checker.WARN && !slices.Contains(flags.severity, "warn") {
				continue
			}
			if rule.Level == checker.INFO && !slices.Contains(flags.severity, "info") {
				continue
			}
		}

		// required / optional
		if flags.required == "false" && rule.Required {
			continue
		}
		if flags.required == "true" && !rule.Required {
			continue
		}

		// tags (experimental, would be better to list the tags in the rule instead of str contains)
		if len(flags.tags) > 0 {
			match := false
			for _, tag := range flags.tags {
				if strings.Contains(rule.Id, tag) {
					match = true
				}
			}

			if !match {
				continue
			}
		}

		filteredRules = append(filteredRules, rule)
	}

	// render
	bytes, err := formatter.RenderChecks(filteredRules, formatters.RenderOpts{})
	if err != nil {
		return getErrFailedPrint("checks "+flags.format, err)
	}

	// print output
	_, _ = fmt.Fprintf(stdout, "%s\n", bytes)

	return nil
}
