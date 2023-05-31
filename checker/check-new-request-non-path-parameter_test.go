package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: new header, query and cookie request params are detected
func TestNewRequestNonPathParameter_DetectsNewPathsAndNewOperations(t *testing.T) {
	s1, err := open("../data/request_params/base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/request_params/optional-request-params.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.NewRequestNonPathParameterCheck), d, osm, checker.INFO)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 3)

	require.Equal(t, "new-optional-request-parameter", errs[0].Id)
	require.Equal(t, "GET", errs[0].Operation)
	require.Equal(t, "/api/test1", errs[0].Path)
	require.Equal(t, checker.INFO, errs[0].Level)
	require.Contains(t, errs[0].Text, "X-NewRequestHeaderParam")

	require.Equal(t, "new-optional-request-parameter", errs[1].Id)
	require.Equal(t, "GET", errs[1].Operation)
	require.Equal(t, "/api/test2", errs[1].Path)
	require.Contains(t, errs[1].Text, "newQueryParam")

	require.Equal(t, "new-optional-request-parameter", errs[2].Id)
	require.Equal(t, "GET", errs[2].Operation)
	require.Equal(t, "/api/test3", errs[2].Path)
	require.Contains(t, errs[2].Text, "csrf-token")
}
