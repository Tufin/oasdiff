package diff_test

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
)

// BC: reducing max length is breaking
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
	require.NotEmpty(t, d)
}

// BC: reducing min length isn't breaking
func TestBreaking_MinLengthSmaller(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MinLength = uint64(13)
	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MinLength = uint64(11)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.Empty(t, d)
}

// BC: increasing max length isn't breaking
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
	require.Empty(t, d)
}

// BC: changing max length from nil to any value is breaking
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
	require.NotEmpty(t, d)
}

// BC: changing max length from any value to nil isn't breaking
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
	require.Empty(t, d)
}

// BC: both max lengths are nil isn't breaking
func TestBreaking_MaxLengthBothNil(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MaxLength = nil
	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MaxLength = nil

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.Empty(t, d)
}

// BC: reducing min items isn't breaking
func TestBreaking_MinItemsSmaller(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MinItems = 13
	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MinItems = 11

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.Empty(t, d)
}

// BC: increasing min items is breaking
func TestBreaking_MinItemsGreater(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MinItems = 13
	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInPath, "domain").Schema.Value.MinItems = 14

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.NotEmpty(t, d)
}

// BC: reducing max in request is breaking
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
	require.NotEmpty(t, d)
}

// BC: reducing max in response is breaking
func TestBreaking_MaxSmallerInReponse(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	maxLengthFrom := uint64(13)
	s1.Paths[installCommandPath].Get.Responses["default"].Value.Headers["X-RateLimit-Limit"].Value.Schema.Value.MaxLength = &maxLengthFrom

	maxLengthTo := uint64(11)
	s2.Paths[installCommandPath].Get.Responses["default"].Value.Headers["X-RateLimit-Limit"].Value.Schema.Value.MaxLength = &maxLengthTo

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.NotEmpty(t, d)
}
