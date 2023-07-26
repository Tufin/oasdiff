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

func l(t *testing.T, v int) load.SpecInfo {
	t.Helper()
	loader := openapi3.NewLoader()
	oas, err := loader.LoadFromFile(fmt.Sprintf("../data/openapi-test%d.yaml", v))
	require.NoError(t, err)
	return load.SpecInfo{Spec: oas, Url: fmt.Sprintf("../data/openapi-test%d.yaml", v)}
}

func d(t *testing.T, config *diff.Config, v1, v2 int) checker.Changes {
	t.Helper()
	l1 := l(t, v1)
	l2 := l(t, v2)
	d, osm, err := diff.GetWithOperationsSourcesMap(config, &l1, &l2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	return errs
}

func getConfig() *diff.Config {
	return diff.NewConfig().WithCheckBreaking()
}

// BC: deleting a path is breaking
func TestBreaking_DeletedPath(t *testing.T) {
	errs := d(t, getConfig(), 1, 701)
	require.Len(t, errs, 1)
	require.Equal(t, "api-path-removed-without-deprecation", errs[0].GetId())
}

// BC: deleting an operation is breaking
func TestBreaking_DeletedOp(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths[installCommandPath].Put = openapi3.NewOperation()

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "api-removed-without-deprecation", errs[0].GetId())
}

// BC: adding a required request body is breaking
func TestBreaking_AddingRequiredRequestBody(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Spec.Paths[installCommandPath].Get.RequestBody = &openapi3.RequestBodyRef{
		Value: openapi3.NewRequestBody().WithRequired(true),
	}

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "added-required-request-body", errs[0].GetId())
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

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "request-body-became-required", errs[0].GetId())
}

// BC: deleting an enum value is breaking
func TestBreaking_DeletedEnum(t *testing.T) {
	errs := d(t, getConfig(), 702, 1)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "request-parameter-enum-value-removed", errs[0].GetId())
}

// BC: added an enum value to response breaking
func TestBreaking_AddedResponseEnum(t *testing.T) {
	errs := d(t, getConfig(), 703, 704)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 2)
	require.Equal(t, "response-property-enum-value-added", errs[0].GetId())
	require.Equal(t, "response-property-enum-value-added", errs[1].GetId())
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

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(),
		&load.SpecInfo{Spec: s1},
		&load.SpecInfo{Spec: s2},
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

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)

	require.Len(t, errs, 1)
	require.Equal(t, "new-request-path-parameter", errs[0].GetId())
}

// BC: new required header param is breaking
func TestBreaking_NewRequiredHeaderParam(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	deleteParam(s1.Spec.Paths[installCommandPath].Get, openapi3.ParameterInHeader, "network-policies")
	s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = true

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "new-required-request-parameter", errs[0].GetId())
}

// BC: changing an existing header param from optional to required is breaking
func TestBreaking_HeaderParamRequiredEnabled(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = false
	s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = true

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Len(t, errs, 1)
	require.Equal(t,
		checker.ApiChange{
			Id:        "request-parameter-became-required",
			Text:      "the 'header' request parameter 'network-policies' became required",
			Comment:   "",
			Level:     checker.ERR,
			Operation: "GET",
			Path:      "/api/{domain}/{project}/install-command",
			Source:    "../data/openapi-test1.yaml",
		}, errs[0])
}

// BC: changing an existing response header from required to optional is breaking
func TestBreaking_ResponseHeaderParamRequiredDisabled(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths[installCommandPath].Get.Responses["default"].Value.Headers["X-RateLimit-Limit"].Value.Required = true
	s2.Spec.Paths[installCommandPath].Get.Responses["default"].Value.Headers["X-RateLimit-Limit"].Value.Required = false

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "response-header-became-optional", errs[0].GetId())
}

// BC: removing an existing required response header is breaking as error
func TestBreaking_ResponseHeaderRemoved(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths[installCommandPath].Get.Responses["default"].Value.Headers["X-RateLimit-Limit"].Value.Required = true
	delete(s2.Spec.Paths[installCommandPath].Get.Responses["default"].Value.Headers, "X-RateLimit-Limit")

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	for _, err := range errs {
		require.Equal(t, checker.ERR, err.GetLevel())
	}
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "required-response-header-removed", errs[0].GetId())
}

// BC: removing an existing response with successful status is breaking
func TestBreaking_ResponseSuccessStatusUpdated(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	delete(s2.Spec.Paths[securityScorePath].Get.Responses, "200")

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	for _, err := range errs {
		require.Equal(t, checker.ERR, err.GetLevel())
	}
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "response-success-status-removed", errs[0].GetId())
}

// BC: removing an existing response with non-successful status is breaking (optional)
func TestBreaking_ResponseNonSuccessStatusUpdated(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	delete(s2.Spec.Paths[securityScorePath].Get.Responses, "400")

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetChecks(utils.StringList{"response-non-success-status-removed"}), d, osm)
	for _, err := range errs {
		require.Equal(t, checker.ERR, err.GetLevel())
	}
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "response-non-success-status-removed", errs[0].GetId())
}

// BC: removing/updating an operation id is breaking (optional)
func TestBreaking_OperationIdRemoved(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Spec.Paths[securityScorePath].Get.OperationID = "newOperationId"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibility(checker.GetChecks(utils.StringList{"api-operation-id-removed"}), d, osm)
	for _, err := range errs {
		require.Equal(t, checker.ERR, err.GetLevel())
	}
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "api-operation-id-removed", errs[0].GetId())
	verifyNonBreakingChangeIsChangelogEntry(t, d, osm, "api-operation-id-removed")
}

// BC: removing/updating an enum in request body is breaking (optional)
func TestBreaking_RequestBodyEnumRemoved(t *testing.T) {
	s1, err := open("../data/enums/request-body-enum.yaml")
	require.NoError(t, err)

	s2, err := open("../data/enums/request-body-enum.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/api/v2/changeOfRequestFieldValueTiedToEnumTest"].Get.RequestBody.Value.Content["application/json"].Schema.Value.Enum = []interface{}{}

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibility(checker.GetChecks(utils.StringList{"request-body-enum-value-removed"}), d, osm)
	for _, err := range errs {
		require.Equal(t, checker.ERR, err.GetLevel())
	}

	require.Len(t, errs, 3)
	require.Equal(t, "request-body-enum-value-removed", errs[0].GetId())
}

// BC: removing/updating a property enum in response is breaking (optional)
func TestBreaking_ResponsePropertyEnumRemoved(t *testing.T) {
	s1 := l(t, 704)
	s2 := l(t, 703)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibility(checker.GetChecks(utils.StringList{"response-property-enum-value-removed"}), d, osm)
	for _, err := range errs {
		require.Equal(t, checker.ERR, err.GetLevel())
	}
	require.NotEmpty(t, errs)
	require.Len(t, errs, 2)
	require.Equal(t, "response-property-enum-value-removed", errs[0].GetId())
}

// BC: removing/updating a tag is breaking (optional)
func TestBreaking_TagRemoved(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Spec.Paths[securityScorePath].Get.Tags[0] = "newTag"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetChecks(utils.StringList{"api-tag-removed"}), d, osm)
	for _, err := range errs {
		require.Equal(t, checker.ERR, err.GetLevel())
	}
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "api-tag-removed", errs[0].GetId())
}

// BC: removing/updating a media type enum in response (optional)
func TestBreaking_ResponseMediaTypeEnumRemoved(t *testing.T) {
	s1, err := open("../data/enums/response-enum.yaml")
	require.NoError(t, err)

	s2, err := open("../data/enums/response-enum-2.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetChecks(utils.StringList{"response-mediatype-enum-value-removed"}), d, osm)
	for _, err := range errs {
		require.Equal(t, checker.ERR, err.GetLevel())
	}
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "response-mediatype-enum-value-removed", errs[0].GetId())
}

// BC: removing an existing response with unparseable status is not breaking
func TestBreaking_ResponseUnparseableStatusRemoved(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	delete(s2.Spec.Paths[installCommandPath].Get.Responses, "default")

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	for _, err := range errs {
		require.Equal(t, checker.ERR, err.GetLevel())
	}
	require.Empty(t, errs)
}

// BC: removing an existing response with error status is not breaking
func TestBreaking_ResponseErrorStatusRemoved(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	delete(s2.Spec.Paths[securityScorePath].Get.Responses, "400")

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	for _, err := range errs {
		require.Equal(t, checker.ERR, err.GetLevel())
	}
	require.Empty(t, errs)
}

// BC: removing an existing optional response header is breaking as warn
func TestBreaking_OptionalResponseHeaderRemoved(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths[installCommandPath].Get.Responses["default"].Value.Headers["X-RateLimit-Limit"].Value.Required = false
	delete(s2.Spec.Paths[installCommandPath].Get.Responses["default"].Value.Headers, "X-RateLimit-Limit")

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	for _, err := range errs {
		require.Equal(t, checker.WARN, err.GetLevel())
	}
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "optional-response-header-removed", errs[0].GetId())
}

// BC: deleting a media-type from response is breaking
func TestBreaking_ResponseDeleteMediaType(t *testing.T) {
	s1, err := open("../data/response-media-type-base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/response-media-type-revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "response-media-type-removed", errs[0].GetId())
}

// BC: deleting a pattern from a schema is not breaking
func TestBreaking_DeletePatten(t *testing.T) {
	s1, err := open("../data/pattern-base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/pattern-revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: adding a pattern to a schema is breaking
func TestBreaking_AddPattern(t *testing.T) {
	s1, err := open("../data/pattern-revision.yaml")
	require.NoError(t, err)

	s2, err := open("../data/pattern-base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "request-property-pattern-added", errs[0].GetId())
}

// BC: adding a pattern to a schema is breaking for recursive properties
func TestBreaking_AddPatternRecursive(t *testing.T) {
	s1, err := open("../data/pattern-revision-recursive.yaml")
	require.NoError(t, err)

	s2, err := open("../data/pattern-base-recursive.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "request-property-pattern-added", errs[0].GetId())
}

// BC: modifying a pattern in a schema is breaking
func TestBreaking_ModifyPattern(t *testing.T) {
	s1, err := open("../data/pattern-base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/pattern-modified-not-anystring.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "request-property-pattern-changed", errs[0].GetId())
}

// BC: modifying a pattern in request parameter is breaking
func TestBreaking_ModifyParameterPattern(t *testing.T) {
	s1, err := open("../data/pattern-parameter-base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/pattern-parameter-modified-not-anystring.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "request-parameter-pattern-changed", errs[0].GetId())
}

// BC: modifying a pattern to ".*" in a schema is not breaking
func TestBreaking_ModifyPatternToAnyString(t *testing.T) {
	s1, err := open("../data/pattern-base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/pattern-modified.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
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

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Len(t, errs, 1)
	require.Equal(t, "request-parameter-default-value-changed", errs[0].GetId())
	require.Equal(t, "for the 'header' request parameter 'network-policies', default value was changed from 'X' to 'Y'", errs[0].GetText())
}

// BC: setting the default value of an optional request parameter is breaking
func TestBreaking_SettingRequiredOptionalParamDefaultValue(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Default = nil
	s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Default = "Y"

	// By default, OpenAPI treats all request parameters as optional

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Len(t, errs, 1)
	require.Equal(t, "request-parameter-default-value-changed", errs[0].GetId())
	require.Equal(t, "for the 'header' request parameter 'network-policies', default value was changed from 'undefined' to 'Y'", errs[0].GetText())
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

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: removing an schema object from components is breaking (optional)
func TestBreaking_SchemaRemoved(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)
	s1.Spec.Paths = map[string]*openapi3.PathItem{}
	s2.Spec.Paths = map[string]*openapi3.PathItem{}

	for k := range s2.Spec.Components.Schemas {
		delete(s2.Spec.Components.Schemas, k)
	}

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	checks := checker.GetChecks(utils.StringList{"api-schema-removed"})
	errs := checker.CheckBackwardCompatibility(checks, d, osm)
	for _, err := range errs {
		require.Equal(t, checker.ERR, err.GetLevel())
	}
	require.NotEmpty(t, errs)
	require.Equal(t, "api-schema-removed", errs[0].GetId())
	require.Equal(t, "removed the schema 'network-policies'", errs[0].GetText())
	require.Equal(t, "api-schema-removed", errs[1].GetId())
	require.Equal(t, "removed the schema 'rules'", errs[1].GetText())
}

// BC: removing a media type from request body is breaking
func TestBreaking_RequestBodyMediaTypeRemoved(t *testing.T) {
	s1, err := open("../data/checker/request_body_media_type_updated_revision.yaml")
	require.NoError(t, err)

	s2, err := open("../data/checker/request_body_media_type_updated_base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Equal(t, "request-body-media-type-removed", errs[0].GetId())
	require.Equal(t, "removed the media type application/json from the request body", errs[0].GetText())
}

// BC: removing 'anyOf' schema from the request body or request body property
func TestBreaking_RequestPropertyAnyOfRemoved(t *testing.T) {
	s1, err := open("../data/checker/request_property_any_of_removed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_any_of_removed_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)

	require.Len(t, errs, 2)

	require.Equal(t, "request-body-any-of-removed", errs[0].GetId())
	require.Equal(t, checker.ERR, errs[0].GetLevel())
	require.Equal(t, "removed 'Rabbit' from the request body 'anyOf' list", errs[0].GetText())

	require.Equal(t, "request-property-any-of-removed", errs[1].GetId())
	require.Equal(t, checker.ERR, errs[1].GetLevel())
	require.Equal(t, "removed 'Breed3' from the '/anyOf[#/components/schemas/Dog]/breed' request property 'anyOf' list", errs[1].GetText())
}

// BC: removing 'oneOf' schema from the request body or request body property
func TestBreaking_RequestPropertyOneOfRemoved(t *testing.T) {
	s1, err := open("../data/checker/request_property_one_of_removed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_one_of_removed_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)

	require.Len(t, errs, 2)

	require.Equal(t, "request-body-one-of-removed", errs[0].GetId())
	require.Equal(t, checker.ERR, errs[0].GetLevel())
	require.Equal(t, "removed 'Rabbit' from the request body 'oneOf' list", errs[0].GetText())

	require.Equal(t, "request-property-one-of-removed", errs[1].GetId())
	require.Equal(t, checker.ERR, errs[1].GetLevel())
	require.Equal(t, "removed 'Breed3' from the '/oneOf[#/components/schemas/Dog]/breed' request property 'oneOf' list", errs[1].GetText())
}
