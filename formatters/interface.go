package formatters

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/checker/localizations"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
	"golang.org/x/exp/slices"
)

// Formatter is a common interface for output formatters
type Formatter interface {
	RenderDiff(diff *diff.Diff, opts RenderOpts) ([]byte, error)
	RenderSummary(diff *diff.Diff, opts RenderOpts) ([]byte, error)
	RenderBreakingChanges(changes checker.Changes, opts RenderOpts) ([]byte, error)
	RenderChangelog(changes checker.Changes, opts RenderOpts, specInfoPair *load.SpecInfoPair) ([]byte, error)
	RenderChecks(checks Checks, opts RenderOpts) ([]byte, error)
	RenderFlatten(spec *openapi3.T, opts RenderOpts) ([]byte, error)
	SupportedOutputs() []Output
}

var formatters = map[Format]Formatter{
	FormatYAML:          YAMLFormatter{},
	FormatJSON:          JSONFormatter{},
	FormatText:          TEXTFormatter{},
	FormatHTML:          HTMLFormatter{},
	FormatGithubActions: GitHubActionsFormatter{},
	FormatJUnit:         JUnitFormatter{},
}

// Lookup returns a formatter by its name
func Lookup(format string, opts FormatterOpts) (Formatter, error) {
	f := Format(format)
	l := checker.NewLocalizer(opts.Language)

	switch f {
	case FormatYAML:
		return YAMLFormatter{
			Localizer: l,
		}, nil
	case FormatJSON:
		return JSONFormatter{
			Localizer: l,
		}, nil
	case FormatText:
		return TEXTFormatter{
			Localizer: l,
		}, nil
	case FormatHTML:
		return HTMLFormatter{
			Localizer: l,
		}, nil
	case FormatGithubActions:
		return GitHubActionsFormatter{
			Localizer: l,
		}, nil
	case FormatJUnit:
		return JUnitFormatter{
			Localizer: l,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", f)
	}
}

func SupportedFormatsByContentType(output Output) []string {
	var formats []string
	for k, v := range formatters {
		if slices.Contains(v.SupportedOutputs(), output) {
			formats = append(formats, string(k))
		}
	}
	return formats
}

// DefaultFormatterOpts returns the default formatter options (e.g. colors, CI mode, etc.)
func DefaultFormatterOpts() FormatterOpts {
	return FormatterOpts{
		Language: localizations.LangDefault,
	}
}
