package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
	"github.com/tufin/oasdiff/utils"
)

// CL: adding discriminator to the request body or request body property: request-body-discriminator-added, request-property-discriminator-added
func TestRequestDiscriminatorUpdatedCheckAdded(t *testing.T) {
	s1, err := open("../data/checker/request_property_discriminator_added_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_discriminator_added_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestDiscriminatorUpdatedCheck), d, osm, checker.INFO)

	require.Len(t, errs, 2)

	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestBodyDiscriminatorAddedId,
		Args:        []any{"application/json"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/pets",
		Source:      load.NewSource("../data/checker/request_property_discriminator_added_revision.yaml"),
		OperationId: "updatePets",
	}, errs[0])
	require.Equal(t, "added discriminator to media-type 'application/json' of request body", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))

	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestPropertyDiscriminatorAddedId,
		Args:        []any{"/oneOf[#/components/schemas/Dog]/breed", "application/json"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/pets",
		Source:      load.NewSource("../data/checker/request_property_discriminator_added_revision.yaml"),
		OperationId: "updatePets",
	}, errs[1])
	require.Equal(t, "added discriminator to property '/oneOf[#/components/schemas/Dog]/breed' of media-type 'application/json' of request body", errs[1].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// CL: removing discriminator from the request body or request body property: request-body-discriminator-removed, request-property-discriminator-removed
func TestRequestDiscriminatorUpdatedCheckRemoved(t *testing.T) {
	s1, err := open("../data/checker/request_property_discriminator_added_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_discriminator_added_base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestDiscriminatorUpdatedCheck), d, osm, checker.INFO)

	require.Len(t, errs, 2)

	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestBodyDiscriminatorRemovedId,
		Args:        []any{"application/json"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/pets",
		Source:      load.NewSource("../data/checker/request_property_discriminator_added_base.yaml"),
		OperationId: "updatePets",
	}, errs[0])
	require.Equal(t, "removed discriminator from media-type 'application/json' of request body", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))

	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestPropertyDiscriminatorRemovedId,
		Args:        []any{"/oneOf[#/components/schemas/Dog]/breed", "application/json"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/pets",
		Source:      load.NewSource("../data/checker/request_property_discriminator_added_base.yaml"),
		OperationId: "updatePets",
	}, errs[1])
	require.Equal(t, "removed discriminator from property '/oneOf[#/components/schemas/Dog]/breed' of media-type 'application/json' of request body", errs[1].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// CL: changing discriminator propertyName in the request body or request body property
func TestRequestDiscriminatorUpdatedCheckPropertyNameChanging(t *testing.T) {
	s1, err := open("../data/checker/request_property_discriminator_added_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_discriminator_added_property_name_changed.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestDiscriminatorUpdatedCheck), d, osm, checker.INFO)

	require.Len(t, errs, 2)

	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestBodyDiscriminatorPropertyNameChangedId,
		Args:        []any{"petType", "petType2", "application/json"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/pets",
		Source:      load.NewSource("../data/checker/request_property_discriminator_added_property_name_changed.yaml"),
		OperationId: "updatePets",
	}, errs[0])
	require.Equal(t, "discriminator property name changed from 'petType' to 'petType2' for media-type 'application/json' of request body", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))

	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestPropertyDiscriminatorPropertyNameChangedId,
		Args:        []any{"/oneOf[#/components/schemas/Dog]/breed", "name", "name2", "application/json"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/pets",
		Source:      load.NewSource("../data/checker/request_property_discriminator_added_property_name_changed.yaml"),
		OperationId: "updatePets",
	}, errs[1])
	require.Equal(t, "request discriminator property name changed from '/oneOf[#/components/schemas/Dog]/breed' to 'name' for 'name2' request property of media-type 'application/json'", errs[1].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// CL: changing discriminator mapping in the request body or request body property
func TestRequestDiscriminatorUpdatedCheckMappingChanging(t *testing.T) {
	s1, err := open("../data/checker/request_property_discriminator_added_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_discriminator_mapping_changed.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestDiscriminatorUpdatedCheck), d, osm, checker.INFO)

	require.Len(t, errs, 5)

	require.ElementsMatch(t, []checker.ApiChange{
		{
			Id:          checker.RequestBodyDiscriminatorMappingAddedId,
			Args:        []any{utils.StringList{"cats"}},
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      load.NewSource("../data/checker/request_property_discriminator_mapping_changed.yaml"),
			OperationId: "updatePets",
		},
		{
			Id:          checker.RequestBodyDiscriminatorMappingDeletedId,
			Args:        []any{utils.StringList{"cat"}},
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      load.NewSource("../data/checker/request_property_discriminator_mapping_changed.yaml"),
			OperationId: "updatePets",
		},
		{
			Id:          checker.RequestPropertyDiscriminatorMappingAddedId,
			Args:        []any{utils.StringList{"breed1Code"}, "/oneOf[#/components/schemas/Dog]/breed"},
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      load.NewSource("../data/checker/request_property_discriminator_mapping_changed.yaml"),
			OperationId: "updatePets",
		},
		{
			Id:          checker.RequestPropertyDiscriminatorMappingChangedId,
			Args:        []any{"breed2", "#/components/schemas/Breed2", "#/components/schemas/Breed3", "/oneOf[#/components/schemas/Dog]/breed"},
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      load.NewSource("../data/checker/request_property_discriminator_mapping_changed.yaml"),
			OperationId: "updatePets",
		},
		{
			Id:          checker.RequestPropertyDiscriminatorMappingDeletedId,
			Args:        []any{utils.StringList{"breed1"}, "/oneOf[#/components/schemas/Dog]/breed"},
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      load.NewSource("../data/checker/request_property_discriminator_mapping_changed.yaml"),
			OperationId: "updatePets",
		}}, errs)
}
