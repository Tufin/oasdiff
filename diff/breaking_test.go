package diff_test

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
)

func TestBreaking_DeletedPaths(t *testing.T) {
	// deleting a path is breaking
	require.NotEmpty(t, d(t, &diff.Config{BreakingOnly: true}, 1, 2))
}

func TestBreaking_DeletedEnum(t *testing.T) {
	// deleting an enum value is breaking
	require.NotEmpty(t,
		d(t, &diff.Config{
			BreakingOnly: true,
		}, 3, 1).PathsDiff.Modified[installCommandPath].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInPath]["project"].SchemaDiff.EnumDiff)
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

func deleteResponseHeader(response *openapi3.Response, name string) {
	delete(response.Headers, name)
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

	// new required path param is breaking
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

	// new required header param is breaking
	require.Contains(t,
		d.PathsDiff.Modified[installCommandPath].OperationsDiff.Modified["GET"].ParametersDiff.Added[openapi3.ParameterInHeader],
		"network-policies")
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

	// changing an existing header param from optional to required is breaking
	require.Equal(t,
		&diff.ValueDiff{
			From: false,
			To:   true,
		},
		d.PathsDiff.Modified[installCommandPath].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInHeader]["network-policies"].RequiredDiff)
}

func TestBreaking_ResponseHeaderParamRequiredDisabled(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Paths[installCommandPath].Get.Responses["default"].Value.Headers["X-RateLimit-Limit"].Value.Required = true
	s2.Paths[installCommandPath].Get.Responses["default"].Value.Headers["X-RateLimit-Limit"].Value.Required = false

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	// changing an existing response header from required to optional is breaking
	require.NotEmpty(t, d.PathsDiff.Modified[installCommandPath].OperationsDiff.Modified["GET"].ResponsesDiff.Modified["default"].HeadersDiff.Modified["X-RateLimit-Limit"].RequiredDiff)
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
