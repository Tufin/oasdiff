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
	require.Len(t, errs, 5)

	require.IsType(t, checker.ApiChange{}, errs[0])
	e0 := errs[0].(checker.ApiChange)
	require.Equal(t, "new-required-request-default-parameter-to-existing-path", e0.Id)
	require.Equal(t, "/api/test1", e0.Path)
	require.Equal(t, checker.ERR, e0.Level)
	require.Contains(t, e0.Text, "version")

	require.IsType(t, checker.ApiChange{}, errs[1])
	e1 := errs[1].(checker.ApiChange)
	require.Equal(t, "new-required-request-default-parameter-to-existing-path", e1.Id)
	require.Equal(t, "/api/test2", e1.Path)
	require.Equal(t, checker.ERR, e1.Level)
	require.Contains(t, e1.Text, "id")

	require.IsType(t, checker.ApiChange{}, errs[2])
	e2 := errs[2].(checker.ApiChange)
	require.Equal(t, "new-required-request-default-parameter-to-existing-path", e2.Id)
	require.Equal(t, "/api/test3", e2.Path)
	require.Equal(t, checker.ERR, e2.Level)
	require.Contains(t, e2.Text, "If-None-Match")

	require.IsType(t, checker.ApiChange{}, errs[3])
	e3 := errs[3].(checker.ApiChange)
	require.Equal(t, "new-optional-request-default-parameter-to-existing-path", e3.Id)
	require.Equal(t, "/api/test1", e3.Path)
	require.Equal(t, checker.INFO, e3.Level)
	require.Contains(t, e3.Text, "optionalQueryParam")

	require.IsType(t, checker.ApiChange{}, errs[4])
	e4 := errs[4].(checker.ApiChange)
	require.Equal(t, "new-optional-request-default-parameter-to-existing-path", e4.Id)
	require.Equal(t, "/api/test2", e4.Path)
	require.Equal(t, checker.INFO, e4.Level)
	require.Contains(t, e4.Text, "optionalHeaderParam")

}
