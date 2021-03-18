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

func TestPatch_ParameterSchemaDescription(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	param := s2.Paths["/api/{domain}/{project}/badges/security-score"].Get.Parameters.GetByInAndName("query", "image")
	param.Schema.Value.Description = "reuven"

	d1, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)

	d1.SpecDiff.Patch(s1)

	d2, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	require.False(t, d2.Summary.Diff)
}
