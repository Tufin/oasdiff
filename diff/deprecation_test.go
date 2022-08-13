package diff_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
)

func getDeprecationFile(file string) string {
	return fmt.Sprintf("../data/deprecation/%s", file)
}

// BC: deleting an operation before sunset is breaking
func TestBreaking_DeprecationFuture(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getDeprecationFile("deprecated-future.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getDeprecationFile("sunset.yaml"))
	require.NoError(t, err)

	dd, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.NotEmpty(t, dd)
}

// BC: deleting an operation without sunset is breaking
func TestBreaking_DeprecationNoSunset(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getDeprecationFile("deprecated-no-sunset.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getDeprecationFile("sunset.yaml"))
	require.NoError(t, err)

	dd, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.NotEmpty(t, dd)
}

// BC: deleting an operation after sunset is not breaking
func TestBreaking_DeprecationPast(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getDeprecationFile("deprecated-past.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getDeprecationFile("sunset.yaml"))
	require.NoError(t, err)

	dd, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.Empty(t, dd)
}

// BC: deprecating an operation with a deprecation policy but without specifying sunset date is breaking
func TestBreaking_DeprecationWithoutSunset(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getDeprecationFile("deprecated-base.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getDeprecationFile("deprecated-no-sunset.yaml"))
	require.NoError(t, err)

	dd, err := diff.Get(&diff.Config{
		BreakingOnly:    true,
		DeprecationDays: 10,
	}, s1, s2)
	require.NoError(t, err)
	require.NotEmpty(t, dd)
}

// BC: deprecating an operation without a deprecation policy and without specifying sunset date is not breaking
func TestBreaking_DeprecationWithoutPoicy(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getDeprecationFile("deprecated-base.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getDeprecationFile("deprecated-no-sunset.yaml"))
	require.NoError(t, err)

	dd, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.Empty(t, dd)
}

func toJson(t *testing.T, value string) json.RawMessage {
	t.Helper()
	data, err := json.Marshal(value)
	require.NoError(t, err)
	return data
}

// BC: deprecating an operation with a deprecation policy and sunset date before required deprecation period is breaking
func TestBreaking_DeprecationWithEarlySunset(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getDeprecationFile("deprecated-base.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getDeprecationFile("deprecated-future.yaml"))
	require.NoError(t, err)
	s2.Paths["/api/test"].Get.ExtensionProps.Extensions[diff.SunsetExtension] = toJson(t, civil.DateOf(time.Now()).AddDays(9).String())

	dd, err := diff.Get(&diff.Config{
		BreakingOnly:    true,
		DeprecationDays: 10,
	}, s1, s2)
	require.NoError(t, err)
	require.NotEmpty(t, dd)
}

// BC: deprecating an operation with a deprecation policy and sunset date after required deprecation period is not breaking
func TestBreaking_DeprecationWithProperSunset(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getDeprecationFile("deprecated-base.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getDeprecationFile("deprecated-future.yaml"))
	require.NoError(t, err)

	s2.Paths["/api/test"].Get.ExtensionProps.Extensions[diff.SunsetExtension] = toJson(t, civil.DateOf(time.Now()).AddDays(10).String())

	dd, err := diff.Get(&diff.Config{
		BreakingOnly:    true,
		DeprecationDays: 10,
	}, s1, s2)
	require.NoError(t, err)
	require.Empty(t, dd)
}

// BC: deleting a path after sunset of all contained operations is not breaking
func TestBreaking_DeprecationPathPast(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getDeprecationFile("deprecated-path-past.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getDeprecationFile("sunset-path.yaml"))
	require.NoError(t, err)

	dd, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.Empty(t, dd)
}
