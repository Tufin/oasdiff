package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

// CL: increasing minimum value of request property
func TestRequestPropertyMinIncreasedCheck(t *testing.T) {
	s1, err := open("../data/checker/request_property_min_increased_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_min_increased_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyMinIncreasedCheck), d, osm, checker.ERR)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestPropertyMinIncreasedId,
		Args:        []any{"age", 15.0},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/pets",
		Source:      load.NewSource("../data/checker/request_property_min_increased_revision.yaml"),
		OperationId: "addPet",
	}, errs[0])
}

// CL: decreasing minimum value of request property
func TestRequestPropertyMinDecreasedCheck(t *testing.T) {
	s1, err := open("../data/checker/request_property_min_increased_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_min_increased_base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyMinIncreasedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestPropertyMinDecreasedId,
		Args:        []any{"age", 15.0, 10.0},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/pets",
		Source:      load.NewSource("../data/checker/request_property_min_increased_base.yaml"),
		OperationId: "addPet",
	}, errs[0])
}
