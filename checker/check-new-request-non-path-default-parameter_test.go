package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// BC: new header, query and cookie required request default param is breaking
func TestNewRequestNonPathParameter_DetectsNewRequiredPathsAndNewOperations(t *testing.T) {
	s1, err := open("../data/request_params/base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/request_params/required-request-params.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.NewRequestNonPathDefaultParameterCheck), d, osm, checker.INFO)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 9)

	require.ElementsMatch(t, []checker.ApiChange{
		{
			Id:          "new-required-request-default-parameter-to-existing-path",
			Text:        "added the new required 'query' request parameter 'version' to all path's operations",
			Comment:     "",
			Level:       3,
			Operation:   "GET",
			OperationId: "getTest",
			Path:        "/api/test1",
			Source:      "../data/request_params/required-request-params.yaml",
		},
		{
			Id:          "new-required-request-default-parameter-to-existing-path",
			Text:        "added the new required 'query' request parameter 'version' to all path's operations",
			Comment:     "",
			Level:       3,
			Operation:   "POST",
			OperationId: "",
			Path:        "/api/test1",
			Source:      "../data/request_params/required-request-params.yaml",
		},
		{
			Id:          "new-required-request-default-parameter-to-existing-path",
			Text:        "added the new required 'query' request parameter 'id' to all path's operations",
			Comment:     "",
			Level:       3,
			Operation:   "GET",
			OperationId: "getTest",
			Path:        "/api/test2",
			Source:      "../data/request_params/required-request-params.yaml",
		},
		{
			Id:          "new-required-request-default-parameter-to-existing-path",
			Text:        "added the new required 'header' request parameter 'If-None-Match' to all path's operations",
			Comment:     "",
			Level:       3,
			Operation:   "GET",
			OperationId: "getTest",
			Path:        "/api/test3",
			Source:      "../data/request_params/required-request-params.yaml",
		},
		{
			Id:          "new-optional-request-default-parameter-to-existing-path",
			Text:        "added the new optional 'query' request parameter 'optionalQueryParam' to all path's operations",
			Comment:     "",
			Level:       1,
			Operation:   "GET",
			OperationId: "getTest",
			Path:        "/api/test1",
			Source:      "../data/request_params/required-request-params.yaml",
		},
		{
			Id:          "new-optional-request-default-parameter-to-existing-path",
			Text:        "added the new optional 'query' request parameter 'optionalQueryParam' to all path's operations",
			Comment:     "",
			Level:       1,
			Operation:   "POST",
			OperationId: "",
			Path:        "/api/test1",
			Source:      "../data/request_params/required-request-params.yaml",
		},
		{
			Id:          "new-optional-request-default-parameter-to-existing-path",
			Text:        "added the new optional 'header' request parameter 'optionalHeaderParam' to all path's operations",
			Comment:     "",
			Level:       1,
			Operation:   "GET",
			OperationId: "getTest",
			Path:        "/api/test2",
			Source:      "../data/request_params/required-request-params.yaml",
		}}, errs)
}
