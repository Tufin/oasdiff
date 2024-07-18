package checker_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
)

func TestStringCond_Info(t *testing.T) {
	level := checker.INFO
	require.Equal(t, level.PrettyString(), level.StringCond(checker.ColorAlways))
	require.Equal(t, level.String(), level.StringCond(checker.ColorNever))
	require.Equal(t, level.String(), level.StringCond(checker.ColorAuto))
	require.Equal(t, level.String(), level.StringCond(checker.ColorInvalid))
}

func TestPrettyString(t *testing.T) {
	require.Equal(t, "\x1b[36minfo\x1b[0m", checker.INFO.PrettyString())
	require.Equal(t, "\x1b[35mwarning\x1b[0m", checker.WARN.PrettyString())
	require.Equal(t, "\x1b[31merror\x1b[0m", checker.ERR.PrettyString())
	require.Equal(t, "\x1b[37missue\x1b[0m", checker.Level(4).PrettyString())
}

func TestProcessSeverityLevels_InvalidFile(t *testing.T) {
	m, err := checker.ProcessSeverityLevels("../data/invalid.txt")
	require.Nil(t, m)
	require.EqualError(t, err, "open ../data/invalid.txt: no such file or directory")
}

func TestProcessSeverityLevels_OK(t *testing.T) {
	m, err := checker.ProcessSeverityLevels("../data/severity-levels.txt")
	require.Equal(t, map[string]checker.Level{"api-security-removed": 2, "request-parameter-enum-value-added": 3}, m)
	require.NoError(t, err)
}

func TestProcessSeverityLevels_WrongFlags(t *testing.T) {
	m, err := checker.ProcessSeverityLevels("../data/severity-levels.txt")
	require.Equal(t, map[string]checker.Level{"api-security-removed": 2, "request-parameter-enum-value-added": 3}, m)
	require.NoError(t, err)
}

func TestGetSeverityLevels_InvalidLine(t *testing.T) {
	m, err := checker.GetSeverityLevels(strings.NewReader("invalid"))
	require.Nil(t, m)
	require.EqualError(t, err, "invalid line #1: invalid")
}

func TestGetSeverityLevels_InvalidRuleId(t *testing.T) {
	m, err := checker.GetSeverityLevels(strings.NewReader("invalid_id err"))
	require.Nil(t, m)
	require.EqualError(t, err, "invalid rule id \"invalid_id\" on line 1")
}

func TestGetSeverityLevels_InvalidLevel(t *testing.T) {
	m, err := checker.GetSeverityLevels(strings.NewReader("request-parameter-enum-value-added invalid_level"))
	require.Nil(t, m)
	require.EqualError(t, err, "invalid level \"invalid_level\" on line 1")
}

func TestGetSeverityLevels_Duplicate(t *testing.T) {
	m, err := checker.GetSeverityLevels(strings.NewReader("request-parameter-enum-value-added info\nrequest-parameter-enum-value-added warn"))
	require.Equal(t, map[string]checker.Level{"request-parameter-enum-value-added": checker.WARN}, m)
	require.NoError(t, err)
}
