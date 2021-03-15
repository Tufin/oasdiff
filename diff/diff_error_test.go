package diff_test

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
)

func TestDiff_SchemaNil(t *testing.T) {
	loader := openapi3.NewSwaggerLoader()
	s1, err := loader.LoadSwaggerFromFile("../data/home-iot-api-1.yaml")
	require.NoError(t, err)

	s1.Paths["/devices"].Post.RequestBody.Value.Content["application/json"].Schema.Value = nil
	_, err = diff.Get(diff.NewConfig(), s1, s1)
	require.Equal(t, "schema reference is nil", err.Error())
}

func TestDiff_MediaTypeNil(t *testing.T) {
	loader := openapi3.NewSwaggerLoader()
	s1, err := loader.LoadSwaggerFromFile("../data/home-iot-api-1.yaml")
	require.NoError(t, err)

	s1.Paths["/devices"].Post.RequestBody.Value.Content["application/json"] = nil
	_, err = diff.Get(diff.NewConfig(), s1, s1)
	require.Equal(t, "media type is nil", err.Error())
}

func TestDiff_PathItemNil(t *testing.T) {
	loader := openapi3.NewSwaggerLoader()
	s1, err := loader.LoadSwaggerFromFile("../data/home-iot-api-1.yaml")
	require.NoError(t, err)

	s1.Paths["/devices"] = nil
	_, err = diff.Get(diff.NewConfig(), s1, s1)
	require.Equal(t, "path item is nil", err.Error())
}

func TestDiff_SpecNil(t *testing.T) {
	loader := openapi3.NewSwaggerLoader()
	s1, err := loader.LoadSwaggerFromFile("../data/home-iot-api-1.yaml")
	_, err = diff.Get(diff.NewConfig(), nil, s1)
	require.Equal(t, "spec is nil", err.Error())
}

func TestSchemaDiff_MediaInvalidMultiEntries(t *testing.T) {

	s5 := l(t, 5)
	s5.Paths[securityScorePath].Get.Parameters.GetByInAndName(openapi3.ParameterInCookie, "test").Content["second/invalid"] = openapi3.NewMediaType()

	s1 := l(t, 1)

	_, err := diff.Get(diff.NewConfig(), s5, s1)
	require.Equal(t, "content map has more than one value", err.Error())
}

func TestFilterByRegex_Invalid(t *testing.T) {
	_, err := diff.Get(&diff.Config{Filter: "["}, l(t, 1), l(t, 2))
	require.Equal(t, "failed to compile filter regex '[' with error parsing regexp: missing closing ]: `[`", err.Error())
}
