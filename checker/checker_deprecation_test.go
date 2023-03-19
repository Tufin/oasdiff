package checker_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

func getDeprecationFile(file string) string {
	return fmt.Sprintf("../data/deprecation/%s", file)
}

// BC: deleting an operation before sunset date is breaking
func TestBreaking_DeprecationEarlySunset(t *testing.T) {

	s1, err := checker.LoadOpenAPISpecInfoFromFile(getDeprecationFile("deprecated-future.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getDeprecationFile("sunset.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "api-removed-before-sunset", errs[0].Id)
}

// BC: deleting an operation without sunset date is breaking
func TestBreaking_DeprecationNoSunset(t *testing.T) {

	s1, err := checker.LoadOpenAPISpecInfoFromFile(getDeprecationFile("deprecated-no-sunset.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getDeprecationFile("sunset.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NoError(t, err)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "api-sunset-parse", errs[0].Id)
}

// BC: deleting an operation after sunset date is not breaking
func TestBreaking_DeprecationPast(t *testing.T) {

	s1, err := checker.LoadOpenAPISpecInfoFromFile(getDeprecationFile("deprecated-past.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getDeprecationFile("sunset.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: deprecating an operation with a deprecation policy but without specifying sunset date is breaking
func TestBreaking_DeprecationWithoutSunset(t *testing.T) {

	s1, err := checker.LoadOpenAPISpecInfoFromFile(getDeprecationFile("base.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getDeprecationFile("deprecated-no-sunset.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	c := checker.GetDefaultChecks()
	c.MinSunsetStableDays = 10
	errs := checker.CheckBackwardCompatibility(c, d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "api-deprecated-sunset-parse", errs[0].Id)
}

// BC: deprecating an operation without a deprecation policy and without specifying sunset date is not breaking
func TestBreaking_DeprecationForAlpha(t *testing.T) {

	s1, err := checker.LoadOpenAPISpecInfoFromFile(getDeprecationFile("base-alpha-stability.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getDeprecationFile("deprecated-no-sunset-alpha-stability.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: removing the path without a deprecation policy and without specifying sunset date is not breaking for alpha level
func TestBreaking_RemovedPathForAlpha(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getDeprecationFile("base-alpha-stability.yaml"))
	require.NoError(t, err)
	alpha := toJson(t, "alpha")
	s1.Spec.Paths["/api/test"].Get.Extensions["x-stability-level"] = alpha
	s1.Spec.Paths["/api/test"].Post.Extensions["x-stability-level"] = alpha

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getDeprecationFile("base-alpha-stability.yaml"))
	require.NoError(t, err)

	delete(s2.Spec.Paths, "/api/test")

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: removing the path without a deprecation policy and without specifying sunset date is breaking if some APIs are not alpha stability level
func TestBreaking_RemovedPathForAlphaBreaking(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getDeprecationFile("base-alpha-stability.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getDeprecationFile("base-alpha-stability.yaml"))
	require.NoError(t, err)

	delete(s2.Spec.Paths, "/api/test")

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
}

// BC: deprecating an operation without a deprecation policy and without specifying sunset date is not breaking for draft level
func TestBreaking_DeprecationForDraft(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getDeprecationFile("base-alpha-stability.yaml"))
	require.NoError(t, err)
	draft := toJson(t, "draft")
	s1.Spec.Paths["/api/test"].Get.Extensions["x-stability-level"] = draft

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getDeprecationFile("deprecated-no-sunset-alpha-stability.yaml"))
	require.NoError(t, err)
	s2.Spec.Paths["/api/test"].Get.Extensions["x-stability-level"] = draft

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: removing the path without a deprecation policy and without specifying sunset date is not breaking for draft level
func TestBreaking_RemovedPathForDraft(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getDeprecationFile("base-alpha-stability.yaml"))
	require.NoError(t, err)
	draft := toJson(t, "draft")
	s1.Spec.Paths["/api/test"].Get.Extensions["x-stability-level"] = draft
	s1.Spec.Paths["/api/test"].Post.Extensions["x-stability-level"] = draft

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getDeprecationFile("base-alpha-stability.yaml"))
	require.NoError(t, err)

	delete(s2.Spec.Paths, "/api/test")

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: removing the path without a deprecation policy and without specifying sunset date is breaking if some APIs are not draft stability level
func TestBreaking_RemovedPathForDraftBreaking(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getDeprecationFile("base-alpha-stability.yaml"))
	require.NoError(t, err)
	draft := toJson(t, "draft")
	s1.Spec.Paths["/api/test"].Get.Extensions["x-stability-level"] = draft

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getDeprecationFile("base-alpha-stability.yaml"))
	require.NoError(t, err)

	delete(s2.Spec.Paths, "/api/test")

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
}

func toJson(t *testing.T, value string) json.RawMessage {
	t.Helper()
	data, err := json.Marshal(value)
	require.NoError(t, err)
	return data
}

// BC: deprecating an operation with a deprecation policy and sunset date before required deprecation period is breaking
func TestBreaking_DeprecationWithEarlySunset(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getDeprecationFile("base.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getDeprecationFile("deprecated-future.yaml"))
	require.NoError(t, err)
	s2.Spec.Paths["/api/test"].Get.Extensions[diff.SunsetExtension] = toJson(t, civil.DateOf(time.Now()).AddDays(9).String())

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	c := checker.GetDefaultChecks()
	c.MinSunsetStableDays = 10
	errs := checker.CheckBackwardCompatibility(c, d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "api-sunset-date-too-small", errs[0].Id)
}

// BC: deprecating an operation with a deprecation policy and sunset date after required deprecation period is not breaking
func TestBreaking_DeprecationWithProperSunset(t *testing.T) {

	s1, err := checker.LoadOpenAPISpecInfoFromFile(getDeprecationFile("base.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getDeprecationFile("deprecated-future.yaml"))
	require.NoError(t, err)

	s2.Spec.Paths["/api/test"].Get.Extensions[diff.SunsetExtension] = toJson(t, civil.DateOf(time.Now()).AddDays(10).String())

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	c := checker.GetDefaultChecks()
	c.MinSunsetStableDays = 10
	errs := checker.CheckBackwardCompatibility(c, d, osm)
	require.Empty(t, errs)
}

// BC: deleting a path after sunset date of all contained operations is not breaking
func TestBreaking_DeprecationPathPast(t *testing.T) {

	s1, err := checker.LoadOpenAPISpecInfoFromFile(getDeprecationFile("deprecated-path-past.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getDeprecationFile("sunset-path.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: deleting a path with some operations having sunset date in the future is breaking
func TestBreaking_DeprecationPathMixed(t *testing.T) {

	s1, err := checker.LoadOpenAPISpecInfoFromFile(getDeprecationFile("deprecated-path-mixed.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getDeprecationFile("sunset-path.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "api-path-removed-before-sunset", errs[0].Id)
}
