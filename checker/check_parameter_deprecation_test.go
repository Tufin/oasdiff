package checker_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

func getParameterDeprecationFile(file string) string {
	return fmt.Sprintf("../data/param-deprecation/%s", file)
}

// // BC: deleting an operation after sunset date is not breaking
// func TestBreaking_DeprecationPast(t *testing.T) {

// 	s1, err := open(getDeprecationFile("deprecated-past.yaml"))
// 	require.NoError(t, err)

// 	s2, err := open(getDeprecationFile("sunset.yaml"))
// 	require.NoError(t, err)

// 	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
// 	require.NoError(t, err)
// 	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.ParameterDeprecationCheck), d, osm)
// 	require.Empty(t, errs)
// }

// // BC: deprecating an operation with a deprecation policy and an invalid sunset date is breaking
// func TestBreaking_DeprecationWithInvalidSunset(t *testing.T) {

// 	s1, err := open(getDeprecationFile("base.yaml"))
// 	require.NoError(t, err)

// 	s2, err := open(getDeprecationFile("deprecated-with-invalid-sunset.yaml"))
// 	require.NoError(t, err)

// 	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
// 	require.NoError(t, err)
// 	c := singleCheckConfig(checker.ParameterDeprecationCheck).WithDeprecation(0, 10)
// 	errs := checker.CheckBackwardCompatibility(c, d, osm)
// 	require.NotEmpty(t, errs)
// 	require.Len(t, errs, 1)
// 	require.Equal(t, checker.APIDeprecatedSunsetParseId, errs[0].GetId())
// 	require.Equal(t, "failed to parse sunset date: 'sunset date doesn't conform with RFC3339: invalid'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
// }

// // BC: deprecating an operation with a deprecation policy and an invalid stability level is breaking
// func TestBreaking_DeprecationWithInvalidStabilityLevel(t *testing.T) {

// 	s1, err := open(getDeprecationFile("base.yaml"))
// 	require.NoError(t, err)

// 	s2, err := open(getDeprecationFile("deprecated-with-invalid-stability.yaml"))
// 	require.NoError(t, err)

// 	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
// 	require.NoError(t, err)
// 	c := singleCheckConfig(checker.ParameterDeprecationCheck).WithDeprecation(0, 10)
// 	errs := checker.CheckBackwardCompatibility(c, d, osm)
// 	require.NotEmpty(t, errs)
// 	require.Len(t, errs, 1)
// 	require.Equal(t, checker.APIInvalidStabilityLevelId, errs[0].GetId())
// 	require.Equal(t, "failed to parse stability level: 'value is not one of draft, alpha, beta or stable: \"invalid\"'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
// 	require.Equal(t, "../data/deprecation/deprecated-with-invalid-stability.yaml", errs[0].GetSource())
// }

// // BC: deprecating an operation without a deprecation policy but without specifying sunset date is not breaking
// func TestBreaking_DeprecationWithoutSunsetNoPolicy(t *testing.T) {

// 	s1, err := open(getDeprecationFile("base.yaml"))
// 	require.NoError(t, err)

// 	s2, err := open(getDeprecationFile("deprecated-no-sunset.yaml"))
// 	require.NoError(t, err)

// 	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
// 	require.NoError(t, err)
// 	c := singleCheckConfig(checker.ParameterDeprecationCheck).WithDeprecation(0, 0)
// 	errs := checker.CheckBackwardCompatibility(c, d, osm)
// 	require.Empty(t, errs)
// }

// // BC: deprecating an operation with a deprecation policy but without specifying sunset date is breaking
// func TestBreaking_DeprecationWithoutSunsetWithPolicy(t *testing.T) {

// 	s1, err := open(getDeprecationFile("base.yaml"))
// 	require.NoError(t, err)

// 	s2, err := open(getDeprecationFile("deprecated-no-sunset.yaml"))
// 	require.NoError(t, err)

// 	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
// 	require.NoError(t, err)
// 	c := singleCheckConfig(checker.ParameterDeprecationCheck).WithDeprecation(30, 100)
// 	errs := checker.CheckBackwardCompatibility(c, d, osm)
// 	require.Len(t, errs, 1)
// 	require.Equal(t, checker.APIDeprecatedSunsetMissingId, errs[0].GetId())
// 	require.Equal(t, "sunset date is missing for deprecated API", errs[0].GetText(checker.NewDefaultLocalizer()))
// }

// // BC: deprecating an operation with a default deprecation policy but without specifying sunset date is not breaking
// func TestBreaking_DeprecationWithoutSunset(t *testing.T) {

// 	s1, err := open(getDeprecationFile("base.yaml"))
// 	require.NoError(t, err)

// 	s2, err := open(getDeprecationFile("deprecated-no-sunset.yaml"))
// 	require.NoError(t, err)

// 	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
// 	require.NoError(t, err)
// 	c := singleCheckConfig(checker.ParameterDeprecationCheck)
// 	errs := checker.CheckBackwardCompatibility(c, d, osm)
// 	require.Empty(t, errs)
// }

// // BC: deprecating an operation without a deprecation policy and without specifying sunset date is not breaking for alpha level
// func TestBreaking_DeprecationForAlpha(t *testing.T) {

// 	s1, err := open(getDeprecationFile("base-alpha-stability.yaml"))
// 	require.NoError(t, err)

// 	s2, err := open(getDeprecationFile("deprecated-no-sunset-alpha-stability.yaml"))
// 	require.NoError(t, err)

// 	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
// 	require.NoError(t, err)
// 	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.ParameterDeprecationCheck), d, osm)
// 	require.Empty(t, errs)
// }

// // BC: removing the path without a deprecation policy and without specifying sunset date is not breaking for alpha level
// func TestBreaking_RemovedPathForAlpha(t *testing.T) {
// 	s1, err := open(getDeprecationFile("base-alpha-stability.yaml"))
// 	require.NoError(t, err)
// 	alpha := toJson(t, checker.STABILITY_ALPHA)
// 	s1.Spec.Paths.Value("/api/test").Get.Extensions["x-stability-level"] = alpha
// 	s1.Spec.Paths.Value("/api/test").Post.Extensions = map[string]interface{}{"x-stability-level": alpha}

// 	s2, err := open(getDeprecationFile("base-alpha-stability.yaml"))
// 	require.NoError(t, err)

// 	s2.Spec.Paths.Delete("/api/test")

// 	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
// 	require.NoError(t, err)
// 	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.ParameterDeprecationCheck), d, osm)
// 	require.Empty(t, errs)
// }

// // BC: deprecating an operation without a deprecation policy and without specifying sunset date is not breaking for draft level
// func TestBreaking_DeprecationForDraft(t *testing.T) {
// 	s1, err := open(getDeprecationFile("base-alpha-stability.yaml"))
// 	require.NoError(t, err)
// 	draft := toJson(t, checker.STABILITY_DRAFT)
// 	s1.Spec.Paths.Value("/api/test").Get.Extensions["x-stability-level"] = draft

// 	s2, err := open(getDeprecationFile("deprecated-no-sunset-alpha-stability.yaml"))
// 	require.NoError(t, err)
// 	s2.Spec.Paths.Value("/api/test").Get.Extensions["x-stability-level"] = draft

// 	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
// 	require.NoError(t, err)
// 	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.ParameterDeprecationCheck), d, osm)
// 	require.Empty(t, errs)
// }

// // BC: removing the path without a deprecation policy and without specifying sunset date is not breaking for draft level
// func TestBreaking_RemovedPathForDraft(t *testing.T) {
// 	s1, err := open(getDeprecationFile("base-alpha-stability.yaml"))
// 	require.NoError(t, err)
// 	draft := toJson(t, checker.STABILITY_DRAFT)
// 	s1.Spec.Paths.Value("/api/test").Get.Extensions["x-stability-level"] = draft
// 	s1.Spec.Paths.Value("/api/test").Post.Extensions = map[string]interface{}{"x-stability-level": draft}

// 	s2, err := open(getDeprecationFile("base-alpha-stability.yaml"))
// 	require.NoError(t, err)

// 	s2.Spec.Paths.Delete("/api/test")

// 	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
// 	require.NoError(t, err)
// 	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.ParameterDeprecationCheck), d, osm)
// 	require.Empty(t, errs)
// }

// func toJson(t *testing.T, value string) json.RawMessage {
// 	t.Helper()
// 	data, err := json.Marshal(value)
// 	require.NoError(t, err)
// 	return data
// }

// // BC: deprecating an operation with a deprecation policy and sunset date before required deprecation period is breaking
// func TestBreaking_DeprecationWithEarlySunset(t *testing.T) {
// 	s1, err := open(getDeprecationFile("base.yaml"))
// 	require.NoError(t, err)

// 	s2, err := open(getDeprecationFile("deprecated-future.yaml"))
// 	require.NoError(t, err)
// 	sunsetDate := civil.DateOf(time.Now()).AddDays(9).String()
// 	s2.Spec.Paths.Value("/api/test").Get.Extensions[diff.SunsetExtension] = toJson(t, sunsetDate)

// 	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
// 	require.NoError(t, err)
// 	c := singleCheckConfig(checker.ParameterDeprecationCheck).WithDeprecation(0, 10)
// 	errs := checker.CheckBackwardCompatibility(c, d, osm)
// 	require.NotEmpty(t, errs)
// 	require.Len(t, errs, 1)
// 	require.Equal(t, checker.APISunsetDateTooSmallId, errs[0].GetId())
// 	require.Equal(t, fmt.Sprintf("sunset date '%s' is too small, must be at least '10' days from now", sunsetDate), errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
// }

// // BC: deprecating an operation with a deprecation policy and sunset date after required deprecation period is not breaking
// func TestBreaking_DeprecationWithProperSunset(t *testing.T) {

// 	s1, err := open(getDeprecationFile("base.yaml"))
// 	require.NoError(t, err)

// 	s2, err := open(getDeprecationFile("deprecated-future.yaml"))
// 	require.NoError(t, err)

// 	s2.Spec.Paths.Value("/api/test").Get.Extensions[diff.SunsetExtension] = toJson(t, civil.DateOf(time.Now()).AddDays(10).String())

// 	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
// 	c := singleCheckConfig(checker.ParameterDeprecationCheck).WithDeprecation(0, 10)
// 	require.NoError(t, err)
// 	errs := checker.CheckBackwardCompatibilityUntilLevel(c, d, osm, checker.INFO)
// 	require.Len(t, errs, 1)
// 	// only a non-breaking change detected
// 	require.Equal(t, checker.INFO, errs[0].GetLevel())
// 	require.Equal(t, "endpoint deprecated", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
// }

// // BC: deleting a path after sunset date of all contained operations is not breaking
// func TestBreaking_DeprecationPathPast(t *testing.T) {

// 	s1, err := open(getDeprecationFile("deprecated-path-past.yaml"))
// 	require.NoError(t, err)

// 	s2, err := open(getDeprecationFile("sunset-path.yaml"))
// 	require.NoError(t, err)

// 	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
// 	require.NoError(t, err)
// 	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.ParameterDeprecationCheck), d, osm)
// 	require.Empty(t, errs)
// }

// CL: parameters that became deprecated
func TestParameterDeprecated_DetectsDeprecated(t *testing.T) {
	s1, err := open(getParameterDeprecationFile("base.yaml"))
	require.NoError(t, err)

	s2, err := open(getParameterDeprecationFile("deprecated-future.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ParameterDeprecationCheck), d, osm, checker.INFO)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)

	require.IsType(t, checker.ApiChange{}, errs[0])
	e0 := errs[0].(checker.ApiChange)
	require.Equal(t, checker.ParameterDeprecatedId, e0.Id)
	require.Equal(t, "GET", e0.Operation)
	require.Equal(t, "/api/test", e0.Path)
	require.Equal(t, "'query' request parameter 'id' was deprecated", e0.GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// CL: parameters that were re-activated
func TestParameterDeprecated_DetectsReactivated(t *testing.T) {
	s1, err := open(getParameterDeprecationFile("deprecated-future.yaml"))
	require.NoError(t, err)

	s2, err := open(getParameterDeprecationFile("base.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ParameterDeprecationCheck), d, osm, checker.INFO)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)

	require.IsType(t, checker.ApiChange{}, errs[0])
	e0 := errs[0].(checker.ApiChange)
	require.Equal(t, checker.ParameterReactivatedId, e0.Id)
	require.Equal(t, "GET", e0.Operation)
	require.Equal(t, "/api/test", e0.Path)
	require.Equal(t, "'query' request parameter 'id' was reactivated", e0.GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// func TestBreaking_InvaidStability(t *testing.T) {

// 	s1, err := open(getDeprecationFile("invalid-stability.yaml"))
// 	require.NoError(t, err)

// 	s2, err := open(getDeprecationFile("base-alpha-stability.yaml"))
// 	require.NoError(t, err)

// 	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
// 	require.NoError(t, err)
// 	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.ParameterDeprecationCheck), d, osm)
// 	require.Len(t, errs, 1)

// 	require.IsType(t, checker.ApiChange{}, errs[0])
// 	e0 := errs[0].(checker.ApiChange)
// 	require.Equal(t, checker.APIInvalidStabilityLevelId, e0.Id)
// 	require.Equal(t, "GET", e0.Operation)
// 	require.Equal(t, "/api/test", e0.Path)
// 	require.Equal(t, "failed to parse stability level: 'value is not one of draft, alpha, beta or stable: \"ga\"'", e0.GetUncolorizedText(checker.NewDefaultLocalizer()))
// 	require.Equal(t, "../data/deprecation/invalid-stability.yaml", errs[0].GetSource())
// }
