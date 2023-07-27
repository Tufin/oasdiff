package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: adding schema discriminator to the response body or response property
func TestResponseDiscriminatorUpdatedCheckAdded(t *testing.T) {
	s1, err := open("../data/checker/response_property_discriminator_added_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_property_discriminator_added_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseDiscriminatorUpdatedCheck), d, osm, checker.INFO)

	require.Len(t, errs, 2)

	require.ElementsMatch(t, []checker.ApiChange{
		{
			Id:          "response-body-discriminator-added",
			Text:        "added response discriminator for the response status 200",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/response_property_discriminator_added_revision.yaml",
			OperationId: "updatePets",
		},
		{
			Id:          "response-property-discriminator-added",
			Text:        "added discriminator to '/oneOf[#/components/schemas/Dog]/breed' response property for the response status 200",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/response_property_discriminator_added_revision.yaml",
			OperationId: "updatePets",
		}}, errs)
}

// CL: removing schema discriminator from the response body or response property
func TestResponseDiscriminatorUpdatedCheckRemoved(t *testing.T) {
	s1, err := open("../data/checker/response_property_discriminator_added_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_property_discriminator_added_base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseDiscriminatorUpdatedCheck), d, osm, checker.INFO)

	require.Len(t, errs, 2)

	require.ElementsMatch(t, []checker.ApiChange{
		{
			Id:          "response-body-discriminator-removed",
			Text:        "removed response discriminator for the response status 200",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/response_property_discriminator_added_base.yaml",
			OperationId: "updatePets",
		},
		{
			Id:          "response-property-discriminator-removed",
			Text:        "removed discriminator to '/oneOf[#/components/schemas/Dog]/breed' response property for the response status 200",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/response_property_discriminator_added_base.yaml",
			OperationId: "updatePets",
		}}, errs)
}

// CL: changing schema discriminator propertyName in the response body or response property
func TestResponseDiscriminatorUpdatedCheckPropertyNameChanging(t *testing.T) {
	s1, err := open("../data/checker/response_property_discriminator_added_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_property_discriminator_added_property_name_changed.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseDiscriminatorUpdatedCheck), d, osm, checker.INFO)

	require.Len(t, errs, 2)

	require.ElementsMatch(t, []checker.ApiChange{
		{
			Id:          "response-body-discriminator-property-name-changed",
			Text:        "response discriminator property name changed from 'petType' to 'petType2' for the response status 200",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/response_property_discriminator_added_property_name_changed.yaml",
			OperationId: "updatePets",
		},
		{
			Id:          "response-property-discriminator-property-name-changed",
			Text:        "response discriminator property name changed for '/oneOf[#/components/schemas/Dog]/breed' response property from 'name' to 'name2' for the response status 200",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/response_property_discriminator_added_property_name_changed.yaml",
			OperationId: "updatePets",
		}}, errs)
}

// CL: changing schema discriminator mapping in the response body or response property
func TestResponseDiscriminatorUpdatedCheckMappingChanging(t *testing.T) {
	s1, err := open("../data/checker/response_property_discriminator_added_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_property_discriminator_mapping_changed.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseDiscriminatorUpdatedCheck), d, osm, checker.INFO)

	require.Len(t, errs, 5)

	require.ElementsMatch(t, []checker.ApiChange{
		{
			Id:          "response-body-discriminator-mapping-added",
			Text:        "added '[cats]' mapping keys to the response discriminator for the response status 200",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/response_property_discriminator_mapping_changed.yaml",
			OperationId: "updatePets",
		},
		{
			Id:          "response-body-discriminator-mapping-deleted",
			Text:        "removed '[cat]' mapping keys from the response discriminator for the response status 200",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/response_property_discriminator_mapping_changed.yaml",
			OperationId: "updatePets",
		},
		{
			Id:          "response-property-discriminator-mapping-added",
			Text:        "added '[breed1Code]' discriminator mapping keys to the '/oneOf[#/components/schemas/Dog]/breed' response property for the response status 200",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/response_property_discriminator_mapping_changed.yaml",
			OperationId: "updatePets",
		},
		{
			Id:          "response-property-discriminator-mapping-changed",
			Text:        "mapped value for discriminator key 'breed2' changed from '#/components/schemas/Breed2' to '#/components/schemas/Breed3' for '/oneOf[#/components/schemas/Dog]/breed' response property for the response status 200",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/response_property_discriminator_mapping_changed.yaml",
			OperationId: "updatePets",
		},
		{
			Id:          "response-property-discriminator-mapping-deleted",
			Text:        "removed '[breed1]' discriminator mapping keys from the '/oneOf[#/components/schemas/Dog]/breed' response property for the response status 200",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/response_property_discriminator_mapping_changed.yaml",
			OperationId: "updatePets",
		}}, errs)
}
