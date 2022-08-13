package diff_test

import (
	"fmt"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
)

func getDeprecationFile(file string) string {
	return fmt.Sprintf("../data/deprecation/%s", file)
}

// BC: deleting an endpoint before sunset is breaking
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

// BC: deleting an endpoint without sunset is breaking
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

// BC: deleting an endpoint after sunset is not breaking
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
