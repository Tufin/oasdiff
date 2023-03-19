package checker_test

import (
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// BC: no change is not breaking
func TestBreaking_Same(t *testing.T) {
	require.Empty(t, d(t, &diff.Config{BreakingOnly: true}, 1, 1))
}

// BC: adding an optional request body is not breaking
func TestBreaking_AddingOptionalRequestBody(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Spec.Paths[installCommandPath].Get.RequestBody = &openapi3.RequestBodyRef{
		Value: openapi3.NewRequestBody().WithRequired(false),
	}

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: changing an existing request body from required to optional is not breaking
func TestBreaking_RequestBodyRequiredDisabled(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths[installCommandPath].Get.RequestBody = &openapi3.RequestBodyRef{
		Value: openapi3.NewRequestBody().WithRequired(true),
	}

	s2.Spec.Paths[installCommandPath].Get.RequestBody = &openapi3.RequestBodyRef{
		Value: openapi3.NewRequestBody().WithRequired(false),
	}

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: deleting a tag is not breaking
func TestBreaking_DeletedTag(t *testing.T) {
	r := d(t, &diff.Config{}, 1, 5)
	require.Len(t, r, 6)
	require.Equal(t, "request-parameter-removed", r[0].Id)
	require.Equal(t, "optional-response-header-removed", r[1].Id)
	require.Equal(t, "response-success-status-removed", r[2].Id)
	require.Equal(t, "response-body-type-changed", r[3].Id)
	require.Equal(t, "api-path-removed-without-deprecation", r[4].Id)
	require.Equal(t, "api-path-removed-without-deprecation", r[5].Id)
}

// BC: adding an enum value is not breaking
func TestBreaking_AddedEnum(t *testing.T) {
	r := d(t, &diff.Config{}, 1, 3)
	require.Len(t, r, 6)
	require.Equal(t, "request-parameter-removed", r[0].Id)
	require.Equal(t, "request-parameter-removed", r[1].Id)
	require.Equal(t, "request-parameter-removed", r[2].Id)
	require.Equal(t, "request-parameter-removed", r[3].Id)
	require.Equal(t, "response-success-status-removed", r[4].Id)
	require.Equal(t, "response-success-status-removed", r[5].Id)
}

// BC: changing extensions is not breaking
func TestBreaking_ModifiedExtension(t *testing.T) {
	r := d(t, &diff.Config{}, 1, 3)
	require.Len(t, r, 6)
	require.Equal(t, "request-parameter-removed", r[0].Id)
	require.Equal(t, "request-parameter-removed", r[1].Id)
	require.Equal(t, "request-parameter-removed", r[2].Id)
	require.Equal(t, "request-parameter-removed", r[3].Id)
	require.Equal(t, "response-success-status-removed", r[4].Id)
	require.Equal(t, "response-success-status-removed", r[5].Id)
}

// BC: changing comments is not breaking
func TestBreaking_Comments(t *testing.T) {
	r := d(t, &diff.Config{}, 1, 3)
	require.Len(t, r, 6)
	require.Equal(t, "request-parameter-removed", r[0].Id)
	require.Equal(t, "request-parameter-removed", r[1].Id)
	require.Equal(t, "request-parameter-removed", r[2].Id)
	require.Equal(t, "request-parameter-removed", r[3].Id)
	require.Equal(t, "response-success-status-removed", r[4].Id)
	require.Equal(t, "response-success-status-removed", r[5].Id)
}

// BC: new optional header param is not breaking
func TestBreaking_NewOptionalHeaderParam(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	deleteParam(s1.Spec.Paths[installCommandPath].Get, openapi3.ParameterInHeader, "network-policies")
	s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = false

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: changing an existing header param to optional is not breaking
func TestBreaking_HeaderParamRequiredDisabled(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = true
	s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = false

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

func deleteResponseHeader(response *openapi3.Response, name string) {
	delete(response.Headers, name)
}

// BC: new required response header param is not breaking
func TestBreaking_NewRequiredResponseHeader(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	deleteResponseHeader(s1.Spec.Paths[installCommandPath].Get.Responses["default"].Value, "X-RateLimit-Limit")
	s2.Spec.Paths[installCommandPath].Get.Responses["default"].Value.Headers["X-RateLimit-Limit"].Value.Required = true

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: changing operation ID is not breaking
func TestBreaking_OperationID(t *testing.T) {
	r := d(t, &diff.Config{}, 3, 1)
	require.Len(t, r, 3)
	require.Equal(t, "request-parameter-pattern-added", r[0].Id)
	require.Equal(t, "request-parameter-max-length-decreased", r[1].Id)
	require.Equal(t, "request-parameter-enum-value-removed", r[2].Id)
}

// BC: changing a link to operation ID is not breaking
func TestBreaking_LinkOperationID(t *testing.T) {
	r := d(t, &diff.Config{}, 3, 1)
	require.Len(t, r, 3)
	require.Equal(t, "request-parameter-pattern-added", r[0].Id)
	require.Equal(t, "request-parameter-max-length-decreased", r[1].Id)
	require.Equal(t, "request-parameter-enum-value-removed", r[2].Id)
}

// BC: adding a media-type to response is not breaking
func TestBreaking_ResponseAddMediaType(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile("../data/response-media-type-revision.yaml")
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile("../data/response-media-type-base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: deprecating an operation is not breaking
func TestBreaking_DeprecatedOperation(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Spec.Paths[installCommandPath].Get.Deprecated = true
	s2.Spec.Paths[installCommandPath].Get.Extensions[diff.SunsetExtension] = toJson(t, civil.DateOf(time.Now()).AddDays(200).String())

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: deprecating a parameter is not breaking
func TestBreaking_DeprecatedParameter(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Deprecated = true

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: deprecating a header is not breaking
func TestBreaking_DeprecatedHeader(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Spec.Paths[installCommandPath].Get.Responses["default"].Value.Headers["X-RateLimit-Limit"].Value.Deprecated = true

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: deprecating a schema is not breaking
func TestBreaking_DeprecatedSchema(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Deprecated = true

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: changing servers is not breaking
func TestBreaking_Servers(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile("../data/servers/baseswagger.json")
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile("../data/servers/revisionswagger.json")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}
