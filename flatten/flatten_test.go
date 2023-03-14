package flatten_test

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/flatten"
)

func TestFlatten_EmptySchema(t *testing.T) {
	schema := openapi3.Schema{}
	flat := flatten.Handle(schema)
	require.Equal(t, &schema, flat)
}

func TestFlatten_NotAllOf(t *testing.T) {
	schema := openapi3.Schema{
		Title: "test",
	}
	flat := flatten.Handle(schema)
	require.Equal(t, &schema, flat)
}

func TestFlatten_OneObjectNoProps(t *testing.T) {

	schema := openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:       "object",
					Properties: openapi3.Schemas{},
				},
			},
		},
	}

	flat := flatten.Handle(schema)
	require.Equal(t, &schema, flat)
}

func TestFlatten_OneObjectOneProp(t *testing.T) {

	object := openapi3.Schemas{}
	object["description"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: "string",
		},
	}

	schema := openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:       "object",
					Properties: object,
				},
			},
		},
	}

	flat := flatten.Handle(schema)
	require.Equal(t, &schema, flat)
}

func TestFlatten_TwoObjects(t *testing.T) {

	obj1 := openapi3.Schemas{}
	obj1["description"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: "string",
		},
	}

	obj2 := openapi3.Schemas{}
	obj2["name"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: "string",
		},
	}

	schema := openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:       "object",
					Properties: obj1,
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:       "object",
					Properties: obj2,
				},
			},
		},
	}

	flat := flatten.Handle(schema)
	require.Len(t, flat.AllOf[0].Value.Properties, 2)
	require.Equal(t, obj1["description"], flat.AllOf[0].Value.Properties["description"])
	require.Equal(t, obj2["name"], flat.AllOf[0].Value.Properties["name"])
}

func TestFlatten_OverlappingProps(t *testing.T) {

	obj1 := openapi3.Schemas{}
	obj1["description"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: "string",
		},
	}

	obj2 := openapi3.Schemas{}
	obj2["description"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type: "int",
		},
	}

	schema := openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:       "object",
					Properties: obj1,
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:       "object",
					Properties: obj2,
				},
			},
		},
	}

	flat := flatten.Handle(schema)
	require.Len(t, flat.AllOf[0].Value.Properties, 1)
	require.Equal(t, obj1["description"], flat.AllOf[0].Value.Properties["description"])
}