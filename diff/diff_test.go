package diff_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

const (
	test1 = "../data/openapi-test1.yaml"
	test2 = "../data/openapi-test2.yaml"
	test3 = "../data/openapi-test3.yaml"
)

func TestDiff_Same(t *testing.T) {
	s, err := load.LoadPath(test1)
	require.NoError(t, err)

	require.Empty(t, diff.Diff(s, s, "").DeletedEndpoints)
}

func TestDiff_DeletedEndpoint(t *testing.T) {
	s1, err := load.LoadPath(test1)
	require.NoError(t, err)

	s2, err := load.LoadPath(test2)
	require.NoError(t, err)

	require.Empty(t, diff.Diff(s2, s1, "").DeletedEndpoints)
	require.EqualValues(t, []string{"/api/{domain}/{project}/install-command"}, diff.Diff(s1, s2, "").DeletedEndpoints)
}

func TestDiff_DeletedOperation(t *testing.T) {
	s1, err := load.LoadPath(test1)
	require.NoError(t, err)

	s2, err := load.LoadPath(test2)
	require.NoError(t, err)

	require.Equal(t,
		diff.DeletedOperations{"POST": struct{}{}},
		diff.Diff(s2, s1, "").ModifiedEndpoints["/api/{domain}/{project}/badges/security-score/"].DeletedOperations)
}

func TestDiff_DeletedParam(t *testing.T) {
	s1, err := load.LoadPath(test1)
	require.NoError(t, err)

	s2, err := load.LoadPath(test2)
	require.NoError(t, err)

	require.Equal(t,
		diff.ParamNames{"X-Auth-Name": struct{}{}},
		diff.Diff(s1, s2, "").ModifiedEndpoints["/api/{domain}/{project}/badges/security-score"].ModifiedOperations["GET"].DeletedParams["header"])
}

func TestDiff_TypeDiff(t *testing.T) {
	s1, err := load.LoadPath(test1)
	require.NoError(t, err)

	s2, err := load.LoadPath(test2)
	require.NoError(t, err)

	require.Equal(t,
		&diff.ValueDiff{
			OldValue: "string",
			NewValue: "integer",
		},
		diff.Diff(s1, s2, "").ModifiedEndpoints["/api/{domain}/{project}/badges/security-score"].ModifiedOperations["GET"].ModifiedParams["path"]["domain"].ShcemaDiff.TypeDiff)
}

func TestDiff_EnumDiff(t *testing.T) {
	s1, err := load.LoadPath(test1)
	require.NoError(t, err)

	s3, err := load.LoadPath(test3)
	require.NoError(t, err)

	require.Equal(t,
		true,
		diff.Diff(s1, s3, "").ModifiedEndpoints["/api/{domain}/{project}/install-command"].ModifiedOperations["GET"].ModifiedParams["path"]["project"].ShcemaDiff.EnumDiff)
}

func TestDiff_NotDiff(t *testing.T) {
	s1, err := load.LoadPath(test1)
	require.NoError(t, err)

	s3, err := load.LoadPath(test3)
	require.NoError(t, err)

	require.Equal(t,
		true,
		diff.Diff(s1, s3, "").ModifiedEndpoints["/api/{domain}/{project}/badges/security-score"].ModifiedOperations["GET"].ModifiedParams["query"]["image"].ShcemaDiff.NotDiff)
}
