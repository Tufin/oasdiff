package diff_test

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
)

func TestDiff_SchemaRefNil(t *testing.T) {
	loader := openapi3.NewSwaggerLoader()
	s1, err := loader.LoadSwaggerFromFile("../data/home-iot-api-1.yaml")
	require.NoError(t, err)

	s1.Components.Schemas["LightingSummary"].Value.Properties["zones"].Value.Items.Value.Properties["deviceId"].Value = nil
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

func TestDiff_InfoNil(t *testing.T) {
	s1 := &openapi3.Swagger{}
	_, err := diff.Get(diff.NewConfig(), s1, s1)
	require.Equal(t, "info is nil", err.Error())
}

func TestDiff_ExampleNil(t *testing.T) {
	s1 := openapi3.Swagger{
		Info: &openapi3.Info{},
		Components: openapi3.Components{
			Examples: openapi3.Examples{"test": &openapi3.ExampleRef{Value: &openapi3.Example{}}},
		},
	}
	s2 := openapi3.Swagger{
		Info: &openapi3.Info{},
		Components: openapi3.Components{
			Examples: openapi3.Examples{"test": &openapi3.ExampleRef{}},
		},
	}
	config := &diff.Config{IncludeExamples: true}
	_, err := diff.Get(config, &s1, &s2)
	require.Equal(t, "example reference is nil", err.Error())

	_, err = diff.Get(config, &s2, &s1)
	require.Equal(t, "example reference is nil", err.Error())
}

func TestDiff_ComponentSchemaNil(t *testing.T) {
	s1 := openapi3.Swagger{
		Info: &openapi3.Info{},
		Components: openapi3.Components{
			Schemas: openapi3.Schemas{"test": &openapi3.SchemaRef{Value: &openapi3.Schema{}}},
		},
	}
	s2 := openapi3.Swagger{
		Info: &openapi3.Info{},
		Components: openapi3.Components{
			Schemas: openapi3.Schemas{"test": &openapi3.SchemaRef{}},
		},
	}
	_, err := diff.Get(diff.NewConfig(), &s1, &s2)
	require.Equal(t, "schema reference is nil", err.Error())

	_, err = diff.Get(diff.NewConfig(), &s2, &s1)
	require.Equal(t, "schema reference is nil", err.Error())
}

func TestDiff_ComponentSchemaDeepNil(t *testing.T) {
	s1 := openapi3.Swagger{
		Info: &openapi3.Info{},
		Components: openapi3.Components{
			Schemas: openapi3.Schemas{
				"test": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						OneOf: []*openapi3.SchemaRef{
							{
								Value: &openapi3.Schema{
									AnyOf: []*openapi3.SchemaRef{
										{
											Value: &openapi3.Schema{
												AllOf: []*openapi3.SchemaRef{
													{
														Value: &openapi3.Schema{
															Not: &openapi3.SchemaRef{
																Value: nil,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := diff.Get(diff.NewConfig(), &s1, &s1)
	require.Equal(t, "schema reference is nil", err.Error())
}

func TestDiff_ComponentParameterNil(t *testing.T) {
	s1 := openapi3.Swagger{
		Info: &openapi3.Info{},
		Components: openapi3.Components{
			Parameters: openapi3.ParametersMap{"test": &openapi3.ParameterRef{Value: &openapi3.Parameter{}}},
		},
	}
	s2 := openapi3.Swagger{
		Info: &openapi3.Info{},
		Components: openapi3.Components{
			Parameters: openapi3.ParametersMap{"test": &openapi3.ParameterRef{}},
		},
	}
	_, err := diff.Get(diff.NewConfig(), &s1, &s2)
	require.Equal(t, "parameter reference is nil", err.Error())

	_, err = diff.Get(diff.NewConfig(), &s2, &s1)
	require.Equal(t, "parameter reference is nil", err.Error())
}

func TestDiff_ComponentHeadersNil(t *testing.T) {
	s1 := openapi3.Swagger{
		Info: &openapi3.Info{},
		Components: openapi3.Components{
			Headers: openapi3.Headers{"test": &openapi3.HeaderRef{Value: &openapi3.Header{}}},
		},
	}
	s2 := openapi3.Swagger{
		Info: &openapi3.Info{},
		Components: openapi3.Components{
			Headers: openapi3.Headers{"test": &openapi3.HeaderRef{}},
		},
	}
	_, err := diff.Get(diff.NewConfig(), &s1, &s2)
	require.Equal(t, "header reference is nil", err.Error())

	_, err = diff.Get(diff.NewConfig(), &s2, &s1)
	require.Equal(t, "header reference is nil", err.Error())
}

func TestDiff_ComponentRequestBodiesNil(t *testing.T) {
	s1 := openapi3.Swagger{
		Info: &openapi3.Info{},
		Components: openapi3.Components{
			RequestBodies: openapi3.RequestBodies{"test": &openapi3.RequestBodyRef{Value: &openapi3.RequestBody{}}},
		},
	}
	s2 := openapi3.Swagger{
		Info: &openapi3.Info{},
		Components: openapi3.Components{
			RequestBodies: openapi3.RequestBodies{"test": &openapi3.RequestBodyRef{}},
		},
	}
	_, err := diff.Get(diff.NewConfig(), &s1, &s2)
	require.Equal(t, "request body reference is nil", err.Error())

	_, err = diff.Get(diff.NewConfig(), &s2, &s1)
	require.Equal(t, "request body reference is nil", err.Error())
}

func TestDiff_ComponentResponsesNil(t *testing.T) {
	s1 := openapi3.Swagger{
		Info: &openapi3.Info{},
		Components: openapi3.Components{
			Responses: openapi3.Responses{"test": &openapi3.ResponseRef{Value: &openapi3.Response{}}},
		},
	}
	s2 := openapi3.Swagger{
		Info: &openapi3.Info{},
		Components: openapi3.Components{
			Responses: openapi3.Responses{"test": &openapi3.ResponseRef{}},
		},
	}
	_, err := diff.Get(diff.NewConfig(), &s1, &s2)
	require.Equal(t, "response reference is nil", err.Error())

	_, err = diff.Get(diff.NewConfig(), &s2, &s1)
	require.Equal(t, "response reference is nil", err.Error())
}

func TestDiff_ComponentSecuritySchemesNil(t *testing.T) {
	s1 := openapi3.Swagger{
		Info: &openapi3.Info{},
		Components: openapi3.Components{
			SecuritySchemes: openapi3.SecuritySchemes{"test": &openapi3.SecuritySchemeRef{Value: &openapi3.SecurityScheme{}}},
		},
	}

	s2 := openapi3.Swagger{
		Info: &openapi3.Info{},
		Components: openapi3.Components{
			SecuritySchemes: openapi3.SecuritySchemes{"test": &openapi3.SecuritySchemeRef{}},
		},
	}
	_, err := diff.Get(diff.NewConfig(), &s1, &s2)
	require.Equal(t, "security scheme reference is nil", err.Error())
	_, err = diff.Get(diff.NewConfig(), &s2, &s1)
	require.Equal(t, "security scheme reference is nil", err.Error())
}

func TestDiff_ComponentCallbacksNil(t *testing.T) {
	s1 := openapi3.Swagger{
		Info: &openapi3.Info{},
		Components: openapi3.Components{
			Callbacks: openapi3.Callbacks{"test": &openapi3.CallbackRef{Value: &openapi3.Callback{}}},
		},
	}
	s2 := openapi3.Swagger{
		Info: &openapi3.Info{},
		Components: openapi3.Components{
			Callbacks: openapi3.Callbacks{"test": &openapi3.CallbackRef{}},
		},
	}
	_, err := diff.Get(diff.NewConfig(), &s1, &s2)
	require.Equal(t, "callback reference is nil", err.Error())

	_, err = diff.Get(diff.NewConfig(), &s2, &s1)
	require.Equal(t, "callback reference is nil", err.Error())
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
