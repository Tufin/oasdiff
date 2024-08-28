package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

// CL: adding 'oneOf' schema to the request body or request body property
func TestRequestPropertyOneOfAdded(t *testing.T) {
	s1, err := open("../data/checker/request_property_one_of_added_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_one_of_added_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyOneOfUpdatedCheck), d, osm, checker.INFO)

	require.Len(t, errs, 2)

	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestBodyOneOfAddedId,
		Args:        []any{"#/components/schemas/Rabbit", "application/json"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/pets",
		Source:      load.NewSource("../data/checker/request_property_one_of_added_revision.yaml"),
		OperationId: "updatePets",
	}, errs[0])
	require.Equal(t, "added '#/components/schemas/Rabbit' to media-type 'application/json' of request body 'oneOf' list", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))

	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestPropertyOneOfAddedId,
		Args:        []any{"#/components/schemas/Breed3", "/oneOf[#/components/schemas/Dog]/breed", "application/json"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/pets",
		Source:      load.NewSource("../data/checker/request_property_one_of_added_revision.yaml"),
		OperationId: "updatePets",
	}, errs[1])
	require.Equal(t, "added '#/components/schemas/Breed3' to '/oneOf[#/components/schemas/Dog]/breed' request property of media-type 'application/json' 'oneOf' list", errs[1].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// CL: removing 'oneOf' schema from the request body or request body property
func TestRequestPropertyOneOfRemoved(t *testing.T) {
	s1, err := open("../data/checker/request_property_one_of_removed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_one_of_removed_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyOneOfUpdatedCheck), d, osm, checker.INFO)

	require.Len(t, errs, 2)

	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestBodyOneOfRemovedId,
		Args:        []any{"#/components/schemas/Rabbit", "application/json"},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/pets",
		Source:      load.NewSource("../data/checker/request_property_one_of_removed_revision.yaml"),
		OperationId: "updatePets",
	}, errs[0])
	require.Equal(t, "removed '#/components/schemas/Rabbit' from media-type 'application/json' of request body 'oneOf' list", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))

	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestPropertyOneOfRemovedId,
		Args:        []any{"#/components/schemas/Breed3", "/oneOf[#/components/schemas/Dog]/breed", "application/json"},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/pets",
		Source:      load.NewSource("../data/checker/request_property_one_of_removed_revision.yaml"),
		OperationId: "updatePets",
	}, errs[1])
	require.Equal(t, "removed '#/components/schemas/Breed3' from '/oneOf[#/components/schemas/Dog]/breed' request property of media-type 'application/json' 'oneOf' list", errs[1].GetUncolorizedText(checker.NewDefaultLocalizer()))
}
