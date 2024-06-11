package checker_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

func open(file string) (*load.SpecInfo, error) {
	return load.NewSpecInfo(openapi3.NewLoader(), load.NewSource(file))
}

func getDeprecationFile(file string) string {
	return fmt.Sprintf("../data/deprecation/%s", file)
}

func singleCheckConfig(c checker.BackwardCompatibilityCheck) *checker.Config {
	return &checker.Config{
		Checks:              []checker.BackwardCompatibilityCheck{c},
		MinSunsetBetaDays:   31,
		MinSunsetStableDays: 180,
	}
}

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

// BC: deleting an operation after sunset date is not breaking
func TestBreaking_DeprecationPast(t *testing.T) {

	s1, err := open(getDeprecationFile("deprecated-past.yaml"))
	require.NoError(t, err)

	s2, err := open(getDeprecationFile("sunset.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.APIDeprecationCheck), d, osm)
	require.Empty(t, errs)
}

// BC: deprecating an operation with a deprecation policy but without specifying sunset date is breaking
func TestBreaking_DeprecationWithoutSunset(t *testing.T) {

	s1, err := open(getDeprecationFile("base.yaml"))
	require.NoError(t, err)

	s2, err := open(getDeprecationFile("deprecated-no-sunset.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	c := singleCheckConfig(checker.APIDeprecationCheck).WithDeprecation(0, 10)
	errs := checker.CheckBackwardCompatibility(c, d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.APIDeprecatedSunsetMissingId, errs[0].GetId())
	require.Equal(t, "sunset date is missing for deprecated API", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: deprecating an operation with a deprecation policy and an invalid sunset date is breaking
func TestBreaking_DeprecationWithInvalidSunset(t *testing.T) {

	s1, err := open(getDeprecationFile("base.yaml"))
	require.NoError(t, err)

	s2, err := open(getDeprecationFile("deprecated-with-invalid-sunset.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	c := singleCheckConfig(checker.APIDeprecationCheck).WithDeprecation(0, 10)
	errs := checker.CheckBackwardCompatibility(c, d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.APIDeprecatedSunsetParseId, errs[0].GetId())
	require.Equal(t, "failed to parse sunset date: 'sunset date doesn't conform with RFC3339: invalid'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: deprecating an operation with a deprecation policy and an invalid stability level is breaking
func TestBreaking_DeprecationWithInvalidStabilityLevel(t *testing.T) {

	s1, err := open(getDeprecationFile("base.yaml"))
	require.NoError(t, err)

	s2, err := open(getDeprecationFile("deprecated-with-invalid-stability.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	c := singleCheckConfig(checker.APIDeprecationCheck).WithDeprecation(0, 10)
	errs := checker.CheckBackwardCompatibility(c, d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.APIInvalidStabilityLevelId, errs[0].GetId())
	require.Equal(t, "failed to parse stability level: 'value is not one of draft, alpha, beta or stable: \"invalid\"'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: deprecating an operation without a deprecation policy but without specifying sunset date is not breaking
func TestBreaking_DeprecationWithoutSunsetNoPolicy(t *testing.T) {

	s1, err := open(getDeprecationFile("base.yaml"))
	require.NoError(t, err)

	s2, err := open(getDeprecationFile("deprecated-no-sunset.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	c := singleCheckConfig(checker.APIDeprecationCheck).WithDeprecation(0, 0)
	errs := checker.CheckBackwardCompatibility(c, d, osm)
	require.Empty(t, errs)
}

// BC: deprecating an operation with a deprecation policy but without specifying sunset date is breaking
func TestBreaking_DeprecationWithoutSunsetWithPolicy(t *testing.T) {

	s1, err := open(getDeprecationFile("base.yaml"))
	require.NoError(t, err)

	s2, err := open(getDeprecationFile("deprecated-no-sunset.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	c := singleCheckConfig(checker.APIDeprecationCheck).WithDeprecation(30, 100)
	errs := checker.CheckBackwardCompatibility(c, d, osm)
	require.Len(t, errs, 1)
	require.Equal(t, checker.APIDeprecatedSunsetMissingId, errs[0].GetId())
	require.Equal(t, "sunset date is missing for deprecated API", errs[0].GetText(checker.NewDefaultLocalizer()))
}

// BC: deprecating an operation without a deprecation policy and without specifying sunset date is not breaking for alpha level
func TestBreaking_DeprecationForAlpha(t *testing.T) {

	s1, err := open(getDeprecationFile("base-alpha-stability.yaml"))
	require.NoError(t, err)

	s2, err := open(getDeprecationFile("deprecated-no-sunset-alpha-stability.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.APIDeprecationCheck), d, osm)
	require.Empty(t, errs)
}

// BC: specifying an invalid stability level is breaking
func TestBreaking_InvalidStabilityLevel(t *testing.T) {
	s1, err := open(getDeprecationFile("base-alpha-stability.yaml"))
	require.NoError(t, err)

	s2, err := open(getDeprecationFile("base-invalid-stability.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.APIDeprecationCheck), d, osm)
	require.Len(t, errs, 1)
	require.Equal(t, checker.APIInvalidStabilityLevelId, errs[0].GetId())
}

// BC: removing the path without a deprecation policy and without specifying sunset date is not breaking for alpha level
func TestBreaking_RemovedPathForAlpha(t *testing.T) {
	s1, err := open(getDeprecationFile("base-alpha-stability.yaml"))
	require.NoError(t, err)
	alpha := toJson(t, checker.STABILITY_ALPHA)
	s1.Spec.Paths.Value("/api/test").Get.Extensions["x-stability-level"] = alpha
	s1.Spec.Paths.Value("/api/test").Post.Extensions = map[string]interface{}{"x-stability-level": alpha}

	s2, err := open(getDeprecationFile("base-alpha-stability.yaml"))
	require.NoError(t, err)

	s2.Spec.Paths.Delete("/api/test")

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.APIDeprecationCheck), d, osm)
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

// BC: deprecating an operation without a deprecation policy and without specifying sunset date is not breaking for draft level
func TestBreaking_DeprecationForDraft(t *testing.T) {
	s1, err := open(getDeprecationFile("base-alpha-stability.yaml"))
	require.NoError(t, err)
	draft := toJson(t, checker.STABILITY_DRAFT)
	s1.Spec.Paths.Value("/api/test").Get.Extensions["x-stability-level"] = draft

	s2, err := open(getDeprecationFile("deprecated-no-sunset-alpha-stability.yaml"))
	require.NoError(t, err)
	s2.Spec.Paths.Value("/api/test").Get.Extensions["x-stability-level"] = draft

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.APIDeprecationCheck), d, osm)
	require.Empty(t, errs)
}

// BC: removing the path without a deprecation policy and without specifying sunset date is not breaking for draft level
func TestBreaking_RemovedPathForDraft(t *testing.T) {
	s1, err := open(getDeprecationFile("base-alpha-stability.yaml"))
	require.NoError(t, err)
	draft := toJson(t, checker.STABILITY_DRAFT)
	s1.Spec.Paths.Value("/api/test").Get.Extensions["x-stability-level"] = draft
	s1.Spec.Paths.Value("/api/test").Post.Extensions = map[string]interface{}{"x-stability-level": draft}

	s2, err := open(getDeprecationFile("base-alpha-stability.yaml"))
	require.NoError(t, err)

	s2.Spec.Paths.Delete("/api/test")

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.APIDeprecationCheck), d, osm)
	require.Empty(t, errs)
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

func toJson(t *testing.T, value string) json.RawMessage {
	t.Helper()
	data, err := json.Marshal(value)
	require.NoError(t, err)
	return data
}

// BC: deprecating an operation with a deprecation policy and sunset date before required deprecation period is breaking
func TestBreaking_DeprecationWithEarlySunset(t *testing.T) {
	s1, err := open(getDeprecationFile("base.yaml"))
	require.NoError(t, err)

	s2, err := open(getDeprecationFile("deprecated-future.yaml"))
	require.NoError(t, err)
	sunsetDate := civil.DateOf(time.Now()).AddDays(9).String()
	s2.Spec.Paths.Value("/api/test").Get.Extensions[diff.SunsetExtension] = toJson(t, sunsetDate)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	c := singleCheckConfig(checker.APIDeprecationCheck).WithDeprecation(0, 10)
	errs := checker.CheckBackwardCompatibility(c, d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.APISunsetDateTooSmallId, errs[0].GetId())
	require.Equal(t, fmt.Sprintf("sunset date '%s' is too small, must be at least '10' days from now", sunsetDate), errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: deprecating an operation with a deprecation policy and sunset date after required deprecation period is not breaking
func TestBreaking_DeprecationWithProperSunset(t *testing.T) {

	s1, err := open(getDeprecationFile("base.yaml"))
	require.NoError(t, err)

	s2, err := open(getDeprecationFile("deprecated-future.yaml"))
	require.NoError(t, err)

	s2.Spec.Paths.Value("/api/test").Get.Extensions[diff.SunsetExtension] = toJson(t, civil.DateOf(time.Now()).AddDays(10).String())

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	c := singleCheckConfig(checker.APIDeprecationCheck).WithDeprecation(0, 10)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(c, d, osm, checker.INFO)
	require.Len(t, errs, 1)
	// only a non-breaking change detected
	require.Equal(t, checker.INFO, errs[0].GetLevel())
	require.Equal(t, "endpoint deprecated", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: deleting a path after sunset date of all contained operations is not breaking
func TestBreaking_DeprecationPathPast(t *testing.T) {

	s1, err := open(getDeprecationFile("deprecated-path-past.yaml"))
	require.NoError(t, err)

	s2, err := open(getDeprecationFile("sunset-path.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.APIDeprecationCheck), d, osm)
	require.Empty(t, errs)
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

// BC: deleting sunset header for a deprecated endpoint is breaking
func TestBreaking_SunsetDeletedForDeprecatedEndpoint(t *testing.T) {

	s1, err := open(getDeprecationFile("deprecated-with-sunset.yaml"))
	require.NoError(t, err)

	s2, err := open(getDeprecationFile("deprecated-no-sunset.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.APISunsetChangedCheck), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.APISunsetDeletedId, errs[0].GetId())
	require.Equal(t, "api sunset date deleted, but deprecated=true kept", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: changing sunset to an earlier date for a deprecated endpoint is breaking
func TestBreaking_SunsetModifiedForDeprecatedEndpoint(t *testing.T) {

	s1, err := open(getDeprecationFile("deprecated-future.yaml"))
	require.NoError(t, err)

	s2, err := open(getDeprecationFile("deprecated-past.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.APISunsetChangedCheck), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.APISunsetDateChangedTooSmallId, errs[0].GetId())
	require.Equal(t, "api sunset date changed to earlier date from '9999-08-10' to '2022-08-10', new sunset date must be not earlier than '9999-08-10' at least '180' days from now", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: changing sunset to an invalid date for a deprecated endpoint is breaking
func TestBreaking_SunsetModifiedToInvalidForDeprecatedEndpoint(t *testing.T) {

	s1, err := open(getDeprecationFile("deprecated-future.yaml"))
	require.NoError(t, err)

	s2, err := open(getDeprecationFile("deprecated-invalid.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.APISunsetChangedCheck), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.APIDeprecatedSunsetParseId, errs[0].GetId())
	require.Equal(t, "failed to parse sunset date: 'sunset date doesn't conform with RFC3339: invalid-date'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: changing sunset from an invalid date for a deprecated endpoint is breaking
func TestBreaking_SunsetModifiedFromInvalidForDeprecatedEndpoint(t *testing.T) {

	s1, err := open(getDeprecationFile("deprecated-invalid.yaml"))
	require.NoError(t, err)

	s2, err := open(getDeprecationFile("deprecated-future.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.APISunsetChangedCheck), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.APIDeprecatedSunsetParseId, errs[0].GetId())
	require.Equal(t, "failed to parse sunset date: 'sunset date doesn't conform with RFC3339: invalid-date'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: deleting other extension (not sunset) header for a deprecated endpoint is not breaking
func TestBreaking_NonSunsetDeletedForDeprecatedEndpoint(t *testing.T) {

	s1, err := open(getDeprecationFile("deprecated-with-other-extension.yaml"))
	require.NoError(t, err)

	s2, err := open(getDeprecationFile("deprecated-no-sunset.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.APISunsetChangedCheck), d, osm)
	require.Empty(t, errs)
}

// BC: no change to headers for a deprecated endpoint is not breaking
func TestBreaking_NoChangeToSunsetDeprecatedEndpoint(t *testing.T) {

	s1, err := open(getDeprecationFile("deprecated-future.yaml"))
	require.NoError(t, err)

	s2, err := open(getDeprecationFile("deprecated-future-2.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.APISunsetChangedCheck), d, osm)
	require.Empty(t, errs)
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

// CL: path operations that became deprecated
func TestApiDeprecated_DetectsDeprecatedOperations(t *testing.T) {
	s1, err := open("../data/deprecation/base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/deprecation/deprecated-future.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APIDeprecationCheck), d, osm, checker.INFO)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)

	require.IsType(t, checker.ApiChange{}, errs[0])
	e0 := errs[0].(checker.ApiChange)
	require.Equal(t, checker.EndpointDeprecatedId, e0.Id)
	require.Equal(t, "GET", e0.Operation)
	require.Equal(t, "/api/test", e0.Path)
	require.Equal(t, "endpoint deprecated", e0.GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// CL: path operations that were re-activated
func TestApiDeprecated_DetectsReactivatedOperations(t *testing.T) {
	s1, err := open("../data/deprecation/deprecated-future.yaml")
	require.NoError(t, err)

	s2, err := open("../data/deprecation/base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APIDeprecationCheck), d, osm, checker.INFO)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)

	require.IsType(t, checker.ApiChange{}, errs[0])
	e0 := errs[0].(checker.ApiChange)
	require.Equal(t, checker.EndpointReactivatedId, e0.Id)
	require.Equal(t, "GET", e0.Operation)
	require.Equal(t, "/api/test", e0.Path)
	require.Equal(t, "endpoint reactivated", e0.GetUncolorizedText(checker.NewDefaultLocalizer()))
}

func TestBreaking_InvaidStability(t *testing.T) {

	s1, err := open(getDeprecationFile("invalid-stability.yaml"))
	require.NoError(t, err)

	s2, err := open(getDeprecationFile("base-alpha-stability.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(singleCheckConfig(checker.APIDeprecationCheck), d, osm)
	require.Len(t, errs, 1)

	require.IsType(t, checker.ApiChange{}, errs[0])
	e0 := errs[0].(checker.ApiChange)
	require.Equal(t, checker.APIInvalidStabilityLevelId, e0.Id)
	require.Equal(t, "GET", e0.Operation)
	require.Equal(t, "/api/test", e0.Path)
	require.Equal(t, "failed to parse stability level: 'value is not one of draft, alpha, beta or stable: \"ga\"'", e0.GetUncolorizedText(checker.NewDefaultLocalizer()))
}
