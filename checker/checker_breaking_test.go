package checker_test

import (
	"fmt"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
	"github.com/tufin/oasdiff/utils"
)

const (
	securityScorePath      = "/api/{domain}/{project}/badges/security-score"
	securityScorePathSlash = securityScorePath + "/"
	installCommandPath     = "/api/{domain}/{project}/install-command"
)

func l(t *testing.T, v int) load.OpenAPISpecInfo {
	t.Helper()
	loader := openapi3.NewLoader()
	oas, err := loader.LoadFromFile(fmt.Sprintf("../data/openapi-test%d.yaml", v))
	require.NoError(t, err)
	return load.OpenAPISpecInfo{Spec: oas, Url: fmt.Sprintf("../data/openapi-test%d.yaml", v)}
}

func d(t *testing.T, config *diff.Config, v1, v2 int) []checker.BackwardCompatibilityError {
	t.Helper()
	l1 := l(t, v1)
	l2 := l(t, v2)
	d, osm, err := diff.GetWithOperationsSourcesMap(config, &l1, &l2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	return errs
}

// BC: deleting a path is breaking
func TestBreaking_DeletedPath(t *testing.T) {
	errs := d(t, &diff.Config{}, 1, 701)
	require.Len(t, errs, 1)
	require.Equal(t, "api-path-removed-without-deprecation", errs[0].Id)
}

// BC: deleting an operation is breaking
func TestBreaking_DeletedOp(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths[installCommandPath].Put = openapi3.NewOperation()

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "api-removed-without-deprecation", errs[0].Id)
}

// BC: adding a required request body is breaking
func TestBreaking_AddingRequiredRequestBody(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Spec.Paths[installCommandPath].Get.RequestBody = &openapi3.RequestBodyRef{
		Value: openapi3.NewRequestBody().WithRequired(true),
	}

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "added-required-request-body", errs[0].Id)
}

// BC: changing an existing request body from optional to required is breaking
func TestBreaking_RequestBodyRequiredEnabled(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths[installCommandPath].Get.RequestBody = &openapi3.RequestBodyRef{
		Value: openapi3.NewRequestBody().WithRequired(false),
	}

	s2.Spec.Paths[installCommandPath].Get.RequestBody = &openapi3.RequestBodyRef{
		Value: openapi3.NewRequestBody().WithRequired(true),
	}

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "request-body-became-required", errs[0].Id)
}

// BC: deleting an enum value is breaking
func TestBreaking_DeletedEnum(t *testing.T) {
	errs := d(t, &diff.Config{}, 702, 1)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "request-parameter-enum-value-removed", errs[0].Id)
}

// BC: added an enum value to response breaking
func TestBreaking_AddedResponseEnum(t *testing.T) {
	errs := d(t, &diff.Config{}, 703, 704)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 2)
	require.Equal(t, "response-property-enum-value-added", errs[0].Id)
	require.Equal(t, "response-property-enum-value-added", errs[1].Id)
}

func deleteParam(op *openapi3.Operation, in string, name string) {

	result := openapi3.NewParameters()

	for _, item := range op.Parameters {
		if v := item.Value; v != nil {
			if v.Name == name && v.In == in {
				continue
			}
			result = append(result, item)
		}
	}
	op.Parameters = result
}

// BC: renaming a path parameter is not breaking
func TestBreaking_PathParamRename(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile("../data/param-rename/method-base.yaml")
	require.NoError(t, err)

	s2, err := loader.LoadFromFile("../data/param-rename/method-revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{},
		&load.OpenAPISpecInfo{Spec: s1},
		&load.OpenAPISpecInfo{Spec: s2},
	)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)

	require.Empty(t, errs)
}

// BC: new required path param is breaking
func TestBreaking_NewPathParam(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	deleteParam(s1.Spec.Paths[installCommandPath].Get, openapi3.ParameterInPath, "project")
	// note: path params are always required

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)

	require.Len(t, errs, 1)
	require.Equal(t, "new-request-path-parameter", errs[0].Id)
}

// BC: new required header param is breaking
func TestBreaking_NewRequiredHeaderParam(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	deleteParam(s1.Spec.Paths[installCommandPath].Get, openapi3.ParameterInHeader, "network-policies")
	s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = true

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "new-required-request-parameter", errs[0].Id)
}

// BC: changing an existing header param from optional to required is breaking
func TestBreaking_HeaderParamRequiredEnabled(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = false
	s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = true

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Equal(t, checker.BackwardCompatibilityErrors{
		{
			Id:        "request-parameter-became-required",
			Text:      "the 'header' request parameter 'network-policies' became required",
			Comment:   "",
			Level:     checker.ERR,
			Operation: "GET",
			Path:      "/api/{domain}/{project}/install-command",
			Source:    "../data/openapi-test1.yaml",
		}}, errs)
}

// BC: changing an existing response header from required to optional is breaking
func TestBreaking_ResponseHeaderParamRequiredDisabled(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths[installCommandPath].Get.Responses["default"].Value.Headers["X-RateLimit-Limit"].Value.Required = true
	s2.Spec.Paths[installCommandPath].Get.Responses["default"].Value.Headers["X-RateLimit-Limit"].Value.Required = false

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "response-header-became-optional", errs[0].Id)
}

// BC: removing an existing required response header is breaking as error
func TestBreaking_ResponseHeaderRemoved(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths[installCommandPath].Get.Responses["default"].Value.Headers["X-RateLimit-Limit"].Value.Required = true
	delete(s2.Spec.Paths[installCommandPath].Get.Responses["default"].Value.Headers, "X-RateLimit-Limit")

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	for _, err := range errs {
		require.Equal(t, checker.ERR, err.Level)
	}
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "required-response-header-removed", errs[0].Id)
}

// BC: removing an existing response with successful status is breaking
func TestBreaking_ResponseSuccessStatusRemoved(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	delete(s2.Spec.Paths[securityScorePath].Get.Responses, "200")

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	for _, err := range errs {
		require.Equal(t, checker.ERR, err.Level)
	}
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "response-success-status-removed", errs[0].Id)
}

// BC: removing an existing response with non-successful status is breaking (optional)
func TestBreaking_ResponseNonSuccessStatusRemoved(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	delete(s2.Spec.Paths[securityScorePath].Get.Responses, "400")

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetChecks(utils.StringList{"response-non-success-status-removed"}), d, osm)
	for _, err := range errs {
		require.Equal(t, checker.ERR, err.Level)
	}
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "response-non-success-status-removed", errs[0].Id)
}

// BC: removing/updating an operation id is breaking (optional)
func TestBreaking_OperationIdRemoved(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Spec.Paths[securityScorePath].Get.OperationID = "newOperationId"

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetChecks(utils.StringList{"api-operation-id-removed"}), d, osm)
	for _, err := range errs {
		require.Equal(t, checker.ERR, err.Level)
	}
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "api-operation-id-removed", errs[0].Id)
}

// BC: removing/updating an enum in request body is breaking (optional)
func TestBreaking_RequestBodyEnumRemoved(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile("../data/enums/request-body-enum.yaml")
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile("../data/enums/request-body-enum.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/api/v2/changeOfRequestFieldValueTiedToEnumTest"].Get.RequestBody.Value.Content["application/json"].Schema.Value.Enum = []interface{}{}

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibility(checker.GetChecks(utils.StringList{"request-body-enum-value-removed"}), d, osm)
	for _, err := range errs {
		require.Equal(t, checker.ERR, err.Level)
	}

	require.Len(t, errs, 3)
	require.Equal(t, "request-body-enum-value-removed", errs[0].Id)
}

// BC: removing/updating a property enum in response is breaking (optional)
func TestBreaking_ResponsePropertyEnumRemoved(t *testing.T) {
	s1 := l(t, 704)
	s2 := l(t, 703)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibility(checker.GetChecks(utils.StringList{"response-property-enum-value-removed"}), d, osm)
	for _, err := range errs {
		require.Equal(t, checker.ERR, err.Level)
	}
	require.NotEmpty(t, errs)
	require.Len(t, errs, 2)
	require.Equal(t, "response-property-enum-value-removed", errs[0].Id)
}

// BC: removing/updating a tag is breaking (optional)
func TestBreaking_TagRemoved(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Spec.Paths[securityScorePath].Get.Tags[0] = "newTag"

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetChecks(utils.StringList{"api-tag-removed"}), d, osm)
	for _, err := range errs {
		require.Equal(t, checker.ERR, err.Level)
	}
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "api-tag-removed", errs[0].Id)
}

// BC: removing/updating a media type enum in response (optional)
func TestBreaking_ResponseMediaTypeEnumRemoved(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile("../data/enums/response-enum.yaml")
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile("../data/enums/response-enum-2.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetChecks(utils.StringList{"response-mediatype-enum-value-removed"}), d, osm)
	for _, err := range errs {
		require.Equal(t, checker.ERR, err.Level)
	}
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "response-mediatype-enum-value-removed", errs[0].Id)
}

// BC: removing an existing response with unparseable status is not breaking
func TestBreaking_ResponseUnparseableStatusRemoved(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	delete(s2.Spec.Paths[installCommandPath].Get.Responses, "default")

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	for _, err := range errs {
		require.Equal(t, checker.ERR, err.Level)
	}
	require.Empty(t, errs)
}

// BC: removing an existing response with error status is not breaking
func TestBreaking_ResponseErrorStatusRemoved(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	delete(s2.Spec.Paths[securityScorePath].Get.Responses, "400")

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	for _, err := range errs {
		require.Equal(t, checker.ERR, err.Level)
	}
	require.Empty(t, errs)
}

// BC: removing an existing optional response header is breaking as warn
func TestBreaking_OptionalResponseHeaderRemoved(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths[installCommandPath].Get.Responses["default"].Value.Headers["X-RateLimit-Limit"].Value.Required = false
	delete(s2.Spec.Paths[installCommandPath].Get.Responses["default"].Value.Headers, "X-RateLimit-Limit")

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	for _, err := range errs {
		require.Equal(t, checker.WARN, err.Level)
	}
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "optional-response-header-removed", errs[0].Id)
}

// BC: deleting a media-type from response is breaking
func TestBreaking_ResponseDeleteMediaType(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile("../data/response-media-type-base.yaml")
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile("../data/response-media-type-revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "response-media-type-removed", errs[0].Id)
}

// BC: deleting a pattern from a schema is not breaking
func TestBreaking_DeletePatten(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile("../data/pattern-base.yaml")
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile("../data/pattern-revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: adding a pattern to a schema is breaking
func TestBreaking_AddPattern(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile("../data/pattern-revision.yaml")
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile("../data/pattern-base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "request-property-pattern-added", errs[0].Id)
}

// BC: adding a pattern to a schema is breaking
/*
func TestBreaking_AddRequestParameterPattern(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile("../data/pattern-parameter-revision.yaml")
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile("../data/pattern-parameter-base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "request-parameter-pattern-added", errs[0].Id)
}
*/

// BC: adding a pattern to a schema is breaking for recursive properties
func TestBreaking_AddPatternRecursive(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile("../data/pattern-revision-recursive.yaml")
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile("../data/pattern-base-recursive.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "request-property-pattern-added", errs[0].Id)
}

// BC: modifying a pattern in a schema is breaking
func TestBreaking_ModifyPattern(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile("../data/pattern-base.yaml")
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile("../data/pattern-modified-not-anystring.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "request-property-pattern-changed", errs[0].Id)
}

// BC: modifying a pattern in request parameter is breaking
func TestBreaking_ModifyParameterPattern(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile("../data/pattern-parameter-base.yaml")
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile("../data/pattern-parameter-modified-not-anystring.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "request-parameter-pattern-changed", errs[0].Id)
}

// BC: modifying a pattern to ".*"" in a schema is not breaking
func TestBreaking_ModifyPatternToAnyString(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile("../data/pattern-base.yaml")
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile("../data/pattern-modified.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: modifying the default value of an optional request parameter is breaking
func TestBreaking_ModifyRequiredOptionalParamDefaultValue(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Default = "X"
	s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Default = "Y"

	// By default, OpenAPI treats all request parameters as optional

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
}

// BC: modifying the default value of a required request parameter is not breaking
func TestBreaking_ModifyRequiredRequiredParamDefaultValue(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	paramBase := s1.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies")
	paramBase.Required = true
	paramBase.Schema.Value.Default = "X"

	paramRevision := s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies")
	paramRevision.Required = true
	paramRevision.Schema.Value.Default = "Y"

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}
