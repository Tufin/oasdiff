package formatters

import (
	"fmt"
	"os"

	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"golang.org/x/exp/slices"
)

// Formatter is a common interface for output formatters
type Formatter interface {
	RenderDiff(diff *diff.Diff, changes checker.Changes, opts RenderOpts) ([]byte, error)
	RenderSummary(checks []checker.BackwardCompatibilityCheck, diff *diff.Diff, changes checker.Changes, opts RenderOpts) ([]byte, error)
	RenderBreakingChanges(checks []checker.BackwardCompatibilityCheck, diff *diff.Diff, changes checker.Changes, opts RenderOpts) ([]byte, error)
	RenderChangelog(changes checker.Changes, opts RenderOpts) ([]byte, error)
	RenderChecks(rules []checker.BackwardCompatibilityRule, opts RenderOpts) ([]byte, error)
	SupportedOutputs() []string
}

var formatters = map[Format]Formatter{
	FormatYAML:          YAMLFormatter{},
	FormatJSON:          JSONFormatter{},
	FormatText:          TEXTFormatter{},
	FormatHTML:          HTMLFormatter{},
	FormatGithubActions: GitHubActionsFormatter{},
	// TODO: FormatSarif:         Sarif{},
}

// Lookup returns a formatter by its name
func Lookup(format string, opts FormatterOpts) (Formatter, error) {
	f := Format(format)

	switch f {
	case FormatYAML:
		return YAMLFormatter{
			Localizer: checker.NewLocalizer(opts.Language, LangDefault),
		}, nil
	case FormatJSON:
		return JSONFormatter{
			Localizer: checker.NewLocalizer(opts.Language, LangDefault),
		}, nil
	case FormatText:
		return TEXTFormatter{
			Localizer: checker.NewLocalizer(opts.Language, LangDefault),
		}, nil
	case FormatHTML:
		return HTMLFormatter{
			Localizer: checker.NewLocalizer(opts.Language, LangDefault),
		}, nil
	case FormatGithubActions:
		return GitHubActionsFormatter{
			Localizer: checker.NewLocalizer(opts.Language, LangDefault),
		}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", f)
	}
}

func SupportedFormatsByContentType(format string) []string {
	var formats []string
	for k, v := range formatters {
		if slices.Contains(v.SupportedOutputs(), format) {
			formats = append(formats, string(k))
		}
	}
	return formats
}

// DefaultFormatterOpts returns the default formatter options (e.g. colors, CI mode, etc.)
func DefaultFormatterOpts() FormatterOpts {
	return FormatterOpts{
		ColorMode: "auto",
		CI:        os.Getenv("CI") == "true",
		Language:  LangDefault,
	}
}
