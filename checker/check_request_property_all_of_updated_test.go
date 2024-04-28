package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

// CL: adding 'allOf' subschema to the request body or request body property
func TestRequestPropertyAllOfAdded(t *testing.T) {
	s1, err := open("../data/checker/request_property_all_of_added_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_all_of_added_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyAllOfUpdatedCheck), d, osm, checker.INFO)

	require.Len(t, errs, 2)

	require.ElementsMatch(t, []checker.ApiChange{
		{
			Id:          checker.RequestBodyAllOfAddedId,
			Args:        []any{"#/components/schemas/Rabbit"},
			Level:       checker.ERR,
			Operation:   "POST",
			Path:        "/pets",
			Source:      load.NewSource("../data/checker/request_property_all_of_added_revision.yaml"),
			OperationId: "updatePets",
		},
		{
			Id:          checker.RequestPropertyAllOfAddedId,
			Args:        []any{"#/components/schemas/Breed3", "/allOf[#/components/schemas/Dog]/breed"},
			Level:       checker.ERR,
			Operation:   "POST",
			Path:        "/pets",
			Source:      load.NewSource("../data/checker/request_property_all_of_added_revision.yaml"),
			OperationId: "updatePets",
		}}, errs)
}

// CL: removing 'allOf' subschema from the request body or request body property
func TestRequestPropertyAllOfRemoved(t *testing.T) {
	s1, err := open("../data/checker/request_property_all_of_removed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_all_of_removed_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyAllOfUpdatedCheck), d, osm, checker.INFO)

	require.Len(t, errs, 2)

	require.ElementsMatch(t, []checker.ApiChange{
		{
			Id:          checker.RequestBodyAllOfRemovedId,
			Args:        []any{"#/components/schemas/Rabbit"},
			Level:       checker.WARN,
			Operation:   "POST",
			Path:        "/pets",
			Source:      load.NewSource("../data/checker/request_property_all_of_removed_revision.yaml"),
			OperationId: "updatePets",
		},
		{
			Id:          checker.RequestPropertyAllOfRemovedId,
			Args:        []any{"#/components/schemas/Breed3", "/allOf[#/components/schemas/Dog]/breed"},
			Level:       checker.WARN,
			Operation:   "POST",
			Path:        "/pets",
			Source:      load.NewSource("../data/checker/request_property_all_of_removed_revision.yaml"),
			OperationId: "updatePets",
		}}, errs)
}
