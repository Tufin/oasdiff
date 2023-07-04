package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: new header, query and cookie request params
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

	require.IsType(t, checker.BackwardCompatibilityError{}, errs[0])
	e0 := errs[0].(checker.BackwardCompatibilityError)
	require.Equal(t, "new-optional-request-parameter", e0.Id)
	require.Equal(t, "GET", e0.Operation)
	require.Equal(t, "/api/test1", e0.Path)
	require.Equal(t, checker.INFO, e0.Level)
	require.Contains(t, e0.Text, "X-NewRequestHeaderParam")

	require.IsType(t, checker.BackwardCompatibilityError{}, errs[1])
	e1 := errs[1].(checker.BackwardCompatibilityError)
	require.Equal(t, "new-optional-request-parameter", e1.Id)
	require.Equal(t, "GET", e1.Operation)
	require.Equal(t, "/api/test2", e1.Path)
	require.Equal(t, checker.INFO, e1.Level)
	require.Contains(t, e1.Text, "newQueryParam")

	require.IsType(t, checker.BackwardCompatibilityError{}, errs[2])
	e2 := errs[2].(checker.BackwardCompatibilityError)
	require.Equal(t, "new-optional-request-parameter", e2.Id)
	require.Equal(t, "GET", e2.Operation)
	require.Equal(t, "/api/test3", e2.Path)
	require.Equal(t, checker.INFO, e2.Level)
	require.Contains(t, e2.Text, "csrf-token")
}
