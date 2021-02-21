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
	require.Empty(t, diff.Run(l(t, 2), l(t, 1), "", "").Diff.PathDiff.Deleted)
}

func TestDiff_DeletedEndpointNotEmpty(t *testing.T) {
	require.EqualValues(t,
		[]string{"/api/{domain}/{project}/install-command"},
		diff.Run(l(t, 1), l(t, 2), "", "").Diff.PathDiff.Deleted)
}

func TestDiff_AddedOperation(t *testing.T) {
	require.Contains(t,
		diff.Run(l(t, 1), l(t, 2), "", "").Diff.PathDiff.Modified["/api/{domain}/{project}/badges/security-score"].Added,
		"POST")
}

func TestDiff_DeletedOperation(t *testing.T) {
	require.Contains(t,
		diff.Run(l(t, 2), l(t, 1), "", "").Diff.PathDiff.Modified["/api/{domain}/{project}/badges/security-score/"].Deleted,
		"POST")
}

func TestDiff_AddedParam(t *testing.T) {
	require.Contains(t,
		diff.Run(l(t, 2), l(t, 1), "", "").Diff.PathDiff.Modified["/api/{domain}/{project}/badges/security-score/"].Modified["GET"].ParamDiff.Added["header"],
		"X-Auth-Name")
}

func TestDiff_DeletedParam(t *testing.T) {
	require.Contains(t,
		diff.Run(l(t, 1), l(t, 2), "", "").Diff.PathDiff.Modified["/api/{domain}/{project}/badges/security-score"].Modified["GET"].ParamDiff.Deleted["header"],
		"X-Auth-Name")
}

func TestDiff_ModifiedParam(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			OldValue: true,
			NewValue: (interface{})(nil),
		},
		diff.Run(l(t, 2), l(t, 1), "", "").Diff.PathDiff.Modified["/api/{domain}/{project}/badges/security-score/"].Modified["GET"].ParamDiff.Modified["query"]["image"].ExplodeDiff)
}

func TestSchemaDiff_TypeDiff(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			OldValue: "string",
			NewValue: "integer",
		},
		diff.Run(l(t, 1), l(t, 2), "", "").Diff.PathDiff.Modified["/api/{domain}/{project}/badges/security-score"].Modified["GET"].ParamDiff.Modified["path"]["domain"].SchemaDiff.TypeDiff)
}

func TestSchemaDiff_EnumDiff(t *testing.T) {
	require.Equal(t,
		&diff.EnumDiff{
			Added:   diff.EnumValues{"test1"},
			Deleted: diff.EnumValues{},
		},
		diff.Run(l(t, 1), l(t, 3), "", "").Diff.PathDiff.Modified["/api/{domain}/{project}/install-command"].Modified["GET"].ParamDiff.Modified["path"]["project"].SchemaDiff.EnumDiff)
}

func TestSchemaDiff_NotDiff(t *testing.T) {
	require.Equal(t,
		true,
		diff.Run(l(t, 1), l(t, 3), "", "").Diff.PathDiff.Modified["/api/{domain}/{project}/badges/security-score"].Modified["GET"].ParamDiff.Modified["query"]["image"].SchemaDiff.NotDiff)
}

func TestSchemaDiff_ContentDiff(t *testing.T) {
	require.Equal(t,
		true,
		diff.Run(l(t, 2), l(t, 1), "", "").Diff.PathDiff.Modified["/api/{domain}/{project}/badges/security-score/"].Modified["GET"].ParamDiff.Modified["query"]["filter"].ContentDiff.SchemaDiff.PropertiesDiff)
}

func TestSchemaDiff_MediaTypeAdded(t *testing.T) {
	require.Equal(t,
		true,
		diff.Run(l(t, 5), l(t, 1), "", "").Diff.PathDiff.Modified["/api/{domain}/{project}/badges/security-score"].Modified["GET"].ParamDiff.Modified["header"]["user"].ContentDiff.MediaTypeAdded)
}

func TestSchemaDiff_MediaTypeDeleted(t *testing.T) {
	require.Equal(t,
		false,
		diff.Run(l(t, 1), l(t, 5), "", "").Diff.PathDiff.Modified["/api/{domain}/{project}/badges/security-score"].Modified["GET"].ParamDiff.Modified["header"]["user"].ContentDiff.MediaTypeAdded)
}

func TestSchemaDiff_MediaTypeModified(t *testing.T) {
	require.Equal(t,
		true,
		diff.Run(l(t, 1), l(t, 5), "", "").Diff.PathDiff.Modified["/api/{domain}/{project}/badges/security-score"].Modified["GET"].ParamDiff.Modified["cookie"]["test"].ContentDiff.MediaTypeDiff)
}

func TestSchemaDiff_AnyOfDiff(t *testing.T) {
	require.Equal(t,
		true,
		diff.Run(l(t, 4), l(t, 2), "/prefix", "").Diff.PathDiff.Modified["/prefix/api/{domain}/{project}/badges/security-score/"].Modified["GET"].ParamDiff.Modified["query"]["token"].SchemaDiff.AnyOfDiff)
}

func TestSchemaDiff_MinDiff(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			OldValue: nil,
			NewValue: float64(7),
		},
		diff.Run(l(t, 4), l(t, 2), "/prefix", "").Diff.PathDiff.Modified["/prefix/api/{domain}/{project}/badges/security-score/"].Modified["GET"].ParamDiff.Modified["path"]["domain"].SchemaDiff.MinDiff)
}

func TestResponseAdded(t *testing.T) {
	require.Contains(t,
		diff.Run(l(t, 1), l(t, 3), "", "").Diff.PathDiff.Modified["/api/{domain}/{project}/badges/security-score"].Modified["GET"].ResponseDiff.Added,
		"default")
}

func TestResponseDeleted(t *testing.T) {
	require.Contains(t,
		diff.Run(l(t, 3), l(t, 1), "", "").Diff.PathDiff.Modified["/api/{domain}/{project}/badges/security-score"].Modified["GET"].ResponseDiff.Deleted,
		"200")
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

func TestFilterByRegex(t *testing.T) {
	require.Nil(t, diff.Run(l(t, 1), l(t, 2), "", "x").Summary.PathSummary)
}
