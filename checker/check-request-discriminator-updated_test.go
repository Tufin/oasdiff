package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: adding schema discriminator to the request body or request body property
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
			Id:          "request-body-discriminator-added",
			Text:        "added request schema discriminator",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_discriminator_added_revision.yaml",
			OperationId: "updatePets",
		},
		{
			Id:          "request-property-discriminator-added",
			Text:        "added discriminator to '/oneOf[#/components/schemas/Dog]/breed' request property",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_discriminator_added_revision.yaml",
			OperationId: "updatePets",
		}}, errs)
}

// CL: removing schema discriminator from the request body or request body property
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
			Id:          "request-body-discriminator-removed",
			Text:        "removed request schema discriminator",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_discriminator_added_base.yaml",
			OperationId: "updatePets",
		},
		{
			Id:          "request-property-discriminator-removed",
			Text:        "removed discriminator to '/oneOf[#/components/schemas/Dog]/breed' request property",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_discriminator_added_base.yaml",
			OperationId: "updatePets",
		}}, errs)
}

// CL: changing schema discriminator propertyName in the request body or request body property
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
			Id:          "request-body-discriminator-property-name-changed",
			Text:        "request schema discriminator property name changed from 'petType' to 'petType2'",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_discriminator_added_property_name_changed.yaml",
			OperationId: "updatePets",
		},
		{
			Id:          "request-property-discriminator-property-name-changed",
			Text:        "request schema discriminator property name changed for '/oneOf[#/components/schemas/Dog]/breed' request property from 'name' to 'name2'",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_discriminator_added_property_name_changed.yaml",
			OperationId: "updatePets",
		}}, errs)
}

// CL: changing schema discriminator mapping in the request body or request body property
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
			Id:          "request-body-discriminator-mapping-added",
			Text:        "added '[cats]' mapping keys to the request schema discriminator",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_discriminator_mapping_changed.yaml",
			OperationId: "updatePets",
		},
		{
			Id:          "request-body-discriminator-mapping-deleted",
			Text:        "removed '[cat]' mapping keys from the request schema discriminator",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_discriminator_mapping_changed.yaml",
			OperationId: "updatePets",
		},
		{
			Id:          "request-property-discriminator-mapping-added",
			Text:        "added '[breed1Code]' discriminator mapping keys to the '/oneOf[#/components/schemas/Dog]/breed' request property",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_discriminator_mapping_changed.yaml",
			OperationId: "updatePets",
		},
		{
			Id:          "request-property-discriminator-mapping-changed",
			Text:        "mapped value for discriminator key 'breed2' changed from '#/components/schemas/Breed2' to '#/components/schemas/Breed3' for '/oneOf[#/components/schemas/Dog]/breed' request property",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_discriminator_mapping_changed.yaml",
			OperationId: "updatePets",
		},
		{
			Id:          "request-property-discriminator-mapping-deleted",
			Text:        "removed '[breed1]' discriminator mapping keys from the '/oneOf[#/components/schemas/Dog]/breed' request property",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_discriminator_mapping_changed.yaml",
			OperationId: "updatePets",
		}}, errs)
}
