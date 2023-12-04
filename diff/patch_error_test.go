package diff_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
)

func TestPatch_StringTypeMismatch_Nil(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Paths.Value("/api/{domain}/{project}/install-command").Get.Parameters.GetByInAndName("path", "domain").Schema.Value.Description = "reuven"

	d1, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)

	d1.PathsDiff.Modified["/api/{domain}/{project}/install-command"].OperationsDiff.Modified["GET"].ParametersDiff.Modified["path"]["domain"].SchemaDiff.DescriptionDiff.To = nil

	require.EqualError(t, d1.Patch(s1), "diff value is nil instead of string")
}

func TestPatch_StringTypeMismatch_Int(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Paths.Value("/api/{domain}/{project}/install-command").Get.Parameters.GetByInAndName("path", "domain").Schema.Value.Description = "reuven"

	d1, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)

	d1.PathsDiff.Modified["/api/{domain}/{project}/install-command"].OperationsDiff.Modified["GET"].ParametersDiff.Modified["path"]["domain"].SchemaDiff.DescriptionDiff.To = 4

	require.EqualError(t, d1.Patch(s1), "diff value type mismatch: string vs. \"int\"")
}

func TestPatch_UINT64TypeMismatch(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	maxLength := uint64(13)
	s2.Paths.Value("/api/{domain}/{project}/install-command").Get.Parameters.GetByInAndName("path", "domain").Schema.Value.MaxLength = &maxLength

	d1, err := diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)

	d1.PathsDiff.Modified["/api/{domain}/{project}/install-command"].OperationsDiff.Modified["GET"].ParametersDiff.Modified["path"]["domain"].SchemaDiff.MaxLengthDiff.To = 13

	require.EqualError(t, d1.Patch(s1), "diff value type mismatch: uint64 vs. \"int\"")
}
