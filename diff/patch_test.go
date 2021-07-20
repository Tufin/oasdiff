package diff_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
)

func TestPatch_MethodDescription(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Paths["/api/{domain}/{project}/badges/security-score"].Get.Description = "reuven"

	d1, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)

	d1.Patch(s1)

	d2, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	require.False(t, d2.GetSummary().Diff)
}

func TestPatch_ParameterDescription(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Paths["/api/{domain}/{project}/badges/security-score"].Get.Parameters.GetByInAndName("query", "filter").Description = "reuven"

	d1, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)

	d1.Patch(s1)

	d2, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	require.False(t, d2.GetSummary().Diff)
}

func TestPatch_ParameterSchemaFormat(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	schema := s2.Paths["/api/{domain}/{project}/badges/security-score"].Get.Parameters.GetByInAndName("query", "image").Schema.Value
	schema.Format = "reuven"

	d1, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)

	d1.Patch(s1)

	d2, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	require.False(t, d2.GetSummary().Diff)
}

func TestPatch_ParameterSchemaEnum(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Paths["/api/{domain}/{project}/install-command"].Get.Parameters.GetByInAndName("path", "domain").Schema.Value.Enum = []interface{}{"reuven", "tufin"}

	d1, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)

	d1.Patch(s1)

	d2, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	require.False(t, d2.GetSummary().Diff)
}

func TestPatch_ParameterSchemaMaxLengthNil(t *testing.T) {
	s1 := l(t, 1)
	maxLength := uint64(13)
	s1.Paths["/api/{domain}/{project}/install-command"].Get.Parameters.GetByInAndName("path", "domain").Schema.Value.MaxLength = &maxLength

	s2 := l(t, 1)
	s2.Paths["/api/{domain}/{project}/install-command"].Get.Parameters.GetByInAndName("path", "domain").Schema.Value.MaxLength = nil

	d1, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)

	d1.Patch(s1)

	d2, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	require.False(t, d2.GetSummary().Diff)
}

func TestPatch_ParameterSchemaMaxLength(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	maxLength := uint64(13)
	s2.Paths["/api/{domain}/{project}/install-command"].Get.Parameters.GetByInAndName("path", "domain").Schema.Value.MaxLength = &maxLength

	d1, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)

	d1.Patch(s1)

	d2, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	require.False(t, d2.GetSummary().Diff)
}

func TestPatch_ValueDiffNil(t *testing.T) {
	valueDiff := &diff.ValueDiff{}
	value := "reuven"
	require.EqualError(t, valueDiff.PatchString(&value), "diff value is nil instead of string")
}

func TestPatch_ValueDiffMismatch(t *testing.T) {
	valueDiff := &diff.ValueDiff{
		To: 4,
	}
	value := "reuven"
	require.EqualError(t, valueDiff.PatchString(&value), "diff value type mismatch: string vs. \"int\"")
}

func TestPatch_ValueDiffInt(t *testing.T) {
	valueDiff := &diff.ValueDiff{
		To: 4,
	}
	value := uint64(3)
	pValue := &value
	require.EqualError(t, valueDiff.PatchUInt64Ref(&pValue), "diff value type mismatch: uint64 vs. \"int\"")
}

func TestPatch_ValueDiff(t *testing.T) {
	v1 := uint64(3)

	valueDiff := &diff.ValueDiff{
		To: v1,
	}

	v2 := uint64(3)
	pV2 := &v2
	require.NoError(t, valueDiff.PatchUInt64Ref(&pV2))
}
