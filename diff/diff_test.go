package diff_test

import (
	"fmt"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
)

const securityScorePath = "/api/{domain}/{project}/badges/security-score"

func l(t *testing.T, v int) *openapi3.Swagger {
	loader := openapi3.NewSwaggerLoader()
	oas, err := loader.LoadSwaggerFromFile(fmt.Sprintf("../data/openapi-test%d.yaml", v))
	require.NoError(t, err)
	return oas
}

func d(t *testing.T, config *diff.Config, v1, v2 int) *diff.Diff {
	d, err := diff.Get(config, l(t, v1), l(t, v2))
	require.NoError(t, err)
	return d
}

func TestDiff_Same(t *testing.T) {
	require.Nil(t, d(t, diff.NewConfig(), 1, 1))
}

func TestDiff_Empty(t *testing.T) {
	require.True(t, (*diff.CallbacksDiff)(nil).Empty())
	require.True(t, (*diff.EncodingsDiff)(nil).Empty())
	require.True(t, (*diff.ExtensionsDiff)(nil).Empty())
	require.True(t, (*diff.HeadersDiff)(nil).Empty())
	require.True(t, (*diff.OperationsDiff)(nil).Empty())
	require.True(t, (*diff.ParametersDiff)(nil).Empty())
	require.True(t, (*diff.RequestBodiesDiff)(nil).Empty())
	require.True(t, (*diff.ResponsesDiff)(nil).Empty())
	require.True(t, (*diff.SchemasDiff)(nil).Empty())
	require.True(t, (*diff.ServersDiff)(nil).Empty())
	require.True(t, (*diff.StringsDiff)(nil).Empty())
	require.True(t, (*diff.StringMapDiff)(nil).Empty())
	require.True(t, (*diff.TagsDiff)(nil).Empty())
	require.True(t, (*diff.SecurityRequirementsDiff)(nil).Empty())
	require.True(t, (*diff.SecuritySchemesDiff)(nil).Empty())
	require.True(t, (*diff.ExamplesDiff)(nil).Empty())
}

func TestDiff_DeletedPaths(t *testing.T) {
	require.ElementsMatch(t,
		[]string{"/api/{domain}/{project}/install-command", "/register", "/subscribe"},
		d(t, diff.NewConfig(), 1, 2).PathsDiff.Deleted)
}

func TestDiff_AddedOperation(t *testing.T) {
	require.Contains(t,
		d(t, diff.NewConfig(), 1, 2).PathsDiff.Modified[securityScorePath].OperationsDiff.Added,
		"POST")
}

func TestDiff_DeletedOperation(t *testing.T) {
	require.Contains(t,
		d(t, diff.NewConfig(), 2, 1).PathsDiff.Modified["/api/{domain}/{project}/badges/security-score/"].OperationsDiff.Deleted,
		"POST")
}

func TestAddedExtension(t *testing.T) {
	config := diff.Config{
		IncludeExtensions: diff.StringSet{"x-extension-test": struct{}{}},
	}

	require.Contains(t,
		d(t, &config, 3, 1).ExtensionsDiff.Added,
		"x-extension-test")
}

func TestDeletedExtension(t *testing.T) {
	config := diff.Config{
		IncludeExtensions: diff.StringSet{"x-extension-test": struct{}{}},
	}

	require.Contains(t,
		d(t, &config, 1, 3).ExtensionsDiff.Deleted,
		"x-extension-test")
}

func TestModifiedExtension(t *testing.T) {
	config := diff.Config{
		IncludeExtensions: diff.StringSet{"x-extension-test2": struct{}{}},
	}

	require.NotNil(t,
		d(t, &config, 1, 3).ExtensionsDiff.Modified["x-extension-test2"])
}

func TestExcludedExtension(t *testing.T) {
	require.Nil(t,
		d(t, diff.NewConfig(), 1, 3).ExtensionsDiff)
}

func TestDiff_AddedGlobalTag(t *testing.T) {
	require.Contains(t,
		d(t, diff.NewConfig(), 3, 1).TagsDiff.Added,
		"security")
}

func TestDiff_ModifiedGlobalTag(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: "Harrison",
			To:   "harrison",
		},
		d(t, diff.NewConfig(), 1, 3).TagsDiff.Modified["reuven"].DescriptionDiff)
}

func TestDiff_AddedTag(t *testing.T) {
	require.Contains(t,
		d(t, diff.NewConfig(), 3, 1).PathsDiff.Modified[securityScorePath].OperationsDiff.Modified["GET"].TagsDiff.Added,
		"security")
}

func TestDiff_DeletedEncoding(t *testing.T) {
	require.Contains(t,
		d(t, diff.NewConfig(), 1, 3).PathsDiff.Modified["/subscribe"].OperationsDiff.Modified["POST"].CallbacksDiff.Modified["myEvent"].Modified["hi"].OperationsDiff.Modified["POST"].RequestBodyDiff.ContentDiff.EncodingsDiff.Deleted,
		"historyMetadata")
}

func TestDiff_ModifiedEncodingHeaders(t *testing.T) {
	require.NotNil(t,
		d(t, diff.NewConfig(), 3, 1).PathsDiff.Modified["/subscribe"].OperationsDiff.Modified["POST"].CallbacksDiff.Modified["myEvent"].Modified["hi"].OperationsDiff.Modified["POST"].RequestBodyDiff.ContentDiff.EncodingsDiff.Modified["profileImage"].HeadersDiff,
		"profileImage")
}

func TestDiff_AddedParam(t *testing.T) {
	require.Contains(t,
		d(t, diff.NewConfig(), 2, 1).PathsDiff.Modified["/api/{domain}/{project}/badges/security-score/"].OperationsDiff.Modified["GET"].ParametersDiff.Added[openapi3.ParameterInHeader],
		"X-Auth-Name")
}

func TestDiff_DeletedParam(t *testing.T) {
	require.Contains(t,
		d(t, diff.NewConfig(), 1, 2).PathsDiff.Modified[securityScorePath].OperationsDiff.Modified["GET"].ParametersDiff.Deleted[openapi3.ParameterInHeader],
		"X-Auth-Name")
}

func TestDiff_ModifiedParam(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: true,
			To:   (interface{})(nil),
		},
		d(t, diff.NewConfig(), 2, 1).PathsDiff.Modified["/api/{domain}/{project}/badges/security-score/"].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInQuery]["image"].ExplodeDiff)
}

func TestSchemaDiff_TypeDiff(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: "string",
			To:   "integer",
		},
		d(t, diff.NewConfig(), 1, 2).PathsDiff.Modified[securityScorePath].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInPath]["domain"].SchemaDiff.TypeDiff)
}

func TestSchemaDiff_EnumDiff(t *testing.T) {
	require.Equal(t,
		&diff.EnumDiff{
			Added:   diff.EnumValues{"test1"},
			Deleted: diff.EnumValues{},
		},
		d(t, diff.NewConfig(), 1, 3).PathsDiff.Modified["/api/{domain}/{project}/install-command"].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInPath]["project"].SchemaDiff.EnumDiff)
}

func TestSchemaDiff_RequiredAdded(t *testing.T) {
	require.Contains(t,
		d(t, diff.NewConfig(), 1, 5).PathsDiff.Modified[securityScorePath].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInQuery]["filter"].ContentDiff.SchemaDiff.RequiredDiff.Added,
		"type")
}

func TestSchemaDiff_RequiredDeleted(t *testing.T) {
	require.Contains(t,
		d(t, diff.NewConfig(), 5, 1).PathsDiff.Modified[securityScorePath].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInQuery]["filter"].ContentDiff.SchemaDiff.RequiredDiff.Deleted,
		"type")
}

func TestSchemaDiff_NotDiff(t *testing.T) {
	require.Equal(t,
		true,
		d(t, diff.NewConfig(), 1, 3).PathsDiff.Modified[securityScorePath].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInQuery]["image"].SchemaDiff.NotDiff.SchemaAdded)
}

func TestSchemaDiff_ContentDiff(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: "number",
			To:   "string",
		},
		d(t, diff.NewConfig(), 2, 1).PathsDiff.Modified["/api/{domain}/{project}/badges/security-score/"].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInQuery]["filter"].ContentDiff.SchemaDiff.PropertiesDiff.Modified["color"].TypeDiff)
}

func TestSchemaDiff_MediaTypeAdded(t *testing.T) {
	require.Equal(t,
		true,
		d(t, diff.NewConfig(), 5, 1).PathsDiff.Modified[securityScorePath].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInHeader]["user"].ContentDiff.MediaTypeAdded)
}

func TestSchemaDiff_MediaTypeDeleted(t *testing.T) {
	require.Equal(t,
		false,
		d(t, diff.NewConfig(), 1, 5).PathsDiff.Modified[securityScorePath].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInHeader]["user"].ContentDiff.MediaTypeAdded)
}

func TestSchemaDiff_MediaTypeModified(t *testing.T) {
	require.Equal(t,
		true,
		d(t, diff.NewConfig(), 1, 5).PathsDiff.Modified[securityScorePath].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInCookie]["test"].ContentDiff.MediaTypeDiff)
}

func TestSchemaDiff_AnyOfModified(t *testing.T) {
	require.False(t,
		d(t, &diff.Config{PathPrefix: "/prefix"}, 4, 2).PathsDiff.Modified["/prefix/api/{domain}/{project}/badges/security-score/"].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInQuery]["token"].SchemaDiff.AnyOfDiff.Empty())
}

func TestSchemaDiff_WithExamples(t *testing.T) {

	require.Equal(t,
		&diff.ValueDiff{
			From: "26734565-dbcc-449a-a370-0beaaf04b0e8",
			To:   "26734565-dbcc-449a-a370-0beaaf04b0e7",
		},
		d(t, diff.NewConfig(), 1, 3).PathsDiff.Modified[securityScorePath].OperationsDiff.Modified["GET"].ParametersDiff.Modified["query"]["token"].SchemaDiff.ExampleDiff)
}

func TestSchemaDiff_MinDiff(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: nil,
			To:   float64(7),
		},
		d(t, &diff.Config{PathPrefix: "/prefix"}, 4, 2).PathsDiff.Modified["/prefix/api/{domain}/{project}/badges/security-score/"].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInPath]["domain"].SchemaDiff.MinDiff)
}

func TestResponseAdded(t *testing.T) {
	require.Contains(t,
		d(t, diff.NewConfig(), 1, 3).PathsDiff.Modified[securityScorePath].OperationsDiff.Modified["GET"].ResponsesDiff.Added,
		"default")
}

func TestResponseDeleted(t *testing.T) {
	require.Contains(t,
		d(t, diff.NewConfig(), 3, 1).PathsDiff.Modified[securityScorePath].OperationsDiff.Modified["GET"].ResponsesDiff.Deleted,
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
		d(t, &config, 3, 1).PathsDiff.Modified["/api/{domain}/{project}/install-command"].OperationsDiff.Modified["GET"].ResponsesDiff.Modified["default"].DescriptionDiff)
}

func TestResponseHeadersModified(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: "Request limit per min.",
			To:   "Request limit per hour.",
		},
		d(t, diff.NewConfig(), 3, 1).PathsDiff.Modified["/api/{domain}/{project}/install-command"].OperationsDiff.Modified["GET"].ResponsesDiff.Modified["default"].HeadersDiff.Modified["X-RateLimit-Limit"].DescriptionDiff)
}

func TestServerAdded(t *testing.T) {
	require.Contains(t,
		d(t, diff.NewConfig(), 5, 3).PathsDiff.Modified["/api/{domain}/{project}/install-command"].OperationsDiff.Modified["GET"].ServersDiff.Added,
		"https://tufin.io/securecloud")
}

func TestServerDeleted(t *testing.T) {
	require.Contains(t,
		d(t, diff.NewConfig(), 3, 5).PathsDiff.Modified["/api/{domain}/{project}/install-command"].OperationsDiff.Modified["GET"].ServersDiff.Deleted,
		"https://tufin.io/securecloud")
}

func TestServerModified(t *testing.T) {
	config := diff.Config{
		IncludeExtensions: diff.StringSet{"x-extension-test": struct{}{}},
	}

	require.Contains(t,
		d(t, &config, 5, 3).PathsDiff.Modified["/api/{domain}/{project}/install-command"].OperationsDiff.Modified["GET"].ServersDiff.Modified,
		"https://www.tufin.io/securecloud")
}

func TestServerAddedToPathItem(t *testing.T) {
	require.Contains(t,
		d(t, diff.NewConfig(), 5, 3).PathsDiff.Modified["/api/{domain}/{project}/install-command"].ServersDiff.Added,
		"https://tufin.io/securecloud")
}

func TestParamAddedToPathItem(t *testing.T) {
	require.Contains(t,
		d(t, diff.NewConfig(), 5, 3).PathsDiff.Modified["/api/{domain}/{project}/install-command"].ParametersDiff.Added[openapi3.ParameterInHeader],
		"name")
}

func TestParamDeletedFromPathItem(t *testing.T) {
	require.Contains(t,
		d(t, diff.NewConfig(), 1, 2).PathsDiff.Modified[securityScorePath].ParametersDiff.Deleted[openapi3.ParameterInPath],
		"domain")
}

func TestHeaderAdded(t *testing.T) {
	require.Contains(t,
		d(t, diff.NewConfig(), 5, 1).HeadersDiff.Added,
		"new")
}

func TestHeaderDeleted(t *testing.T) {
	require.Contains(t,
		d(t, diff.NewConfig(), 1, 5).HeadersDiff.Deleted,
		"new")
}

func TestRequestBodyModified(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: "number",
			To:   "integer",
		},
		d(t, diff.NewConfig(), 1, 3).RequestBodiesDiff.Modified["reuven"].ContentDiff.SchemaDiff.PropertiesDiff.Modified["meter_value"].TypeDiff,
	)
}

func TestHeaderModifiedSchema(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: false,
			To:   true,
		},
		d(t, diff.NewConfig(), 5, 1).HeadersDiff.Modified["test"].SchemaDiff.AdditionalPropertiesAllowedDiff)
}

func TestHeaderModifiedContent(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: "string",
			To:   "object",
		},
		d(t, diff.NewConfig(), 5, 1).HeadersDiff.Modified["testc"].ContentDiff.SchemaDiff.TypeDiff)
}

func TestResponseContentModified(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: "object",
			To:   "string",
		},
		d(t, diff.NewConfig(), 5, 1).PathsDiff.Modified[securityScorePath].OperationsDiff.Modified["GET"].ResponsesDiff.Modified["201"].ContentDiff.SchemaDiff.TypeDiff)
}

func TestResponseDespcriptionNil(t *testing.T) {
	s3 := l(t, 3)
	s3.Paths["/api/{domain}/{project}/install-command"].Get.Responses["default"].Value.Description = nil

	d, err := diff.Get(diff.NewConfig(), s3, l(t, 1))
	require.NoError(t, err)

	require.Equal(t,
		&diff.ValueDiff{
			From: interface{}(nil),
			To:   "Tufin1",
		},
		d.PathsDiff.Modified["/api/{domain}/{project}/install-command"].OperationsDiff.Modified["GET"].ResponsesDiff.Modified["default"].DescriptionDiff)
}

func TestSchemaDiff_DeletedCallback(t *testing.T) {
	require.Contains(t,
		d(t, diff.NewConfig(), 3, 1).PathsDiff.Modified["/register"].OperationsDiff.Modified["POST"].CallbacksDiff.Deleted,
		"myEvent")
}

func TestSchemaDiff_ModifiedCallback(t *testing.T) {
	require.Contains(t,
		d(t, diff.NewConfig(), 3, 1).PathsDiff.Modified["/subscribe"].OperationsDiff.Modified["POST"].CallbacksDiff.Modified["myEvent"].Deleted,
		"{$request.body#/callbackUrl}")
}

func TestSchemaDiff_AddedRequestBody(t *testing.T) {
	require.True(t,
		d(t, diff.NewConfig(), 3, 1).PathsDiff.Modified["/subscribe"].OperationsDiff.Modified["POST"].CallbacksDiff.Modified["myEvent"].Modified["bye"].OperationsDiff.Modified["POST"].RequestBodyDiff.Added)
}

func TestSchemaDiff_AddedSchemas(t *testing.T) {
	require.Contains(t,
		d(t, diff.NewConfig(), 1, 5).SchemasDiff.Added,
		"requests")
}

func TestSchemaDiff_DeletedSchemas(t *testing.T) {
	require.Contains(t,
		d(t, diff.NewConfig(), 5, 1).SchemasDiff.Deleted,
		"requests")
}

func TestSchemaDiff_ModifiedSchemas(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: true,
			To:   false,
		},
		d(t, diff.NewConfig(), 1, 5).SchemasDiff.Modified["network-policies"].AdditionalPropertiesAllowedDiff)
}

func TestSchemaDiff_ModifiedSchemasOldNil(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: nil,
			To:   false,
		},
		d(t, diff.NewConfig(), 1, 5).SchemasDiff.Modified["rules"].AdditionalPropertiesAllowedDiff)
}

func TestSchemaDiff_ModifiedSchemasNewNil(t *testing.T) {
	require.Equal(t,
		&diff.ValueDiff{
			From: false,
			To:   nil,
		},
		d(t, diff.NewConfig(), 5, 1).SchemasDiff.Modified["rules"].AdditionalPropertiesAllowedDiff)
}

func TestSummary(t *testing.T) {

	d := d(t, diff.NewConfig(), 1, 2).GetSummary()

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
	require.Equal(t, diff.SummaryDetails{}, d.GetSummaryDetails(diff.ExamplesDetail))
	require.Equal(t, diff.SummaryDetails{}, d.GetSummaryDetails(diff.LinksDetail))
	require.Equal(t, diff.SummaryDetails{}, d.GetSummaryDetails(diff.CallbacksDetail))
}

func TestSummaryInvalidComponent(t *testing.T) {
	require.Equal(t, diff.SummaryDetails{
		Added:    0,
		Deleted:  0,
		Modified: 0,
	}, d(t, diff.NewConfig(), 1, 2).GetSummary().GetSummaryDetails("invalid"))
}

func TestFilterByRegex(t *testing.T) {
	d, err := diff.Get(&diff.Config{PathFilter: "x"}, l(t, 1), l(t, 2))
	require.NoError(t, err)
	require.Nil(t, d.GetSummary().Details[diff.PathsDetail])
}

func TestAddedSecurityRequirement(t *testing.T) {
	require.Contains(t,
		d(t, diff.NewConfig(), 3, 1).PathsDiff.Modified["/register"].OperationsDiff.Modified["POST"].SecurityDiff.Added,
		"bearerAuth")
}

func TestModifiedSecurityRequirement(t *testing.T) {
	securityScopesDiff := d(t, diff.NewConfig(), 3, 1).PathsDiff.Modified["/register"].OperationsDiff.Modified["POST"].SecurityDiff.Modified["OAuth"]
	require.NotEmpty(t, securityScopesDiff)

	require.Contains(t,
		securityScopesDiff["OAuth"].Deleted,
		"write:pets")
}
