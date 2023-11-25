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
			Text:        "added request discriminator",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_discriminator_added_revision.yaml",
			OperationId: "updatePets",
		},
		{
			Id:          checker.RequestPropertyDiscriminatorAddedId,
			Text:        "added discriminator to '/oneOf[#/components/schemas/Dog]/breed' request property",
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
			Text:        "removed request discriminator",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_discriminator_added_base.yaml",
			OperationId: "updatePets",
		},
		{
			Id:          checker.RequestPropertyDiscriminatorRemovedId,
			Text:        "removed discriminator from '/oneOf[#/components/schemas/Dog]/breed' request property",
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
			Text:        "request discriminator property name changed from 'petType' to 'petType2'",
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
			Text:        "request discriminator property name changed for '/oneOf[#/components/schemas/Dog]/breed' request property from 'name' to 'name2'",
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
			Text:        "added '[cats]' mapping keys to the request discriminator",
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
			Text:        "removed '[cat]' mapping keys from the request discriminator",
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
			Text:        "added '[breed1Code]' discriminator mapping keys to the '/oneOf[#/components/schemas/Dog]/breed' request property",
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
			Text:        "mapped value for discriminator key 'breed2' changed from '#/components/schemas/Breed2' to '#/components/schemas/Breed3' for '/oneOf[#/components/schemas/Dog]/breed' request property",
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
			Text:        "removed '[breed1]' discriminator mapping keys from the '/oneOf[#/components/schemas/Dog]/breed' request property",
			Args:        []any{utils.StringList{"breed1"}, "/oneOf[#/components/schemas/Dog]/breed"},
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_discriminator_mapping_changed.yaml",
			OperationId: "updatePets",
		}}, errs)
}
