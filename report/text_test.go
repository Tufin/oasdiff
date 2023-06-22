package report_test

import (
	"fmt"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/report"
)

func l(t *testing.T, v int) *openapi3.T {
	t.Helper()

	loader := openapi3.NewLoader()
	oas, err := loader.LoadFromFile(fmt.Sprintf("../data/openapi-test%d.yaml", v))
	require.NoError(t, err)
	return oas
}

func d(t *testing.T, config *diff.Config, v1, v2 int) *diff.Diff {
	t.Helper()

	d, err := diff.Get(config, l(t, v1), l(t, v2))
	require.NoError(t, err)
	return d
}

func Test_NoChanges(t *testing.T) {
	require.Equal(t, report.GetTextReportAsString(d(t, &diff.Config{}, 3, 3)), "No changes\n")
}

func Test_NoEndpointChanges(t *testing.T) {
	s1 := openapi3.T{
		Info: &openapi3.Info{},
	}
	s2 := openapi3.T{
		Info: &openapi3.Info{
			Title: "reuven",
		},
	}

	dd, err := diff.Get(diff.NewConfig(), &s1, &s2)
	require.NoError(t, err)

	require.Equal(t, report.GetTextReportAsString(dd), "No endpoint changes\n")
}

func TestText1(t *testing.T) {
	require.Contains(t, report.GetTextReportAsString(d(t, &diff.Config{}, 3, 5)), "GET /api/{domain}/{project}/install-command")
}

func TestText2(t *testing.T) {
	require.Contains(t, report.GetTextReportAsString(d(t, &diff.Config{}, 5, 3)), "Deleted response: 201")
}

func TestText3(t *testing.T) {
	textReport := report.GetTextReportAsString(d(t, &diff.Config{}, 1, 3))

	require.Contains(t, textReport, "New enum values: [test1]")
	require.Contains(t, textReport, "Scheme OAuth Added scopes: [write:pets]")
}

func TestText4(t *testing.T) {
	textReport := report.GetTextReportAsString(d(t, &diff.Config{}, 3, 1))

	require.Contains(t, textReport, "New security requirements: bearerAuth")
	require.Contains(t, textReport, "Scheme OAuth Deleted scopes: [write:pets]")
}

func TestText5(t *testing.T) {
	textReport := report.GetTextReportAsString(d(t, &diff.Config{}, 2, 1))
	require.Contains(t, textReport, "Type changed from 'integer' to 'string'")
}

func TestText6(t *testing.T) {
	textReport := report.GetTextReportAsString(d(t, &diff.Config{}, 1, 5))
	require.Contains(t, textReport, "Type changed from 'string' to 'object'")
}

func TestText_DerefUint64(t *testing.T) {
	textReport := report.GetTextReportAsString(d(t, &diff.Config{}, 1, 3))
	require.Contains(t, textReport, "MaxLength changed from 29 to 30")
}

func TestText_DerefUint64Nil(t *testing.T) {
	textReport := report.GetTextReportAsString(d(t, &diff.Config{}, 1, 5))
	require.Contains(t, textReport, "MaxLength changed from 29 to null")
}
