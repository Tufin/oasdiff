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
	failOnDiff, diffEmpty, returnErr := internal.Run([]string{}, io.Discard)
	require.False(t, failOnDiff)
	require.False(t, diffEmpty)
	require.Equal(t, 101, returnErr.Code)
}

func Test_OneArg(t *testing.T) {
	failOnDiff, diffEmpty, returnErr := internal.Run([]string{"oasdiff"}, io.Discard)
	require.False(t, failOnDiff)
	require.False(t, diffEmpty)
	require.Equal(t, 101, returnErr.Code)
}

func Test_NoRevision(t *testing.T) {
	failOnDiff, diffEmpty, returnErr := internal.Run([]string{"oasdiff", "-base", "base.yaml"}, io.Discard)
	require.False(t, failOnDiff)
	require.False(t, diffEmpty)
	require.Equal(t, 101, returnErr.Code)
}

func Test_BasicDiff(t *testing.T) {
	failOnDiff, diffEmpty, returnErr := internal.Run([]string{"oasdiff", "-base", "../data/openapi-test1.yaml", "-revision", "../data/openapi-test3.yaml"}, io.Discard)
	require.False(t, failOnDiff)
	require.False(t, diffEmpty)
	require.Nil(t, returnErr)
}

func Test_DiffInvalidFormat(t *testing.T) {
	failOnDiff, diffEmpty, returnErr := internal.Run([]string{"oasdiff", "-base", "../data/openapi-test1.yaml", "-revision", "../data/openapi-test3.yaml", "-format", "xxx"}, io.Discard)
	require.False(t, failOnDiff)
	require.False(t, diffEmpty)
	require.Equal(t, 108, returnErr.Code)
}

func Test_BasicBreakingChanges(t *testing.T) {
	failOnDiff, diffEmpty, returnErr := internal.Run([]string{"oasdiff", "-base", "../data/openapi-test1.yaml", "-revision", "../data/openapi-test3.yaml", "-check-breaking"}, io.Discard)
	require.False(t, failOnDiff)
	require.False(t, diffEmpty)
	require.Nil(t, returnErr)
}

func Test_BreakingChangesInvalidFormat(t *testing.T) {
	failOnDiff, diffEmpty, returnErr := internal.Run([]string{"oasdiff", "-base", "../data/openapi-test1.yaml", "-revision", "../data/openapi-test3.yaml", "-check-breaking", "-format", "html"}, io.Discard)
	require.False(t, failOnDiff)
	require.False(t, diffEmpty)
	require.Equal(t, 108, returnErr.Code)
}

func Test_BreakingChangesYaml(t *testing.T) {
	var stdout bytes.Buffer
	failOnDiff, diffEmpty, returnErr := internal.Run([]string{"oasdiff", "-base", "../data/openapi-test1.yaml", "-revision", "../data/openapi-test3.yaml", "-check-breaking", "-format", "yaml"}, &stdout)
	require.False(t, failOnDiff)
	require.False(t, diffEmpty)
	require.Nil(t, returnErr)

	out := stdout.Bytes()
	var bc interface{}
	require.Nil(t, yaml.Unmarshal(out, &bc))
}

func Test_BreakingChangesJson(t *testing.T) {
	var stdout bytes.Buffer
	failOnDiff, diffEmpty, returnErr := internal.Run([]string{"oasdiff", "-base", "../data/openapi-test1.yaml", "-revision", "../data/openapi-test3.yaml", "-check-breaking", "-format", "json"}, &stdout)
	require.False(t, failOnDiff)
	require.False(t, diffEmpty)
	require.Nil(t, returnErr)

	out := stdout.Bytes()
	var bc interface{}
	require.Nil(t, json.Unmarshal(out, &bc))
}

func Test_BreakingChangesText(t *testing.T) {
	var stdout bytes.Buffer
	failOnDiff, diffEmpty, returnErr := internal.Run([]string{"oasdiff", "-base", "../data/openapi-test1.yaml", "-revision", "../data/openapi-test3.yaml", "-check-breaking"}, &stdout)
	require.False(t, failOnDiff)
	require.False(t, diffEmpty)
	require.Nil(t, returnErr)

	out := stdout.Bytes()
	var bc interface{}
	require.NotNil(t, json.Unmarshal(out, &bc))
	require.NotNil(t, yaml.Unmarshal(out, &bc))
}
