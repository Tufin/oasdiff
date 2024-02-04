package load_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/load"
)

func TestSpecInfoPair(t *testing.T) {
	spec, err := load.NewSpecInfo(MockLoader{}, load.NewSource("../data/openapi-test1.yaml"))
	require.NoError(t, err)

	pair := load.NewSpecInfoPair(spec, spec)
	require.Equal(t, "1.0.0", pair.GetBaseVersion())
	require.Equal(t, "1.0.0", pair.GetRevisionVersion())
}

func TestSpecInfoPair_NA(t *testing.T) {
	var pair *load.SpecInfoPair
	require.Equal(t, "n/a", pair.GetBaseVersion())
	require.Equal(t, "n/a", pair.GetRevisionVersion())
}

func TestSpecInfoPair_Nil(t *testing.T) {
	var spec *load.SpecInfo
	pair := load.NewSpecInfoPair(spec, spec)

	require.Equal(t, "n/a", pair.GetBaseVersion())
	require.Equal(t, "n/a", pair.GetRevisionVersion())
}
