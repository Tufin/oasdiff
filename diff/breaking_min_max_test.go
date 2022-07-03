package diff_test

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
)

// TODO: check whether this logic applies to response too

func TestBreaking_MaxLengthSmaller(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	maxLengthFrom := uint64(13)
	s1.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MaxLength = &maxLengthFrom

	maxLengthTo := uint64(11)
	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MaxLength = &maxLengthTo

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	// BC: reducing max length is breaking
	require.NotEmpty(t, d)
}

func TestBreaking_MinLengthSmaller(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MinLength = uint64(13)
	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MinLength = uint64(11)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	// BC: reducing min length isn't breaking
	require.Empty(t, d)
}

func TestBreaking_MaxLengthGreater(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	maxLengthFrom := uint64(13)
	s1.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MaxLength = &maxLengthFrom

	maxLengthTo := uint64(14)
	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MaxLength = &maxLengthTo

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	// BC: increasing max length isn't breaking
	require.Empty(t, d)
}

func TestBreaking_MaxLengthFromNil(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MaxLength = nil

	maxLengthTo := uint64(14)
	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MaxLength = &maxLengthTo

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	// BC: changing max length from nil to any value is breaking
	require.NotEmpty(t, d)
}

func TestBreaking_MaxLengthToNil(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	maxLengthFrom := uint64(13)
	s1.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MaxLength = &maxLengthFrom

	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MaxLength = nil

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	// BC: changing max length from any value to nil isn't breaking
	require.Empty(t, d)
}

func TestBreaking_MaxLengthBothNil(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MaxLength = nil
	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MaxLength = nil

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	// BC: both max lengths are nil isn't breaking
	require.Empty(t, d)
}

func TestBreaking_MinItemsSmaller(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MinItems = 13
	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MinItems = 11

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	// BC: reducing min items isn't breaking
	require.Empty(t, d)
}

func TestBreaking_MinItemsGreater(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MinItems = 13
	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MinItems = 14

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	// BC: increasing min items is breaking
	require.NotEmpty(t, d)
}

func TestBreaking_MaxSmaller(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	maxFrom := float64(13)
	s1.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.Max = &maxFrom

	maxTo := float64(11)
	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.Max = &maxTo

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	// BC: reducing max is breaking
	require.NotEmpty(t, d)
}
