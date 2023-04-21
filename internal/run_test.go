package internal_test

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/internal"
	"gopkg.in/yaml.v3"
)

func cmdToArgs(cmd string) []string {
	return strings.Split(cmd, " ")
}

func Test_NoArgs(t *testing.T) {
	require.Equal(t, 101, internal.Run([]string{}, io.Discard, io.Discard))
}

func Test_OneArg(t *testing.T) {
	require.Equal(t, 101, internal.Run(cmdToArgs("oasdiff"), io.Discard, io.Discard))
}

func Test_NoRevision(t *testing.T) {
	require.Equal(t, 101, internal.Run(cmdToArgs("oasdiff -base base.yaml"), io.Discard, io.Discard))
}

func Test_InvalidArg(t *testing.T) {
	require.Equal(t, 100, internal.Run(cmdToArgs("oasdiff -base data/openapi-test1.yaml -revision data/openapi-test1.yaml -deprecation-days 23s"), io.Discard, io.Discard))
}

func Test_BasicDiff(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff -base ../data/openapi-test1.yaml -revision ../data/openapi-test3.yaml"), io.Discard, io.Discard))
	var bc interface{}
	require.Nil(t, yaml.Unmarshal(stdout.Bytes(), &bc))
}

func Test_InvalidFile(t *testing.T) {
	require.Equal(t, 102, internal.Run(cmdToArgs("oasdiff -base ../data/no-file.yaml -revision ../data/openapi-test3.yaml"), io.Discard, io.Discard))
}

func Test_DiffInvalidFormat(t *testing.T) {
	require.Equal(t, 108, internal.Run(cmdToArgs("oasdiff -base ../data/openapi-test1.yaml -revision ../data/openapi-test3.yaml -format xxx"), io.Discard, io.Discard))
}

func Test_BasicBreakingChanges(t *testing.T) {
	require.Zero(t, internal.Run(cmdToArgs("oasdiff -base ../data/openapi-test1.yaml -revision ../data/openapi-test3.yaml -check-breaking"), io.Discard, io.Discard))
}

func Test_BreakingChangesInvalidFormat(t *testing.T) {
	require.Equal(t, 108, internal.Run(cmdToArgs("oasdiff -base ../data/openapi-test1.yaml -revision ../data/openapi-test3.yaml -check-breaking -format html"), io.Discard, io.Discard))
}

func Test_BreakingChangesYaml(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff -base ../data/openapi-test1.yaml -revision ../data/openapi-test3.yaml -check-breaking -format yaml"), &stdout, io.Discard))
	var bc interface{}
	require.Nil(t, yaml.Unmarshal(stdout.Bytes(), &bc))
}

func Test_BreakingChangesJson(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff -base ../data/openapi-test1.yaml -revision ../data/openapi-test3.yaml -check-breaking -format json"), &stdout, io.Discard))
	var bc interface{}
	require.Nil(t, json.Unmarshal(stdout.Bytes(), &bc))
}

func Test_BreakingChangesText(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff -base ../data/openapi-test1.yaml -revision ../data/openapi-test3.yaml -check-breaking"), &stdout, io.Discard))
	var bc interface{}
	require.Error(t, json.Unmarshal(stdout.Bytes(), &bc))
	require.Error(t, yaml.Unmarshal(stdout.Bytes(), &bc))
}

func Test_BreakingChangesFailOnDiff(t *testing.T) {
	require.Equal(t, 1, internal.Run(cmdToArgs("oasdiff -base ../data/openapi-test1.yaml -revision ../data/openapi-test3.yaml -check-breaking -fail-on-diff"), io.Discard, io.Discard))
}

func Test_BreakingChangesFailOnWarns(t *testing.T) {
	require.Equal(t, 1, internal.Run(cmdToArgs("oasdiff -base ../data/openapi-test1.yaml -revision ../data/openapi-test3.yaml -check-breaking -fail-on-diff -fail-on-warns"), io.Discard, io.Discard))
}

func Test_BreakingChangesFailOnWarnsErrsOnly(t *testing.T) {
	require.Equal(t, 1, internal.Run(cmdToArgs("oasdiff -base ../data/openapi-test2.yaml -revision ../data/openapi-test4.yaml -check-breaking -fail-on-diff -fail-on-warns"), io.Discard, io.Discard))
}

func Test_BreakingChangesFailOnDiffNoDiff(t *testing.T) {
	require.Zero(t, internal.Run(cmdToArgs("oasdiff -base ../data/openapi-test1.yaml -revision ../data/openapi-test1.yaml -check-breaking -fail-on-diff"), io.Discard, io.Discard))
}

func Test_BreakingChangesFailOnWarnsNoDiff(t *testing.T) {
	require.Zero(t, internal.Run(cmdToArgs("oasdiff -base ../data/openapi-test1.yaml -revision ../data/openapi-test1.yaml -check-breaking -fail-on-diff -fail-on-warns"), io.Discard, io.Discard))
}

func Test_Help(t *testing.T) {
	var stdout bytes.Buffer
	internal.Run(cmdToArgs("oasdiff -help"), &stdout, io.Discard)
	require.Contains(t, stdout.String(), "Usage of oasdiff")
}

func Test_HelpShortcut(t *testing.T) {
	var stdout bytes.Buffer
	internal.Run(cmdToArgs("oasdiff -h"), &stdout, io.Discard)
	require.Contains(t, stdout.String(), "Usage of oasdiff")
}

func Test_Version(t *testing.T) {
	var stdout bytes.Buffer
	internal.Run(cmdToArgs("oasdiff -version"), &stdout, io.Discard)
	require.Contains(t, stdout.String(), "oasdiff version:")
}
