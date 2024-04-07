package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
	"github.com/tufin/oasdiff/utils"
)

// CL: adding discriminator to the response body or response property
func TestResponseDiscriminatorUpdatedCheckAdded(t *testing.T) {
	s1, err := open("../data/checker/response_property_discriminator_added_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_property_discriminator_added_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseDiscriminatorUpdatedCheck), d, osm, checker.INFO)

	require.Len(t, errs, 2)

	require.ElementsMatch(t, []checker.ApiChange{
		{
			Id:          checker.ResponseBodyDiscriminatorAddedId,
			Args:        []any{"200"},
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      load.NewSource("../data/checker/response_property_discriminator_added_revision.yaml"),
			OperationId: "updatePets",
		},
		{
			Id:          checker.ResponsePropertyDiscriminatorAddedId,
			Args:        []any{"/oneOf[#/components/schemas/Dog]/breed", "200"},
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      load.NewSource("../data/checker/response_property_discriminator_added_revision.yaml"),
			OperationId: "updatePets",
		}}, errs)
}

// CL: removing discriminator from the response body or response property
func TestResponseDiscriminatorUpdatedCheckRemoved(t *testing.T) {
	s1, err := open("../data/checker/response_property_discriminator_added_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_property_discriminator_added_base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseDiscriminatorUpdatedCheck), d, osm, checker.INFO)

	require.Len(t, errs, 2)

	require.ElementsMatch(t, []checker.ApiChange{
		{
			Id:          checker.ResponseBodyDiscriminatorRemovedId,
			Args:        []any{"200"},
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      load.NewSource("../data/checker/response_property_discriminator_added_base.yaml"),
			OperationId: "updatePets",
		},
		{
			Id:          checker.ResponsePropertyDiscriminatorRemovedId,
			Args:        []any{"/oneOf[#/components/schemas/Dog]/breed", "200"},
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      load.NewSource("../data/checker/response_property_discriminator_added_base.yaml"),
			OperationId: "updatePets",
		}}, errs)
}

// CL: changing discriminator propertyName in the response body or response property
func TestResponseDiscriminatorUpdatedCheckPropertyNameChanging(t *testing.T) {
	s1, err := open("../data/checker/response_property_discriminator_added_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_property_discriminator_added_property_name_changed.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseDiscriminatorUpdatedCheck), d, osm, checker.INFO)

	require.Len(t, errs, 2)

	require.ElementsMatch(t, []checker.ApiChange{
		{
			Id:          checker.ResponseBodyDiscriminatorPropertyNameChangedId,
			Args:        []any{"petType", "petType2", "200"},
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      load.NewSource("../data/checker/response_property_discriminator_added_property_name_changed.yaml"),
			OperationId: "updatePets",
		},
		{
			Id:          checker.ResponsePropertyDiscriminatorPropertyNameChangedId,
			Args:        []any{"/oneOf[#/components/schemas/Dog]/breed", "name", "name2", "200"},
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      load.NewSource("../data/checker/response_property_discriminator_added_property_name_changed.yaml"),
			OperationId: "updatePets",
		}}, errs)
}

// CL: changing discriminator mapping in the response body or response property
func TestResponseDiscriminatorUpdatedCheckMappingChanging(t *testing.T) {
	s1, err := open("../data/checker/response_property_discriminator_added_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_property_discriminator_mapping_changed.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseDiscriminatorUpdatedCheck), d, osm, checker.INFO)

	require.Len(t, errs, 5)

	require.ElementsMatch(t, []checker.ApiChange{
		{
			Id:          checker.ResponseBodyDiscriminatorMappingAddedId,
			Args:        []any{utils.StringList{"cats"}, "200"},
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      load.NewSource("../data/checker/response_property_discriminator_mapping_changed.yaml"),
			OperationId: "updatePets",
		},
		{
			Id:          checker.ResponseBodyDiscriminatorMappingDeletedId,
			Args:        []any{utils.StringList{"cat"}, "200"},
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      load.NewSource("../data/checker/response_property_discriminator_mapping_changed.yaml"),
			OperationId: "updatePets",
		},
		{
			Id:          checker.ResponsePropertyDiscriminatorMappingAddedId,
			Args:        []any{utils.StringList{"breed1Code"}, "/oneOf[#/components/schemas/Dog]/breed", "200"},
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      load.NewSource("../data/checker/response_property_discriminator_mapping_changed.yaml"),
			OperationId: "updatePets",
		},
		{
			Id:          checker.ResponsePropertyDiscriminatorMappingChangedId,
			Args:        []any{"breed2", "#/components/schemas/Breed2", "#/components/schemas/Breed3", "/oneOf[#/components/schemas/Dog]/breed", "200"},
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      load.NewSource("../data/checker/response_property_discriminator_mapping_changed.yaml"),
			OperationId: "updatePets",
		},
		{
			Id:          checker.ResponsePropertyDiscriminatorMappingDeletedId,
			Args:        []any{utils.StringList{"breed1"}, "/oneOf[#/components/schemas/Dog]/breed", "200"},
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      load.NewSource("../data/checker/response_property_discriminator_mapping_changed.yaml"),
			OperationId: "updatePets",
		}}, errs)
}
