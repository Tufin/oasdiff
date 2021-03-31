package report_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/report"
)

func TestHTML(t *testing.T) {
	d, err := diff.Get(diff.NewConfig(), l(t, 1), l(t, 3))
	require.NoError(t, err)

	html, err := report.GetHTMLReportAsString(d)
	require.NoError(t, err)
	require.NotEmpty(t, html)
}
