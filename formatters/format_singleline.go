package formatters

import (
	"bytes"
	"fmt"

	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/load"
)

type SingleLineFormatter struct {
	notImplementedFormatter
	Localizer checker.Localizer
}

func newSingleLineFormatter(l checker.Localizer) SingleLineFormatter {
	return SingleLineFormatter{
		Localizer: l,
	}
}

func (f SingleLineFormatter) RenderBreakingChanges(changes checker.Changes, opts RenderOpts) ([]byte, error) {
	result := bytes.NewBuffer(nil)

	if len(changes) > 0 {
		_, _ = fmt.Fprint(result, getBreakingTitle(changes, f.Localizer, opts.ColorMode))
	}

	for _, c := range changes {
		_, _ = fmt.Fprintf(result, "%s\n\n", c.SingleLineError(f.Localizer, opts.ColorMode))
	}

	return result.Bytes(), nil
}

func (f SingleLineFormatter) RenderChangelog(changes checker.Changes, tags []string, opts RenderOpts, specInfoPair *load.SpecInfoPair) ([]byte, error) {
	result := bytes.NewBuffer(nil)

	if len(changes) > 0 {
		_, _ = fmt.Fprint(result, getChangelogTitle(changes, f.Localizer, opts.ColorMode))
	}

	for _, c := range changes {
		_, _ = fmt.Fprintf(result, "%s\n\n", c.SingleLineError(f.Localizer, opts.ColorMode))
	}

	return result.Bytes(), nil
}

func (f SingleLineFormatter) SupportedOutputs() []Output {
	return []Output{OutputBreaking, OutputChangelog}
}

func getBreakingTitle(changes checker.Changes, l checker.Localizer, colorMode checker.ColorMode) string {
	count := changes.GetLevelCount()
	return l(
		"total-errors",
		len(changes),
		count[checker.ERR],
		checker.ERR.StringCond(colorMode),
		count[checker.WARN],
		checker.WARN.StringCond(colorMode),
	)
}

func getChangelogTitle(changes checker.Changes, l checker.Localizer, colorMode checker.ColorMode) string {
	count := changes.GetLevelCount()
	return l(
		"total-changes",
		len(changes),
		count[checker.ERR],
		checker.ERR.StringCond(colorMode),
		count[checker.WARN],
		checker.WARN.StringCond(colorMode),
		count[checker.INFO],
		checker.INFO.StringCond(colorMode),
	)
}
