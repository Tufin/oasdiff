package checker_test

import (
	"fmt"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
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
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	return errs
}

// BC: deleting a path is breaking
func TestBreaking_DeletedPath(t *testing.T) {
	d := d(t, &diff.Config{}, 1, 2)
	require.NotEmpty(t, d)
}

// BC: deleting an operation is breaking
func TestBreaking_DeletedOp(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths[installCommandPath].Put = openapi3.NewOperation()

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
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
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
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
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
}

// BC: deleting an enum value is breaking
func TestBreaking_DeletedEnum(t *testing.T) {
	require.NotEmpty(t, d(t, &diff.Config{}, 3, 1))
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

// BC: new required path param is breaking
func TestBreaking_NewPathParam(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	deleteParam(s1.Spec.Paths[installCommandPath].Get, openapi3.ParameterInPath, "project")
	// note: path params are always required

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)

	// FIXME
	require.Contains(t, errs, checker.BackwardCompatibilityError{})
}

// BC: new required header param is breaking
func TestBreaking_NewRequiredHeaderParam(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	deleteParam(s1.Spec.Paths[installCommandPath].Get, openapi3.ParameterInHeader, "network-policies")
	s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = true

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	// FIXME
	require.Contains(t, errs, checker.BackwardCompatibilityError{})
}

// BC: changing an existing header param from optional to required is breaking
func TestBreaking_HeaderParamRequiredEnabled(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = false
	s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = true

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	// FIXME
	require.Equal(t, []checker.BackwardCompatibilityError{
		{
			Id:        "request-parameter-became-required",
			Text:      "the header request parameter network-policies became required",
			Comment:   "",
			Level:     checker.ERR,
			Operation: "GET",
			Path:      "/api/{domain}/{project}/install-command",
			Source:    "../data/openapi-test1.yaml",
			ToDo:      "Add to exceptions-list.md",
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
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
}

// BC: deleting a media-type from response is breaking
func TestBreaking_ResponseDeleteMediaType(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile("../data/response-media-type-base.yaml")
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile("../data/response-media-type-revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
}

// BC: deleting a pattern from a schema is not breaking
func TestBreaking_DeletePatten(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile("../data/pattern-base.yaml")
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile("../data/pattern-revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: adding a pattern to a schema is breaking
func TestBreaking_AddPatten(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile("../data/pattern-revision.yaml")
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile("../data/pattern-base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
}

// BC: modifying a pattern in a schema is breaking
func TestBreaking_ModifyPatten(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile("../data/pattern-base.yaml")
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile("../data/pattern-modified-not-anystring.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
}

// BC: modifying a pattern to ".*"" in a schema is not breaking
func TestBreaking_ModifyPattenToAnyString(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile("../data/pattern-base.yaml")
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile("../data/pattern-modified.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.DefaultChecks(), d, osm)
	require.Empty(t, errs)
}
