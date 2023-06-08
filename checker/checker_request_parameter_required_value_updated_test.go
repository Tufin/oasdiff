package checker_test

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// BC: changing an existing header param from optional to required is breaking
func TestBreaking_HeaderParamBecameRequired(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = false
	s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = true

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.RequestParameterRequiredValueUpdatedCheck), d, osm)
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

// CL: changing an existing header param from required to optional is not breaking
func TestBreaking_HeaderParamBecameOptional(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = true
	s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = false

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterRequiredValueUpdatedCheck), d, osm, checker.INFO)
	require.NotEmpty(t, errs)
	require.Equal(t, checker.BackwardCompatibilityErrors{
		{
			Id:        "request-parameter-became-optional",
			Text:      "the 'header' request parameter 'network-policies' became optional",
			Comment:   "",
			Level:     checker.INFO,
			Operation: "GET",
			Path:      "/api/{domain}/{project}/install-command",
			Source:    "../data/openapi-test1.yaml",
		}}, errs)
}
