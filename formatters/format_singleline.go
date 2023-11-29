package formatters

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
	"github.com/tufin/oasdiff/report"
)

type SingleLineFormatter struct {
	notImplementedFormatter
	Localizer checker.Localizer
}

func (f SingleLineFormatter) RenderDiff(diff *diff.Diff, opts RenderOpts) ([]byte, error) {
	return []byte(report.GetTextReportAsString(diff)), nil
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

func (f SingleLineFormatter) RenderChangelog(changes checker.Changes, opts RenderOpts, specInfoPair *load.SpecInfoPair) ([]byte, error) {
	result := bytes.NewBuffer(nil)

	if len(changes) > 0 {
		_, _ = fmt.Fprint(result, getChangelogTitle(changes, f.Localizer, opts.ColorMode))
	}

	for _, c := range changes {
		_, _ = fmt.Fprintf(result, "%s\n\n", c.SingleLineError(f.Localizer, opts.ColorMode))
	}

	return result.Bytes(), nil
}

func (f SingleLineFormatter) RenderChecks(checks Checks, opts RenderOpts) ([]byte, error) {
	result := bytes.NewBuffer(nil)

	w := tabwriter.NewWriter(result, 1, 1, 1, ' ', 0)
	_, _ = fmt.Fprintln(w, "ID\tDESCRIPTION\tLEVEL")
	for _, check := range checks {
		_, _ = fmt.Fprintln(w, check.Id+"\t"+f.Localizer(check.Description)+"\t"+check.Level)
	}
	_ = w.Flush()

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
