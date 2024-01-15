package delta_test

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/delta"
	"github.com/tufin/oasdiff/diff"
)

func TestCalc_EndpointAdded(t *testing.T) {
	loader := openapi3.NewLoader()
	s1, err := loader.LoadFromFile("../data/simple1.yaml")
	require.NoError(t, err)

	s2, err := loader.LoadFromFile("../data/simple3.yaml")
	require.NoError(t, err)

	d, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	require.Equal(t, delta.Get(false, d, s1, s2), 0.25)
}

func TestCalc_EndpointAddedAndDeleted(t *testing.T) {
	loader := openapi3.NewLoader()
	s1, err := loader.LoadFromFile("../data/simple1.yaml")
	require.NoError(t, err)

	s2, err := loader.LoadFromFile("../data/simple2.yaml")
	require.NoError(t, err)

	d, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	require.Equal(t, delta.Get(false, d, s1, s2), 1.0)
}
