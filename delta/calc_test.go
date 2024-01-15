package delta_test

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/delta"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/utils"
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

func TestSymmetric(t *testing.T) {
	specs := utils.StringList{"../data/simple.yaml", "../data/simple1.yaml", "../data/simple2.yaml", "../data/simple3.yaml"}
	specPairs := specs.CartesianProduct(specs)

	loader := openapi3.NewLoader()
	for _, pair := range specPairs {
		s1, err := loader.LoadFromFile(pair.X)
		require.NoError(t, err)

		s2, err := loader.LoadFromFile(pair.Y)
		require.NoError(t, err)

		d1, err := diff.Get(diff.NewConfig(), s1, s2)
		require.NoError(t, err)

		d2, err := diff.Get(diff.NewConfig(), s2, s1)
		require.NoError(t, err)

		require.Equal(t, delta.Get(false, d1, s1, s2), delta.Get(false, d2, s2, s1), pair)
	}
}

func TestAsymmetric(t *testing.T) {
	specs := utils.StringList{"../data/simple.yaml", "../data/simple1.yaml", "../data/simple2.yaml", "../data/simple3.yaml"}
	specPairs := specs.CartesianProduct(specs)

	loader := openapi3.NewLoader()
	for _, pair := range specPairs {
		s1, err := loader.LoadFromFile(pair.X)
		require.NoError(t, err)

		s2, err := loader.LoadFromFile(pair.Y)
		require.NoError(t, err)

		d1, err := diff.Get(diff.NewConfig(), s1, s2)
		require.NoError(t, err)
		asymmetric1 := delta.Get(true, d1, s1, s2)

		d2, err := diff.Get(diff.NewConfig(), s2, s1)
		require.NoError(t, err)
		asymmetric2 := delta.Get(true, d2, s2, s1)

		symmetric := delta.Get(false, d2, s2, s1)

		require.Equal(t, utils.Average(asymmetric1, asymmetric2), symmetric, pair)
	}
}
