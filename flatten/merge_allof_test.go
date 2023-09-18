package flatten_test

import (
	"context"
	"testing"

	"github.com/tufin/oasdiff/flatten"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
)

func TestMerge_NestedAllOfInProperties(t *testing.T) {
	merged, err := flatten.Merge(openapi3.Schema{
		Properties: openapi3.Schemas{
			"prop1": &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: "object",
					AllOf: openapi3.SchemaRefs{
						&openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type:     "object",
								MinProps: 10,
								MaxProps: openapi3.Uint64Ptr(40),
							},
						},
						&openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type:     "object",
								MinProps: 5,
								MaxProps: openapi3.Uint64Ptr(25),
							},
						},
					},
				},
			},
		},
	},
	)
	require.NoError(t, err)
	require.Nil(t, merged.Properties["prop1"].Value.AllOf)
	require.Equal(t, uint64(10), merged.Properties["prop1"].Value.MinProps)
	require.Equal(t, uint64(25), *merged.Properties["prop1"].Value.MaxProps)
}

func TestMerge_NestedAllOfInNot(t *testing.T) {
	merged, err := flatten.Merge(openapi3.Schema{
		Not: &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type: "object",
				AllOf: openapi3.SchemaRefs{
					&openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Type:     "object",
							MinProps: 10,
							MaxProps: openapi3.Uint64Ptr(40),
						},
					},
					&openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Type:     "object",
							MinProps: 5,
							MaxProps: openapi3.Uint64Ptr(25),
						},
					},
				},
			},
		},
	},
	)
	require.NoError(t, err)
	require.Nil(t, merged.Not.Value.AllOf)
	require.Equal(t, uint64(10), merged.Not.Value.MinProps)
	require.Equal(t, uint64(25), *merged.Not.Value.MaxProps)
}

func TestMerge_NestedAllOfInOneOf(t *testing.T) {
	merged, err := flatten.Merge(openapi3.Schema{
		OneOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: "object",
					AllOf: openapi3.SchemaRefs{
						&openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type:     "object",
								MinProps: 10,
								MaxProps: openapi3.Uint64Ptr(40),
							},
						},
						&openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type:     "object",
								MinProps: 5,
								MaxProps: openapi3.Uint64Ptr(25),
							},
						},
					},
				},
			},
		},
	})
	require.NoError(t, err)
	require.Nil(t, merged.OneOf[0].Value.AllOf)
	require.Equal(t, uint64(10), merged.OneOf[0].Value.MinProps)
	require.Equal(t, uint64(25), *merged.OneOf[0].Value.MaxProps)
}

// AllOf is empty in base schema, but there is nested non-empty AllOf in base schema.
func TestMerge_NestedAllOfInAnyOf(t *testing.T) {
	merged, err := flatten.Merge(openapi3.Schema{
		AnyOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: "object",
					AllOf: openapi3.SchemaRefs{
						&openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type:     "object",
								MinProps: 10,
								MaxProps: openapi3.Uint64Ptr(40),
							},
						},
						&openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type:     "object",
								MinProps: 5,
								MaxProps: openapi3.Uint64Ptr(25),
							},
						},
					},
				},
			},
		},
	})
	require.NoError(t, err)
	require.Nil(t, merged.AnyOf[0].Value.AllOf)
	require.Equal(t, uint64(10), merged.AnyOf[0].Value.MinProps)
	require.Equal(t, uint64(25), *merged.AnyOf[0].Value.MaxProps)
}

// identical numeric types are merged successfully
func TestMerge_TypeNumeric(t *testing.T) {
	merged, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Properties: openapi3.Schemas{
						"prop1": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: "number",
							},
						},
					},
					Type: "object",
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Properties: openapi3.Schemas{
						"prop1": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: "number",
							},
						},
					},
					Type: "object",
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, "number", merged.Properties["prop1"].Value.Type)

	merged, err = flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Properties: openapi3.Schemas{
						"prop1": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: "integer",
							},
						},
					},
					Type: "object",
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Properties: openapi3.Schemas{
						"prop1": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: "integer",
							},
						},
					},
					Type: "object",
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, "integer", merged.Properties["prop1"].Value.Type)
}

// Conflicting numeric types are merged successfully
func TestMerge_TypeNumericConflictResolved(t *testing.T) {
	merged, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Properties: openapi3.Schemas{
						"prop1": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: "integer",
							},
						},
					},
					Type: "object",
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Properties: openapi3.Schemas{
						"prop1": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: "number",
							},
						},
					},
					Type: "object",
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, "integer", merged.Properties["prop1"].Value.Type)
}

// Conflicting types cannot be resolved
func TestMerge_TypeFailure(t *testing.T) {
	_, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Properties: openapi3.Schemas{
						"prop1": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: "integer",
							},
						},
					},
					Type: "object",
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Properties: openapi3.Schemas{
						"prop1": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: "string",
							},
						},
					},
					Type: "object",
				},
			},
		},
	})

	require.EqualError(t, err, flatten.TypeErrorMessage)
}

// if ExclusiveMax is true on the minimum Max value, then ExclusiveMax is true in the merged schema.
func TestMerge_ExclusiveMaxIsTrue(t *testing.T) {
	merged, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:         "object",
					ExclusiveMax: true,
					Max:          openapi3.Float64Ptr(1),
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:         "object",
					ExclusiveMax: false,
					Max:          openapi3.Float64Ptr(2),
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, true, merged.ExclusiveMax)
}

// if ExclusiveMax is false on the minimum Max value, then ExclusiveMax is false in the merged schema.
func TestMerge_ExclusiveMaxIsFalse(t *testing.T) {
	merged, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:         "object",
					ExclusiveMax: false,
					Max:          openapi3.Float64Ptr(1),
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:         "object",
					ExclusiveMax: true,
					Max:          openapi3.Float64Ptr(2),
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, false, merged.ExclusiveMax)
}

// if ExclusiveMin is false on the highest Min value, then ExclusiveMin is false in the merged schema.
func TestMerge_ExclusiveMinIsFalse(t *testing.T) {
	merged, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:         "object",
					ExclusiveMin: false,
					Min:          openapi3.Float64Ptr(40),
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:         "object",
					ExclusiveMin: true,
					Min:          openapi3.Float64Ptr(5),
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, false, merged.ExclusiveMin)
}

// if ExclusiveMin is true on the highest Min value, then ExclusiveMin is true in the merged schema.
func TestMerge_ExclusiveMinIsTrue(t *testing.T) {
	merged, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:         "object",
					ExclusiveMin: true,
					Min:          openapi3.Float64Ptr(40),
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:         "object",
					ExclusiveMin: false,
					Min:          openapi3.Float64Ptr(5),
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, true, merged.ExclusiveMin)
}

// merge multiple Not inside AllOf
func TestMerge_Not(t *testing.T) {
	merged, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: "object",
					Not: &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Type: "string",
						},
					},
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: "object",
					Not: &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Type: "integer",
						},
					},
				},
			},
		},
	})

	require.NoError(t, err)
	require.Equal(t, "string", merged.Not.Value.AnyOf[0].Value.Type)
	require.Equal(t, "integer", merged.Not.Value.AnyOf[1].Value.Type)
}

// merge multiple OneOf inside AllOf
func TestMerge_OneOf(t *testing.T) {
	merged, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: "object",
					OneOf: openapi3.SchemaRefs{
						&openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type:     "object",
								Required: []string{"prop1"},
							},
						},
						&openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type:     "object",
								Required: []string{"prop2"},
							},
						},
					},
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: "object",
					OneOf: openapi3.SchemaRefs{
						&openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type:     "object",
								Required: []string{"prop2"},
							},
						},
					},
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, []string{"prop1", "prop2"}, merged.OneOf[0].Value.Required)
	require.Equal(t, []string{"prop2"}, merged.OneOf[1].Value.Required)
}

// merge multiple AnyOf inside AllOf
func TestMerge_AnyOf(t *testing.T) {
	merged, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: "object",
					AnyOf: openapi3.SchemaRefs{
						&openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type:     "object",
								Required: []string{"string"},
							},
						},
						&openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type:     "object",
								Required: []string{"boolean"},
							},
						},
					},
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: "object",
					AnyOf: openapi3.SchemaRefs{
						&openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type:     "object",
								Required: []string{"boolean"},
							},
						},
					},
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, []string{"string", "boolean"}, merged.AnyOf[0].Value.Required)
}

// conflicting uniqueItems values are merged successfully
func TestMerge_UniqueItemsTrue(t *testing.T) {
	merged, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:        "object",
					UniqueItems: true,
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:        "object",
					UniqueItems: false,
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, true, merged.UniqueItems)
}

// non-conflicting uniqueItems values are merged successfully
func TestMerge_UniqueItemsFalse(t *testing.T) {
	merged, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:        "object",
					UniqueItems: false,
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:        "object",
					UniqueItems: false,
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, false, merged.UniqueItems)
}

// Item merge fails due to conflicting item types.
func TestMerge_Items_Failure(t *testing.T) {
	_, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: "object",
					Properties: openapi3.Schemas{
						"test": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: "array",
								Items: &openapi3.SchemaRef{
									Value: &openapi3.Schema{
										Type: "integer",
									},
								},
							},
						},
					},
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: "object",
					Properties: openapi3.Schemas{
						"test": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: "array",
								Items: &openapi3.SchemaRef{
									Value: &openapi3.Schema{
										Type: "string",
									},
								},
							},
						},
					},
				},
			},
		},
	})
	require.EqualError(t, err, flatten.TypeErrorMessage)
}

// items are merged successfully when there are no conflicts
func TestMerge_Items(t *testing.T) {
	merged, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: "object",
					Properties: openapi3.Schemas{
						"test": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: "array",
								Items: &openapi3.SchemaRef{
									Value: &openapi3.Schema{
										Type: "integer",
									},
								},
							},
						},
					},
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: "object",
					Properties: openapi3.Schemas{
						"test": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: "array",
								Items: &openapi3.SchemaRef{
									Value: &openapi3.Schema{
										Type: "integer",
									},
								},
							},
						},
					},
				},
			},
		},
	})
	require.NoError(t, err)
	require.Nil(t, merged.AllOf)
	require.Equal(t, "array", merged.Properties["test"].Value.Type)
	require.Equal(t, "integer", merged.Properties["test"].Value.Items.Value.Type)
}

func TestMerge_MultipleOfContained(t *testing.T) {

	//todo - more tests
	merged, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:       "object",
					MultipleOf: openapi3.Float64Ptr(10.0),
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:       "object",
					MultipleOf: openapi3.Float64Ptr(2.0),
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, float64(10), *merged.MultipleOf)
}

func TestMerge_MultipleOfDecimal(t *testing.T) {
	merged, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:       "object",
					MultipleOf: openapi3.Float64Ptr(11.0),
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:       "object",
					MultipleOf: openapi3.Float64Ptr(0.7),
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, float64(77), *merged.MultipleOf)
}

func TestMerge_EnumContained(t *testing.T) {
	merged, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: "object",
					Enum: []interface{}{"1", nil, 1},
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: "object",
					Enum: []interface{}{"1"},
				},
			},
		},
	})
	require.NoError(t, err)
	require.ElementsMatch(t, []interface{}{"1"}, merged.Enum)
}

// enum merge fails if the intersection of enum values is empty.
func TestMerge_EnumNoIntersection(t *testing.T) {
	_, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: "object",
					Enum: []interface{}{"1", nil},
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: "object",
					Enum: []interface{}{"2"},
				},
			},
		},
	})
	require.Error(t, err)
}

// Properties range is the most restrictive
func TestMerge_RangeProperties(t *testing.T) {
	merged, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:     "object",
					MinProps: 10,
					MaxProps: openapi3.Uint64Ptr(40),
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:     "object",
					MinProps: 5,
					MaxProps: openapi3.Uint64Ptr(25),
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, uint64(10), merged.MinProps)
	require.Equal(t, uint64(25), *merged.MaxProps)
}

// Items range is the most restrictive
func TestMerge_RangeItems(t *testing.T) {

	merged, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:     "object",
					MinItems: 10,
					MaxItems: openapi3.Uint64Ptr(40),
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:     "object",
					MinItems: 5,
					MaxItems: openapi3.Uint64Ptr(25),
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, uint64(10), merged.MinItems)
	require.Equal(t, uint64(25), *merged.MaxItems)
}

func TestMerge_Range(t *testing.T) {
	merged, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: "object",
					Min:  openapi3.Float64Ptr(10),
					Max:  openapi3.Float64Ptr(40),
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: "object",
					Min:  openapi3.Float64Ptr(5),
					Max:  openapi3.Float64Ptr(25),
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, float64(10), *merged.Min)
	require.Equal(t, float64(25), *merged.Max)
}

func TestMerge_MaxLength(t *testing.T) {
	merged, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:      "object",
					MaxLength: openapi3.Uint64Ptr(10),
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:      "object",
					MaxLength: openapi3.Uint64Ptr(20),
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, uint64(10), *merged.MaxLength)
}

func TestMerge_MinLength(t *testing.T) {
	merged, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:      "object",
					MinLength: 10,
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:      "object",
					MinLength: 20,
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, uint64(20), merged.MinLength)
}

func TestMerge_Description(t *testing.T) {
	merged, err := flatten.Merge(openapi3.Schema{
		Description: "desc0",
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:        "object",
					Description: "desc1",
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:        "object",
					Description: "desc2",
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, "desc0", merged.Description)
}

// non-conflicting types are merged successfully
func TestMerge_NonConflictingType(t *testing.T) {
	merged, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: "object",
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: "object",
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, "object", merged.Type)
}

// schema cannot be merged if types are conflicting
func TestMerge_FailsOnConflictingTypes(t *testing.T) {
	_, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: "object",
					Properties: openapi3.Schemas{
						"name": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: "string",
							},
						},
					},
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: "object",
					Properties: openapi3.Schemas{
						"name": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: "object",
							},
						},
					},
				},
			},
		},
	})
	require.Error(t, err)
}

func TestMerge_Title(t *testing.T) {
	merged, err := flatten.Merge(openapi3.Schema{
		Title: "base schema",
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:  "object",
					Title: "first",
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:  "object",
					Title: "second",
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, "base schema", merged.Title)
}

// merge conflicting integer formats
func TestMerge_FormatInteger(t *testing.T) {
	merged, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Properties: openapi3.Schemas{
						"prop1": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Format: flatten.FormatInt32,
							},
						},
					},
					Type: "object",
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Properties: openapi3.Schemas{
						"prop1": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Format: flatten.FormatInt64,
							},
						},
					},
					Type: "object",
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, flatten.FormatInt32, merged.Properties["prop1"].Value.Format)
}

// merge conflicting float formats
func TestMerge_FormatFloat(t *testing.T) {
	merged, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Properties: openapi3.Schemas{
						"prop1": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Format: flatten.FormatFloat,
							},
						},
					},
					Type: "object",
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Properties: openapi3.Schemas{
						"prop1": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Format: flatten.FormatDouble,
							},
						},
					},
					Type: "object",
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, flatten.FormatFloat, merged.Properties["prop1"].Value.Format)
}

// merge conflicting integer and float formats
func TestMerge_NumericFormat(t *testing.T) {
	merged, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Properties: openapi3.Schemas{
						"prop1": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Format: flatten.FormatFloat,
							},
						},
					},
					Type: "object",
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Properties: openapi3.Schemas{
						"prop1": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Format: flatten.FormatDouble,
							},
						},
					},
					Type: "object",
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Properties: openapi3.Schemas{
						"prop1": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Format: flatten.FormatInt32,
							},
						},
					},
					Type: "object",
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, flatten.FormatInt32, merged.Properties["prop1"].Value.Format)
}

func TestMerge_Format(t *testing.T) {
	merged, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:   "object",
					Format: "date",
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:   "object",
					Format: "date",
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, "date", merged.Format)
}

func TestMerge_Format_Failure(t *testing.T) {
	_, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:   "object",
					Format: "date",
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:   "object",
					Format: "byte",
				},
			},
		},
	})
	require.EqualError(t, err, flatten.FormatErrorMessage)
}

func TestMerge_EmptySchema(t *testing.T) {
	schema := openapi3.Schema{}
	merged, err := flatten.Merge(schema)
	require.NoError(t, err)
	require.Equal(t, &schema, merged)
}

func TestMerge_NoAllOf(t *testing.T) {
	schema := openapi3.Schema{
		Title: "test",
	}
	merged, err := flatten.Merge(schema)
	require.NoError(t, err)
	require.Equal(t, &schema, merged)
}

func TestMerge_TwoObjects(t *testing.T) {

	obj1 := openapi3.Schemas{
		"description": &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type: "string",
			},
		},
	}

	obj2 := openapi3.Schemas{
		"name": &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type: "string",
			},
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

	merged, err := flatten.Merge(schema)
	require.NoError(t, err)
	require.Len(t, merged.AllOf, 0)
	require.Len(t, merged.Properties, 2)
	require.Equal(t, obj1["description"].Value.Type, merged.Properties["description"].Value.Type)
	require.Equal(t, obj2["name"].Value.Type, merged.Properties["name"].Value.Type)
}

func TestMerge_OneObjectOneProp(t *testing.T) {

	object := openapi3.Schemas{
		"description": &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type: "string",
			},
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

	merged, err := flatten.Merge(schema)
	require.NoError(t, err)
	require.Len(t, merged.Properties, 1)
	require.Equal(t, object["description"].Value.Type, merged.Properties["description"].Value.Type)
}

func TestMerge_OneObjectNoProps(t *testing.T) {

	merged, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:       "object",
					Properties: openapi3.Schemas{},
				},
			},
		},
	})
	require.NoError(t, err)
	require.Len(t, merged.Properties, 0)
}

func TestMerge_OverlappingProps(t *testing.T) {

	obj1 := openapi3.Schemas{
		"description": &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Title: "first",
			},
		},
	}

	obj2 := openapi3.Schemas{
		"description": &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Title: "second",
			},
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
	merged, err := flatten.Merge(schema)
	require.NoError(t, err)
	require.Len(t, merged.AllOf, 0)
	require.Len(t, merged.Properties, 1)
	require.Equal(t, (*obj1["description"].Value), (*merged.Properties["description"].Value))
}

// if additionalProperties is false, then the merged additionalProperties is the intersection of relevant properties.
func TestMerge_AdditionalProperties_False(t *testing.T) {
	apFalse := false
	apTrue := true

	var firstPropEnum []interface{}
	var secondPropEnum []interface{}
	var thirdPropEnum []interface{}

	firstPropEnum = append(firstPropEnum, "1", "5", "3")
	secondPropEnum = append(secondPropEnum, "1", "8", "7")
	thirdPropEnum = append(thirdPropEnum, "3", "8", "5")

	merged, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: "object",
					Properties: openapi3.Schemas{
						"prop1": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Enum: firstPropEnum,
							},
						},
						"name": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: "string",
							},
						},
					},
					AdditionalProperties: openapi3.AdditionalProperties{
						Has: &apTrue,
					},
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: "object",
					Properties: openapi3.Schemas{
						"prop2": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Enum: secondPropEnum,
							},
						},
					},
					AdditionalProperties: openapi3.AdditionalProperties{
						Has: &apFalse,
					},
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: "object",
					Properties: openapi3.Schemas{
						"prop2": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Enum: thirdPropEnum,
							},
						},
					},
					AdditionalProperties: openapi3.AdditionalProperties{
						Has: &apFalse,
					},
				},
			},
		}})
	require.NoError(t, err)
	require.Equal(t, "8", merged.Properties["prop2"].Value.Enum[0])
}

// if additionalProperties is true, then the merged additionalProperties is the intersection of all properties.
func TestMerge_AdditionalProperties_True(t *testing.T) {
	apTrue := true

	var firstPropEnum []interface{}
	var secondPropEnum []interface{}
	var thirdPropEnum []interface{}

	firstPropEnum = append(firstPropEnum, "1", "5", "3")
	secondPropEnum = append(secondPropEnum, "1", "8", "7")
	thirdPropEnum = append(thirdPropEnum, "3", "8", "5")

	merged, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: "object",
					Properties: openapi3.Schemas{
						"prop1": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Enum: firstPropEnum,
							},
						},
						"name": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: "string",
							},
						},
					},
					AdditionalProperties: openapi3.AdditionalProperties{
						Has: &apTrue,
					},
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: "object",
					Properties: openapi3.Schemas{
						"prop2": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Enum: secondPropEnum,
							},
						},
					},
					AdditionalProperties: openapi3.AdditionalProperties{
						Has: &apTrue,
					},
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: "object",
					Properties: openapi3.Schemas{
						"prop2": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Enum: thirdPropEnum,
							},
						},
					},
					AdditionalProperties: openapi3.AdditionalProperties{
						Has: &apTrue,
					},
				},
			},
		}})
	require.NoError(t, err)
	require.Equal(t, "string", merged.Properties["name"].Value.Type)
}

func TestMergeAllOf_Pattern(t *testing.T) {

	merged, err := flatten.Merge(openapi3.Schema{
		AllOf: openapi3.SchemaRefs{
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:    "object",
					Pattern: "foo",
				},
			},
			&openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type:    "object",
					Pattern: "bar",
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, "(?=foo)(?=bar)", merged.Pattern)
}

func TestMerge_Required(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		filename string
	}{
		{"testdata/properties.yml"},
	}

	for _, test := range tests {
		t.Run(test.filename, func(t *testing.T) {
			// Load in the reference spec from the testdata
			sl := openapi3.NewLoader()
			sl.IsExternalRefsAllowed = true
			doc, err := sl.LoadFromFile(test.filename)
			require.NoError(t, err, "loading test file")
			err = doc.Validate(ctx)
			require.NoError(t, err, "validating spec")
			merged, err := flatten.Merge(*doc.Paths["/products"].Get.Responses["200"].Value.Content["application/json"].Schema.Value)
			require.NoError(t, err)

			props := merged.Properties
			require.Len(t, props, 3)
			require.Contains(t, props, "id")
			require.Contains(t, props, "createdAt")
			require.Contains(t, props, "otherId")

			required := merged.Required
			require.Len(t, required, 2)
			require.Contains(t, required, "id")
			require.Contains(t, required, "otherId")
		})
	}
}
