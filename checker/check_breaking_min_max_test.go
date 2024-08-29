package checker_test

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// BC: reducing max length in request is breaking: request-parameter-max-length-decreased
func TestBreaking_RequestMaxLengthSmaller(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	maxLengthFrom := uint64(13)
	s1.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MaxLength = &maxLengthFrom

	maxLengthTo := uint64(11)
	s2.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MaxLength = &maxLengthTo

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.RequestParameterMaxLengthDecreasedId, errs[0].GetId())
	require.Equal(t, "for 'path' request parameter 'domain', maxLength was decreased from '13' to '11'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: reducing max length in response is not breaking
func TestBreaking_ResponseMaxLengthSmaller(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	maxLengthFrom := uint64(13)
	s1.Spec.Paths.Value(securityScorePath).Get.Responses.Value("201").Value.Content["application/xml"].Schema.Value.MaxLength = &maxLengthFrom

	maxLengthTo := uint64(11)
	s2.Spec.Paths.Value(securityScorePath).Get.Responses.Value("201").Value.Content["application/xml"].Schema.Value.MaxLength = &maxLengthTo

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.Empty(t, errs)
}

// BC: reducing min length in request is not breaking
func TestBreaking_RequestMinLengthSmaller(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MinLength = uint64(13)
	s2.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MinLength = uint64(11)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.Empty(t, errs)
}

// BC: reducing min length in response is breaking: response-body-min-length-decreased
func TestBreaking_MinLengthSmaller(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths.Value(securityScorePath).Get.Responses.Value("201").Value.Content["application/xml"].Schema.Value.MinLength = uint64(13)
	s2.Spec.Paths.Value(securityScorePath).Get.Responses.Value("201").Value.Content["application/xml"].Schema.Value.MinLength = uint64(11)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.Equal(t, checker.ResponseBodyMinLengthDecreasedId, errs[0].GetId())
	require.Equal(t, "minLength value of response body was decreased from '13' to '11'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: increasing max length in request is not breaking
func TestBreaking_RequestMaxLengthGreater(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	maxLengthFrom := uint64(13)
	s1.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MaxLength = &maxLengthFrom

	maxLengthTo := uint64(14)
	s2.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MaxLength = &maxLengthTo

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.Empty(t, errs)
}

// BC: increasing max length in response is breaking: response-body-max-length-increased
func TestBreaking_ResponseMaxLengthGreater(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	maxLengthFrom := uint64(13)
	s1.Spec.Paths.Value(securityScorePath).Get.Responses.Value("201").Value.Content["application/xml"].Schema.Value.MaxLength = &maxLengthFrom

	maxLengthTo := uint64(14)
	s2.Spec.Paths.Value(securityScorePath).Get.Responses.Value("201").Value.Content["application/xml"].Schema.Value.MaxLength = &maxLengthTo

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ResponseBodyMaxLengthIncreasedId, errs[0].GetId())
	require.Equal(t, "maxLength value of response body was increased from '13' to '14'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: changing max length in request from nil to any value is breaking: request-parameter-max-length-set
func TestBreaking_MaxLengthFromNil(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MaxLength = nil

	maxLengthTo := uint64(14)
	s2.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MaxLength = &maxLengthTo

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.RequestParameterMaxLengthSetId, errs[0].GetId())
	require.Equal(t, "for 'path' request parameter 'domain', maxLength was set to '14'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: changing max length in response from nil to any value is not breaking
func TestBreaking_ResponseMaxLengthFromNil(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths.Value(securityScorePath).Get.Responses.Value("201").Value.Content["application/xml"].Schema.Value.MaxLength = nil

	maxLengthTo := uint64(14)
	s2.Spec.Paths.Value(securityScorePath).Get.Responses.Value("201").Value.Content["application/xml"].Schema.Value.MaxLength = &maxLengthTo

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.Empty(t, errs)
}

// BC: changing max length in request from any value to nil is not breaking
func TestBreaking_RequestMaxLengthToNil(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	maxLengthFrom := uint64(13)
	s1.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MaxLength = &maxLengthFrom

	s2.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MaxLength = nil

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.Empty(t, errs)
}

// BC: changing max length in response from any value to nil is breaking: response-body-max-length-unset
func TestBreaking_ResponseMaxLengthToNil(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	maxLengthFrom := uint64(13)
	s1.Spec.Paths.Value(securityScorePath).Get.Responses.Value("201").Value.Content["application/xml"].Schema.Value.MaxLength = &maxLengthFrom

	s2.Spec.Paths.Value(securityScorePath).Get.Responses.Value("201").Value.Content["application/xml"].Schema.Value.MaxLength = nil

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ResponseBodyMaxLengthUnsetId, errs[0].GetId())
	require.Equal(t, "maxLength value of response body was unset from '13'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: both max lengths in request are nil is not breaking
func TestBreaking_MaxLengthBothNil(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MaxLength = nil
	s2.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MaxLength = nil

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.Empty(t, errs)
}

// BC: both max lengths in response are nil is not breaking
func TestBreaking_ResponseMaxLengthBothNil(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths.Value(securityScorePath).Get.Responses.Value("201").Value.Content["application/xml"].Schema.Value.MaxLength = nil
	s2.Spec.Paths.Value(securityScorePath).Get.Responses.Value("201").Value.Content["application/xml"].Schema.Value.MaxLength = nil

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.Empty(t, errs)
}

// BC: reducing min items in request is not breaking
func TestBreaking_RequestMinItemsSmaller(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MinItems = 13
	s2.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MinItems = 11

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.Empty(t, errs)
}

// BC: reducing min items in response is breaking: response-body-min-items-decreased
func TestBreaking_ResponseMinItemsSmaller(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths.Value(securityScorePath).Get.Responses.Value("201").Value.Content["application/xml"].Schema.Value.MinItems = 13
	s2.Spec.Paths.Value(securityScorePath).Get.Responses.Value("201").Value.Content["application/xml"].Schema.Value.MinItems = 11

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ResponseBodyMinItemsDecreasedId, errs[0].GetId())
	require.Equal(t, "minItems value of response body was decreased from '13' to '11'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: increasing min items in request is breaking: request-parameter-min-items-increased
func TestBreaking_RequeatMinItemsGreater(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MinItems = 13
	s2.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MinItems = 14

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.Len(t, errs, 1)
	require.Equal(t, checker.RequestParameterMinItemsIncreasedId, errs[0].GetId())
	require.Equal(t, "for 'path' request parameter 'domain', minItems was increased from '13' to '14'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: increasing min items in response is not breaking
func TestBreaking_ResponseMinItemsGreater(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Spec.Paths.Value(securityScorePath).Get.Responses.Value("201").Value.Content["application/xml"].Schema.Value.MinItems = 13
	s2.Spec.Paths.Value(securityScorePath).Get.Responses.Value("201").Value.Content["application/xml"].Schema.Value.MinItems = 14

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.Empty(t, errs)
}

// BC: reducing max in request is breaking: request-parameter-max-decreased
func TestBreaking_MaxSmaller(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	maxFrom := float64(13)
	s1.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.Max = &maxFrom

	maxTo := float64(11)
	s2.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.Max = &maxTo

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.Len(t, errs, 1)
	require.Equal(t, checker.RequestParameterMaxDecreasedId, errs[0].GetId())
	require.Equal(t, "for 'path' request parameter 'domain', max was decreased from '13.00' to '11.00'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: reducing max in response is not breaking
func TestBreaking_MaxSmallerInResponse(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	maxFrom := float64(13)
	s1.Spec.Paths.Value(installCommandPath).Get.Responses.Value("default").Value.Headers["X-RateLimit-Limit"].Value.Schema.Value.Max = &maxFrom

	maxTo := float64(11)
	s2.Spec.Paths.Value(installCommandPath).Get.Responses.Value("default").Value.Headers["X-RateLimit-Limit"].Value.Schema.Value.Max = &maxTo

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.Empty(t, errs)
}
