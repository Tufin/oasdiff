package report_test

import (
	"testing"

	"github.com/oasdiff/oasdiff/diff"
	"github.com/oasdiff/oasdiff/report"
	"github.com/stretchr/testify/require"
)

func TestHTML(t *testing.T) {
	d, err := diff.Get(diff.NewConfig(), l(t, 1), l(t, 3))
	require.NoError(t, err)

	html, err := report.GetHTMLReportAsString(d)
	require.NoError(t, err)
	require.NotEmpty(t, html)
}
