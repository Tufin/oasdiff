package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

// CL: adding 'anyOf' schema to the response body or response body property
func TestResponsePropertyAnyOfAdded(t *testing.T) {
	s1, err := open("../data/checker/response_property_any_of_added_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_property_any_of_added_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePropertyAnyOfUpdatedCheck), d, osm, checker.INFO)

	require.Len(t, errs, 2)

	require.ElementsMatch(t, []checker.ApiChange{
		{
			Id:          checker.ResponseBodyAnyOfAddedId,
			Args:        []any{"#/components/schemas/Rabbit", "200"},
			Level:       checker.INFO,
			Operation:   "GET",
			Path:        "/pets",
			Source:      load.NewSource("../data/checker/response_property_any_of_added_revision.yaml"),
			OperationId: "listPets",
		},
		{
			Id:          checker.ResponsePropertyAnyOfAddedId,
			Args:        []any{"#/components/schemas/Breed3", "/anyOf[#/components/schemas/Dog]/breed", "200"},
			Level:       checker.INFO,
			Operation:   "GET",
			Path:        "/pets",
			Source:      load.NewSource("../data/checker/response_property_any_of_added_revision.yaml"),
			OperationId: "listPets",
		}}, errs)
}

// CL: removing 'anyOf' schema from the response body or response body property
func TestResponsePropertyAnyOfRemoved(t *testing.T) {
	s1, err := open("../data/checker/response_property_any_of_removed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_property_any_of_removed_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePropertyAnyOfUpdatedCheck), d, osm, checker.INFO)

	require.Len(t, errs, 2)

	require.ElementsMatch(t, []checker.ApiChange{
		{
			Id:          checker.ResponseBodyAnyOfRemovedId,
			Args:        []any{"#/components/schemas/Rabbit", "200"},
			Level:       checker.INFO,
			Operation:   "GET",
			Path:        "/pets",
			Source:      load.NewSource("../data/checker/response_property_any_of_removed_revision.yaml"),
			OperationId: "listPets",
		},
		{
			Id:          checker.ResponsePropertyAnyOfRemovedId,
			Args:        []any{"#/components/schemas/Breed3", "/anyOf[#/components/schemas/Dog]/breed", "200"},
			Level:       checker.INFO,
			Operation:   "GET",
			Path:        "/pets",
			Source:      load.NewSource("../data/checker/response_property_any_of_removed_revision.yaml"),
			OperationId: "listPets",
		}}, errs)
}
