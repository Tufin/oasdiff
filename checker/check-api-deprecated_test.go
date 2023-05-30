package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/checker/localizations"
	"github.com/tufin/oasdiff/diff"
)

// CL: path operations that became deprecated are detected
func TestApiDeprecated_DetectsDeprecatedOperations(t *testing.T) {
	s1, err := open("../data/deprecation/base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/deprecation/deprecated-no-sunset.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)

	config := checker.BackwardCompatibilityCheckConfig{
		Checks:              []checker.BackwardCompatibilityCheck{checker.APIDeprecatedCheck},
		MinSunsetBetaDays:   31,
		MinSunsetStableDays: 180,
		Localizer:           *localizations.New("en", "en"),
	}

	errs := checker.CheckBackwardCompatibility(config, d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)

	require.Equal(t, "api-path-deprecated", errs[0].Id)
	require.Equal(t, "GET", errs[0].Operation)
	require.Equal(t, "/api/test", errs[0].Path)
}

// CL: path operations that was re-activated are detected
func TestApiDeprecated_DetectsReactivatedOperations(t *testing.T) {
	s1, err := open("../data/deprecation/deprecated-no-sunset.yaml")
	require.NoError(t, err)

	s2, err := open("../data/deprecation/base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)

	config := checker.BackwardCompatibilityCheckConfig{
		Checks:              []checker.BackwardCompatibilityCheck{checker.APIDeprecatedCheck},
		MinSunsetBetaDays:   31,
		MinSunsetStableDays: 180,
		Localizer:           *localizations.New("en", "en"),
	}

	errs := checker.CheckBackwardCompatibility(config, d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)

	require.Equal(t, "api-path-reactivated", errs[0].Id)
	require.Equal(t, "GET", errs[0].Operation)
	require.Equal(t, "/api/test", errs[0].Path)
}
