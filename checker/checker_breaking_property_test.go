package checker_test

import (
	"fmt"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

func getReqPropFile(file string) string {
	return fmt.Sprintf("../data/required-properties/%s", file)
}

// BC: new required property in request header is breaking
func TestBreaking_NewRequiredProperty(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Properties["courseId"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type:        "string",
			Description: "Unique ID of the course",
		},
	}
	s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Required = []string{"courseId"}

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "new-required-request-header-property", errs[0].Id)
}

// BC: new optional property in request header is not breaking
func TestBreaking_NewNonRequiredProperty(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Properties["courseId"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type:        "string",
			Description: "Unique ID of the course",
		},
	}

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
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

	s1.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Properties["courseId"] = &sr
	s1.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Required = []string{}

	s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Properties["courseId"] = &sr
	s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Required = []string{"courseId"}

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "request-header-property-became-required", errs[0].Id)
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

	s1.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Properties["courseId"] = &sr
	s1.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Required = []string{"courseId"}

	s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Properties["courseId"] = &sr
	s2.Spec.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Required = []string{}

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, &s1, &s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: changing an existing property in response body to optional is breaking
func TestBreaking_RespBodyRequiredPropertyDisabled(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("response-base.json"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("response-revision.json"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "response-property-became-optional", errs[0].Id)
}

// BC: changing a request property to enum is breaking
func TestBreaking_ReqBodyBecameEnum(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile("../data/enum/base-body.yaml")
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile("../data/enum/revision-body.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "request-body-became-enum", errs[0].Id)
}

// BC: changing a response body to nullable is breaking
func TestBreaking_RespBodyNullable(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile("../data/nullable/base-body.yaml")
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile("../data/nullable/revision-body.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "response-body-became-nullable", errs[0].Id)
}

// BC: changing a response property to nullable is breaking
func TestBreaking_RespBodyPropertyNullable(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile("../data/nullable/base-property.yaml")
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile("../data/nullable/revision-property.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "response-property-became-nullable", errs[0].Id)
}

// BC: changing an embedded reponse property to nullable is breaking
func TestBreaking_RespBodyEmbeddedPropertyNullable(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile("../data/nullable/base-embedded-property.yaml")
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile("../data/nullable/revision-embedded-property.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "response-property-became-nullable", errs[0].Id)
}

// BC: changing a required property in response body to optional and also deleting it is breaking
func TestBreaking_RespBodyDeleteAndDisableRequiredProperty(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("response-del-required-prop-base.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("response-del-required-prop-revision.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
}

// BC: adding a non-existent required property in request body is not breaking
func TestBreaking_ReqBodyNewRequiredPropertyNew(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("request-new-required-prop-base.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("request-new-required-prop-revision.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: changing an existing property in response body to required is not breaking
func TestBreaking_RespBodyRequiredPropertyEnabled(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("response-revision.json"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("response-base.json"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: changing an existing property in request body to optional is not breaking
func TestBreaking_ReqBodyRequiredPropertyDisabled(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("request-base.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("request-revision.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: changing an existing property in request body to required is breaking
func TestBreaking_ReqBodyRequiredPropertyEnabled(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("request-revision.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("request-base.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "request-property-became-required", errs[0].Id)
}

// BC: adding a new required property in request body is breaking
func TestBreaking_ReqBodyNewRequiredProperty(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("request-new-base.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("request-new-revision.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "new-required-request-property", errs[0].Id)
}

// BC: deleting a required property in request is breaking with warn
func TestBreaking_ReqBodyDeleteRequiredProperty(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("request-new-revision.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("request-new-base.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "request-property-removed", errs[0].Id)
	require.Equal(t, checker.WARN, errs[0].Level)
}

// BC: deleting a required property within another property in request is breaking with warn
func TestBreaking_ReqBodyDeleteRequiredProperty2(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("request-property-items.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("request-property-items-2.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "request-property-removed", errs[0].Id)
	require.Equal(t, checker.WARN, errs[0].Level)
}

// BC: adding a new required property in response body is not breaking
func TestBreaking_RespBodyNewRequiredProperty(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("response-new-base.json"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("response-new-revision.json"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: deleting a required property in response body is breaking
func TestBreaking_RespBodyDeleteRequiredProperty(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("response-new-revision.json"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("response-new-base.json"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "response-required-property-removed", errs[0].Id)
}

// BC: adding a new required property under AllOf in response body is not breaking
func TestBreaking_RespBodyNewAllOfRequiredProperty(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("response-allof-base.json"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("response-allof-revision.json"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: deleting a required property under AllOf in response body is breaking
func TestBreaking_RespBodyDeleteAllOfRequiredProperty(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("response-allof-revision.json"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("response-allof-base.json"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "response-required-property-removed", errs[0].Id)
}

// Old BC: adding a new required property under AllOf in response body is not breaking but when multiple inline (without $ref) schemas under AllOf are modified simultaneously, we detect is as breaking
// explanation: when multiple inline (without $ref) schemas under AllOf are modified we can't correlate schemas across base and revision
// as a result we can't determine that the change was "a new required property" and the change appears as breaking
// New BC: this "breaking change" is only the tip of the iceberg. The real problem is that allOf / anyOf must be processed specifically.
// For allOf all subSchemas must be merged before diff checking
// For anyOf schemas must be compared one-by-one
// I am going to change the behaviour in this case because currently it is false-positive case which can't correctly be checked
// At least now oit is a warn
func TestBreaking_RespBodyNewAllOfMultiRequiredProperty(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("response-allof-multi-base.json"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("response-allof-multi-revision.json"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "response-allOf-modified", errs[0].Id)
	require.Equal(t, checker.WARN, errs[0].Level)
}

// BC: adding a new required read-only property in request body is not breaking
func TestBreaking_ReadOnlyNewRequiredProperty(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("read-only-new-base.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("read-only-new-revision.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: changing an existing read-only property in request body to required is not breaking
func TestBreaking_ReadOnlyPropertyRequiredEnabled(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("read-only-base.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("read-only-revision.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: deleting a required write-only property in response body is not breaking
func TestBreaking_WriteOnlyDeleteRequiredProperty(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("write-only-delete-base.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("write-only-delete-revision.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "request-property-removed", errs[0].Id)
	require.Equal(t, checker.WARN, errs[0].Level)
}

// BC: deleting a non-required non-write-only property in response body is not breaking
func TestBreaking_WriteOnlyDeleteNonRequiredProperty(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("write-only-delete-partial-base.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("write-only-delete-partial-revision.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, "request-property-removed", errs[0].Id)
	require.Equal(t, checker.WARN, errs[0].Level)
}

// BC: changing an existing write-only property in response body to optional is not breaking
func TestBreaking_WriteOnlyPropertyRequiredDisabled(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("write-only-base.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("write-only-revision.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: changing an existing required property in response body to write-only is not breaking
func TestBreaking_RequiredPropertyWriteOnlyEnabled(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("write-only-changed-base.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("write-only-changed-revision.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.Empty(t, errs)
}

// BC: changing an existing required property in response body to not-write-only is breaking
func TestBreaking_RequiredPropertyWriteOnlyDisabled(t *testing.T) {
	s1, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("write-only-changed-revision.yaml"))
	require.NoError(t, err)

	s2, err := checker.LoadOpenAPISpecInfoFromFile(getReqPropFile("write-only-changed-base.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(&diff.Config{}, s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(checker.GetDefaultChecks(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 2)
	require.Equal(t, "response-required-property-became-not-write-only", errs[0].Id)
	require.Equal(t, checker.WARN, errs[0].Level)
	require.Equal(t, "response-required-property-became-not-write-only", errs[1].Id)
	require.Equal(t, checker.WARN, errs[1].Level)
}
