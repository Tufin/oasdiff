package load_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/load"
)

func TestLoadPath_NoError(t *testing.T) {
	_, err := load.LoadPath("../data/openapi-test1.yaml")
	require.NoError(t, err)
}

func TestLoadPath_Error(t *testing.T) {
	_, err := load.LoadPath("../data/null")
	require.Error(t, err)
}

func TestLoadURI_Error(t *testing.T) {
	path, err := url.ParseRequestURI("http://null")
	require.NoError(t, err)

	_, err = load.LoadURI(path)
	require.Error(t, err)
}

func TestLoad_NoError(t *testing.T) {
	_, err := load.Load("../data/openapi-test1.yaml")
	require.NoError(t, err)
}
