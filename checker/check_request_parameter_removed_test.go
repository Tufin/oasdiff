package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// BC: deleting a parameter without deprecation is breaking
func TestBreaking_DeletedParameter(t *testing.T) {
	s1, err := open(getParameterDeprecationFile("base.yaml"))
	require.NoError(t, err)

	s2, err := open(getParameterDeprecationFile("sunset.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.RequestParameterRemovedCheck), d, osm)
	require.Len(t, errs, 1)
	require.Equal(t, checker.RequestParameterRemovedId, errs[0].GetId())
	require.Equal(t, "deleted the 'query' request parameter 'id'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
	require.Equal(t, "This is a warning because some apps may return an error when receiving a parameter that they do not expect. It is recommended to deprecate the parameter first.", errs[0].GetComment(checker.NewDefaultLocalizer()))
}

// BC: deleting a parameter after sunset date is not breaking
func TestBreaking_ParameterDeprecationPast(t *testing.T) {

	s1, err := open(getParameterDeprecationFile("deprecated-past.yaml"))
	require.NoError(t, err)

	s2, err := open(getParameterDeprecationFile("sunset.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.RequestParameterRemovedCheck), d, osm)
	require.Empty(t, errs)
}

// BC: deleting a parameter before sunset date is breaking
func TestBreaking_ParameterDeprecationFuture(t *testing.T) {

	s1, err := open(getParameterDeprecationFile("deprecated-future.yaml"))
	require.NoError(t, err)

	s2, err := open(getParameterDeprecationFile("sunset.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.RequestParameterRemovedCheck), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ParameterRemovedBeforeSunsetId, errs[0].GetId())
	require.Equal(t, "deleted the 'query' request parameter 'id' before the sunset date '9999-08-10'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: removing a deprecated parameter with an invalid date is breaking
func TestBreaking_RemoveParameterWithInvalidSunset(t *testing.T) {

	s1, err := open(getParameterDeprecationFile("deprecated-invalid.yaml"))
	require.NoError(t, err)

	s2, err := open(getParameterDeprecationFile("sunset.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.RequestParameterRemovedCheck), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.RequestParameterSunsetParseId, errs[0].GetId())
	require.Equal(t, "failed to parse sunset date for the 'query' request parameter 'id': 'sunset date doesn't conform with RFC3339: invalid-date'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
	require.Equal(t, "../data/param-deprecation/sunset.yaml", errs[0].GetSource())
}

// // BC: removing the path without a deprecation policy and without specifying sunset date is not breaking for alpha level
// func TestBreaking_RemovedPathForAlpha(t *testing.T) {
// 	s1, err := open(getParameterDeprecationFile("base-alpha-stability.yaml"))
// 	require.NoError(t, err)
// 	alpha := toJson(t, checker.STABILITY_ALPHA)
// 	s1.Spec.Paths.Value("/api/test").Get.Extensions["x-stability-level"] = alpha
// 	s1.Spec.Paths.Value("/api/test").Post.Extensions = map[string]interface{}{"x-stability-level": alpha}

// 	s2, err := open(getParameterDeprecationFile("base-alpha-stability.yaml"))
// 	require.NoError(t, err)

// 	s2.Spec.Paths.Delete("/api/test")

// 	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
// 	require.NoError(t, err)
// 	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.ParameterDeprecationCheck), d, osm)
// 	require.Empty(t, errs)
// }

// // BC: removing the path without a deprecation policy and without specifying sunset date is not breaking for draft level
// func TestBreaking_RemovedPathForDraft(t *testing.T) {
// 	s1, err := open(getParameterDeprecationFile("base-alpha-stability.yaml"))
// 	require.NoError(t, err)
// 	draft := toJson(t, checker.STABILITY_DRAFT)
// 	s1.Spec.Paths.Value("/api/test").Get.Extensions["x-stability-level"] = draft
// 	s1.Spec.Paths.Value("/api/test").Post.Extensions = map[string]interface{}{"x-stability-level": draft}

// 	s2, err := open(getParameterDeprecationFile("base-alpha-stability.yaml"))
// 	require.NoError(t, err)

// 	s2.Spec.Paths.Delete("/api/test")

// 	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
// 	require.NoError(t, err)
// 	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.ParameterDeprecationCheck), d, osm)
// 	require.Empty(t, errs)
// }

// // BC: deleting a path after sunset date of all contained operations is not breaking
// func TestBreaking_DeprecationPathPast(t *testing.T) {

// 	s1, err := open(getParameterDeprecationFile("deprecated-path-past.yaml"))
// 	require.NoError(t, err)

// 	s2, err := open(getParameterDeprecationFile("sunset-path.yaml"))
// 	require.NoError(t, err)

// 	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
// 	require.NoError(t, err)
// 	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.ParameterDeprecationCheck), d, osm)
// 	require.Empty(t, errs)
// }
