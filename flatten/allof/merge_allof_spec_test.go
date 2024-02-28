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
