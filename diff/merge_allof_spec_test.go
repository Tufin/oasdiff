package diff_test

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
)

func Test_MergeSpec(t *testing.T) {
	loader := openapi3.NewLoader()
	s, err := loader.LoadFromFile("../data/allof.yaml")
	require.NoError(t, err)
	merged, err := diff.MergeSpec(*s)
	require.NoError(t, err)
	require.Equal(t, "string", merged.Components.Schemas["GroupView"].Value.Properties["created"].Value.Type)
	require.Equal(t, "string", merged.Components.Parameters["groupId"].Value.Schema.Value.Properties["prop1"].Value.Type)
	require.Equal(t, "boolean", merged.Components.Parameters["groupId"].Value.Schema.Value.Properties["prop2"].Value.Type)
	require.Empty(t, merged.Components.Parameters["groupId"].Value.Schema.Value.AllOf)
	require.Equal(t, "string", merged.Paths["/api/v1.0/groups"].Patch.RequestBody.Value.Content["application/json"].Schema.Value.Properties["prop1"].Value.Type)
	require.Equal(t, "boolean", merged.Paths["/api/v1.0/groups"].Patch.RequestBody.Value.Content["application/json"].Schema.Value.Properties["prop2"].Value.Type)
	require.Empty(t, merged.Paths["/api/v1.0/groups"].Patch.RequestBody.Value.Content["application/json"].Schema.Value.AllOf)
}
