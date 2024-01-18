package delta_test

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/delta"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/utils"
)

func TestEndpointAdded(t *testing.T) {
	loader := openapi3.NewLoader()
	s1, err := loader.LoadFromFile("../data/simple1.yaml")
	require.NoError(t, err)

	s2, err := loader.LoadFromFile("../data/simple3.yaml")
	require.NoError(t, err)

	d, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	require.Equal(t, 0.5, delta.Get(false, d))
}

func TestEndpointDeletedAsym(t *testing.T) {
	loader := openapi3.NewLoader()
	s1, err := loader.LoadFromFile("../data/simple3.yaml")
	require.NoError(t, err)

	s2, err := loader.LoadFromFile("../data/simple1.yaml")
	require.NoError(t, err)

	d, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	require.Equal(t, 0.5, delta.Get(true, d))
}

func TestEndpointAddedAndDeleted(t *testing.T) {
	loader := openapi3.NewLoader()
	s1, err := loader.LoadFromFile("../data/simple1.yaml")
	require.NoError(t, err)

	s2, err := loader.LoadFromFile("../data/simple2.yaml")
	require.NoError(t, err)

	d, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	require.Equal(t, 1.0, delta.Get(false, d))
}

func TestSymmetric(t *testing.T) {
	specs := utils.StringList{"../data/simple.yaml", "../data/simple1.yaml", "../data/simple2.yaml", "../data/simple3.yaml", "../data/simple4.yaml"}
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

		require.Equal(t, delta.Get(false, d1), delta.Get(false, d2), pair)
	}
}

func TestAsymmetric(t *testing.T) {
	specs := utils.StringList{"../data/simple.yaml", "../data/simple1.yaml", "../data/simple2.yaml", "../data/simple3.yaml", "../data/simple4.yaml"}
	specPairs := specs.CartesianProduct(specs)

	loader := openapi3.NewLoader()
	for _, pair := range specPairs {
		s1, err := loader.LoadFromFile(pair.X)
		require.NoError(t, err)

		s2, err := loader.LoadFromFile(pair.Y)
		require.NoError(t, err)

		d1, err := diff.Get(diff.NewConfig(), s1, s2)
		require.NoError(t, err)
		asymmetric1 := delta.Get(true, d1)

		d2, err := diff.Get(diff.NewConfig(), s2, s1)
		require.NoError(t, err)
		asymmetric2 := delta.Get(true, d2)

		symmetric := delta.Get(false, d2)

		require.Equal(t, asymmetric1+asymmetric2, symmetric, pair)
	}
}

func TestParameters(t *testing.T) {
	loader := openapi3.NewLoader()
	s1, err := loader.LoadFromFile("../data/simple4.yaml")
	require.NoError(t, err)

	s2, err := loader.LoadFromFile("../data/simple3.yaml")
	require.NoError(t, err)

	d, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	require.Equal(t, 0.25, delta.Get(true, d))
}
