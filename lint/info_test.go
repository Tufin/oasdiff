package lint_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/lint"
)

func TestInfo_NoInfo(t *testing.T) {

	const source = "../data/lint/info/no-info.yaml"
	errs := lint.Run(lint.NewConfig([]lint.Check{lint.InfoCheck}), source, loadFrom(t, source))
	require.Len(t, errs, 1)
	require.Equal(t, "info-missing", errs[0].Id)
	require.Equal(t, lint.LEVEL_ERROR, errs[0].Level)
	require.Equal(t, source, errs[0].Source)
}

func TestInfo_TitleMissing(t *testing.T) {

	const source = "../data/lint/info/title-missing.yaml"
	errs := lint.Run(lint.NewConfig([]lint.Check{lint.InfoCheck}), source, loadFrom(t, source))
	require.Len(t, errs, 1)
	require.Equal(t, "info-title-missing", errs[0].Id)
	require.Equal(t, lint.LEVEL_ERROR, errs[0].Level)
	require.Equal(t, source, errs[0].Source)
}

func TestInfo_VersionMissing(t *testing.T) {

	const source = "../data/lint/info/version-missing.yaml"
	errs := lint.Run(lint.NewConfig([]lint.Check{lint.InfoCheck}), source, loadFrom(t, source))
	require.Len(t, errs, 1)
	require.Equal(t, "info-version-missing", errs[0].Id)
	require.Equal(t, lint.LEVEL_ERROR, errs[0].Level)
	require.Equal(t, source, errs[0].Source)
}

func TestInfo_InvalidTOS(t *testing.T) {

	const source = "../data/lint/info/invalid-terms-of-service.yaml"
	errs := lint.Run(lint.NewConfig([]lint.Check{lint.InfoCheck}), source, loadFrom(t, source))
	require.Len(t, errs, 1)
	require.Equal(t, "info-invalid-terms-of-service", errs[0].Id)
	require.Equal(t, lint.LEVEL_ERROR, errs[0].Level)
	require.Equal(t, "terms of service must be in the format of a URL: bla", errs[0].Text)
	require.Equal(t, source, errs[0].Source)
}
