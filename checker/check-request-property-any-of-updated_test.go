package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: adding 'anyOf' schema to the request body or request body property
func TestRequestPropertyAnyOfAdded(t *testing.T) {
	s1, err := open("../data/checker/request_property_any_of_added_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_any_of_added_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyAnyOfUpdatedCheck), d, osm, checker.INFO)

	require.Len(t, errs, 2)

	require.ElementsMatch(t, []checker.ApiChange{
		{
			Id:          checker.RequestBodyAnyOfAddedId,
			Text:        "added 'Rabbit' to the request body 'anyOf' list",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_any_of_added_revision.yaml",
			OperationId: "updatePets",
		},
		{
			Id:          "request-property-any-of-added",
			Text:        "added 'Breed3' to the '/anyOf[#/components/schemas/Dog]/breed' request property 'anyOf' list",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_any_of_added_revision.yaml",
			OperationId: "updatePets",
		}}, errs)
}

// CL: removing 'anyOf' schema from the request body or request body property
func TestRequestPropertyAnyOfRemoved(t *testing.T) {
	s1, err := open("../data/checker/request_property_any_of_removed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_any_of_removed_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyAnyOfUpdatedCheck), d, osm, checker.INFO)

	require.Len(t, errs, 2)

	require.ElementsMatch(t, []checker.ApiChange{
		{
			Id:          checker.RequestBodyAnyOfRemovedId,
			Text:        "removed 'Rabbit' from the request body 'anyOf' list",
			Comment:     "",
			Level:       checker.ERR,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_any_of_removed_revision.yaml",
			OperationId: "updatePets",
		},
		{
			Id:          "request-property-any-of-removed",
			Text:        "removed 'Breed3' from the '/anyOf[#/components/schemas/Dog]/breed' request property 'anyOf' list",
			Comment:     "",
			Level:       checker.ERR,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_any_of_removed_revision.yaml",
			OperationId: "updatePets",
		}}, errs)
}
