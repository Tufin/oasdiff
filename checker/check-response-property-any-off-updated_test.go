package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: adding 'anyOf' schema to the response body or response body property
func TestResponsePropertyAnyOffAdded(t *testing.T) {
	s1, err := open("../data/checker/response_property_any_of_added_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_property_any_of_added_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePropertyAnyOffUpdated), d, osm, checker.INFO)

	require.Len(t, errs, 2)

	require.Equal(t, checker.ApiChange{
		Id:          "response-body-any-of-added",
		Text:        "added ''Rabbit'' to the response body 'anyOf' list for the response status 200",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "GET",
		Path:        "/pets",
		Source:      "../data/checker/response_property_any_of_added_revision.yaml",
		OperationId: "listPets",
	}, errs[0])

	require.Equal(t, checker.ApiChange{
		Id:          "response-property-any-of-added",
		Text:        "added ''Breed3'' to the '/anyOf[#/components/schemas/Dog]/breed' response property 'anyOf' list for the response status 200",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "GET",
		Path:        "/pets",
		Source:      "../data/checker/response_property_any_of_added_revision.yaml",
		OperationId: "listPets",
	}, errs[1])
}

// CL: removing 'anyOf' schema from the response body or response body property
func TestResponsePropertyAnyOffRemoved(t *testing.T) {
	s1, err := open("../data/checker/response_property_any_of_removed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_property_any_of_removed_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePropertyAnyOffUpdated), d, osm, checker.INFO)

	require.Len(t, errs, 2)

	require.Equal(t, checker.ApiChange{
		Id:          "response-body-any-of-removed",
		Text:        "removed ''Rabbit'' from the response body 'anyOf' list for the response status 200",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "GET",
		Path:        "/pets",
		Source:      "../data/checker/response_property_any_of_removed_revision.yaml",
		OperationId: "listPets",
	}, errs[0])

	require.Equal(t, checker.ApiChange{
		Id:          "response-property-any-of-removed",
		Text:        "removed ''Breed3'' from the '/anyOf[#/components/schemas/Dog]/breed' response property 'anyOf' list for the response status 200",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "GET",
		Path:        "/pets",
		Source:      "../data/checker/response_property_any_of_removed_revision.yaml",
		OperationId: "listPets",
	}, errs[1])
}
