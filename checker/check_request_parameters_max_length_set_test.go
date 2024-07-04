package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

// CL: setting maxLength of request parameters
func TestRequestParameterMaxLengthSetCheck(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_max_length_set_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_max_length_set_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterMaxLengthSetCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:        checker.RequestParameterMaxLengthSetId,
		Args:      []any{"query", "category", uint64(15)},
		Level:     checker.WARN,
		Comment:   checker.RequestParameterMaxLengthSetId + "-comment",
		Operation: "POST",
		Path:      "/test",
		Source:    load.NewSource("../data/checker/request_parameter_max_length_set_revision.yaml"),
	}, errs[0])
}
