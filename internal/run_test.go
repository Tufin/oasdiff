package internal_test

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/oasdiff/telemetry/model"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/formatters"
	"github.com/tufin/oasdiff/internal"
	"gopkg.in/yaml.v3"
)

func cmdToArgs(cmd string) []string {
	return strings.Fields(cmd)
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	os.Setenv(model.EnvNoTelemetry, "1")
}

func Test_InvalidCmd(t *testing.T) {
	require.Equal(t, 100, internal.Run(cmdToArgs("oasdiff invalid"), io.Discard, io.Discard))
}

func Test_NoRevision(t *testing.T) {
	require.Equal(t, 100, internal.Run(cmdToArgs("oasdiff diff base.yaml"), io.Discard, io.Discard))
}

func Test_InvalidFlag(t *testing.T) {
	require.Equal(t, 100, internal.Run(cmdToArgs("oasdiff diff data/openapi-test1.yaml data/openapi-test1.yaml --invalid"), io.Discard, io.Discard))
}

func Test_BasicDiff(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff diff ../data/openapi-test1.yaml ../data/openapi-test3.yaml --exclude-elements endpoints"), &stdout, io.Discard))
	var bc interface{}
	require.Nil(t, yaml.Unmarshal(stdout.Bytes(), &bc))
}

func Test_DiffJson(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff diff ../data/openapi-test1.yaml ../data/openapi-test3.yaml -f json --exclude-elements endpoints"), &stdout, io.Discard))
	var bc interface{}
	require.Nil(t, json.Unmarshal(stdout.Bytes(), &bc))
}

func Test_DiffHtml(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff diff ../data/openapi-test1.yaml ../data/openapi-test3.yaml -f html"), &stdout, io.Discard))
	require.Contains(t, stdout.String(), `<h3 id="new-endpoints-none">New Endpoints: None</h3>`)
}

func Test_DiffText(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff diff ../data/openapi-test1.yaml ../data/openapi-test3.yaml -f text"), &stdout, io.Discard))
	require.Contains(t, stdout.String(), `### New Endpoints: None`)
}

func Test_Summary(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff summary ../data/openapi-test1.yaml ../data/openapi-test3.yaml"), &stdout, io.Discard))
	require.Contains(t, stdout.String(), `diff: true`)
}

func Test_InvalidFile(t *testing.T) {
	var stderr bytes.Buffer
	require.Equal(t, 102, internal.Run(cmdToArgs("oasdiff diff no-file ../data/openapi-test3.yaml"), io.Discard, &stderr))
	require.Condition(t, func() (success bool) {
		return stderr.String() == "Error: failed to load base spec from \"no-file\": open no-file: no such file or directory\n" ||
			stderr.String() == "Error: failed to load base spec from \"no-file\": open no-file: The system cannot find the file specified.\n" // windows
	}, stderr.String())
}

func Test_InvalidGlob(t *testing.T) {
	var stderr bytes.Buffer
	require.Equal(t, 103, internal.Run(cmdToArgs(`oasdiff diff -c "a[" ../data/openapi-test3.yaml`), io.Discard, &stderr))
	require.Equal(t, "Error: failed to load base specs from glob \"\\\"a[\\\"\": syntax error in pattern\n", stderr.String())
}

func Test_GlobNoFiles(t *testing.T) {
	var stderr bytes.Buffer
	require.Equal(t, 103, internal.Run(cmdToArgs("oasdiff diff -c no-file ../data/openapi-test3.yaml"), io.Discard, &stderr))
	require.Equal(t, "Error: failed to load base specs from glob \"no-file\": no matching files\n", stderr.String())
}

func Test_GlobWithUrl(t *testing.T) {
	var stderr bytes.Buffer
	require.Equal(t, 103, internal.Run(cmdToArgs("oasdiff diff -c ../data/openapi-test1.yaml https://"), io.Discard, &stderr))
	require.Equal(t, "Error: failed to load revision specs from glob \"https://\": no matching files (should be a glob, not a URL)\n", stderr.String())
}

func Test_DiffInvalidFormat(t *testing.T) {
	require.Equal(t, 100, internal.Run(cmdToArgs("oasdiff diff ../data/openapi-test1.yaml ../data/openapi-test3.yaml --format xxx"), io.Discard, io.Discard))
}

func Test_BreakingChangesIncludeChecks(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff breaking ../data/run_test/breaking_changes_include_checks_base.yaml ../data/run_test/breaking_changes_include_checks_revision.yaml --include-checks response-non-success-status-removed,api-tag-removed --format json"), &stdout, io.Discard))
	bc := formatters.Changes{}
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &bc))
	require.Len(t, bc, 2)
	for _, c := range bc {
		require.Equal(t, c.Level, checker.ERR)
	}
}

func Test_BasicBreakingChanges(t *testing.T) {
	require.Zero(t, internal.Run(cmdToArgs("oasdiff breaking ../data/openapi-test1.yaml ../data/openapi-test3.yaml"), io.Discard, io.Discard))
}

func Test_BreakingChangesInvalidFormat(t *testing.T) {
	require.Equal(t, 100, internal.Run(cmdToArgs("oasdiff breaking ../data/openapi-test1.yaml ../data/openapi-test3.yaml --format html"), io.Discard, io.Discard))
}

func Test_BreakingChangesYaml(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff breaking ../data/openapi-test1.yaml ../data/openapi-test3.yaml --format yaml"), &stdout, io.Discard))
	var bc interface{}
	require.Nil(t, yaml.Unmarshal(stdout.Bytes(), &bc))
}

func Test_BreakingChangesJson(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff breaking ../data/openapi-test1.yaml ../data/openapi-test3.yaml --format json"), &stdout, io.Discard))
	var bc interface{}
	require.Nil(t, json.Unmarshal(stdout.Bytes(), &bc))
}

func Test_BreakingChangesText(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff breaking ../data/openapi-test1.yaml ../data/openapi-test3.yaml"), &stdout, io.Discard))
	var bc interface{}
	require.Error(t, json.Unmarshal(stdout.Bytes(), &bc))
	require.Error(t, yaml.Unmarshal(stdout.Bytes(), &bc))
}

func Test_BreakingChangesFailOnErr(t *testing.T) {
	require.Equal(t, 1, internal.Run(cmdToArgs("oasdiff breaking ../data/openapi-test1.yaml ../data/openapi-test3.yaml --fail-on ERR"), io.Discard, io.Discard))
}

func Test_BreakingChangesFailOnWarn(t *testing.T) {
	require.Equal(t, 1, internal.Run(cmdToArgs("oasdiff breaking ../data/openapi-test1.yaml ../data/openapi-test3.yaml --fail-on WARN"), io.Discard, io.Discard))
}

func Test_BreakingChangesFailOnWarnsErrsOnly(t *testing.T) {
	require.Equal(t, 1, internal.Run(cmdToArgs("oasdiff breaking ../data/openapi-test2.yaml ../data/openapi-test4.yaml --fail-on WARN"), io.Discard, io.Discard))
}

func Test_BreakingChangesFailOnDiffNoDiff(t *testing.T) {
	require.Zero(t, internal.Run(cmdToArgs("oasdiff breaking ../data/openapi-test1.yaml ../data/openapi-test1.yaml --fail-on ERR"), io.Discard, io.Discard))
}

func Test_BreakingChangesFailOnWarnsNoDiff(t *testing.T) {
	require.Zero(t, internal.Run(cmdToArgs("oasdiff breaking ../data/openapi-test1.yaml ../data/openapi-test1.yaml --fail-on WARN"), io.Discard, io.Discard))
}

func Test_BreakingChangesIgnoreErrs(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff breaking ../data/openapi-test1.yaml ../data/openapi-test3.yaml --err-ignore ../data/ignore-err-example.txt --format json"), &stdout, io.Discard))
	bc := formatters.Changes{}
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &bc))
	require.Len(t, bc, 5)
}

func Test_BreakingChangesIgnoreErrsAndWarns(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff breaking ../data/openapi-test1.yaml ../data/openapi-test3.yaml --err-ignore ../data/ignore-err-example.txt --warn-ignore ../data/ignore-warn-example.txt --format json"), &stdout, io.Discard))
	bc := formatters.Changes{}
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &bc))
	require.Len(t, bc, 4)
}

func Test_BreakingChangesIgnoreErrsApiSchemaOptional(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff breaking ../data/openapi-test1.yaml ../data/openapi-test3.yaml --err-ignore ../data/ignore-err-example.txt --warn-ignore ../data/ignore-warn-example.txt --include-checks api-schema-removed --format json"), &stdout, io.Discard))
	bc := formatters.Changes{}
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &bc))
	require.Len(t, bc, 4)
}

func Test_BreakingChangesInvalidIgnoreFile(t *testing.T) {
	require.Equal(t, 121, internal.Run(cmdToArgs("oasdiff breaking ../data/openapi-test1.yaml ../data/openapi-test3.yaml --err-ignore no-file"), io.Discard, io.Discard))
}

func Test_ComposedMode(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff diff ../data/composed/base/*.yaml ../data/composed/revision/*.yaml --composed --exclude-elements endpoints"), &stdout, io.Discard))
	var bc interface{}
	require.NoError(t, yaml.Unmarshal(stdout.Bytes(), &bc))
	require.Equal(t, map[string]interface{}{"paths": map[string]interface{}{"deleted": []interface{}{"/api/old-test"}}}, bc)
}

func Test_ComposedModeInvalidFile(t *testing.T) {
	var stderr bytes.Buffer
	require.Equal(t, 103, internal.Run(cmdToArgs("oasdiff diff ../data/allof/* ../data/allof/* --composed --flatten"), io.Discard, &stderr))

	require.Condition(t, func() (success bool) {
		return stderr.String() == "Error: failed to load base specs from glob \"../data/allof/*\": failed to flatten allOf in \"../data/allof/invalid.yaml\": unable to resolve Type conflict: all Type values must be identical\n" ||
			stderr.String() == "Error: failed to load base specs from glob \"../data/allof/*\": failed to flatten allOf in \"..\\data\\allof\\invalid.yaml\": unable to resolve Type conflict: all Type values must be identical" // windows
	}, stderr.String())
}

func Test_Help(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff --help"), &stdout, io.Discard))
	require.Contains(t, stdout.String(), "Usage:")
}

func Test_HelpShortcut(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff -h"), &stdout, io.Discard))
	require.Contains(t, stdout.String(), "Usage:")
}

func Test_Version(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff -v"), &stdout, io.Discard))
	require.Contains(t, stdout.String(), "oasdiff version")
}

func Test_StripPrefixBase(t *testing.T) {
	require.Zero(t, internal.Run(cmdToArgs("oasdiff breaking ../data/simple.yaml ../data/simple.yaml --strip-prefix-base /partner-api"), io.Discard, io.Discard))
}

func Test_DuplicatePathsFail(t *testing.T) {
	require.NotZero(t, internal.Run(cmdToArgs("oasdiff breaking ../data/duplicate_endpoints/base.yaml ../data/duplicate_endpoints/revision.yaml"), io.Discard, io.Discard))
}

func Test_DuplicatePathsOK(t *testing.T) {
	require.Zero(t, internal.Run(cmdToArgs("oasdiff breaking ../data/duplicate_endpoints/base.yaml ../data/duplicate_endpoints/revision.yaml --include-path-params"), io.Discard, io.Discard))
}

func Test_Changelog(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff changelog ../data/run_test/changelog_base.yaml ../data/run_test/changelog_revision.yaml --format json"), &stdout, io.Discard))
	cl := formatters.Changes{}
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &cl))
	require.Len(t, cl, 1)
}

func Test_BreakingChangesChangelogOptionalCheckersAreInfoLevel(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff changelog ../data/run_test/changelog_include_checks_base.yaml ../data/run_test/changelog_include_checks_revision.yaml --format json"), &stdout, io.Discard))
	cl := formatters.Changes{}
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &cl))
	require.Len(t, cl, 2)
	for _, c := range cl {
		require.Equal(t, c.Level, checker.INFO)
	}
}

func Test_BreakingChangesChangelogOptionalCheckersAreErrorLevelWhenSpecified(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run(cmdToArgs("oasdiff changelog ../data/run_test/changelog_include_checks_base.yaml ../data/run_test/changelog_include_checks_revision.yaml --format json --include-checks api-tag-removed,response-non-success-status-removed"), &stdout, io.Discard))
	cl := formatters.Changes{}
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &cl))
	require.Len(t, cl, 2)
	for _, c := range cl {
		require.Equal(t, c.Level, checker.ERR)
	}
}

func Test_BreakingChangesFlattenDeprecated(t *testing.T) {
	require.Zero(t, internal.Run(cmdToArgs("oasdiff breaking ../data/allof/simple.yaml ../data/allof/revision.yaml --flatten --fail-on ERR"), io.Discard, io.Discard))
}

func Test_BreakingChangesFlattenAllOf(t *testing.T) {
	require.Zero(t, internal.Run(cmdToArgs("oasdiff breaking ../data/allof/simple.yaml ../data/allof/revision.yaml --flatten-allof --fail-on ERR"), io.Discard, io.Discard))
}

func Test_BreakingChangesFlattenCommonParams(t *testing.T) {
	require.Zero(t, internal.Run(cmdToArgs("oasdiff breaking ../data/common-params/params_in_path.yaml ../data/common-params/params_in_op.yaml --flatten-params --fail-on ERR"), io.Discard, io.Discard))
}

func Test_BreakingChangesCaseInsensitiveHeaders(t *testing.T) {
	require.Zero(t, internal.Run(cmdToArgs("oasdiff breaking ../data/header-case/base.yaml ../data/header-case/revision.yaml --case-insensitive-headers --fail-on ERR"), io.Discard, io.Discard))
}

func Test_FlattenCmdOK(t *testing.T) {
	require.Zero(t, internal.Run(cmdToArgs("oasdiff flatten ../data/allof/simple.yaml"), io.Discard, io.Discard))
}

func Test_FlattenCmdInvalid(t *testing.T) {
	var stderr bytes.Buffer
	require.Equal(t, 102, internal.Run(cmdToArgs("oasdiff flatten ../data/allof/invalid.yaml"), io.Discard, &stderr))
	require.Equal(t, `Error: failed to load original spec from "../data/allof/invalid.yaml": failed to flatten allOf in "../data/allof/invalid.yaml": unable to resolve Type conflict: all Type values must be identical
`, stderr.String())
}

func Test_Checks(t *testing.T) {
	require.Zero(t, internal.Run(cmdToArgs("oasdiff checks -l ru"), io.Discard, io.Discard))
}

func Test_Color(t *testing.T) {
	require.Zero(t, internal.Run(cmdToArgs("oasdiff breaking ../data/allof/simple.yaml ../data/allof/revision.yaml --color always"), io.Discard, io.Discard))
}

func Test_ColorWithNonTextFormat(t *testing.T) {
	var stderr bytes.Buffer
	require.Equal(t, 100, internal.Run(cmdToArgs("oasdiff changelog ../data/allof/simple.yaml ../data/allof/revision.yaml -f yaml --color always"), io.Discard, &stderr))
	require.Equal(t, "Error: --color flag is only relevant with 'text' or 'singleline' formats\n", stderr.String())
}

func Test_Delta(t *testing.T) {
	require.Zero(t, internal.Run(cmdToArgs("oasdiff delta ../data/simple1.yaml ../data/simple2.yaml"), io.Discard, io.Discard))
}
