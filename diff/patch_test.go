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

	d1.SpecDiff.Patch(s1)

	d2, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	require.False(t, d2.Summary.Diff)
}

func TestPatch_ParameterDescription(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Paths["/api/{domain}/{project}/badges/security-score"].Get.Parameters.GetByInAndName("query", "filter").Description = "reuven"

	d1, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)

	d1.SpecDiff.Patch(s1)

	d2, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	require.False(t, d2.Summary.Diff)
}

func TestPatch_ParameterSchemaFormat(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	schema := s2.Paths["/api/{domain}/{project}/badges/security-score"].Get.Parameters.GetByInAndName("query", "image").Schema.Value
	schema.Format = "reuven"

	d1, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)

	d1.SpecDiff.Patch(s1)

	d2, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	require.False(t, d2.Summary.Diff)
}

func TestPatch_ParameterSchemaEnum(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	schema := s2.Paths["/api/{domain}/{project}/install-command"].Get.Parameters.GetByInAndName("path", "domain").Schema.Value
	schema.Enum = []interface{}{"reuven", "tufin"}

	d1, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)

	d1.SpecDiff.Patch(s1)

	d2, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	require.False(t, d2.Summary.Diff)
}
