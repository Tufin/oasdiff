package checker_test

import (
	"testing"

	"github.com/oasdiff/oasdiff/checker"
	"github.com/oasdiff/oasdiff/diff"
	"github.com/oasdiff/oasdiff/load"
	"github.com/stretchr/testify/require"
)

// CL: setting max of request body
func TestRequestBodyMaxSetCheck(t *testing.T) {
	s1, err := open("../data/checker/request_body_max_set_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_body_max_set_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyMaxSetCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestBodyMaxSetId,
		Args:        []any{float64(15)},
		Level:       checker.WARN,
		Comment:     checker.RequestBodyMaxSetId + "-comment",
		Operation:   "POST",
		OperationId: "addPet",
		Path:        "/pets",
		Source:      load.NewSource("../data/checker/request_body_max_set_revision.yaml"),
	}, errs[0])
}

// CL: setting max of request propreties
func TestRequestPropertyMaxSetCheck(t *testing.T) {
	s1, err := open("../data/checker/request_property_max_set_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_max_set_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyMaxSetCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestPropertyMaxSetId,
		Args:        []any{"age", float64(15)},
		Level:       checker.WARN,
		Comment:     checker.RequestPropertyMaxSetId + "-comment",
		Operation:   "POST",
		OperationId: "addPet",
		Path:        "/pets",
		Source:      load.NewSource("../data/checker/request_property_max_set_revision.yaml"),
	}, errs[0])
}
