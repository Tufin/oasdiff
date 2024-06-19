package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// BC: deleting an operation before sunset date is breaking
func TestBreaking_RemoveBeforeSunset(t *testing.T) {

	s1, err := open(getDeprecationFile("deprecated-future.yaml"))
	require.NoError(t, err)

	s2, err := open(getDeprecationFile("sunset.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.APIRemovedCheck), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.APIRemovedBeforeSunsetId, errs[0].GetId())
	require.Equal(t, "api removed before the sunset date '9999-08-10'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: deleting an operation without sunset date is not breaking
func TestBreaking_DeprecationNoSunset(t *testing.T) {

	s1, err := open(getDeprecationFile("deprecated-no-sunset.yaml"))
	require.NoError(t, err)

	s2, err := open(getDeprecationFile("sunset.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.APIRemovedCheck), d, osm)
	require.NoError(t, err)
	require.Empty(t, errs)
}

// BC: removing the path without a deprecation policy and without specifying sunset date is breaking if some APIs are not alpha stability level
func TestBreaking_RemovedPathForAlphaBreaking(t *testing.T) {
	s1, err := open(getDeprecationFile("base-alpha-stability.yaml"))
	require.NoError(t, err)

	s2, err := open(getDeprecationFile("base-alpha-stability.yaml"))
	require.NoError(t, err)

	s2.Spec.Paths.Delete("/api/test")

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.APIRemovedCheck), d, osm)
	require.Len(t, errs, 2)
	require.Equal(t, checker.APIPathRemovedWithoutDeprecationId, errs[0].GetId())
	require.Equal(t, "api path removed without deprecation", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
	require.Equal(t, checker.APIPathRemovedWithoutDeprecationId, errs[1].GetId())
	require.Equal(t, "api path removed without deprecation", errs[1].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: removing the path without a deprecation policy and without specifying sunset date is breaking if some APIs are not draft stability level
func TestBreaking_RemovedPathForDraftBreaking(t *testing.T) {
	s1, err := open(getDeprecationFile("base-alpha-stability.yaml"))
	require.NoError(t, err)
	draft := toJson(t, checker.STABILITY_DRAFT)
	s1.Spec.Paths.Value("/api/test").Get.Extensions["x-stability-level"] = draft

	s2, err := open(getDeprecationFile("base-alpha-stability.yaml"))
	require.NoError(t, err)

	s2.Spec.Paths.Delete("/api/test")

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.APIRemovedCheck), d, osm)
	require.Len(t, errs, 2)
	require.Equal(t, checker.APIPathRemovedWithoutDeprecationId, errs[0].GetId())
	require.Equal(t, "api path removed without deprecation", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
	require.Equal(t, checker.APIPathRemovedWithoutDeprecationId, errs[1].GetId())
	require.Equal(t, "api path removed without deprecation", errs[1].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: deleting a path with some operations having sunset date in the future is breaking
func TestBreaking_DeprecationPathMixed(t *testing.T) {

	s1, err := open(getDeprecationFile("deprecated-path-mixed.yaml"))
	require.NoError(t, err)

	s2, err := open(getDeprecationFile("sunset-path.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.APIRemovedCheck), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.APIPathRemovedBeforeSunsetId, errs[0].GetId())
	require.Equal(t, "api path removed before the sunset date '9999-08-10'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: removing a deprecated enpoint with an invalid date is breaking
func TestBreaking_RemoveEndpointWithInvalidSunset(t *testing.T) {

	s1, err := open(getDeprecationFile("deprecated-invalid.yaml"))
	require.NoError(t, err)

	s2, err := open(getDeprecationFile("deprecated-invalid.yaml"))
	require.NoError(t, err)

	s2.Spec.Paths.Find("/api/test").SetOperation("GET", nil)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.APIRemovedCheck), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.APIPathSunsetParseId, errs[0].GetId())
	require.Equal(t, "failed to parse sunset date: 'sunset date doesn't conform with RFC3339: invalid-date'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
	require.Equal(t, "../data/deprecation/deprecated-invalid.yaml", errs[0].GetSource())
}

// test sunset date without double quotes, see https://github.com/Tufin/oasdiff/pull/198/files
func TestBreaking_DeprecationPathMixed_RFC3339_Sunset(t *testing.T) {

	s1, err := open(getDeprecationFile("deprecated-path-mixed-rfc3339-sunset.yaml"))
	require.NoError(t, err)

	s2, err := open(getDeprecationFile("sunset-path.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.APIRemovedCheck), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.APIPathRemovedBeforeSunsetId, errs[0].GetId())
	require.Equal(t, "api path removed before the sunset date '9999-08-10'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}
