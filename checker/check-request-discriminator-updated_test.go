package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/utils"
)

// CL: adding discriminator to the request body or request body property
func TestRequestDiscriminatorUpdatedCheckAdded(t *testing.T) {
	s1, err := open("../data/checker/request_property_discriminator_added_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_discriminator_added_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestDiscriminatorUpdatedCheck), d, osm, checker.INFO)

	require.Len(t, errs, 2)

	require.ElementsMatch(t, []checker.ApiChange{
		{
			Id:          checker.RequestBodyDiscriminatorAddedId,
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_discriminator_added_revision.yaml",
			OperationId: "updatePets",
		},
		{
			Id:          checker.RequestPropertyDiscriminatorAddedId,
			Args:        []any{"/oneOf[#/components/schemas/Dog]/breed"},
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_discriminator_added_revision.yaml",
			OperationId: "updatePets",
		}}, errs)
}

// CL: removing discriminator from the request body or request body property
func TestRequestDiscriminatorUpdatedCheckRemoved(t *testing.T) {
	s1, err := open("../data/checker/request_property_discriminator_added_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_discriminator_added_base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestDiscriminatorUpdatedCheck), d, osm, checker.INFO)

	require.Len(t, errs, 2)

	require.ElementsMatch(t, []checker.ApiChange{
		{
			Id:          checker.RequestBodyDiscriminatorRemovedId,
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_discriminator_added_base.yaml",
			OperationId: "updatePets",
		},
		{
			Id:          checker.RequestPropertyDiscriminatorRemovedId,
			Args:        []any{"/oneOf[#/components/schemas/Dog]/breed"},
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_discriminator_added_base.yaml",
			OperationId: "updatePets",
		}}, errs)
}

// CL: changing discriminator propertyName in the request body or request body property
func TestRequestDiscriminatorUpdatedCheckPropertyNameChanging(t *testing.T) {
	s1, err := open("../data/checker/request_property_discriminator_added_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_discriminator_added_property_name_changed.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestDiscriminatorUpdatedCheck), d, osm, checker.INFO)

	require.Len(t, errs, 2)

	require.ElementsMatch(t, []checker.ApiChange{
		{
			Id:          checker.RequestBodyDiscriminatorPropertyNameChangedId,
			Args:        []any{"petType", "petType2"},
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_discriminator_added_property_name_changed.yaml",
			OperationId: "updatePets",
		},
		{
			Id:          checker.RequestPropertyDiscriminatorPropertyNameChangedId,
			Args:        []any{"/oneOf[#/components/schemas/Dog]/breed", "name", "name2"},
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_discriminator_added_property_name_changed.yaml",
			OperationId: "updatePets",
		}}, errs)
}

// CL: changing discriminator mapping in the request body or request body property
func TestRequestDiscriminatorUpdatedCheckMappingChanging(t *testing.T) {
	s1, err := open("../data/checker/request_property_discriminator_added_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_discriminator_mapping_changed.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestDiscriminatorUpdatedCheck), d, osm, checker.INFO)

	require.Len(t, errs, 5)

	require.ElementsMatch(t, []checker.ApiChange{
		{
			Id:          checker.RequestBodyDiscriminatorMappingAddedId,
			Args:        []any{utils.StringList{"cats"}},
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_discriminator_mapping_changed.yaml",
			OperationId: "updatePets",
		},
		{
			Id:          checker.RequestBodyDiscriminatorMappingDeletedId,
			Args:        []any{utils.StringList{"cat"}},
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_discriminator_mapping_changed.yaml",
			OperationId: "updatePets",
		},
		{
			Id:          checker.RequestPropertyDiscriminatorMappingAddedId,
			Args:        []any{utils.StringList{"breed1Code"}, "/oneOf[#/components/schemas/Dog]/breed"},
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_discriminator_mapping_changed.yaml",
			OperationId: "updatePets",
		},
		{
			Id:          checker.RequestPropertyDiscriminatorMappingChangedId,
			Args:        []any{"breed2", "#/components/schemas/Breed2", "#/components/schemas/Breed3", "/oneOf[#/components/schemas/Dog]/breed"},
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_discriminator_mapping_changed.yaml",
			OperationId: "updatePets",
		},
		{
			Id:          checker.RequestPropertyDiscriminatorMappingDeletedId,
			Args:        []any{utils.StringList{"breed1"}, "/oneOf[#/components/schemas/Dog]/breed"},
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_discriminator_mapping_changed.yaml",
			OperationId: "updatePets",
		}}, errs)
}
