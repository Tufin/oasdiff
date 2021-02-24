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
	s, err := load.NewOASLoader().FromPath(fmt.Sprintf("../data/openapi-test%d.yaml", v))
	require.NoError(t, err)
	return s
}

func TestDiff_Same(t *testing.T) {
	s := l(t, 1)
	require.Empty(t, diff.Get(s, s, "", "").SpecDiff.PathsDiff)
}

func TestDiff_DeletedPaths(t *testing.T) {
	require.ElementsMatch(t,
		[]string{"/api/{domain}/{project}/install-command", "/register", "/subscribe"},
		diff.Get(l(t, 1), l(t, 2), "", "").SpecDiff.PathsDiff.Deleted)
}

func TestDiff_AddedOperation(t *testing.T) {
	require.Contains(t,
		diff.Get(l(t, 1), l(t, 2), "", "").SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score"].Added,
		"POST")
}

func TestDiff_DeletedOperation(t *testing.T) {
	require.Contains(t,
		diff.Get(l(t, 2), l(t, 1), "", "").SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score/"].Deleted,
		"POST")
}

func TestDiff_AddedGlobalTag(t *testing.T) {
	require.Contains(t,
		diff.Get(l(t, 3), l(t, 1), "", "").SpecDiff.TagsDiff.Added,
		"security")
}

func TestDiff_ModifiedGlobalTag(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: "Harrison",
			To:   "harrison",
		},
		diff.Get(l(t, 1), l(t, 3), "", "").SpecDiff.TagsDiff.Modified["reuven"].DescriptionDiff)
}

func TestDiff_AddedTag(t *testing.T) {
	require.Contains(t,
		diff.Get(l(t, 3), l(t, 1), "", "").SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score"].Modified["GET"].TagsDiff.Added,
		"security")
}

func TestDiff_AddedParam(t *testing.T) {
	require.Contains(t,
		diff.Get(l(t, 2), l(t, 1), "", "").SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score/"].Modified["GET"].ParamDiff.Added["header"],
		"X-Auth-Name")
}

func TestDiff_DeletedParam(t *testing.T) {
	require.Contains(t,
		diff.Get(l(t, 1), l(t, 2), "", "").SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score"].Modified["GET"].ParamDiff.Deleted["header"],
		"X-Auth-Name")
}

func TestDiff_ModifiedParam(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: true,
			To:   (interface{})(nil),
		},
		diff.Get(l(t, 2), l(t, 1), "", "").SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score/"].Modified["GET"].ParamDiff.Modified["query"]["image"].ExplodeDiff)
}

func TestSchemaDiff_TypeDiff(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: "string",
			To:   "integer",
		},
		diff.Get(l(t, 1), l(t, 2), "", "").SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score"].Modified["GET"].ParamDiff.Modified["path"]["domain"].SchemaDiff.TypeDiff)
}

func TestSchemaDiff_EnumDiff(t *testing.T) {
	require.Equal(t,
		&diff.EnumDiff{
			Added:   diff.EnumValues{"test1"},
			Deleted: diff.EnumValues{},
		},
		diff.Get(l(t, 1), l(t, 3), "", "").SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/install-command"].Modified["GET"].ParamDiff.Modified["path"]["project"].SchemaDiff.EnumDiff)
}

func TestSchemaDiff_RequiredAdded(t *testing.T) {
	require.Contains(t,
		diff.Get(l(t, 1), l(t, 5), "", "").SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score"].Modified["GET"].ParamDiff.Modified["query"]["filter"].ContentDiff.SchemaDiff.Required.Added,
		"type")
}

func TestSchemaDiff_RequiredDeleted(t *testing.T) {
	require.Contains(t,
		diff.Get(l(t, 5), l(t, 1), "", "").SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score"].Modified["GET"].ParamDiff.Modified["query"]["filter"].ContentDiff.SchemaDiff.Required.Deleted,
		"type")
}

func TestSchemaDiff_NotDiff(t *testing.T) {
	require.Equal(t,
		true,
		diff.Get(l(t, 1), l(t, 3), "", "").SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score"].Modified["GET"].ParamDiff.Modified["query"]["image"].SchemaDiff.NotDiff)
}

func TestSchemaDiff_ContentDiff(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: "number",
			To:   "string",
		},
		diff.Get(l(t, 2), l(t, 1), "", "").SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score/"].Modified["GET"].ParamDiff.Modified["query"]["filter"].ContentDiff.SchemaDiff.PropertiesDiff.Modified["color"].TypeDiff)
}

func TestSchemaDiff_MediaTypeAdded(t *testing.T) {
	require.Equal(t,
		true,
		diff.Get(l(t, 5), l(t, 1), "", "").SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score"].Modified["GET"].ParamDiff.Modified["header"]["user"].ContentDiff.MediaTypeAdded)
}

func TestSchemaDiff_MediaTypeDeleted(t *testing.T) {
	require.Equal(t,
		false,
		diff.Get(l(t, 1), l(t, 5), "", "").SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score"].Modified["GET"].ParamDiff.Modified["header"]["user"].ContentDiff.MediaTypeAdded)
}

func TestSchemaDiff_MediaTypeModified(t *testing.T) {
	require.Equal(t,
		true,
		diff.Get(l(t, 1), l(t, 5), "", "").SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score"].Modified["GET"].ParamDiff.Modified["cookie"]["test"].ContentDiff.MediaTypeDiff)
}

func TestSchemaDiff_MediaInvalidMultiEntries(t *testing.T) {
	s5 := l(t, 5)
	s5.Paths["/api/{domain}/{project}/badges/security-score"].Get.Parameters.GetByInAndName("cookie", "test").Content["second/invalid"] = openapi3.NewMediaType()

	s1 := l(t, 1)

	require.Nil(t,
		diff.Get(s1, s5, "", "").SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score"].Modified["GET"].ParamDiff.Modified["cookie"]["test"].ContentDiff)

	require.Nil(t,
		diff.Get(s5, s1, "", "").SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score"].Modified["GET"].ParamDiff.Modified["cookie"]["test"].ContentDiff)
}

func TestSchemaDiff_AnyOfDiff(t *testing.T) {
	require.Equal(t,
		true,
		diff.Get(l(t, 4), l(t, 2), "/prefix", "").SpecDiff.PathsDiff.Modified["/prefix/api/{domain}/{project}/badges/security-score/"].Modified["GET"].ParamDiff.Modified["query"]["token"].SchemaDiff.AnyOfDiff)
}

func TestSchemaDiff_MinDiff(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: nil,
			To:   float64(7),
		},
		diff.Get(l(t, 4), l(t, 2), "/prefix", "").SpecDiff.PathsDiff.Modified["/prefix/api/{domain}/{project}/badges/security-score/"].Modified["GET"].ParamDiff.Modified["path"]["domain"].SchemaDiff.MinDiff)
}

func TestResponseAdded(t *testing.T) {
	require.Contains(t,
		diff.Get(l(t, 1), l(t, 3), "", "").SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score"].Modified["GET"].ResponseDiff.Added,
		"default")
}

func TestResponseDeleted(t *testing.T) {
	require.Contains(t,
		diff.Get(l(t, 3), l(t, 1), "", "").SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score"].Modified["GET"].ResponseDiff.Deleted,
		"default")
}

func TestResponseDescriptionModified(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: "Tufin",
			To:   "Tufin1",
		},
		diff.Get(l(t, 3), l(t, 1), "", "").SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/install-command"].Modified["GET"].ResponseDiff.Modified["default"].DescriptionDiff)
}

func TestResponseHeadersModified(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: "Request limit per min.",
			To:   "Request limit per hour.",
		},
		diff.Get(l(t, 3), l(t, 1), "", "").SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/install-command"].Modified["GET"].ResponseDiff.Modified["default"].HeadersDiff.Modified["X-RateLimit-Limit"].DescriptionDiff)
}

func TestHeaderAdded(t *testing.T) {
	require.Contains(t,
		diff.Get(l(t, 5), l(t, 1), "", "").SpecDiff.HeadersDiff.Added,
		"new")
}

func TestHeaderDeleted(t *testing.T) {
	require.Contains(t,
		diff.Get(l(t, 1), l(t, 5), "", "").SpecDiff.HeadersDiff.Deleted,
		"new")
}

func TestHeaderModifiedSchema(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: false,
			To:   true,
		},
		diff.Get(l(t, 5), l(t, 1), "", "").SpecDiff.HeadersDiff.Modified["test"].SchemaDiff.AdditionalPropertiesAllowedDiff)
}

func TestHeaderModifiedContent(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: "string",
			To:   "object",
		},
		diff.Get(l(t, 5), l(t, 1), "", "").SpecDiff.HeadersDiff.Modified["testc"].ContentDiff.SchemaDiff.TypeDiff)
}

func TestResponseContentModified(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: "object",
			To:   "string",
		},
		diff.Get(l(t, 5), l(t, 1), "", "").SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score"].Modified["GET"].ResponseDiff.Modified["201"].ContentDiff.SchemaDiff.TypeDiff)
}

func TestResponseDespcriptionNil(t *testing.T) {

	s3 := l(t, 3)
	s3.Paths["/api/{domain}/{project}/install-command"].Get.Responses["default"].Value.Description = nil

	require.Equal(t,
		&diff.ValueDiff{
			From: interface{}(nil),
			To:   "Tufin1",
		},
		diff.Get(s3, l(t, 1), "", "").SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/install-command"].Modified["GET"].ResponseDiff.Modified["default"].DescriptionDiff)
}

func TestSchemaDiff_DeletedCallback(t *testing.T) {
	require.Contains(t,
		diff.Get(l(t, 3), l(t, 1), "", "").SpecDiff.PathsDiff.Modified["/register"].Modified["POST"].CallbacksDiff.Deleted,
		"myEvent")
}

func TestSchemaDiff_ModifiedCallback(t *testing.T) {
	require.Contains(t,
		diff.Get(l(t, 3), l(t, 1), "", "").SpecDiff.PathsDiff.Modified["/subscribe"].Modified["POST"].CallbacksDiff.Modified["myEvent"].Deleted,
		"{$request.body#/callbackUrl}")
}

func TestSchemaDiff_AddedSchemas(t *testing.T) {
	require.Contains(t,
		diff.Get(l(t, 1), l(t, 5), "", "").SpecDiff.SchemasDiff.Added,
		"requests")
}

func TestSchemaDiff_DeletedSchemas(t *testing.T) {
	require.Contains(t,
		diff.Get(l(t, 5), l(t, 1), "", "").SpecDiff.SchemasDiff.Deleted,
		"requests")
}

func TestSchemaDiff_ModifiedSchemas(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: true,
			To:   false,
		},
		diff.Get(l(t, 1), l(t, 5), "", "").SpecDiff.SchemasDiff.Modified["network-policies"].AdditionalPropertiesAllowedDiff)
}

func TestSchemaDiff_ModifiedSchemasOldNil(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: nil,
			To:   false,
		},
		diff.Get(l(t, 1), l(t, 5), "", "").SpecDiff.SchemasDiff.Modified["rules"].AdditionalPropertiesAllowedDiff)
}

func TestSchemaDiff_ModifiedSchemasNewNil(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: false,
			To:   nil,
		},
		diff.Get(l(t, 5), l(t, 1), "", "").SpecDiff.SchemasDiff.Modified["rules"].AdditionalPropertiesAllowedDiff)
}

func TestSchemaDiff_ModifiedSchemasValueDeleted(t *testing.T) {
	s5 := l(t, 5)
	s5.Components.Schemas["network-policies"].Value = nil

	require.Equal(t,
		true,
		diff.Get(l(t, 1), s5, "", "").SpecDiff.SchemasDiff.Modified["network-policies"].ValueDeleted)
}

func TestSchemaDiff_ModifiedSchemasValueAdded(t *testing.T) {
	s5 := l(t, 5)
	s5.Components.Schemas["network-policies"].Value = nil

	require.Equal(t,
		true,
		diff.Get(s5, l(t, 1), "", "").SpecDiff.SchemasDiff.Modified["network-policies"].ValueAdded)
}

func TestSchemaDiff_ModifiedSchemasBothValuesNil(t *testing.T) {
	s5 := l(t, 5)
	s5.Components.Schemas["network-policies"].Value = nil

	require.False(t, diff.Get(s5, s5, "", "").Summary.Diff)
}

func TestSummary(t *testing.T) {

	require.Equal(t,
		&diff.Summary{
			Diff: true,
			Components: map[string]*diff.SummaryDetails{
				diff.PathsComponent: {
					Added:    0,
					Deleted:  3,
					Modified: 1,
				},
				diff.SchemasComponent: {
					Deleted: 2,
				},
				diff.ParametersComponent: {
					Deleted: 1,
				},
				diff.HeadersComponent: {
					Deleted: 3,
				},
				diff.TagsComponent: {
					Deleted: 2,
				},
			},
		},
		diff.Get(l(t, 1), l(t, 2), "", "").Summary)
}

func TestSummary_NoDiff(t *testing.T) {
	s := l(t, 1)

	require.Equal(t,
		&diff.Summary{
			Diff:       false,
			Components: map[string]*diff.SummaryDetails{},
		},
		diff.Get(s, s, "", "").Summary)
}

func TestFilterByRegex(t *testing.T) {
	require.Nil(t, diff.Get(l(t, 1), l(t, 2), "", "x").Summary.Components[diff.PathsComponent])
}

func TestFilterByRegex_Invalid(t *testing.T) {
	require.Equal(t, true, diff.Get(l(t, 1), l(t, 2), "", "[").Summary.Diff)
}
