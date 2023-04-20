package internal_test

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/internal"
	"gopkg.in/yaml.v3"
)

func Test_NoArgs(t *testing.T) {
	require.Equal(t, 101, internal.Run([]string{}, io.Discard, io.Discard))
}

func Test_OneArg(t *testing.T) {
	require.Equal(t, 101, internal.Run([]string{"oasdiff"}, io.Discard, io.Discard))
}

func Test_NoRevision(t *testing.T) {
	require.Equal(t, 101, internal.Run([]string{"oasdiff", "-base", "base.yaml"}, io.Discard, io.Discard))
}

func Test_BasicDiff(t *testing.T) {
	require.Zero(t, internal.Run([]string{"oasdiff", "-base", "../data/openapi-test1.yaml", "-revision", "../data/openapi-test3.yaml"}, io.Discard, io.Discard))
}

func Test_DiffInvalidFormat(t *testing.T) {
	require.Equal(t, 108, internal.Run([]string{"oasdiff", "-base", "../data/openapi-test1.yaml", "-revision", "../data/openapi-test3.yaml", "-format", "xxx"}, io.Discard, io.Discard))
}

func Test_BasicBreakingChanges(t *testing.T) {
	require.Zero(t, internal.Run([]string{"oasdiff", "-base", "../data/openapi-test1.yaml", "-revision", "../data/openapi-test3.yaml", "-check-breaking"}, io.Discard, io.Discard))
}

func Test_BreakingChangesInvalidFormat(t *testing.T) {
	require.Equal(t, 108, internal.Run([]string{"oasdiff", "-base", "../data/openapi-test1.yaml", "-revision", "../data/openapi-test3.yaml", "-check-breaking", "-format", "html"}, io.Discard, io.Discard))
}

func Test_BreakingChangesYaml(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run([]string{"oasdiff", "-base", "../data/openapi-test1.yaml", "-revision", "../data/openapi-test3.yaml", "-check-breaking", "-format", "yaml"}, &stdout, io.Discard))
	out := stdout.Bytes()
	var bc interface{}
	require.Nil(t, yaml.Unmarshal(out, &bc))
}

func Test_BreakingChangesJson(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run([]string{"oasdiff", "-base", "../data/openapi-test1.yaml", "-revision", "../data/openapi-test3.yaml", "-check-breaking", "-format", "json"}, &stdout, io.Discard))
	out := stdout.Bytes()
	var bc interface{}
	require.Nil(t, json.Unmarshal(out, &bc))
}

func Test_BreakingChangesText(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run([]string{"oasdiff", "-base", "../data/openapi-test1.yaml", "-revision", "../data/openapi-test3.yaml", "-check-breaking"}, &stdout, io.Discard))
	out := stdout.Bytes()
	var bc interface{}
	require.NotNil(t, json.Unmarshal(out, &bc))
	require.NotNil(t, yaml.Unmarshal(out, &bc))
}

func Test_BreakingChangesFailOnDiff(t *testing.T) {
	var stdout bytes.Buffer
	require.Equal(t, 1, internal.Run([]string{"oasdiff", "-base", "../data/openapi-test1.yaml", "-revision", "../data/openapi-test3.yaml", "-check-breaking", "-fail-on-diff"}, &stdout, io.Discard))
}

func Test_BreakingChangesFailOnWarns(t *testing.T) {
	var stdout bytes.Buffer
	require.Equal(t, 1, internal.Run([]string{"oasdiff", "-base", "../data/openapi-test1.yaml", "-revision", "../data/openapi-test3.yaml", "-check-breaking", "-fail-on-diff", "-fail-on-warns"}, &stdout, io.Discard))
}

func Test_BreakingChangesFailOnWarnsErrsOnly(t *testing.T) {
	var stdout bytes.Buffer
	require.Equal(t, 1, internal.Run([]string{"oasdiff", "-base", "../data/openapi-test2.yaml", "-revision", "../data/openapi-test4.yaml", "-check-breaking", "-fail-on-diff", "-fail-on-warns"}, &stdout, io.Discard))
}

func Test_BreakingChangesFailOnDiffNoDiff(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run([]string{"oasdiff", "-base", "../data/openapi-test1.yaml", "-revision", "../data/openapi-test1.yaml", "-check-breaking", "-fail-on-diff"}, &stdout, io.Discard))
}

func Test_BreakingChangesFailOnWarnsNoDiff(t *testing.T) {
	var stdout bytes.Buffer
	require.Zero(t, internal.Run([]string{"oasdiff", "-base", "../data/openapi-test1.yaml", "-revision", "../data/openapi-test1.yaml", "-check-breaking", "-fail-on-diff", "-fail-on-warns"}, &stdout, io.Discard))
}
