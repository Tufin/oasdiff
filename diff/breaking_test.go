package diff_test

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
)

// BC: deleting a path is breaking
func TestBreaking_DeletedPath(t *testing.T) {
	d := d(t, &diff.Config{BreakingOnly: true}, 1, 2)
	require.NotEmpty(t, d.PathsDiff.Deleted)
}

// BC: deleting an operation is breaking
func TestBreaking_DeletedOp(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Paths[installCommandPath].Put = openapi3.NewOperation()

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.NotEmpty(t, d.PathsDiff.Modified[installCommandPath].OperationsDiff.Deleted)
}

// BC: adding a required request body is breaking
func TestBreaking_AddingRequiredRequestBody(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Paths[installCommandPath].Get.RequestBody = &openapi3.RequestBodyRef{
		Value: openapi3.NewRequestBody().WithRequired(true),
	}

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.NotEmpty(t, d.PathsDiff.Modified[installCommandPath].OperationsDiff.Modified["GET"].RequestBodyDiff)
}

// BC: changing an existing request body from optional to required is breaking
func TestBreaking_RequestBodyRequiredEnabled(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Paths[installCommandPath].Get.RequestBody = &openapi3.RequestBodyRef{
		Value: openapi3.NewRequestBody().WithRequired(false),
	}

	s2.Paths[installCommandPath].Get.RequestBody = &openapi3.RequestBodyRef{
		Value: openapi3.NewRequestBody().WithRequired(true),
	}

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.NotEmpty(t, d.PathsDiff.Modified[installCommandPath].OperationsDiff.Modified["GET"].RequestBodyDiff)
}

// BC: deleting an enum value is breaking
func TestBreaking_DeletedEnum(t *testing.T) {
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

// BC: new required path param is breaking
func TestBreaking_NewPathParam(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	deleteParam(s1.Paths[installCommandPath].Get, openapi3.ParameterInPath, "project")
	// note: path params are always required

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	require.Contains(t,
		d.PathsDiff.Modified[installCommandPath].OperationsDiff.Modified["GET"].ParametersDiff.Added[openapi3.ParameterInPath],
		"project")
}

// BC: new required header param is breaking
func TestBreaking_NewRequiredHeaderParam(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	deleteParam(s1.Paths[installCommandPath].Get, openapi3.ParameterInHeader, "network-policies")
	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = true

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	require.Contains(t,
		d.PathsDiff.Modified[installCommandPath].OperationsDiff.Modified["GET"].ParametersDiff.Added[openapi3.ParameterInHeader],
		"network-policies")
}

// BC: changing an existing header param from optional to required is breaking
func TestBreaking_HeaderParamRequiredEnabled(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = false
	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = true

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	require.Equal(t,
		&diff.ValueDiff{
			From: false,
			To:   true,
		},
		d.PathsDiff.Modified[installCommandPath].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInHeader]["network-policies"].RequiredDiff)
}

// BC: changing an existing response header from required to optional is breaking
func TestBreaking_ResponseHeaderParamRequiredDisabled(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Paths[installCommandPath].Get.Responses["default"].Value.Headers["X-RateLimit-Limit"].Value.Required = true
	s2.Paths[installCommandPath].Get.Responses["default"].Value.Headers["X-RateLimit-Limit"].Value.Required = false

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.NotEmpty(t, d.PathsDiff.Modified[installCommandPath].OperationsDiff.Modified["GET"].ResponsesDiff.Modified["default"].HeadersDiff.Modified["X-RateLimit-Limit"].RequiredDiff)
}

// BC: deleting a media-type from response is breaking
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
	require.NotEmpty(t, d)
}

// BC: deleting a pattern from a schema is not breaking
func TestBreaking_DeletePatten(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile("../data/pattern-base.yaml")
	require.NoError(t, err)

	s2, err := loader.LoadFromFile("../data/pattern-revision.yaml")
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.Empty(t, d)
}

// BC: adding a pattern to a schema is breaking
func TestBreaking_AddPatten(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile("../data/pattern-revision.yaml")
	require.NoError(t, err)

	s2, err := loader.LoadFromFile("../data/pattern-base.yaml")
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.NotEmpty(t, d)
}

// BC: modifying a pattern in a schema is breaking
func TestBreaking_ModifyPatten(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile("../data/pattern-base.yaml")
	require.NoError(t, err)

	s2, err := loader.LoadFromFile("../data/pattern-modified.yaml")
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.NotEmpty(t, d)
}

// BC: removing an existing required response header is breaking
func TestBreaking_ResponseHeaderRemoved(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Paths[installCommandPath].Get.Responses["default"].Value.Headers["X-RateLimit-Limit"].Value.Required = true
	delete(s2.Paths[installCommandPath].Get.Responses["default"].Value.Headers, "X-RateLimit-Limit")

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.NotEmpty(t, d)
}

// BC: removing an existing optional response header is breaking
func TestBreaking_OptionalResponseHeaderRemoved(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Paths[installCommandPath].Get.Responses["default"].Value.Headers["X-RateLimit-Limit"].Value.Required = false
	delete(s2.Paths[installCommandPath].Get.Responses["default"].Value.Headers, "X-RateLimit-Limit")

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.NotEmpty(t, d)
}