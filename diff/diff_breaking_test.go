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

func TestBreaking_Ref(t *testing.T) {
	require.True(t,
		d(t, diff.NewConfig(), 1, 3).RequestBodiesDiff.Modified["reuven"].ContentDiff.MediaTypeModified["application/json"].SchemaDiff.PropertiesDiff.Modified["meter_value"].TypeDiff.Breaking(),
	)
}

func TestCompareWithDefault(t *testing.T) {
	require.True(t,
		d(t, diff.NewConfig(), 1, 3).TagsDiff.Modified["reuven"].DescriptionDiff.CompareWithDefault("Harrison", "harrison", ""),
	)
}

func TestCompareWithDefault_Nil(t *testing.T) {
	require.True(t,
		d(t, diff.NewConfig(), 2, 1).PathsDiff.Modified[securityScorePathSlash].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInQuery]["image"].ExplodeDiff.CompareWithDefault(true, false, false),
	)
}

func TestBreaking_MaxLengthSmaller(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	maxLengthFrom := uint64(13)
	s1.Paths["/api/{domain}/{project}/install-command"].Get.Parameters.GetByInAndName("path", "domain").Schema.Value.MaxLength = &maxLengthFrom

	maxLengthTo := uint64(11)
	s2.Paths["/api/{domain}/{project}/install-command"].Get.Parameters.GetByInAndName("path", "domain").Schema.Value.MaxLength = &maxLengthTo

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.True(t, d.Breaking())
}

func TestBreaking_MaxLengthGreater(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	maxLengthFrom := uint64(13)
	s1.Paths["/api/{domain}/{project}/install-command"].Get.Parameters.GetByInAndName("path", "domain").Schema.Value.MaxLength = &maxLengthFrom

	maxLengthTo := uint64(14)
	s2.Paths["/api/{domain}/{project}/install-command"].Get.Parameters.GetByInAndName("path", "domain").Schema.Value.MaxLength = &maxLengthTo

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.False(t, d.Breaking())
}

func TestBreaking_MaxLengthFromNil(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Paths["/api/{domain}/{project}/install-command"].Get.Parameters.GetByInAndName("path", "domain").Schema.Value.MaxLength = nil

	maxLengthTo := uint64(14)
	s2.Paths["/api/{domain}/{project}/install-command"].Get.Parameters.GetByInAndName("path", "domain").Schema.Value.MaxLength = &maxLengthTo

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.True(t, d.Breaking())
}

func TestBreaking_MaxLengthToNil(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	maxLengthFrom := uint64(13)
	s1.Paths["/api/{domain}/{project}/install-command"].Get.Parameters.GetByInAndName("path", "domain").Schema.Value.MaxLength = &maxLengthFrom

	s2.Paths["/api/{domain}/{project}/install-command"].Get.Parameters.GetByInAndName("path", "domain").Schema.Value.MaxLength = nil

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.False(t, d.Breaking())
}

func TestBreaking_MaxLengthBothNil(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Paths["/api/{domain}/{project}/install-command"].Get.Parameters.GetByInAndName("path", "domain").Schema.Value.MaxLength = nil
	s2.Paths["/api/{domain}/{project}/install-command"].Get.Parameters.GetByInAndName("path", "domain").Schema.Value.MaxLength = nil

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.False(t, d.Breaking())
}

func TestBreaking_MinItemsSmaller(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Paths["/api/{domain}/{project}/install-command"].Get.Parameters.GetByInAndName("path", "domain").Schema.Value.MinItems = 13
	s2.Paths["/api/{domain}/{project}/install-command"].Get.Parameters.GetByInAndName("path", "domain").Schema.Value.MinItems = 11

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.False(t, d.Breaking())
}

func TestBreaking_MinItemsGreater(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Paths["/api/{domain}/{project}/install-command"].Get.Parameters.GetByInAndName("path", "domain").Schema.Value.MinItems = 13
	s2.Paths["/api/{domain}/{project}/install-command"].Get.Parameters.GetByInAndName("path", "domain").Schema.Value.MinItems = 14

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.True(t, d.Breaking())
}
