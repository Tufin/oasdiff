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

type TEXTFormatter struct {
	notImplementedFormatter
	Localizer checker.Localizer
}

func (f TEXTFormatter) RenderDiff(diff *diff.Diff, opts RenderOpts) ([]byte, error) {
	return []byte(report.GetTextReportAsString(diff)), nil
}

func (f TEXTFormatter) RenderBreakingChanges(changes checker.Changes, opts RenderOpts) ([]byte, error) {
	result := bytes.NewBuffer(nil)

	if len(changes) > 0 {
		count := changes.GetLevelCount()
		title := f.Localizer(
			"total-errors",
			len(changes),
			count[checker.ERR],
			checker.ERR.PrettyString(),
			count[checker.WARN],
			checker.WARN.PrettyString(),
		)

		_, _ = fmt.Fprint(result, title)
	}

	for _, c := range changes {
		_, _ = fmt.Fprintf(result, "%s\n\n", c.PrettyErrorText(f.Localizer))
	}

	return result.Bytes(), nil
}

func (f TEXTFormatter) RenderChangelog(changes checker.Changes, opts RenderOpts, specInfoPair *load.SpecInfoPair) ([]byte, error) {
	result := bytes.NewBuffer(nil)

	if len(changes) > 0 {
		count := changes.GetLevelCount()
		title := f.Localizer(
			"total-changes",
			len(changes),
			count[checker.ERR],
			checker.ERR.PrettyString(),
			count[checker.WARN],
			checker.WARN.PrettyString(),
			count[checker.INFO],
			checker.INFO.PrettyString(),
		)

		_, _ = fmt.Fprint(result, title)
	}

	for _, c := range changes {
		_, _ = fmt.Fprintf(result, "%s\n\n", c.PrettyErrorText(f.Localizer))
	}

	return result.Bytes(), nil
}

func (f TEXTFormatter) RenderChecks(checks Checks, opts RenderOpts) ([]byte, error) {
	result := bytes.NewBuffer(nil)

	w := tabwriter.NewWriter(result, 1, 1, 1, ' ', 0)
	_, _ = fmt.Fprintln(w, "ID\tDESCRIPTION\tLEVEL")
	for _, check := range checks {
		_, _ = fmt.Fprintln(w, check.Id+"\t"+check.Description+"\t"+check.Level)
	}
	_ = w.Flush()

	return result.Bytes(), nil
}

func (f TEXTFormatter) SupportedOutputs() []Output {
	return []Output{OutputDiff, OutputBreaking, OutputChangelog, OutputChecks}
}
