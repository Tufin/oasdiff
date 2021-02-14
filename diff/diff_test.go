package diff_test

import (
	"fmt"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

func l(t *testing.T, v int) *openapi3.Swagger {
	s, err := load.LoadPath(fmt.Sprintf("../data/openapi-test%d.yaml", v))
	require.NoError(t, err)
	return s
}

func TestDiff_Same(t *testing.T) {
	s := l(t, 1)
	require.Empty(t, diff.Diff(s, s, "").DeletedEndpoints)
}

func TestDiff_DeletedEndpointEmpty(t *testing.T) {
	require.Empty(t, diff.Diff(l(t, 2), l(t, 1), "").DeletedEndpoints)
}

func TestDiff_DeletedEndpointNotEmpty(t *testing.T) {
	require.EqualValues(t, []string{"/api/{domain}/{project}/install-command"}, diff.Diff(l(t, 1), l(t, 2), "").DeletedEndpoints)
}

func TestDiff_AddedOperation(t *testing.T) {
	require.Equal(t,
		diff.OperationMap{"POST": struct{}{}},
		diff.Diff(l(t, 1), l(t, 2), "").ModifiedEndpoints["/api/{domain}/{project}/badges/security-score"].AddedOperations)
}

func TestDiff_DeletedOperation(t *testing.T) {
	require.Equal(t,
		diff.OperationMap{"POST": struct{}{}},
		diff.Diff(l(t, 2), l(t, 1), "").ModifiedEndpoints["/api/{domain}/{project}/badges/security-score/"].DeletedOperations)
}

func TestDiff_AddedParam(t *testing.T) {
	require.Equal(t,
		diff.ParamNames{"X-Auth-Name": struct{}{}},
		diff.Diff(l(t, 2), l(t, 1), "").ModifiedEndpoints["/api/{domain}/{project}/badges/security-score/"].ModifiedOperations["GET"].AddedParams["header"])
}

func TestDiff_DeletedParam(t *testing.T) {
	require.Equal(t,
		diff.ParamNames{"X-Auth-Name": struct{}{}},
		diff.Diff(l(t, 1), l(t, 2), "").ModifiedEndpoints["/api/{domain}/{project}/badges/security-score"].ModifiedOperations["GET"].DeletedParams["header"])
}

func TestSchemaDiff_TypeDiff(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			OldValue: "string",
			NewValue: "integer",
		},
		diff.Diff(l(t, 1), l(t, 2), "").ModifiedEndpoints["/api/{domain}/{project}/badges/security-score"].ModifiedOperations["GET"].ModifiedParams["path"]["domain"].SchemaDiff.TypeDiff)
}

func TestSchemaDiff_EnumDiff(t *testing.T) {
	require.Equal(t,
		true,
		diff.Diff(l(t, 1), l(t, 3), "").ModifiedEndpoints["/api/{domain}/{project}/install-command"].ModifiedOperations["GET"].ModifiedParams["path"]["project"].SchemaDiff.EnumDiff)
}

func TestSchemaDiff_NotDiff(t *testing.T) {
	require.Equal(t,
		true,
		diff.Diff(l(t, 1), l(t, 3), "").ModifiedEndpoints["/api/{domain}/{project}/badges/security-score"].ModifiedOperations["GET"].ModifiedParams["query"]["image"].SchemaDiff.NotDiff)
}

func TestSchemaDiff_ContentDiff(t *testing.T) {
	require.Equal(t,
		true,
		diff.Diff(l(t, 2), l(t, 1), "").ModifiedEndpoints["/api/{domain}/{project}/badges/security-score/"].ModifiedOperations["GET"].ModifiedParams["query"]["filter"].ContentDiff.SchemaDiff.PropertiesDiff)
}

func TestSchemaDiff_AnyOfDiff(t *testing.T) {
	require.Equal(t,
		true,
		diff.Diff(l(t, 4), l(t, 2), "/prefix").ModifiedEndpoints["/prefix/api/{domain}/{project}/badges/security-score/"].ModifiedOperations["GET"].ModifiedParams["query"]["token"].SchemaDiff.AnyOfDiff)
}

func TestSchemaDiff_MinDiff(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			OldValue: nil,
			NewValue: float64(7),
		},
		diff.Diff(l(t, 4), l(t, 2), "/prefix").ModifiedEndpoints["/prefix/api/{domain}/{project}/badges/security-score/"].ModifiedOperations["GET"].ModifiedParams["path"]["domain"].SchemaDiff.MinDiff)
}
