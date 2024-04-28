package checker_test

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

// BC: changing an existing header param from optional to required is breaking
func TestBreaking_HeaderParamBecameRequired(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = false
	s2.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = true

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.RequestParameterRequiredValueUpdatedCheck), d, osm)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:        checker.RequestParameterBecomeRequiredId,
		Args:      []any{"header", "network-policies"},
		Level:     checker.ERR,
		Operation: "GET",
		Path:      "/api/{domain}/{project}/install-command",
		Source:    load.NewSource("../data/openapi-test1.yaml"),
	}, errs[0])
}

// CL: changing an existing header param from required to optional
func TestBreaking_HeaderParamBecameOptional(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = true
	s2.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = false

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterRequiredValueUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:        checker.RequestParameterBecomeOptionalId,
		Args:      []any{"header", "network-policies"},
		Level:     checker.INFO,
		Operation: "GET",
		Path:      "/api/{domain}/{project}/install-command",
		Source:    load.NewSource("../data/openapi-test1.yaml"),
	}, errs[0])
}
