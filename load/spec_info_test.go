package load_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/load"
)

func TestSpecInfo(t *testing.T) {
	_, err := load.LoadSpecInfo(MockLoader{}, load.NewSource("../data/openapi-test1.yaml"))
	require.NoError(t, err)
}

func TestSpecInfo_GlobOK(t *testing.T) {
	_, err := load.FromGlob(MockLoader{}, "../data/*.yaml")
	require.NoError(t, err)
}

func TestSpecInfo_InvalidSpec(t *testing.T) {
	_, err := load.FromGlob(MockLoader{}, "../data/ignore-err-example.txt")
	require.EqualError(t, err, "error unmarshaling JSON: while decoding JSON: json: cannot unmarshal string into Go value of type openapi3.TBis")
}

func TestSpecInfo_InvalidGlob(t *testing.T) {
	_, err := load.FromGlob(MockLoader{}, "[*")
	require.EqualError(t, err, "syntax error in pattern")
}

func TestSpecInfo_URL(t *testing.T) {
	_, err := load.FromGlob(MockLoader{}, "http://localhost/openapi-test1.yaml")
	require.EqualError(t, err, "no matching files (should be a glob, not a URL)")
}

func TestSpecInfo_GlobNoFiles(t *testing.T) {
	_, err := load.FromGlob(MockLoader{}, "../data/*.xxx")
	require.EqualError(t, err, "no matching files")
}
