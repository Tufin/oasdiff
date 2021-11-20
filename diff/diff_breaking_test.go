package diff_test

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
)

func TestBreaking_Same(t *testing.T) {
	require.False(t, d(t, diff.NewConfig(), 1, 1).Breaking())
}

func TestBreaking_DeletedPaths(t *testing.T) {
	require.True(t, d(t, diff.NewConfig(), 1, 2).Breaking())
}

func TestBreaking_DeletedTagAllChanges(t *testing.T) {
	require.False(t, d(t, &diff.Config{
		BreakingOnly: false,
	}, 1, 5).PathsDiff.Modified[securityScorePath].OperationsDiff.Modified["GET"].TagsDiff.Empty())
}

func TestBreaking_DeletedTag(t *testing.T) {
	require.True(t, d(t, &diff.Config{
		BreakingOnly: true,
	}, 1, 5).PathsDiff.Modified[securityScorePath].OperationsDiff.Modified["GET"].TagsDiff.Empty())
}

func TestBreaking_DeletedEnum(t *testing.T) {
	require.False(t,
		d(t, &diff.Config{
			BreakingOnly: true,
		}, 3, 1).PathsDiff.Modified[installCommandPath].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInPath]["project"].SchemaDiff.EnumDiff.Empty())
}

func TestBreaking_AddedEnum(t *testing.T) {
	require.False(t,
		d(t, &diff.Config{
			BreakingOnly: true,
		}, 1, 3).PathsDiff.Modified[installCommandPath].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInPath].Breaking())
}

func TestBreaking_ModifiedExtension(t *testing.T) {
	config := diff.Config{
		IncludeExtensions: diff.StringSet{"x-extension-test2": struct{}{}},
	}

	require.False(t, d(t, &config, 1, 3).ExtensionsDiff.Breaking())
}
