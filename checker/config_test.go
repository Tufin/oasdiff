package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
)

func TestNewConfig(t *testing.T) {
	config := checker.NewConfig()
	require.NotEmpty(t, config.Checks)
	require.Empty(t, config.LogLevelOverrides)
	require.Equal(t, checker.BetaDeprecationDays, config.MinSunsetBetaDays)
	require.Equal(t, checker.StableDeprecationDays, config.MinSunsetStableDays)
}

func TestNewConfigWithDeprecation(t *testing.T) {
	config := checker.NewConfig().WithDeprecation(10, 20)
	require.NotEmpty(t, config.Checks)
	require.Empty(t, config.LogLevelOverrides)
	require.Equal(t, 10, config.MinSunsetBetaDays)
	require.Equal(t, 20, config.MinSunsetStableDays)
}

func TestNewConfigWithOptionalCheck(t *testing.T) {
	config := checker.NewConfig().WithOptionalCheck("id")
	require.Equal(t, checker.ERR, config.LogLevelOverrides["id"])
}
