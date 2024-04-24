package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

// CL: changing a response schema type
func TestResponseSchemaTypeChangedCheck(t *testing.T) {
	s1, err := open("../data/checker/response_schema_type_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_schema_type_changed_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePropertyTypeChangedCheck), d, osm, checker.ERR)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.ResponseBodyTypeChangedId,
		Args:        []any{"string", "", "object", "", "200"},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/response_schema_type_changed_revision.yaml"),
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: changing a response property schema type
func TestResponsePropertyTypeChangedCheck(t *testing.T) {
	s1, err := open("../data/checker/response_schema_type_changed_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_schema_type_changed_revision.yaml")
	require.NoError(t, err)

	s2.Spec.Paths.Value("/api/v1.0/groups").Post.Responses.Value("200").Value.Content["application/json"].Schema.Value.Properties["data"].Value.Properties["name"].Value.Type = "integer"

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePropertyTypeChangedCheck), d, osm, checker.ERR)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.ResponsePropertyTypeChangedId,
		Args:        []any{"data/name", "string", "", "integer", "", "200"},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/response_schema_type_changed_revision.yaml"),
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: changing a response property schema format
func TestResponsePropertyFormatChangedCheck(t *testing.T) {
	s1, err := open("../data/checker/response_schema_format_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_schema_format_changed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths.Value("/api/v1.0/groups").Post.Responses.Value("200").Value.Content["application/json"].Schema.Value.Properties["data"].Value.Properties["name"].Value.Format = "uuid"

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePropertyTypeChangedCheck), d, osm, checker.ERR)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.ResponsePropertyTypeChangedId,
		Args:        []any{"data/name", "string", "hostname", "string", "uuid", "200"},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/response_schema_format_changed_base.yaml"),
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: changing properties of subschemas under allOf
func TestResponsePropertyAnyOfModified(t *testing.T) {
	s1, err := open("../data/checker/response_property_any_of_complex_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_property_any_of_complex_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePropertyTypeChangedCheck), d, osm, checker.INFO)

	require.Len(t, errs, 3)
	require.ElementsMatch(t, []checker.ApiChange{
		{
			Id:          checker.ResponsePropertyTypeChangedId,
			Args:        []any{"/anyOf[#/components/schemas/Dog]/breed/anyOf[#/components/schemas/Breed2]/name", "string", "", "number", "", "200"},
			Level:       checker.ERR,
			Operation:   "GET",
			Path:        "/pets",
			Source:      load.NewSource("../data/checker/response_property_any_of_complex_revision.yaml"),
			OperationId: "listPets",
		},
		{
			Id:          checker.ResponsePropertyTypeChangedId,
			Args:        []any{"/anyOf[subschema #3: Rabbit]/", "string", "", "number", "", "200"},
			Level:       checker.ERR,
			Operation:   "GET",
			Path:        "/pets",
			Source:      load.NewSource("../data/checker/response_property_any_of_complex_revision.yaml"),
			OperationId: "listPets",
		},
		{
			Id:          checker.ResponsePropertyTypeChangedId,
			Args:        []any{"/anyOf[subschema #4 -> subschema #5]/", "string", "", "number", "", "200"},
			Level:       checker.ERR,
			Operation:   "GET",
			Path:        "/pets",
			Source:      load.NewSource("../data/checker/response_property_any_of_complex_revision.yaml"),
			OperationId: "listPets",
		}}, errs)
}
