package diff_test

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
)

// BC: no change isn't breaking
func TestBreaking_Same(t *testing.T) {
	require.Empty(t, d(t, &diff.Config{BreakingOnly: true}, 1, 1))
}

// BC: adding an optional request body isn't breaking
func TestBreaking_AddingOptionalRequestBody(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Paths[installCommandPath].Get.RequestBody = &openapi3.RequestBodyRef{
		Value: openapi3.NewRequestBody().WithRequired(false),
	}

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.Empty(t, d)
}

// BC: changing an existing request body from required to optional isn't breaking
func TestBreaking_RequestBodyRequiredDisabled(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Paths[installCommandPath].Get.RequestBody = &openapi3.RequestBodyRef{
		Value: openapi3.NewRequestBody().WithRequired(true),
	}

	s2.Paths[installCommandPath].Get.RequestBody = &openapi3.RequestBodyRef{
		Value: openapi3.NewRequestBody().WithRequired(false),
	}

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.Empty(t, d)
}

// BC: deleting a tag isn't breaking
func TestBreaking_DeletedTag(t *testing.T) {
	require.Empty(t, d(t, &diff.Config{
		BreakingOnly: true,
	}, 1, 5).PathsDiff.Modified[securityScorePath].OperationsDiff.Modified["GET"].TagsDiff)
}

// BC: adding an enum value isn't breaking
func TestBreaking_AddedEnum(t *testing.T) {
	require.Empty(t,
		d(t, &diff.Config{
			BreakingOnly: true,
		}, 1, 3).PathsDiff.Modified[installCommandPath].OperationsDiff.Modified["GET"].ParametersDiff.Modified[openapi3.ParameterInPath])
}

// BC: changing extensions isn't breaking
func TestBreaking_ModifiedExtension(t *testing.T) {
	config := diff.Config{
		BreakingOnly:      true,
		IncludeExtensions: diff.StringSet{"x-extension-test2": struct{}{}},
	}

	require.Empty(t, d(t, &config, 1, 3).ExtensionsDiff)
}

// BC: changing comments isn't breaking
func TestBreaking_Components(t *testing.T) {
	require.Empty(t, d(t, &diff.Config{BreakingOnly: true},
		1, 3).ComponentsDiff)
}

// BC: new optional header param isn't breaking
func TestBreaking_NewOptionalHeaderParam(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	deleteParam(s1.Paths[installCommandPath].Get, openapi3.ParameterInHeader, "network-policies")
	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = false

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.Empty(t, d)
}

// BC: changing an existing header param to optional isn't breaking
func TestBreaking_HeaderParamRequiredDisabled(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s1.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = true
	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Required = false

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.Empty(t, d)
}

func deleteResponseHeader(response *openapi3.Response, name string) {
	delete(response.Headers, name)
}

// BC: new required response header param isn't breaking
func TestBreaking_NewRequiredResponseHeader(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	deleteResponseHeader(s1.Paths[installCommandPath].Get.Responses["default"].Value, "X-RateLimit-Limit")
	s2.Paths[installCommandPath].Get.Responses["default"].Value.Headers["X-RateLimit-Limit"].Value.Required = true

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.Empty(t, d)
}

// BC: changing operation ID isn't breaking
func TestBreaking_OperationID(t *testing.T) {
	require.Empty(t,
		d(t, &diff.Config{
			BreakingOnly: true,
		}, 3, 1).PathsDiff.Modified[securityScorePath].OperationsDiff.Modified["GET"].OperationIDDiff)
}

// BC: changing a link to operation ID isn't breaking
func TestBreaking_LinkOperationID(t *testing.T) {
	require.Empty(t,
		d(t, &diff.Config{
			BreakingOnly: true,
		}, 3, 1).PathsDiff.Modified["/subscribe"].OperationsDiff.Modified["POST"].CallbacksDiff.Modified["myEvent"].Modified["hi"].OperationsDiff.Modified["POST"].ResponsesDiff.Modified["200"].LinksDiff.Modified)
}

// BC: adding a media-type to response isn't breaking
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
	require.Empty(t, d)
}

// BC: deprecating an operation isn't breaking
func TestBreaking_DeprecatedOperation(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Paths[installCommandPath].Get.Deprecated = true

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.Empty(t, d)
}

// BC: deprecating a parameter isn't breaking
func TestBreaking_DeprecatedParameter(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Deprecated = true

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.Empty(t, d)
}

// BC: deprecating a header isn't breaking
func TestBreaking_DeprecatedHeader(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Paths[installCommandPath].Get.Responses["default"].Value.Headers["X-RateLimit-Limit"].Value.Deprecated = true

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.Empty(t, d)
}

// BC: deprecating a schema isn't breaking
func TestBreaking_DeprecatedSchema(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Deprecated = true

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.Empty(t, d)
}

// BC: changing servers isn't breaking
func TestBreaking_Servers(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile("../data/servers/baseswagger.json")
	require.NoError(t, err)

	s2, err := loader.LoadFromFile("../data/servers/revisionswagger.json")
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.Empty(t, d)
}
