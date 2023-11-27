package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: adding 'oneOf' schema to the response body or response body property
func TestResponsePropertyOneOfAdded(t *testing.T) {
	s1, err := open("../data/checker/response_property_one_of_added_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_property_one_of_added_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePropertyOneOfUpdated), d, osm, checker.INFO)

	require.Len(t, errs, 3)

	require.ElementsMatch(t, []checker.ApiChange{
		{
			Id:          checker.ResponseBodyOneOfAddedId,
			Args:        []any{"Rabbit", "200"},
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "GET",
			Path:        "/pets",
			Source:      "../data/checker/response_property_one_of_added_revision.yaml",
			OperationId: "listPets",
		},
		{
			Id:          checker.ResponsePropertyOneOfAddedId,
			Args:        []any{"Breed3", "/oneOf[#/components/schemas/Dog]/breed", "200"},
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "GET",
			Path:        "/pets",
			Source:      "../data/checker/response_property_one_of_added_revision.yaml",
			OperationId: "listPets",
		},
		{
			Id:          checker.ResponsePropertyOneOfAddedId,
			Args:        []any{"RevisionSchema[1]:Dark brown types", "/oneOf[#/components/schemas/Fox]/breed", "200"},
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "GET",
			Path:        "/pets",
			Source:      "../data/checker/response_property_one_of_added_revision.yaml",
			OperationId: "listPets",
		}}, errs)
}

// CL: removing 'oneOf' schema from the response body or response body property
func TestResponsePropertyOneOfRemoved(t *testing.T) {
	s1, err := open("../data/checker/response_property_one_of_removed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_property_one_of_removed_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePropertyOneOfUpdated), d, osm, checker.INFO)

	require.Len(t, errs, 3)

	require.ElementsMatch(t, []checker.ApiChange{
		{
			Id:          checker.ResponseBodyOneOfRemovedId,
			Args:        []any{"Rabbit", "200"},
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "GET",
			Path:        "/pets",
			Source:      "../data/checker/response_property_one_of_removed_revision.yaml",
			OperationId: "listPets",
		},
		{
			Id:          checker.ResponsePropertyOneOfRemovedId,
			Args:        []any{"Breed3", "/oneOf[#/components/schemas/Dog]/breed", "200"},
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "GET",
			Path:        "/pets",
			Source:      "../data/checker/response_property_one_of_removed_revision.yaml",
			OperationId: "listPets",
		},
		{
			Id:          checker.ResponsePropertyOneOfRemovedId,
			Args:        []any{"BaseSchema[1]:Dark brown types", "/oneOf[#/components/schemas/Fox]/breed", "200"},
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "GET",
			Path:        "/pets",
			Source:      "../data/checker/response_property_one_of_removed_revision.yaml",
			OperationId: "listPets",
		}}, errs)
}
