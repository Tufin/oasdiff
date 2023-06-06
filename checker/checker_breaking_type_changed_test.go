package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// BC: changing response's body schema type from string to number is breaking
func TestBreaking_TypeStringToNumber(t *testing.T) {
	s1, err := open("../data/type-change/simple.yaml")
	require.NoError(t, err)
	s1.Spec.Paths["/test"].Get.Responses["200"].Value.Content["application/json"].Schema.Value.Type = "string"

	s2, err := open("../data/type-change/simple.yaml")
	require.NoError(t, err)
	s2.Spec.Paths["/test"].Get.Responses["200"].Value.Content["application/json"].Schema.Value.Type = "number"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Len(t, errs, 1)
	require.Equal(t, "response-body-type-changed", errs[0].Id)
	require.Equal(t, "the response's body type/format changed from 'string'/'none' to 'number'/'none' for status '200'", errs[0].Text)
}

// BC: changing response's body schema type from number to string is breaking
func TestBreaking_TypeNumberToString(t *testing.T) {
	s1, err := open("../data/type-change/simple.yaml")
	require.NoError(t, err)
	s1.Spec.Paths["/test"].Get.Responses["200"].Value.Content["application/json"].Schema.Value.Type = "number"

	s2, err := open("../data/type-change/simple.yaml")
	require.NoError(t, err)
	s2.Spec.Paths["/test"].Get.Responses["200"].Value.Content["application/json"].Schema.Value.Type = "string"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Len(t, errs, 1)
	require.Equal(t, "response-body-type-changed", errs[0].Id)
	require.Equal(t, "the response's body type/format changed from 'number'/'none' to 'string'/'none' for status '200'", errs[0].Text)
}

// BC: changing response's body schema type from number to integer is not breaking
func TestBreaking_TypeNumberToInteger(t *testing.T) {
	s1, err := open("../data/type-change/simple.yaml")
	require.NoError(t, err)
	s1.Spec.Paths["/test"].Get.Responses["200"].Value.Content["application/json"].Schema.Value.Type = "number"

	s2, err := open("../data/type-change/simple.yaml")
	require.NoError(t, err)
	s2.Spec.Paths["/test"].Get.Responses["200"].Value.Content["application/json"].Schema.Value.Type = "integer"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: changing response's body schema type from integer to number is breaking
func TestBreaking_TypeIntegerToNumber(t *testing.T) {
	s1, err := open("../data/type-change/simple.yaml")
	require.NoError(t, err)
	s1.Spec.Paths["/test"].Get.Responses["200"].Value.Content["application/json"].Schema.Value.Type = "integer"

	s2, err := open("../data/type-change/simple.yaml")
	require.NoError(t, err)
	s2.Spec.Paths["/test"].Get.Responses["200"].Value.Content["application/json"].Schema.Value.Type = "number"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Len(t, errs, 1)
	require.Equal(t, "response-body-type-changed", errs[0].Id)
	require.Equal(t, "the response's body type/format changed from 'integer'/'none' to 'number'/'none' for status '200'", errs[0].Text)
}

// BC: changing response's body schema type from number/none to integer/int32 is not breaking
func TestBreaking_TypeNumberToInt32(t *testing.T) {
	s1, err := open("../data/type-change/simple.yaml")
	require.NoError(t, err)
	s1.Spec.Paths["/test"].Get.Responses["200"].Value.Content["application/json"].Schema.Value.Type = "number"

	s2, err := open("../data/type-change/simple.yaml")
	require.NoError(t, err)
	s2.Spec.Paths["/test"].Get.Responses["200"].Value.Content["application/json"].Schema.Value.Type = "integer"
	s2.Spec.Paths["/test"].Get.Responses["200"].Value.Content["application/json"].Schema.Value.Format = "int32"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: changing response's embedded property schema type from string/none to integer/int32 is breaking
func TestBreaking_RespTypeChanged(t *testing.T) {
	s1, err := open("../data/type-change/base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/type-change/revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Len(t, errs, 1)
	require.Equal(t, "response-property-type-changed", errs[0].Id)
	require.Equal(t, "the response's property type/format changed from 'string'/'none' to 'integer'/'int32' for status '200'", errs[0].Text)
}
