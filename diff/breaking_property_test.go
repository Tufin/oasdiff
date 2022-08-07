package diff_test

import (
	"fmt"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
)

func getReqPropFile(file string) string {
	return fmt.Sprintf("../data/required-properties/%s", file)
}

// BC: new required property in request header is breaking
func TestBreaking_NewRequiredProperty(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Properties["courseId"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type:        "string",
			Description: "Unique ID of the course",
		},
	}
	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Required = []string{"courseId"}

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.NotEmpty(t, d)
}

// BC: new optional property in request header is not breaking
func TestBreaking_NewNonRequiredProperty(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Properties["courseId"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type:        "string",
			Description: "Unique ID of the course",
		},
	}

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.Empty(t, d)
}

// BC: changing an existing property in request header to required is breaking
func TestBreaking_PropertyRequiredEnabled(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	sr := openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type:        "string",
			Description: "Unique ID of the course",
		},
	}

	s1.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Properties["courseId"] = &sr
	s1.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Required = []string{}

	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Properties["courseId"] = &sr
	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Required = []string{"courseId"}

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.NotEmpty(t, d)
}

// BC: changing an existing property in request header to optional is not breaking
func TestBreaking_PropertyRequiredDisabled(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	sr := openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type:        "string",
			Description: "Unique ID of the course",
		},
	}

	s1.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Properties["courseId"] = &sr
	s1.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Required = []string{"courseId"}

	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Properties["courseId"] = &sr
	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Required = []string{}

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.Empty(t, d)
}

// BC: changing an existing property in response body to optional is breaking
func TestBreaking_RespBodyRequiredPropertyDisabled(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getReqPropFile("response-base.json"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getReqPropFile("response-revision.json"))
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.NotEmpty(t, d)
}

// BC: changing an existing property in response body to required is not breaking
func TestBreaking_RespBodyRequiredPropertyEnabled(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getReqPropFile("response-revision.json"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getReqPropFile("response-base.json"))
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.Empty(t, d)
}

// BC: changing an existing property in request body to optional is not breaking
func TestBreaking_ReqBodyRequiredPropertyDisabled(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getReqPropFile("request-base.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getReqPropFile("request-revision.yaml"))
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.Empty(t, d)
}

// BC: changing an existing property in request body to required is breaking
func TestBreaking_ReqBodyRequiredPropertyEnabled(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getReqPropFile("request-revision.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getReqPropFile("request-base.yaml"))
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.NotEmpty(t, d)
}

// BC: adding a new required property in request body is breaking
func TestBreaking_ReqBodyNewRequiredProperty(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getReqPropFile("request-new-base.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getReqPropFile("request-new-revision.yaml"))
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.NotEmpty(t, d)
}

// BC: deleting a required property in request is not breaking
func TestBreaking_ReqBodyDeleteRequiredProperty(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getReqPropFile("request-new-revision.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getReqPropFile("request-new-base.yaml"))
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.Empty(t, d)
}

// BC: adding a new required property in response body is not breaking
func TestBreaking_RespBodyNewRequiredProperty(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getReqPropFile("response-new-base.json"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getReqPropFile("response-new-revision.json"))
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.Empty(t, d)
}

// BC: deleting a required property in response body is breaking
func TestBreaking_RespBodyDeleteRequiredProperty(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getReqPropFile("response-new-revision.json"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getReqPropFile("response-new-base.json"))
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.NotEmpty(t, d)
}

// BC: adding a new required property under AllOf in response body is not breaking
func TestBreaking_RespBodyNewAllOfRequiredProperty(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getReqPropFile("response-allof-base.json"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getReqPropFile("response-allof-revision.json"))
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.Empty(t, d)
}

// BC: deleting a required property under AllOf in response body is breaking
func TestBreaking_RespBodyDeleteAllOfRequiredProperty(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getReqPropFile("response-allof-revision.json"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getReqPropFile("response-allof-base.json"))
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.NotEmpty(t, d)
}

// BC: adding a new required property under AllOf in response body is not breaking but when multiple inline (without $ref) schemas under AllOf are modified simultaneously, we detect is as breaking
// explanation: when multiple inline (without $ref) schemas under AllOf are modified we can't correlate schemas across base and revision
// as a result we can't determine that the change was "a new required property" and the change appears as breaking
func TestBreaking_RespBodyNewAllOfMultiRequiredProperty(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getReqPropFile("response-allof-multi-base.json"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getReqPropFile("response-allof-multi-revision.json"))
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.NotEmpty(t, d)
}

// BC: adding a new required read-only property in request body is not breaking
func TestBreaking_ReadOnlyNewRequiredProperty(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getReqPropFile("read-only-new-base.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getReqPropFile("read-only-new-revision.yaml"))
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.Empty(t, d)
}

// BG: changing an existing read-only property in request body to required is not breaking
func TestBreaking_ReadOnlyPropertyRequiredEnabled(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getReqPropFile("read-only-base.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getReqPropFile("read-only-revision.yaml"))
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.Empty(t, d)
}

// BC: deleting a required write-only property in response body is not breaking
func TestBreaking_WriteOnlyDeleteRequiredProperty(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getReqPropFile("write-only-delete-base.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getReqPropFile("write-only-delete-revision.yaml"))
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.Empty(t, d)
}

// BC: deleting a non-required non-write-only property in response body is not breaking
func TestBreaking_WriteOnlyDeleteNonRequiredProperty(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getReqPropFile("write-only-delete-partial-base.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getReqPropFile("write-only-delete-partial-revision.yaml"))
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.Empty(t, d)
}

// BC: changing an existing write-only property in response body to optional is not breaking
func TestBreaking_WriteOnlyPropertyRequiredDisabled(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getReqPropFile("write-only-base.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getReqPropFile("write-only-revision.yaml"))
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.Empty(t, d)
}

// BC: changing an existing required property in response body to write-only is not breaking
func TestBreaking_RequiredPropertyWriteOnlyEnabled(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getReqPropFile("write-only-changed-base.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getReqPropFile("write-only-changed-revision.yaml"))
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.Empty(t, d)
}

// BC: changing an existing required property in response body to not-write-only is breaking
func TestBreaking_RequiredPropertyWriteOnlyDisabled(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getReqPropFile("write-only-changed-revision.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getReqPropFile("write-only-changed-base.yaml"))
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)
	require.NotEmpty(t, d)
}
