package internal

import (
	"fmt"
	"io"

	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/report"
)

func handleDiff(stdout io.Writer, diffReport *diff.Diff, format string) *ReturnError {
	switch format {
	case FormatYAML:
		if err := printYAML(stdout, diffReport); err != nil {
			return getErrFailedPrint("diff YAML", err)
		}
	case FormatJSON:
		if err := printJSON(stdout, diffReport); err != nil {
			return getErrFailedPrint("diff JSON", err)
		}
	case FormatText:
		fmt.Fprintf(stdout, "%s", report.GetTextReportAsString(diffReport))
	case FormatHTML:
		html, err := report.GetHTMLReportAsString(diffReport)
		if err != nil {
			return getErrFailedGenerateHTML(err)
		}
		fmt.Fprintf(stdout, "%s", html)
	default:
		return getErrUnsupportedDiffFormat(format)
	}

	return nil
}
