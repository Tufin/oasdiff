package diff_test

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
)

func TestDiff_NilPointers(t *testing.T) {
	loader := openapi3.NewSwaggerLoader()
	s1, err := loader.LoadSwaggerFromFile("../data/home-iot-api-1.yaml")
	require.NoError(t, err)
	s2, err := loader.LoadSwaggerFromFile("../data/home-iot-api-2.yaml")
	require.NoError(t, err)

	s1.Paths["/devices"].Post.RequestBody.Value.Content["application/json"].Schema.Value = nil
	_, err = diff.Get(diff.NewConfig(), s1, s2)
	require.Error(t, err)

	s1.Paths["/devices"].Post.RequestBody.Value.Content["application/json"].Schema = nil
	_, err = diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)

	s1.Paths["/devices"].Post.RequestBody.Value.Content["application/json"] = nil
	_, err = diff.Get(diff.NewConfig(), s1, s2)
	require.Error(t, err)

	s1.Paths["/devices"].Post.RequestBody.Value.Content = nil
	_, err = diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)

	s1.Paths["/devices"].Post.RequestBody = nil
	_, err = diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)

	s1.Paths["/devices"].Post = nil
	_, err = diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)

	s1.Paths["/devices"] = nil
	_, err = diff.Get(diff.NewConfig(), s1, s2)
	require.Error(t, err)

	s1.Paths = nil
	_, err = diff.Get(diff.NewConfig(), s1, s2)
	require.NoError(t, err)

	s1 = nil
	_, err = diff.Get(diff.NewConfig(), s1, s2)
	require.Error(t, err)
}

func TestSchemaDiff_MediaInvalidMultiEntries1(t *testing.T) {

	s5 := l(t, 5)
	s5.Paths[securityScorePath].Get.Parameters.GetByInAndName(openapi3.ParameterInCookie, "test").Content["second/invalid"] = openapi3.NewMediaType()

	s1 := l(t, 1)

	_, err := diff.Get(diff.NewConfig(), s1, s5)
	require.Error(t, err)
}

func TestSchemaDiff_MediaInvalidMultiEntries2(t *testing.T) {

	s5 := l(t, 5)
	s5.Paths[securityScorePath].Get.Parameters.GetByInAndName(openapi3.ParameterInCookie, "test").Content["second/invalid"] = openapi3.NewMediaType()

	s1 := l(t, 1)

	_, err := diff.Get(diff.NewConfig(), s5, s1)
	require.Error(t, err)
}

func TestFilterByRegex_Invalid(t *testing.T) {
	_, err := diff.Get(&diff.Config{Filter: "["}, l(t, 1), l(t, 2))
	require.Error(t, err)
}
