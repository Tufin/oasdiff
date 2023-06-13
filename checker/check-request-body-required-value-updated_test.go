package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// BC: changing response's body schema type from string to number is breaking
func TestRequestBodyBecameRequired(t *testing.T) {
	s1, err := open("../data/checker/request_body_became_required_base.yaml")
	s2, err := open("../data/checker/request_body_became_required_base.yaml")

	s2.Spec.Paths["/api/v1.0/groups"].Post.RequestBody.Value.Required = true

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.RequestBodyRequiredUpdatedCheck), d, osm)
	require.NotEmpty(t, errs)
	require.Equal(t, checker.BackwardCompatibilityErrors{
		{
			Id:          "request-body-became-required",
			Text:        "request body became required",
			Comment:     "",
			Level:       checker.ERR,
			Operation:   "POST",
			Path:        "/api/v1.0/groups",
			Source:      "../data/checker/request_body_became_required_base.yaml",
			OperationId: "createOneGroup",
		}}, errs)
}

func TestRequestBodyBecameOptional(t *testing.T) {
	s1, err := open("../data/checker/request_body_became_optional_base.yaml")
	s2, err := open("../data/checker/request_body_became_optional_base.yaml")
	require.Empty(t, err)

	s2.Spec.Paths["/api/v1.0/groups"].Post.RequestBody.Value.Required = false

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig().WithCheckBreaking(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestBodyRequiredUpdatedCheck), d, osm, checker.INFO)
	require.NotEmpty(t, errs)
	require.Equal(t, checker.BackwardCompatibilityErrors{
		{
			Id:          "request-body-became-optional",
			Text:        "request body became optional",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/api/v1.0/groups",
			Source:      "../data/checker/request_body_became_optional_base.yaml",
			OperationId: "createOneGroup",
		}}, errs)
}
