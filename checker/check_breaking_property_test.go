package checker_test

import (
	"fmt"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

func getReqPropFile(file string) string {
	return fmt.Sprintf("../data/required-properties/%s", file)
}

// BC: new required property in request header is breaking: new-required-request-header-property
func TestBreaking_NewRequiredProperty(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Properties["courseId"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type:        &openapi3.Types{"string"},
			Description: "Unique ID of the course",
		},
	}
	s2.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Required = []string{"courseId"}

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.NewRequiredRequestHeaderPropertyId, errs[0].GetId())
	require.Equal(t, "added required 'network-policies' request header property 'courseId'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: new optional property in request header is not breaking
func TestBreaking_NewNonRequiredProperty(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Properties["courseId"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type:        &openapi3.Types{"string"},
			Description: "Unique ID of the course",
		},
	}

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.Empty(t, errs)
}

// BC: changing an existing property in request header to required is breaking: request-header-property-became-required
func TestBreaking_PropertyRequiredEnabled(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	sr := openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type:        &openapi3.Types{"string"},
			Description: "Unique ID of the course",
		},
	}

	s1.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Properties["courseId"] = &sr
	s1.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Required = []string{}

	s2.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Properties["courseId"] = &sr
	s2.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Required = []string{"courseId"}

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.RequestHeaderPropertyBecameRequiredId, errs[0].GetId())
	require.Equal(t, "property 'courseId' of request header 'network-policies' became required", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: changing an existing property in request header to optional is not breaking
func TestBreaking_PropertyRequiredDisabled(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	sr := openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type:        &openapi3.Types{"string"},
			Description: "Unique ID of the course",
		},
	}

	s1.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Properties["courseId"] = &sr
	s1.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Required = []string{"courseId"}

	s2.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Properties["courseId"] = &sr
	s2.Spec.Paths.Value(installCommandPath).Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Required = []string{}

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.Empty(t, errs)
}

// BC: changing an existing property in response body to optional is breaking: response-property-became-optional
func TestBreaking_RespBodyRequiredPropertyDisabled(t *testing.T) {
	s1, err := open(getReqPropFile("response-base.json"))
	require.NoError(t, err)

	s2, err := open(getReqPropFile("response-revision.json"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ResponsePropertyBecameOptionalId, errs[0].GetId())
	require.Equal(t, "property 'helpAndSupport/title' became optional for response status '200'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: changing a request body to enum is breaking: request-body-became-enum
func TestBreaking_ReqBodyBecameEnum(t *testing.T) {
	s1, err := open("../data/enums/request-body-no-enum.yaml")
	require.NoError(t, err)

	s2, err := open("../data/enums/request-body-enum.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.RequestBodyBecameEnumId, errs[0].GetId())
	require.Equal(t, "media-type 'application/json' of request body was restricted to a list of enum values", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: adding an enum value to request body is not breaking
func TestBreaking_ReqBodyEnumValueAdded(t *testing.T) {
	s1, err := open("../data/enums/request-body-enum.yaml")
	require.NoError(t, err)

	s2, err := open("../data/enums/request-body-enum-revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.Empty(t, errs)
}

// BC: changing a request body type and changing it to enum simultaneously is breaking: request-body-became-enum, request-body-type-changed
func TestBreaking_ReqBodyBecameEnumAndTypeChanged(t *testing.T) {
	s1, err := open("../data/enums/request-body-no-enum.yaml")
	require.NoError(t, err)

	s2, err := open("../data/enums/request-body-enum-int.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.Len(t, errs, 2)

	require.Equal(t, checker.RequestBodyBecameEnumId, errs[0].GetId())
	require.Equal(t, "media-type 'application/json' of request body was restricted to a list of enum values", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))

	require.Equal(t, checker.RequestBodyTypeChangedId, errs[1].GetId())
	require.Equal(t, "type/format of media-type 'application/json' of request body changed from 'string'/'' to 'int'/''", errs[1].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: changing an existing property in request body to enum is breaking: request-property-became-enum
func TestBreaking_ReqPropertyBecameEnum(t *testing.T) {
	s1, err := open("../data/enums/request-property-no-enum.yaml")
	require.NoError(t, err)

	s2, err := open("../data/enums/request-property-enum.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.RequestPropertyBecameEnumId, errs[0].GetId())
	require.Equal(t, "request property 'name' of media-type 'application/json' was restricted to a list of enum values", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: changing an existing path param to enum is breaking: request-parameter-became-enum
func TestBreaking_ReqParameterBecameEnum(t *testing.T) {
	s1, err := open("../data/enums/request-parameter-op-no-enum.yaml")
	require.NoError(t, err)

	s2, err := open("../data/enums/request-parameter-op-enum.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.RequestParameterBecameEnumId, errs[0].GetId())
	require.Equal(t, "'path' request parameter 'bookId' was restricted to a list of enum values", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: changing an existing property in request header to enum is breaking: request-header-property-became-enum
func TestBreaking_ReqParameterHeaderPropertyBecameEnum(t *testing.T) {
	s1, err := open("../data/enums/request-parameter-property-no-enum.yaml")
	require.NoError(t, err)

	s2, err := open("../data/enums/request-parameter-property-enum.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.RequestHeaderPropertyBecameEnumId, errs[0].GetId())
	require.Equal(t, "property 'name' of request header 'bookId' was restricted to a list of enum values", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: changing a response body to nullable is breaking: response-body-became-nullable
func TestBreaking_RespBodyNullable(t *testing.T) {
	s1, err := open("../data/nullable/base-body.yaml")
	require.NoError(t, err)

	s2, err := open("../data/nullable/revision-body.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ResponseBodyBecameNullableId, errs[0].GetId())
	require.Equal(t, "response body became nullable", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: changing a request property to not nullable is breaking: request-property-became-not-nullable
func TestBreaking_ReqBodyPropertyNotNullable(t *testing.T) {
	s1, err := open("../data/nullable/base-req.yaml")
	require.NoError(t, err)

	s2, err := open("../data/nullable/revision-req.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.RequestPropertyBecomeNotNullableId, errs[0].GetId())
	require.Equal(t, "request property 'id' of media-type 'application/json' became not nullable", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: changing a response property to nullable is breaking: response-property-became-nullable
func TestBreaking_RespBodyPropertyNullable(t *testing.T) {
	s1, err := open("../data/nullable/base-property.yaml")
	require.NoError(t, err)

	s2, err := open("../data/nullable/revision-property.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ResponsePropertyBecameNullableId, errs[0].GetId())
	require.Equal(t, "property 'name' became nullable for response status '201'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: changing an embedded response property to nullable is breaking: response-property-became-nullable
func TestBreaking_RespBodyEmbeddedPropertyNullable(t *testing.T) {
	s1, err := open("../data/nullable/base-embedded-property.yaml")
	require.NoError(t, err)

	s2, err := open("../data/nullable/revision-embedded-property.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ResponsePropertyBecameNullableId, errs[0].GetId())
	require.Equal(t, "property 'name/name' became nullable for response status '201'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: changing a required property in response body to optional and also deleting it is breaking: response-required-property-removed
func TestBreaking_RespBodyDeleteAndDisableRequiredProperty(t *testing.T) {
	s1, err := open(getReqPropFile("response-del-required-prop-base.yaml"))
	require.NoError(t, err)

	s2, err := open(getReqPropFile("response-del-required-prop-revision.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ResponseRequiredPropertyRemovedId, errs[0].GetId())
	require.Equal(t, "removed required property 'name' from response status '200'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: adding a non-existent required property in request body is not breaking
func TestBreaking_ReqBodyNewRequiredPropertyNew(t *testing.T) {
	s1, err := open(getReqPropFile("request-new-required-prop-base.yaml"))
	require.NoError(t, err)

	s2, err := open(getReqPropFile("request-new-required-prop-revision.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.Empty(t, errs)
}

// BC: changing an existing property in response body to required is not breaking
func TestBreaking_RespBodyRequiredPropertyEnabled(t *testing.T) {
	s1, err := open(getReqPropFile("response-revision.json"))
	require.NoError(t, err)

	s2, err := open(getReqPropFile("response-base.json"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.Empty(t, errs)
}

// BC: changing an existing property in request body to optional is not breaking
func TestBreaking_ReqBodyRequiredPropertyDisabled(t *testing.T) {
	s1, err := open(getReqPropFile("request-base.yaml"))
	require.NoError(t, err)

	s2, err := open(getReqPropFile("request-revision.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.Empty(t, errs)
}

// BC: changing an existing property in request body to required is breaking: request-property-became-required
func TestBreaking_ReqBodyRequiredPropertyEnabled(t *testing.T) {
	s1, err := open(getReqPropFile("request-revision.yaml"))
	require.NoError(t, err)

	s2, err := open(getReqPropFile("request-base.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.RequestPropertyBecameRequiredId, errs[0].GetId())
	require.Equal(t, "request property 'email' of media-type 'application/json' became required", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: adding a new required property in request body is breaking: new-required-request-property
func TestBreaking_ReqBodyNewRequiredProperty(t *testing.T) {
	s1, err := open(getReqPropFile("request-new-base.yaml"))
	require.NoError(t, err)

	s2, err := open(getReqPropFile("request-new-revision.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.NewRequiredRequestPropertyId, errs[0].GetId())
	require.Equal(t, "added required request property 'email' to media-type 'application/json'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: deleting a required property in request is breaking with warn: request-property-removed
func TestBreaking_ReqBodyDeleteRequiredProperty(t *testing.T) {
	s1, err := open(getReqPropFile("request-new-revision.yaml"))
	require.NoError(t, err)

	s2, err := open(getReqPropFile("request-new-base.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.RequestPropertyRemovedId, errs[0].GetId())
	require.Equal(t, checker.WARN, errs[0].GetLevel())
	require.Equal(t, "removed request property 'email' of media-type 'application/json'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: deleting an embedded optional property in request is breaking with warn: request-property-removed
func TestBreaking_ReqBodyDeleteRequiredProperty2(t *testing.T) {
	s1, err := open(getReqPropFile("request-property-items.yaml"))
	require.NoError(t, err)

	s2, err := open(getReqPropFile("request-property-items-2.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.Contains(t, errs, checker.ApiChange{
		Id:        checker.RequestPropertyRemovedId,
		Args:      []any{"roleAssignments/items/role", "application/json"},
		Level:     checker.WARN,
		Operation: "POST",
		Path:      "/api/roleMappings",
		Source:    load.NewSource("../data/required-properties/request-property-items-2.yaml"),
	})
	require.Equal(t, "removed request property 'roleAssignments/items/role' of media-type 'application/json'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: adding a new required property in response body is not breaking
func TestBreaking_RespBodyNewRequiredProperty(t *testing.T) {
	s1, err := open(getReqPropFile("response-new-base.json"))
	require.NoError(t, err)

	s2, err := open(getReqPropFile("response-new-revision.json"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.Empty(t, errs)
}

// BC: deleting a required property in response body is breaking: response-required-property-removed
func TestBreaking_RespBodyDeleteRequiredProperty(t *testing.T) {
	s1, err := open(getReqPropFile("response-new-revision.json"))
	require.NoError(t, err)

	s2, err := open(getReqPropFile("response-new-base.json"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ResponseRequiredPropertyRemovedId, errs[0].GetId())
	require.Equal(t, "removed required property 'helpAndSupport/title' from response status '200'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: adding a new required property under AllOf in response body is not breaking
func TestBreaking_RespBodyNewAllOfRequiredProperty(t *testing.T) {
	s1, err := open(getReqPropFile("response-allof-base.json"))
	require.NoError(t, err)

	s2, err := open(getReqPropFile("response-allof-revision.json"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.Empty(t, errs)
}

// BC: deleting a required property under AllOf in response body is breaking: response-required-property-removed
func TestBreaking_RespBodyDeleteAllOfRequiredProperty(t *testing.T) {
	s1, err := open(getReqPropFile("response-allof-revision.json"))
	require.NoError(t, err)

	s2, err := open(getReqPropFile("response-allof-base.json"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ResponseRequiredPropertyRemovedId, errs[0].GetId())
	require.Equal(t, "removed required property '/allOf[subschema #1]/bazqux' from response status '200'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: adding a new required read-only property in request body is not breaking
func TestBreaking_ReadOnlyNewRequiredProperty(t *testing.T) {
	s1, err := open(getReqPropFile("read-only-new-base.yaml"))
	require.NoError(t, err)

	s2, err := open(getReqPropFile("read-only-new-revision.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.Empty(t, errs)
}

// BC: changing an existing read-only property in request body to required is not breaking
func TestBreaking_ReadOnlyPropertyRequiredEnabled(t *testing.T) {
	s1, err := open(getReqPropFile("read-only-base.yaml"))
	require.NoError(t, err)

	s2, err := open(getReqPropFile("read-only-revision.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.Empty(t, errs)
}

// BC: deleting a required write-only property in response body is not breaking: request-property-removed
func TestBreaking_WriteOnlyDeleteRequiredProperty(t *testing.T) {
	s1, err := open(getReqPropFile("write-only-delete-base.yaml"))
	require.NoError(t, err)

	s2, err := open(getReqPropFile("write-only-delete-revision.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.RequestPropertyRemovedId, errs[0].GetId())
	require.Equal(t, checker.WARN, errs[0].GetLevel())
	require.Equal(t, "removed request property 'test' of media-type 'application/json'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: deleting a non-required non-write-only property in response body is breaking with warning: request-property-removed, response-optional-property-removed
func TestBreaking_WriteOnlyDeleteNonRequiredProperty(t *testing.T) {
	s1, err := open(getReqPropFile("write-only-delete-partial-base.yaml"))
	require.NoError(t, err)

	s2, err := open(getReqPropFile("write-only-delete-partial-revision.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 3)

	require.Equal(t, checker.RequestPropertyRemovedId, errs[0].GetId())
	require.Equal(t, checker.WARN, errs[0].GetLevel())
	require.Equal(t, "/api/atlas/v1.0/groups", errs[0].GetPath())
	require.Equal(t, "removed request property 'test' of media-type 'application/json'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))

	require.Equal(t, checker.WARN, errs[1].GetLevel())
	require.Equal(t, checker.ResponseOptionalPropertyRemovedId, errs[1].GetId())
	require.Equal(t, "/api/atlas/v1.0/groups", errs[1].GetPath())
	require.Equal(t, "removed optional property 'test' from response status '200'", errs[1].GetUncolorizedText(checker.NewDefaultLocalizer()))

	require.Equal(t, checker.WARN, errs[2].GetLevel())
	require.Equal(t, checker.ResponseOptionalPropertyRemovedId, errs[2].GetId())
	require.Equal(t, "/api/atlas/v1.0/groups/{groupId}", errs[2].GetPath())
	require.Equal(t, "removed optional property 'test' from response status '200'", errs[2].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: changing an existing write-only property in response body to optional is not breaking
func TestBreaking_WriteOnlyPropertyRequiredDisabled(t *testing.T) {
	s1, err := open(getReqPropFile("write-only-base.yaml"))
	require.NoError(t, err)

	s2, err := open(getReqPropFile("write-only-revision.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.Empty(t, errs)
}

// BC: changing an existing required property in response body to write-only is not breaking
func TestBreaking_RequiredPropertyWriteOnlyEnabled(t *testing.T) {
	s1, err := open(getReqPropFile("write-only-changed-base.yaml"))
	require.NoError(t, err)

	s2, err := open(getReqPropFile("write-only-changed-revision.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.Empty(t, errs)
}

// BC: changing an existing required property in response body to not-write-only is breaking: response-required-property-became-not-write-only
func TestBreaking_RequiredPropertyWriteOnlyDisabled(t *testing.T) {
	s1, err := open(getReqPropFile("write-only-changed-revision.yaml"))
	require.NoError(t, err)

	s2, err := open(getReqPropFile("write-only-changed-base.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 2)

	require.Equal(t, checker.ResponseRequiredPropertyBecameNonWriteOnlyId, errs[0].GetId())
	require.Equal(t, checker.WARN, errs[0].GetLevel())
	require.Equal(t, "/api/atlas/v1.0/groups", errs[0].GetPath())
	require.Equal(t, "response required property 'test' became not write-only for status '200'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))

	require.Equal(t, checker.ResponseRequiredPropertyBecameNonWriteOnlyId, errs[1].GetId())
	require.Equal(t, checker.WARN, errs[1].GetLevel())
	require.Equal(t, "/api/atlas/v1.0/groups/{groupId}", errs[1].GetPath())
	require.Equal(t, "response required property 'test' became not write-only for status '200'", errs[1].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: changing an existing property in request body to required is breaking: request-property-became-required
func TestBreaking_Body(t *testing.T) {
	s1, err := open(getReqPropFile("body1.yaml"))
	require.NoError(t, err)

	s2, err := open(getReqPropFile("body2.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.RequestPropertyBecameRequiredId, errs[0].GetId())
	require.Equal(t, []interface{}{"id", "application/json"}, errs[0].GetArgs())
	require.Equal(t, "request property 'id' of media-type 'application/json' became required", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: changing an existing property in request body items to required is breaking: request-property-became-required
func TestBreaking_Items(t *testing.T) {
	s1, err := open(getReqPropFile("items1.yaml"))
	require.NoError(t, err)

	s2, err := open(getReqPropFile("items2.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.RequestPropertyBecameRequiredId, errs[0].GetId())
	require.Equal(t, []interface{}{"/items/id", "application/json"}, errs[0].GetArgs())
	require.Equal(t, "request property '/items/id' of media-type 'application/json' became required", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: changing an existing property in request body items to required with a default value is not breaking
func TestBreaking_ItemsWithDefault(t *testing.T) {
	s1, err := open(getReqPropFile("items1.yaml"))
	require.NoError(t, err)

	s2, err := open(getReqPropFile("items3.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.Empty(t, errs)
}

// BC: changing an existing property in request body anyOf to required is breaking: request-property-became-required
func TestBreaking_AnyOf(t *testing.T) {
	s1, err := open(getReqPropFile("anyOf1.yaml"))
	require.NoError(t, err)

	s2, err := open(getReqPropFile("anyOf2.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.RequestPropertyBecameRequiredId, errs[0].GetId())
	require.Equal(t, "request property '/anyOf[subschema #1]/id' of media-type 'application/json' became required", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: changing an existing property under another property in request body to required is breaking: request-property-became-required
func TestBreaking_NestedProp(t *testing.T) {
	s1, err := open(getReqPropFile("nested-property1.yaml"))
	require.NoError(t, err)

	s2, err := open(getReqPropFile("nested-property2.yaml"))
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.RequestPropertyBecameRequiredId, errs[0].GetId())
	require.Equal(t, "request property 'id/userId' of media-type 'application/json' became required", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// BC: changing a response property to optional under AllOf, AnyOf or OneOf is breaking: response-property-became-optional
func TestBreaking_OneOf(t *testing.T) {
	s1, err := open("../data/x-of/base.json")
	require.NoError(t, err)

	s2, err := open("../data/x-of/revision.json")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibility(allChecksConfig(), d, osm)
	require.Len(t, errs, 3)

	require.Equal(t, checker.ResponsePropertyBecameOptionalId, errs[0].GetId())
	require.Equal(t, "property '/allOf[#/components/schemas/ProblemSchema]/changedProperty' became optional for response status '200'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))

	require.Equal(t, checker.ResponsePropertyBecameOptionalId, errs[1].GetId())
	require.Equal(t, "property '/anyOf[#/components/schemas/ProblemSchema]/changedProperty' became optional for response status '200'", errs[1].GetUncolorizedText(checker.NewDefaultLocalizer()))

	require.Equal(t, checker.ResponsePropertyBecameOptionalId, errs[2].GetId())
	require.Equal(t, "property '/oneOf[#/components/schemas/ProblemSchema]/changedProperty' became optional for response status '200'", errs[2].GetUncolorizedText(checker.NewDefaultLocalizer()))
}
