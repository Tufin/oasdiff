package diff_test

import (
	"fmt"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
)

func l(t *testing.T, v int) *openapi3.Swagger {
	loader := openapi3.NewSwaggerLoader()
	oas, err := loader.LoadSwaggerFromFile(fmt.Sprintf("../data/openapi-test%d.yaml", v))
	require.NoError(t, err)
	return oas
}

func TestDiff_Same(t *testing.T) {
	s := l(t, 1)
	require.Nil(t, diff.Get(diff.NewConfig(), s, s).SpecDiff)
}

func TestDiff_Empty(t *testing.T) {
	require.True(t, (&diff.CallbacksDiff{}).Empty())
	require.True(t, (&diff.EncodingsDiff{}).Empty())
	require.True(t, (&diff.ExtensionsDiff{}).Empty())
	require.True(t, (&diff.HeadersDiff{}).Empty())
	require.True(t, (&diff.OperationsDiff{}).Empty())
	require.True(t, (&diff.ParametersDiff{}).Empty())
	require.True(t, (&diff.RequestBodiesDiff{}).Empty())
	require.True(t, (&diff.ResponsesDiff{}).Empty())
	require.True(t, (&diff.SchemasDiff{}).Empty())
	require.True(t, (&diff.ServersDiff{}).Empty())
	require.True(t, (&diff.StringsDiff{}).Empty())
	require.True(t, (&diff.StringMapDiff{}).Empty())
	require.True(t, (&diff.TagsDiff{}).Empty())
}

func TestDiff_DeletedPaths(t *testing.T) {
	require.ElementsMatch(t,
		[]string{"/api/{domain}/{project}/install-command", "/register", "/subscribe"},
		diff.Get(diff.NewConfig(), l(t, 1), l(t, 2)).SpecDiff.PathsDiff.Deleted)
}

func TestDiff_AddedOperation(t *testing.T) {
	require.Contains(t,
		diff.Get(diff.NewConfig(), l(t, 1), l(t, 2)).SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score"].OperationsDiff.Added,
		"POST")
}

func TestDiff_DeletedOperation(t *testing.T) {
	require.Contains(t,
		diff.Get(diff.NewConfig(), l(t, 2), l(t, 1)).SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score/"].OperationsDiff.Deleted,
		"POST")
}

func TestAddedExtension(t *testing.T) {
	config := diff.Config{
		IncludeExtensions: diff.StringSet{"x-extension-test": struct{}{}},
	}

	require.Contains(t,
		diff.Get(&config, l(t, 3), l(t, 1)).SpecDiff.ExtensionsDiff.Added,
		"x-extension-test")
}

func TestDeletedExtension(t *testing.T) {
	config := diff.Config{
		IncludeExtensions: diff.StringSet{"x-extension-test": struct{}{}},
	}

	require.Contains(t,
		diff.Get(&config, l(t, 1), l(t, 3)).SpecDiff.ExtensionsDiff.Deleted,
		"x-extension-test")
}

func TestModifiedExtension(t *testing.T) {
	config := diff.Config{
		IncludeExtensions: diff.StringSet{"x-extension-test2": struct{}{}},
	}
	require.NotNil(t,
		diff.Get(&config, l(t, 1), l(t, 3)).SpecDiff.ExtensionsDiff.Modified["x-extension-test2"])
}

func TestExcludedExtension(t *testing.T) {
	require.Nil(t,
		diff.Get(diff.NewConfig(), l(t, 1), l(t, 3)).SpecDiff.ExtensionsDiff)
}

func TestDiff_AddedGlobalTag(t *testing.T) {
	require.Contains(t,
		diff.Get(diff.NewConfig(), l(t, 3), l(t, 1)).SpecDiff.TagsDiff.Added,
		"security")
}

func TestDiff_ModifiedGlobalTag(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: "Harrison",
			To:   "harrison",
		},
		diff.Get(diff.NewConfig(), l(t, 1), l(t, 3)).SpecDiff.TagsDiff.Modified["reuven"].DescriptionDiff)
}

func TestDiff_AddedTag(t *testing.T) {
	require.Contains(t,
		diff.Get(diff.NewConfig(), l(t, 3), l(t, 1)).SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score"].OperationsDiff.Modified["GET"].TagsDiff.Added,
		"security")
}

func TestDiff_DeletedEncoding(t *testing.T) {
	require.Contains(t,
		diff.Get(diff.NewConfig(), l(t, 1), l(t, 3)).SpecDiff.PathsDiff.Modified["/subscribe"].OperationsDiff.Modified["POST"].CallbacksDiff.Modified["myEvent"].Modified["hi"].OperationsDiff.Modified["POST"].RequestBodyDiff.ContentDiff.EncodingsDiff.Deleted,
		"historyMetadata")
}

func TestDiff_ModifiedEncodingHeaders(t *testing.T) {
	require.NotNil(t,
		diff.Get(diff.NewConfig(), l(t, 3), l(t, 1)).SpecDiff.PathsDiff.Modified["/subscribe"].OperationsDiff.Modified["POST"].CallbacksDiff.Modified["myEvent"].Modified["hi"].OperationsDiff.Modified["POST"].RequestBodyDiff.ContentDiff.EncodingsDiff.Modified["profileImage"].HeadersDiff,
		"profileImage")
}

func TestDiff_AddedParam(t *testing.T) {
	require.Contains(t,
		diff.Get(diff.NewConfig(), l(t, 2), l(t, 1)).SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score/"].OperationsDiff.Modified["GET"].ParametersDiff.Added[openapi3.ParameterInHeader],
		"X-Auth-Name")
}

func TestDiff_DeletedParam(t *testing.T) {
	require.Contains(t,
		diff.Get(diff.NewConfig(), l(t, 1), l(t, 2)).SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score"].OperationsDiff.Modified["GET"].ParametersDiff.Deleted[openapi3.ParameterInHeader],
		"X-Auth-Name")
}

func TestDiff_ModifiedParam(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: true,
			To:   (interface{})(nil),
		},
		diff.Get(diff.NewConfig(), l(t, 2), l(t, 1)).SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score/"].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInQuery]["image"].ExplodeDiff)
}

func TestSchemaDiff_TypeDiff(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: "string",
			To:   "integer",
		},
		diff.Get(diff.NewConfig(), l(t, 1), l(t, 2)).SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score"].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInPath]["domain"].SchemaDiff.TypeDiff)
}

func TestSchemaDiff_EnumDiff(t *testing.T) {
	require.Equal(t,
		&diff.EnumDiff{
			Added:   diff.EnumValues{"test1"},
			Deleted: diff.EnumValues{},
		},
		diff.Get(diff.NewConfig(), l(t, 1), l(t, 3)).SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/install-command"].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInPath]["project"].SchemaDiff.EnumDiff)
}

func TestSchemaDiff_RequiredAdded(t *testing.T) {
	require.Contains(t,
		diff.Get(diff.NewConfig(), l(t, 1), l(t, 5)).SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score"].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInQuery]["filter"].ContentDiff.SchemaDiff.Required.Added,
		"type")
}

func TestSchemaDiff_RequiredDeleted(t *testing.T) {
	require.Contains(t,
		diff.Get(diff.NewConfig(), l(t, 5), l(t, 1)).SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score"].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInQuery]["filter"].ContentDiff.SchemaDiff.Required.Deleted,
		"type")
}

func TestSchemaDiff_NotDiff(t *testing.T) {
	require.Equal(t,
		true,
		diff.Get(diff.NewConfig(), l(t, 1), l(t, 3)).SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score"].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInQuery]["image"].SchemaDiff.NotDiff.SchemaAdded)
}

func TestSchemaDiff_ContentDiff(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: "number",
			To:   "string",
		},
		diff.Get(diff.NewConfig(), l(t, 2), l(t, 1)).SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score/"].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInQuery]["filter"].ContentDiff.SchemaDiff.PropertiesDiff.Modified["color"].TypeDiff)
}

func TestSchemaDiff_MediaTypeAdded(t *testing.T) {
	require.Equal(t,
		true,
		diff.Get(diff.NewConfig(), l(t, 5), l(t, 1)).SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score"].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInHeader]["user"].ContentDiff.MediaTypeAdded)
}

func TestSchemaDiff_MediaTypeDeleted(t *testing.T) {
	require.Equal(t,
		false,
		diff.Get(diff.NewConfig(), l(t, 1), l(t, 5)).SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score"].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInHeader]["user"].ContentDiff.MediaTypeAdded)
}

func TestSchemaDiff_MediaTypeModified(t *testing.T) {
	require.Equal(t,
		true,
		diff.Get(diff.NewConfig(), l(t, 1), l(t, 5)).SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score"].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInCookie]["test"].ContentDiff.MediaTypeDiff)
}

func TestSchemaDiff_MediaInvalidMultiEntries(t *testing.T) {
	s5 := l(t, 5)
	s5.Paths["/api/{domain}/{project}/badges/security-score"].Get.Parameters.GetByInAndName(openapi3.ParameterInCookie, "test").Content["second/invalid"] = openapi3.NewMediaType()

	s1 := l(t, 1)

	require.Nil(t,
		diff.Get(diff.NewConfig(), s1, s5).SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score"].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInCookie])

	require.Nil(t,
		diff.Get(diff.NewConfig(), s5, s1).SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score"].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInCookie])
}

func TestSchemaDiff_AnyOfDiff(t *testing.T) {
	require.Equal(t,
		true,
		diff.Get(&diff.Config{Prefix: "/prefix"}, l(t, 4), l(t, 2)).SpecDiff.PathsDiff.Modified["/prefix/api/{domain}/{project}/badges/security-score/"].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInQuery]["token"].SchemaDiff.AnyOfDiff)
}

func TestSchemaDiff_WithExamples(t *testing.T) {

	require.Equal(t,
		&diff.ValueDiff{
			From: "26734565-dbcc-449a-a370-0beaaf04b0e8",
			To:   "26734565-dbcc-449a-a370-0beaaf04b0e7",
		},
		diff.Get(&diff.Config{IncludeExamples: true}, l(t, 1), l(t, 3)).SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score"].OperationsDiff.Modified["GET"].ParametersDiff.Modified["query"]["token"].SchemaDiff.ExampleDiff)
}

func TestSchemaDiff_MinDiff(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: nil,
			To:   float64(7),
		},
		diff.Get(&diff.Config{Prefix: "/prefix"}, l(t, 4), l(t, 2)).SpecDiff.PathsDiff.Modified["/prefix/api/{domain}/{project}/badges/security-score/"].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInPath]["domain"].SchemaDiff.MinDiff)
}

func TestResponseAdded(t *testing.T) {
	require.Contains(t,
		diff.Get(diff.NewConfig(), l(t, 1), l(t, 3)).SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score"].OperationsDiff.Modified["GET"].ResponsesDiff.Added,
		"default")
}

func TestResponseDeleted(t *testing.T) {
	require.Contains(t,
		diff.Get(diff.NewConfig(), l(t, 3), l(t, 1)).SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score"].OperationsDiff.Modified["GET"].ResponsesDiff.Deleted,
		"default")
}

func TestResponseDescriptionModified(t *testing.T) {
	config := diff.Config{
		IncludeExtensions: diff.StringSet{"x-extension-test": struct{}{}},
	}

	require.Equal(t,
		&diff.ValueDiff{
			From: "Tufin",
			To:   "Tufin1",
		},
		diff.Get(&config, l(t, 3), l(t, 1)).SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/install-command"].OperationsDiff.Modified["GET"].ResponsesDiff.Modified["default"].DescriptionDiff)
}

func TestResponseHeadersModified(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: "Request limit per min.",
			To:   "Request limit per hour.",
		},
		diff.Get(diff.NewConfig(), l(t, 3), l(t, 1)).SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/install-command"].OperationsDiff.Modified["GET"].ResponsesDiff.Modified["default"].HeadersDiff.Modified["X-RateLimit-Limit"].DescriptionDiff)
}

func TestServerAdded(t *testing.T) {
	require.Contains(t,
		diff.Get(diff.NewConfig(), l(t, 5), l(t, 3)).SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/install-command"].OperationsDiff.Modified["GET"].ServersDiff.Added,
		"https://tufin.io/securecloud")
}

func TestServerDeleted(t *testing.T) {
	require.Contains(t,
		diff.Get(diff.NewConfig(), l(t, 3), l(t, 5)).SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/install-command"].OperationsDiff.Modified["GET"].ServersDiff.Deleted,
		"https://tufin.io/securecloud")
}

func TestServerModified(t *testing.T) {
	config := diff.Config{
		IncludeExtensions: diff.StringSet{"x-extension-test": struct{}{}},
	}

	require.Contains(t,
		diff.Get(&config, l(t, 5), l(t, 3)).SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/install-command"].OperationsDiff.Modified["GET"].ServersDiff.Modified,
		"https://www.tufin.io/securecloud")
}

func TestServerAddedToPathItem(t *testing.T) {
	require.Contains(t,
		diff.Get(diff.NewConfig(), l(t, 5), l(t, 3)).SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/install-command"].ServersDiff.Added,
		"https://tufin.io/securecloud")
}

func TestParamAddedToPathItem(t *testing.T) {
	require.Contains(t,
		diff.Get(diff.NewConfig(), l(t, 5), l(t, 3)).SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/install-command"].ParametersDiff.Added[openapi3.ParameterInHeader],
		"name")
}

func TestParamDeletedFromPathItem(t *testing.T) {
	require.Contains(t,
		diff.Get(diff.NewConfig(), l(t, 1), l(t, 2)).SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score"].ParametersDiff.Deleted[openapi3.ParameterInPath],
		"domain")
}

func TestHeaderAdded(t *testing.T) {
	require.Contains(t,
		diff.Get(diff.NewConfig(), l(t, 5), l(t, 1)).SpecDiff.HeadersDiff.Added,
		"new")
}

func TestHeaderDeleted(t *testing.T) {
	require.Contains(t,
		diff.Get(diff.NewConfig(), l(t, 1), l(t, 5)).SpecDiff.HeadersDiff.Deleted,
		"new")
}

func TestRequestBodyModified(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: "number",
			To:   "integer",
		},
		diff.Get(diff.NewConfig(), l(t, 1), l(t, 3)).SpecDiff.RequestBodiesDiff.Modified["reuven"].ContentDiff.SchemaDiff.PropertiesDiff.Modified["meter_value"].TypeDiff,
	)
}

func TestHeaderModifiedSchema(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: false,
			To:   true,
		},
		diff.Get(diff.NewConfig(), l(t, 5), l(t, 1)).SpecDiff.HeadersDiff.Modified["test"].SchemaDiff.AdditionalPropertiesAllowedDiff)
}

func TestHeaderModifiedContent(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: "string",
			To:   "object",
		},
		diff.Get(diff.NewConfig(), l(t, 5), l(t, 1)).SpecDiff.HeadersDiff.Modified["testc"].ContentDiff.SchemaDiff.TypeDiff)
}

func TestResponseContentModified(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: "object",
			To:   "string",
		},
		diff.Get(diff.NewConfig(), l(t, 5), l(t, 1)).SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/badges/security-score"].OperationsDiff.Modified["GET"].ResponsesDiff.Modified["201"].ContentDiff.SchemaDiff.TypeDiff)
}

func TestResponseDespcriptionNil(t *testing.T) {

	s3 := l(t, 3)
	s3.Paths["/api/{domain}/{project}/install-command"].Get.Responses["default"].Value.Description = nil

	require.Equal(t,
		&diff.ValueDiff{
			From: interface{}(nil),
			To:   "Tufin1",
		},
		diff.Get(diff.NewConfig(), s3, l(t, 1)).SpecDiff.PathsDiff.Modified["/api/{domain}/{project}/install-command"].OperationsDiff.Modified["GET"].ResponsesDiff.Modified["default"].DescriptionDiff)
}

func TestSchemaDiff_DeletedCallback(t *testing.T) {
	require.Contains(t,
		diff.Get(diff.NewConfig(), l(t, 3), l(t, 1)).SpecDiff.PathsDiff.Modified["/register"].OperationsDiff.Modified["POST"].CallbacksDiff.Deleted,
		"myEvent")
}

func TestSchemaDiff_ModifiedCallback(t *testing.T) {
	require.Contains(t,
		diff.Get(diff.NewConfig(), l(t, 3), l(t, 1)).SpecDiff.PathsDiff.Modified["/subscribe"].OperationsDiff.Modified["POST"].CallbacksDiff.Modified["myEvent"].Deleted,
		"{$request.body#/callbackUrl}")
}

func TestSchemaDiff_AddedSchemas(t *testing.T) {
	require.Contains(t,
		diff.Get(diff.NewConfig(), l(t, 1), l(t, 5)).SpecDiff.SchemasDiff.Added,
		"requests")
}

func TestSchemaDiff_DeletedSchemas(t *testing.T) {
	require.Contains(t,
		diff.Get(diff.NewConfig(), l(t, 5), l(t, 1)).SpecDiff.SchemasDiff.Deleted,
		"requests")
}

func TestSchemaDiff_ModifiedSchemas(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: true,
			To:   false,
		},
		diff.Get(diff.NewConfig(), l(t, 1), l(t, 5)).SpecDiff.SchemasDiff.Modified["network-policies"].AdditionalPropertiesAllowedDiff)
}

func TestSchemaDiff_ModifiedSchemasOldNil(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: nil,
			To:   false,
		},
		diff.Get(diff.NewConfig(), l(t, 1), l(t, 5)).SpecDiff.SchemasDiff.Modified["rules"].AdditionalPropertiesAllowedDiff)
}

func TestSchemaDiff_ModifiedSchemasNewNil(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: false,
			To:   nil,
		},
		diff.Get(diff.NewConfig(), l(t, 5), l(t, 1)).SpecDiff.SchemasDiff.Modified["rules"].AdditionalPropertiesAllowedDiff)
}

func TestSchemaDiff_ModifiedSchemasValueDeleted(t *testing.T) {
	s5 := l(t, 5)
	s5.Components.Schemas["network-policies"].Value = nil

	require.Equal(t,
		true,
		diff.Get(diff.NewConfig(), l(t, 1), s5).SpecDiff.SchemasDiff.Modified["network-policies"].ValueDeleted)
}

func TestSchemaDiff_ModifiedSchemasValueAdded(t *testing.T) {
	s5 := l(t, 5)
	s5.Components.Schemas["network-policies"].Value = nil

	require.Equal(t,
		true,
		diff.Get(diff.NewConfig(), s5, l(t, 1)).SpecDiff.SchemasDiff.Modified["network-policies"].ValueAdded)
}

func TestSchemaDiff_ModifiedSchemasBothValuesNil(t *testing.T) {
	s5 := l(t, 5)
	s5.Components.Schemas["network-policies"].Value = nil

	require.False(t, diff.Get(diff.NewConfig(), s5, s5).Summary.Diff)
}

func TestSummary(t *testing.T) {

	d := diff.Get(diff.NewConfig(), l(t, 1), l(t, 2)).Summary

	require.Equal(t, diff.SummaryDetails{0, 3, 1}, d.GetSummaryDetails(diff.PathsDetail))
	require.Equal(t, diff.SummaryDetails{0, 1, 0}, d.GetSummaryDetails(diff.SecurityDetail))
	require.Equal(t, diff.SummaryDetails{0, 1, 0}, d.GetSummaryDetails(diff.ServersDetail))
	require.Equal(t, diff.SummaryDetails{0, 2, 0}, d.GetSummaryDetails(diff.TagsDetail))
	require.Equal(t, diff.SummaryDetails{0, 2, 0}, d.GetSummaryDetails(diff.SchemasDetail))
	require.Equal(t, diff.SummaryDetails{0, 1, 0}, d.GetSummaryDetails(diff.ParametersDetail))
	require.Equal(t, diff.SummaryDetails{0, 3, 0}, d.GetSummaryDetails(diff.HeadersDetail))
	require.Equal(t, diff.SummaryDetails{0, 1, 0}, d.GetSummaryDetails(diff.RequestBodiesDetail))
	require.Equal(t, diff.SummaryDetails{0, 1, 0}, d.GetSummaryDetails(diff.ResponsesDetail))
	require.Equal(t, diff.SummaryDetails{0, 2, 0}, d.GetSummaryDetails(diff.SecuritySchemesDetail))
	require.Equal(t, diff.SummaryDetails{}, d.GetSummaryDetails(diff.CallbacksDetail))
}

func TestSummaryInvalidComponent(t *testing.T) {

	require.Equal(t, diff.SummaryDetails{
		Added:    0,
		Deleted:  0,
		Modified: 0,
	}, diff.Get(diff.NewConfig(), l(t, 1), l(t, 2)).Summary.GetSummaryDetails("invalid"))
}

func TestFilterByRegex(t *testing.T) {
	require.Nil(t, diff.Get(&diff.Config{Filter: "x"}, l(t, 1), l(t, 2)).Summary.Details[diff.PathsDetail])
}

func TestFilterByRegex_Invalid(t *testing.T) {
	require.Equal(t, true, diff.Get(&diff.Config{Filter: "["}, l(t, 1), l(t, 2)).Summary.Diff)
}

func TestAddedSecurityRequireent(t *testing.T) {
	require.Contains(t,
		diff.Get(&diff.Config{}, l(t, 3), l(t, 1)).SpecDiff.PathsDiff.Modified["/register"].OperationsDiff.Modified["POST"].SecurityDiff.Added,
		"bearerAuth")
}
