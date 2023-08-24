package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: adding 'allOf' subschema to the response body or response body property
func TestResponsePropertyAllOfAdded(t *testing.T) {
	s1, err := open("../data/checker/response_property_all_of_added_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_property_all_of_added_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePropertyAllOfUpdated), d, osm, checker.INFO)

	require.Len(t, errs, 2)

	require.ElementsMatch(t, []checker.ApiChange{
		{
			Id:          "response-body-all-of-added",
			Text:        "added 'Rabbit' to the response body 'allOf' list for the response status 200",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "GET",
			Path:        "/pets",
			Source:      "../data/checker/response_property_all_of_added_revision.yaml",
			OperationId: "listPets",
		},
		{
			Id:          "response-property-all-of-added",
			Text:        "added 'Breed3' to the '/allOf[#/components/schemas/Dog]/breed' response property 'allOf' list for the response status 200",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "GET",
			Path:        "/pets",
			Source:      "../data/checker/response_property_all_of_added_revision.yaml",
			OperationId: "listPets",
		}}, errs)
}

// CL: removing 'allOf' subschema from the response body or response body property
func TestResponsePropertyAllOfRemoved(t *testing.T) {
	s1, err := open("../data/checker/response_property_all_of_removed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_property_all_of_removed_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePropertyAllOfUpdated), d, osm, checker.INFO)

	require.Len(t, errs, 2)

	require.ElementsMatch(t, []checker.ApiChange{
		{
			Id:          "response-body-all-of-removed",
			Text:        "removed 'Rabbit' from the response body 'allOf' list for the response status 200",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "GET",
			Path:        "/pets",
			Source:      "../data/checker/response_property_all_of_removed_revision.yaml",
			OperationId: "listPets",
		},
		{
			Id:          "response-property-all-of-removed",
			Text:        "removed 'Breed3' from the '/allOf[#/components/schemas/Dog]/breed' response property 'allOf' list for the response status 200",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "GET",
			Path:        "/pets",
			Source:      "../data/checker/response_property_all_of_removed_revision.yaml",
			OperationId: "listPets",
		}}, errs)
}
