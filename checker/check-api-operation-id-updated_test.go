package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: removing an existing operation id
func TestOperationIdRemoved(t *testing.T) {
	s1, err := open("../data/checker/operation_id_removed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/operation_id_removed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/api/v1.0/groups"].Post.OperationID = ""

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APIOperationIdUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "api-operation-id-removed",
		Text:        "api operation id 'createOneGroup' removed and replaced with ''",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/operation_id_removed_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: updating an existing operation id
func TestOperationIdUpdated(t *testing.T) {
	s1, err := open("../data/checker/operation_id_removed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/operation_id_removed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/api/v1.0/groups"].Post.OperationID = "newOperationId"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APIOperationIdUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "api-operation-id-removed",
		Text:        "api operation id 'createOneGroup' removed and replaced with 'newOperationId'",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/operation_id_removed_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: adding a new operation id
func TestOperationIdAdded(t *testing.T) {
	s1, err := open("../data/checker/operation_id_added_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/operation_id_added_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/api/v1.0/groups"].Post.OperationID = "NewOperationId"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig().WithCheckBreaking(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APIOperationIdUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "api-operation-id-added",
		Text:        "api operation id 'NewOperationId' was added",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/operation_id_added_base.yaml",
		OperationId: "NewOperationId",
	}, errs[0])
}
