package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// BC: changing request's body schema type from string to number is breaking
func TestBreaking_ReqTypeStringToNumber(t *testing.T) {
	file := "../data/type-change/simple-request.yaml"

	s1, err := open(file)
	require.NoError(t, err)
	s1.Spec.Paths["/test"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Type = "string"

	s2, err := open(file)
	require.NoError(t, err)
	s2.Spec.Paths["/test"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Type = "number"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Len(t, errs, 1)
	require.Equal(t, checker.RequestBodyTypeChangedId, errs[0].GetId())
	require.Equal(t, "the request's body type/format changed from 'string'/'' to 'number'/''", errs[0].GetText(checker.NewDefaultLocalizer()))
}

// BC: changing request's body schema type from number to string is breaking
func TestBreaking_ReqTypeNumberToString(t *testing.T) {
	file := "../data/type-change/simple-request.yaml"

	s1, err := open(file)
	require.NoError(t, err)
	s1.Spec.Paths["/test"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Type = "number"

	s2, err := open(file)
	require.NoError(t, err)
	s2.Spec.Paths["/test"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Type = "string"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Len(t, errs, 1)
	require.Equal(t, checker.RequestBodyTypeChangedId, errs[0].GetId())
	require.Equal(t, "the request's body type/format changed from 'number'/'' to 'string'/''", errs[0].GetText(checker.NewDefaultLocalizer()))
}

// BC: changing request's body schema type from number to integer is breaking
func TestBreaking_ReqTypeNumberToInteger(t *testing.T) {
	file := "../data/type-change/simple-request.yaml"

	s1, err := open(file)
	require.NoError(t, err)
	s1.Spec.Paths["/test"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Type = "number"

	s2, err := open(file)
	require.NoError(t, err)
	s2.Spec.Paths["/test"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Type = "integer"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Len(t, errs, 1)
	require.Equal(t, checker.RequestBodyTypeChangedId, errs[0].GetId())
	require.Equal(t, "the request's body type/format changed from 'number'/'' to 'integer'/''", errs[0].GetText(checker.NewDefaultLocalizer()))
}

// BC: changing request's body schema type from integer to number is not breaking
func TestBreaking_ReqTypeIntegerToNumber(t *testing.T) {
	file := "../data/type-change/simple-request.yaml"

	s1, err := open(file)
	require.NoError(t, err)
	s1.Spec.Paths["/test"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Type = "integer"

	s2, err := open(file)
	require.NoError(t, err)
	s2.Spec.Paths["/test"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Type = "number"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: changing request's body schema type from number/none to integer/int32 is breaking
func TestBreaking_ReqTypeNumberToInt32(t *testing.T) {
	file := "../data/type-change/simple-request.yaml"

	s1, err := open(file)
	require.NoError(t, err)
	s1.Spec.Paths["/test"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Type = "number"

	s2, err := open(file)
	require.NoError(t, err)
	s2.Spec.Paths["/test"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Type = "integer"
	s2.Spec.Paths["/test"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Format = "int32"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Len(t, errs, 1)
	require.Equal(t, checker.RequestBodyTypeChangedId, errs[0].GetId())
	require.Equal(t, "the request's body type/format changed from 'number'/'' to 'integer'/'int32'", errs[0].GetText(checker.NewDefaultLocalizer()))
}
