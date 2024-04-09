package diff_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
)

func TestPatch_MethodDescription(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Paths.Value("/api/{domain}/{project}/badges/security-score").Get.Description = "reuven"

	d1, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)

	require.NoError(t, d1.Patch(s1))

	d2, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	require.False(t, d2.GetSummary().Diff)
}

func TestPatch_ParameterDescription(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Paths.Value("/api/{domain}/{project}/badges/security-score").Get.Parameters.GetByInAndName("query", "filter").Description = "reuven"

	d1, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)

	require.NoError(t, d1.Patch(s1))

	d2, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	require.False(t, d2.GetSummary().Diff)
}

func TestPatch_ParameterSchemaFormat(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	schema := s2.Paths.Value("/api/{domain}/{project}/badges/security-score").Get.Parameters.GetByInAndName("query", "image").Schema.Value
	schema.Format = "reuven"

	d1, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)

	require.NoError(t, d1.Patch(s1))

	d2, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	require.False(t, d2.GetSummary().Diff)
}

func TestPatch_ParameterSchemaEnum(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Paths.Value("/api/{domain}/{project}/install-command").Get.Parameters.GetByInAndName("path", "domain").Schema.Value.Enum = []interface{}{"reuven", "tufin"}

	d1, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)

	require.NoError(t, d1.Patch(s1))

	d2, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	require.False(t, d2.GetSummary().Diff)
}

func TestPatch_ParameterSchemaMaxLengthNil(t *testing.T) {
	s1 := l(t, 1)
	maxLength := uint64(13)
	s1.Paths.Value("/api/{domain}/{project}/install-command").Get.Parameters.GetByInAndName("path", "domain").Schema.Value.MaxLength = &maxLength

	s2 := l(t, 1)
	s2.Paths.Value("/api/{domain}/{project}/install-command").Get.Parameters.GetByInAndName("path", "domain").Schema.Value.MaxLength = nil

	d1, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)

	require.NoError(t, d1.Patch(s1))

	d2, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	require.False(t, d2.GetSummary().Diff)
}

func TestPatch_ParameterSchemaMaxLength(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	maxLength := uint64(13)
	s2.Paths.Value("/api/{domain}/{project}/install-command").Get.Parameters.GetByInAndName("path", "domain").Schema.Value.MaxLength = &maxLength

	d1, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)

	require.NoError(t, d1.Patch(s1))

	d2, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	require.False(t, d2.GetSummary().Diff)
}
