package allof_test

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/load"
)

func Test_MergeSpecOK(t *testing.T) {
	spec, err := load.NewSpecInfo(openapi3.NewLoader(), load.NewSource("../../data/allof/simple.yaml"), load.WithFlattenAllOf())
	require.NoError(t, err)
	merged := spec.Spec
	require.NoError(t, err)
	require.Equal(t, "string", merged.Components.Schemas["GroupView"].Value.Properties["created"].Value.Type)
	require.Equal(t, "string", merged.Components.Parameters["groupId"].Value.Schema.Value.Properties["prop1"].Value.Type)
	require.Equal(t, "boolean", merged.Components.Parameters["groupId"].Value.Schema.Value.Properties["prop2"].Value.Type)
	require.Empty(t, merged.Components.Parameters["groupId"].Value.Schema.Value.AllOf)
	require.Equal(t, "string", merged.Paths.Value("/api/v1.0/groups").Patch.RequestBody.Value.Content["application/json"].Schema.Value.Properties["prop1"].Value.Type)
	require.Equal(t, "boolean", merged.Paths.Value("/api/v1.0/groups").Patch.RequestBody.Value.Content["application/json"].Schema.Value.Properties["prop2"].Value.Type)
	require.Empty(t, merged.Paths.Value("/api/v1.0/groups").Patch.RequestBody.Value.Content["application/json"].Schema.Value.AllOf)
}

func Test_MergeSpecInvalid(t *testing.T) {
	_, err := load.NewSpecInfo(openapi3.NewLoader(), load.NewSource("../../data/allof/invalid.yaml"), load.WithFlattenAllOf())
	require.EqualError(t, err, "failed to flatten allOf in \"../../data/allof/invalid.yaml\": unable to resolve Type conflict: all Type values must be identical")
}

func TestMergeSpec_CircularAdditionalPropsWithoutAllOf(t *testing.T) {
	spec, err := load.NewSpecInfo(openapi3.NewLoader(), load.NewSource("testdata/circular_additional_props1.yaml"), load.WithFlattenAllOf())
	require.NoError(t, err)

	merged := spec.Spec
	require.Equal(t, "object", merged.Components.Schemas["BaseSchema"].Value.Properties["prop1"].Value.Type)
	require.NotNil(t, merged.Components.Schemas["BaseSchema"].Value.Properties["prop1"].Value.AdditionalProperties.Schema)
	require.NotNil(t, merged.Components.Schemas["BaseSchema"].Value.Properties["prop1"].Value.AdditionalProperties.Schema.Value)

	baseSchema := merged.Components.Schemas["BaseSchema"].Value
	referencedAdditionalPropSchema := merged.Components.Schemas["BaseSchema"].Value.Properties["prop1"].Value.AdditionalProperties.Schema.Value
	require.Equal(t, baseSchema, referencedAdditionalPropSchema)
}

func TestMergeSpec_MergeCircularAdditionalPropsWithAllOf(t *testing.T) {
	spec, err := load.NewSpecInfo(openapi3.NewLoader(), load.NewSource("testdata/circular_additional_props2.yaml"), load.WithFlattenAllOf())
	require.NoError(t, err)

	merged := spec.Spec
	require.Nil(t, merged.Components.Schemas["BaseSchema"].Value.AllOf)
	require.Equal(t, "string", merged.Components.Schemas["BaseSchema"].Value.Properties["fixedProperty"].Value.Type)
	require.Equal(t, "object", merged.Components.Schemas["BaseSchema"].Value.Properties["prop1"].Value.Type)
	require.NotNil(t, merged.Components.Schemas["BaseSchema"].Value.Properties["prop1"].Value.AdditionalProperties.Schema)
	require.NotNil(t, merged.Components.Schemas["BaseSchema"].Value.Properties["prop1"].Value.AdditionalProperties.Schema.Value)

	baseSchema := merged.Components.Schemas["BaseSchema"].Value
	referencedAdditionalPropSchema := merged.Components.Schemas["BaseSchema"].Value.Properties["prop1"].Value.AdditionalProperties.Schema.Value
	require.Equal(t, baseSchema, referencedAdditionalPropSchema)
}

func TestMergeSpec_MergeCircularAdditionalPropsNestedWithinAllOf(t *testing.T) {
	spec, err := load.NewSpecInfo(openapi3.NewLoader(), load.NewSource("testdata/circular_additional_props3.yaml"), load.WithFlattenAllOf())
	require.NoError(t, err)

	merged := spec.Spec
	require.Nil(t, merged.Components.Schemas["BaseSchema"].Value.AllOf)
	require.Equal(t, "object", merged.Components.Schemas["BaseSchema"].Value.Properties["prop1"].Value.Type)
	require.NotNil(t, merged.Components.Schemas["BaseSchema"].Value.Properties["prop1"].Value.AdditionalProperties.Schema)
	require.NotNil(t, merged.Components.Schemas["BaseSchema"].Value.Properties["prop1"].Value.AdditionalProperties.Schema.Value)

	baseSchemaReferencedAdditionalPropSchema := merged.Components.Schemas["BaseSchema"].Value.Properties["prop1"].Value.AdditionalProperties.Schema.Value
	NestedSelfReferentialSchema := merged.Components.Schemas["NestedSelfReferentialSchema"].Value
	require.Equal(t, baseSchemaReferencedAdditionalPropSchema, NestedSelfReferentialSchema)
}
