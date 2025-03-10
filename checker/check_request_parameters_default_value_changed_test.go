package checker_test

import (
	"testing"

	"github.com/oasdiff/oasdiff/checker"
	"github.com/oasdiff/oasdiff/diff"
	"github.com/oasdiff/oasdiff/load"
	"github.com/stretchr/testify/require"
)

// CL: changing request parameter default value
func TestRequestParameterDefaultValueChanged(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_default_value_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_default_value_changed_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterDefaultValueChangedCheck), d, osm, checker.ERR)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestParameterDefaultValueChangedId,
		Args:        []any{"query", "category", "default_category", "updated_category"},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/request_parameter_default_value_changed_revision.yaml"),
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: changing request parameter default value, while the param is also renamed
func TestRequestParameterDefaultValueChangedAndRenamedParameter(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_default_value_changed_base_renamed.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_default_value_changed_revision_renamed.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterDefaultValueChangedCheck), d, osm, checker.ERR)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestParameterDefaultValueChangedId,
		Args:        []any{"path", "group_id", "2", "1"},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups/{group_id}",
		Source:      load.NewSource("../data/checker/request_parameter_default_value_changed_revision_renamed.yaml"),
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: adding request parameter default value
func TestRequestParameterDefaultValueAdded(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_default_value_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_default_value_changed_base.yaml")
	require.NoError(t, err)

	s1.Spec.Paths.Value("/api/v1.0/groups").Post.Parameters[1].Value.Schema.Value.Default = nil

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterDefaultValueChangedCheck), d, osm, checker.ERR)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestParameterDefaultValueAddedId,
		Args:        []any{"query", "category", "default_category"},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/request_parameter_default_value_changed_base.yaml"),
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: removing request parameter default value
func TestRequestParameterDefaultValueRemoved(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_default_value_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_default_value_changed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths.Value("/api/v1.0/groups").Post.Parameters[1].Value.Schema.Value.Default = nil

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterDefaultValueChangedCheck), d, osm, checker.ERR)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestParameterDefaultValueRemovedId,
		Args:        []any{"query", "category", "default_category"},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/request_parameter_default_value_changed_base.yaml"),
		OperationId: "createOneGroup",
	}, errs[0])
}
