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

func verifyNonBreakingChangeIsChangelogEntry(t *testing.T, d *diff.Diff, osm *diff.OperationsSourcesMap, changeId string) {
	t.Helper()

	// Check no breaking change is detected
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
	// Check changelog captures the change
	errs = checker.CheckBackwardCompatibilityUntilLevel(checker.GetDefaultChecks(), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.INFO, errs[0].Level)
	require.Equal(t, changeId, errs[0].Id)
}

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

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
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

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	verifyNonBreakingChangeIsChangelogEntry(t, d, osm, "request-body-became-optional")
}

// BC: deleting a tag is not breaking
func TestBreaking_DeletedTag(t *testing.T) {
	r := d(t, getConfig(), 1, 5)
	require.Len(t, r, 6)
	require.Equal(t, "response-body-type-changed", r[0].Id)
	require.Equal(t, "response-success-status-removed", r[1].Id)
	require.Equal(t, "api-path-removed-without-deprecation", r[2].Id)
	require.Equal(t, "api-path-removed-without-deprecation", r[3].Id)
	require.Equal(t, "optional-response-header-removed", r[4].Id)
	require.Equal(t, "request-parameter-removed", r[5].Id)
}

// BC: adding an enum value is not breaking
func TestBreaking_AddedEnum(t *testing.T) {
	r := d(t, getConfig(), 1, 3)
	require.Len(t, r, 6)
	require.Equal(t, "response-success-status-removed", r[0].Id)
	require.Equal(t, "response-success-status-removed", r[1].Id)
	require.Equal(t, "request-parameter-removed", r[2].Id)
	require.Equal(t, "request-parameter-removed", r[3].Id)
	require.Equal(t, "request-parameter-removed", r[4].Id)
	require.Equal(t, "request-parameter-removed", r[5].Id)
}

// BC: changing extensions is not breaking
func TestBreaking_ModifiedExtension(t *testing.T) {
	r := d(t, getConfig(), 1, 3)
	require.Len(t, r, 6)
	require.Equal(t, "response-success-status-removed", r[0].Id)
	require.Equal(t, "response-success-status-removed", r[1].Id)
	require.Equal(t, "request-parameter-removed", r[2].Id)
	require.Equal(t, "request-parameter-removed", r[3].Id)
	require.Equal(t, "request-parameter-removed", r[4].Id)
	require.Equal(t, "request-parameter-removed", r[5].Id)
}

// BC: changing comments is not breaking
func TestBreaking_Comments(t *testing.T) {
	r := d(t, getConfig(), 1, 3)
	require.Len(t, r, 6)
	require.Equal(t, "response-success-status-removed", r[0].Id)
	require.Equal(t, "response-success-status-removed", r[1].Id)
	require.Equal(t, "request-parameter-removed", r[2].Id)
	require.Equal(t, "request-parameter-removed", r[3].Id)
	require.Equal(t, "request-parameter-removed", r[4].Id)
	require.Equal(t, "request-parameter-removed", r[5].Id)
}

// BC: new optional header param is not breaking
func TestBreaking_NewOptionalHeaderParam(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	deleteParam(s1.Spec.Paths[installCommandPath].Get, openapi3.ParameterInHeader, "network-policies")
	s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = false

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// CL: changing an existing header param to optional
func TestBreaking_HeaderParamRequiredDisabled(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = true
	s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = false

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	changes := checker.CheckBackwardCompatibilityUntilLevel(checker.GetDefaultChecks(), d, osm, checker.INFO)
	require.NotEmpty(t, changes)
	require.Equal(t, "request-parameter-became-optional", changes[0].Id)
	require.Len(t, changes, 1)
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

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: changing operation ID is not breaking
func TestBreaking_OperationID(t *testing.T) {
	r := d(t, getConfig(), 3, 1)
	require.Len(t, r, 3)
	require.Equal(t, "request-parameter-max-length-decreased", r[0].Id)
	require.Equal(t, "request-parameter-enum-value-removed", r[1].Id)
	require.Equal(t, "request-parameter-pattern-added", r[2].Id)
}

// BC: changing a link to operation ID is not breaking
func TestBreaking_LinkOperationID(t *testing.T) {
	r := d(t, getConfig(), 3, 1)
	require.Len(t, r, 3)
	require.Equal(t, "request-parameter-max-length-decreased", r[0].Id)
	require.Equal(t, "request-parameter-enum-value-removed", r[1].Id)
	require.Equal(t, "request-parameter-pattern-added", r[2].Id)
}

// BC: adding a media-type to response is not breaking
func TestBreaking_ResponseAddMediaType(t *testing.T) {
	s1, err := open("../data/response-media-type-revision.yaml")
	require.NoError(t, err)

	s2, err := open("../data/response-media-type-base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// CL: deprecating an operation with sunset greater than min
func TestBreaking_DeprecatedOperation(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Spec.Paths[installCommandPath].Get.Deprecated = true
	s2.Spec.Paths[installCommandPath].Get.Extensions[diff.SunsetExtension] = toJson(t, civil.DateOf(time.Now()).AddDays(180).String())

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(checker.GetDefaultChecks(), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, errs[0].Level, checker.INFO)
}

// BC: deprecating a parameter is not breaking
func TestBreaking_DeprecatedParameter(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Deprecated = true

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: deprecating a header is not breaking
func TestBreaking_DeprecatedHeader(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Spec.Paths[installCommandPath].Get.Responses["default"].Value.Headers["X-RateLimit-Limit"].Value.Deprecated = true

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: deprecating a schema is not breaking
func TestBreaking_DeprecatedSchema(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Deprecated = true

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: changing servers is not breaking
func TestBreaking_Servers(t *testing.T) {
	s1, err := open("../data/servers/baseswagger.json")
	require.NoError(t, err)

	s2, err := open("../data/servers/revisionswagger.json")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: adding a tag is not breaking
func TestBreaking_TagAdded(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Spec.Paths[securityScorePath].Get.Tags = append(s2.Spec.Paths[securityScorePath].Get.Tags, "newTag")
	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	for _, err := range errs {
		require.Equal(t, checker.ERR, err.Level)
	}
	require.Empty(t, errs)
}

// BC: adding a tag is not breaking with "api-tag-removed" check
func TestBreaking_TagAddedWithCustomCheck(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Spec.Paths[securityScorePath].Get.Tags = append(s2.Spec.Paths[securityScorePath].Get.Tags, "newTag")
	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	verifyNonBreakingChangeIsChangelogEntry(t, d, osm, "api-tag-added")
}

// BC: adding an operation ID is not breaking with "api-operation-id-removed" check
func TestBreaking_OperationIdAdded(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths[securityScorePath].Get.OperationID = ""

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	verifyNonBreakingChangeIsChangelogEntry(t, d, osm, "api-operation-id-added")
}
