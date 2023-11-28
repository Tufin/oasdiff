package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: changing request's body to required is breaking
func TestRequestBodyBecameRequired(t *testing.T) {
	s1, err := open("../data/checker/request_body_became_required_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_body_became_required_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/api/v1.0/groups"].Post.RequestBody.Value.Required = true

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.RequestBodyRequiredUpdatedCheck), d, osm)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestBodyBecameRequiredId,
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/request_body_became_required_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: changing request's body to optional
func TestRequestBodyBecameOptional(t *testing.T) {
	s1, err := open("../data/checker/request_body_became_optional_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_body_became_optional_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/api/v1.0/groups"].Post.RequestBody.Value.Required = false

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig().WithCheckBreaking(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestBodyRequiredUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestBodyBecameOptionalId,
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/request_body_became_optional_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}
