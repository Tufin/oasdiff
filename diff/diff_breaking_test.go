package diff_test

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
)

func TestBreaking_Same(t *testing.T) {
	require.Empty(t, d(t, &diff.Config{BreakingOnly: true}, 1, 1))
}

func TestBreaking_DeletedPaths(t *testing.T) {
	require.NotEmpty(t, d(t, &diff.Config{BreakingOnly: true}, 1, 2))
}

func TestBreaking_DeletedTagAllChanges(t *testing.T) {
	require.NotEmpty(t, d(t, &diff.Config{
		BreakingOnly: false,
	}, 1, 5).PathsDiff.Modified[securityScorePath].OperationsDiff.Modified["GET"].TagsDiff)
}

func TestBreaking_DeletedTag(t *testing.T) {
	require.Empty(t, d(t, &diff.Config{
		BreakingOnly: true,
	}, 1, 5).PathsDiff.Modified[securityScorePath].OperationsDiff.Modified["GET"].TagsDiff)
}

func TestBreaking_DeletedEnum(t *testing.T) {
	require.NotEmpty(t,
		d(t, &diff.Config{
			BreakingOnly: true,
		}, 3, 1).PathsDiff.Modified[installCommandPath].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInPath]["project"].SchemaDiff.EnumDiff)
}

func TestBreaking_AddedEnum(t *testing.T) {
	require.Empty(t,
		d(t, &diff.Config{
			BreakingOnly: true,
		}, 1, 3).PathsDiff.Modified[installCommandPath].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInPath])
}

func TestBreaking_ModifiedExtension(t *testing.T) {
	config := diff.Config{
		BreakingOnly:      true,
		IncludeExtensions: diff.StringSet{"x-extension-test2": struct{}{}},
	}

	require.Empty(t, d(t, &config, 1, 3).ExtensionsDiff)
}

func TestBreaking_Components(t *testing.T) {
	require.Empty(t, d(t, &diff.Config{BreakingOnly: true},
		1, 3).ComponentsDiff)
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

func deleteParam(op *openapi3.Operation, in string, name string) {

	result := openapi3.NewParameters()

	for _, item := range op.Parameters {
		if v := item.Value; v != nil {
			if v.Name == name && v.In == in {
				continue
			}
			result = append(result, item)
		}
	}
	op.Parameters = result
}

func TestBreaking_NewPathParam(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	deleteParam(s1.Paths[installCommandPath].Get, openapi3.ParameterInPath, "project")
	// note: path params are always required

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	// new required path param breaks client
	require.Contains(t,
		d.PathsDiff.Modified[installCommandPath].OperationsDiff.Modified["GET"].ParametersDiff.Added[openapi3.ParameterInPath],
		"project")
}

func TestBreaking_NewRequiredHeaderParam(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	deleteParam(s1.Paths[installCommandPath].Get, openapi3.ParameterInHeader, "network-policies")
	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = true

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	// new required header param breaks client
	require.Contains(t,
		d.PathsDiff.Modified[installCommandPath].OperationsDiff.Modified["GET"].ParametersDiff.Added[openapi3.ParameterInHeader],
		"network-policies")
}

func TestBreaking_NewNonRequiredHeaderParam(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	deleteParam(s1.Paths[installCommandPath].Get, openapi3.ParameterInHeader, "network-policies")
	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = false

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	// new optional header param doesn't break client
	require.Empty(t, d)
}

func TestBreaking_HeaderParamRequiredEnabled(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = false
	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = true

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	// changing an existing header param to required breaks client
	require.Equal(t,
		&diff.ValueDiff{
			From: false,
			To:   true,
		},
		d.PathsDiff.Modified[installCommandPath].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInHeader]["network-policies"].RequiredDiff)
}

func TestBreaking_HeaderParamRequiredDisabled(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = true
	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = false

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	// changing an existing header param to optional doesn't break client
	require.Empty(t, d)
}

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

func TestBreaking_OperationID(t *testing.T) {
	require.Empty(t,
		d(t, &diff.Config{
			BreakingOnly: true,
		}, 3, 1).PathsDiff.Modified[securityScorePath].OperationsDiff.Modified["GET"].OperationIDDiff)
}

func TestBreaking_LinkOperationID(t *testing.T) {
	require.Empty(t,
		d(t, &diff.Config{
			BreakingOnly: true,
		}, 3, 1).PathsDiff.Modified["/subscribe"].OperationsDiff.Modified["POST"].CallbacksDiff.Modified["myEvent"].Modified["hi"].OperationsDiff.Modified["POST"].ResponsesDiff.Modified["200"].LinksDiff.Modified)
}

func TestBreaking_ResponseDeleteMediaType(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile("../data/response-media-type-base.yaml")
	require.NoError(t, err)

	s2, err := loader.LoadFromFile("../data/response-media-type-revision.yaml")
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	// deleting a media-type from response is breaking
	require.NotEmpty(t, d)
}

func TestBreaking_ResponseAddMediaType(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile("../data/response-media-type-revision.yaml")
	require.NoError(t, err)

	s2, err := loader.LoadFromFile("../data/response-media-type-base.yaml")
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	// adding a media-type to response is not breaking
	require.Empty(t, d)
}
