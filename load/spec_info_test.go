package load_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/load"
)

func TestLoadSpecInfo(t *testing.T) {
	_, err := load.LoadSpecInfo(MockLoader{}, load.GetSource("openapi-test1.yaml"))
	require.NoError(t, err)
}

func TestLoadGlob_OK(t *testing.T) {
	_, err := load.FromGlob(MockLoader{}, RelativeDataPath+"*.yaml")
	require.NoError(t, err)
}

func TestLoadGlob_InvalidSpec(t *testing.T) {
	_, err := load.FromGlob(MockLoader{}, RelativeDataPath+"ignore-err-example.txt")
	require.Error(t, err)
	require.Equal(t, "error unmarshaling JSON: while decoding JSON: json: cannot unmarshal string into Go value of type openapi3.TBis", err.Error())
}

func TestLoadGlob_Invalid(t *testing.T) {
	_, err := load.FromGlob(MockLoader{}, "[*")
	require.Error(t, err)
	require.Equal(t, "syntax error in pattern", err.Error())
}

func TestLoadGlob_URL(t *testing.T) {
	_, err := load.FromGlob(MockLoader{}, "http://localhost/openapi-test1.yaml")
	require.Error(t, err)
	require.Equal(t, "no matching files (should be a glob, not a URL)", err.Error())
}

func TestLoadGlob_NoFiles(t *testing.T) {
	_, err := load.FromGlob(MockLoader{}, RelativeDataPath+"*.xxx")
	require.Error(t, err)
	require.Equal(t, "no matching files", err.Error())
}
