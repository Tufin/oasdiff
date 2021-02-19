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
	s, err := load.NewSwaggerLoader().FromPath(fmt.Sprintf("../data/openapi-test%d.yaml", v))
	require.NoError(t, err)
	return s
}

func TestDiff_Same(t *testing.T) {
	s := l(t, 1)
	require.Empty(t, diff.Run(s, s, "", "").Diff.PathDiff)
}

func TestDiff_DeletedEndpointEmpty(t *testing.T) {
	require.Empty(t, diff.Run(l(t, 2), l(t, 1), "", "").Diff.PathDiff.DeletedEndpoints)
}

func TestDiff_DeletedEndpointNotEmpty(t *testing.T) {
	require.EqualValues(t, []string{"/api/{domain}/{project}/install-command"}, diff.Run(l(t, 1), l(t, 2), "", "").Diff.PathDiff.DeletedEndpoints)
}

func TestDiff_AddedOperation(t *testing.T) {
	require.Contains(t,
		diff.Run(l(t, 1), l(t, 2), "", "").Diff.PathDiff.ModifiedEndpoints["/api/{domain}/{project}/badges/security-score"].AddedOperations,
		"POST")
}

func TestDiff_DeletedOperation(t *testing.T) {
	require.Contains(t,
		diff.Run(l(t, 2), l(t, 1), "", "").Diff.PathDiff.ModifiedEndpoints["/api/{domain}/{project}/badges/security-score/"].DeletedOperations,
		"POST")
}

func TestDiff_AddedParam(t *testing.T) {
	require.Contains(t,
		diff.Run(l(t, 2), l(t, 1), "", "").Diff.PathDiff.ModifiedEndpoints["/api/{domain}/{project}/badges/security-score/"].ModifiedOperations["GET"].AddedParams["header"],
		"X-Auth-Name")
}

func TestDiff_DeletedParam(t *testing.T) {
	require.Contains(t,
		diff.Run(l(t, 1), l(t, 2), "", "").Diff.PathDiff.ModifiedEndpoints["/api/{domain}/{project}/badges/security-score"].ModifiedOperations["GET"].DeletedParams["header"],
		"X-Auth-Name")
}

func TestSchemaDiff_TypeDiff(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			OldValue: "string",
			NewValue: "integer",
		},
		diff.Run(l(t, 1), l(t, 2), "", "").Diff.PathDiff.ModifiedEndpoints["/api/{domain}/{project}/badges/security-score"].ModifiedOperations["GET"].ModifiedParams["path"]["domain"].SchemaDiff.TypeDiff)
}

func TestSchemaDiff_EnumDiff(t *testing.T) {
	require.Equal(t,
		true,
		diff.Run(l(t, 1), l(t, 3), "", "").Diff.PathDiff.ModifiedEndpoints["/api/{domain}/{project}/install-command"].ModifiedOperations["GET"].ModifiedParams["path"]["project"].SchemaDiff.EnumDiff)
}

func TestSchemaDiff_NotDiff(t *testing.T) {
	require.Equal(t,
		true,
		diff.Run(l(t, 1), l(t, 3), "", "").Diff.PathDiff.ModifiedEndpoints["/api/{domain}/{project}/badges/security-score"].ModifiedOperations["GET"].ModifiedParams["query"]["image"].SchemaDiff.NotDiff)
}

func TestSchemaDiff_ContentDiff(t *testing.T) {
	require.Equal(t,
		true,
		diff.Run(l(t, 2), l(t, 1), "", "").Diff.PathDiff.ModifiedEndpoints["/api/{domain}/{project}/badges/security-score/"].ModifiedOperations["GET"].ModifiedParams["query"]["filter"].ContentDiff.SchemaDiff.PropertiesDiff)
}

func TestSchemaDiff_MediaTypeAdded(t *testing.T) {
	require.Equal(t,
		true,
		diff.Run(l(t, 5), l(t, 1), "", "").Diff.PathDiff.ModifiedEndpoints["/api/{domain}/{project}/badges/security-score"].ModifiedOperations["GET"].ModifiedParams["header"]["user"].ContentDiff.MediaTypeAdded)
}

func TestSchemaDiff_MediaTypeDeleted(t *testing.T) {
	require.Equal(t,
		false,
		diff.Run(l(t, 1), l(t, 5), "", "").Diff.PathDiff.ModifiedEndpoints["/api/{domain}/{project}/badges/security-score"].ModifiedOperations["GET"].ModifiedParams["header"]["user"].ContentDiff.MediaTypeAdded)
}

func TestSchemaDiff_MediaTypeModified(t *testing.T) {
	require.Equal(t,
		true,
		diff.Run(l(t, 1), l(t, 5), "", "").Diff.PathDiff.ModifiedEndpoints["/api/{domain}/{project}/badges/security-score"].ModifiedOperations["GET"].ModifiedParams["cookie"]["test"].ContentDiff.MediaTypeDiff)
}

func TestSchemaDiff_AnyOfDiff(t *testing.T) {
	require.Equal(t,
		true,
		diff.Run(l(t, 4), l(t, 2), "/prefix", "").Diff.PathDiff.ModifiedEndpoints["/prefix/api/{domain}/{project}/badges/security-score/"].ModifiedOperations["GET"].ModifiedParams["query"]["token"].SchemaDiff.AnyOfDiff)
}

func TestSchemaDiff_MinDiff(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			OldValue: nil,
			NewValue: float64(7),
		},
		diff.Run(l(t, 4), l(t, 2), "/prefix", "").Diff.PathDiff.ModifiedEndpoints["/prefix/api/{domain}/{project}/badges/security-score/"].ModifiedOperations["GET"].ModifiedParams["path"]["domain"].SchemaDiff.MinDiff)
}

func TestSchemaDiff_AddedSchemas(t *testing.T) {
	require.Contains(t,
		diff.Run(l(t, 1), l(t, 5), "", "").Diff.SchemaDiff.AddedSchemas,
		"requests")
}

func TestSchemaDiff_DeletedSchemas(t *testing.T) {
	require.Contains(t,
		diff.Run(l(t, 5), l(t, 1), "", "").Diff.SchemaDiff.DeletedSchemas,
		"requests")
}

func TestSchemaDiff_ModifiedSchemas(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			OldValue: true,
			NewValue: false,
		},
		diff.Run(l(t, 1), l(t, 5), "", "").Diff.SchemaDiff.ModifiedSchemas["network-policies"].AdditionalPropertiesAllowedDiff,
		"requests")
}

func TestSchemaDiff_ModifiedSchemasOldNil(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			OldValue: nil,
			NewValue: false,
		},
		diff.Run(l(t, 1), l(t, 5), "", "").Diff.SchemaDiff.ModifiedSchemas["rules"].AdditionalPropertiesAllowedDiff,
		"requests")
}

func TestSchemaDiff_ModifiedSchemasNewNil(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			OldValue: false,
			NewValue: nil,
		},
		diff.Run(l(t, 5), l(t, 1), "", "").Diff.SchemaDiff.ModifiedSchemas["rules"].AdditionalPropertiesAllowedDiff,
		"requests")
}

func TestSummary(t *testing.T) {
	require.Equal(t,
		&diff.Summary{
			Diff: true,
			PathSummary: &diff.PathSummary{
				Added:    0,
				Deleted:  1,
				Modified: 1,
			},
			SchemaSummary: &diff.SchemaSummary{
				Deleted: 2,
			},
		},
		diff.Run(l(t, 1), l(t, 2), "", "").Summary)
}

func TestSummary_NoDiff(t *testing.T) {
	s := l(t, 1)

	require.Equal(t,
		&diff.Summary{
			Diff: false,
		},
		diff.Run(s, s, "", "").Summary)
}

func TestPrefix(t *testing.T) {
	require.Equal(t,
		&diff.Summary{
			Diff: true,
			PathSummary: &diff.PathSummary{
				Deleted:  0,
				Modified: 1,
			},
		},
		diff.Run(l(t, 4), l(t, 2), "/prefix", "").Summary)
}
