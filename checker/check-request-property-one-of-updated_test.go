package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: adding 'oneOf' schema to the request body or request body property
func TestRequestPropertyOneOfAdded(t *testing.T) {
	s1, err := open("../data/checker/request_property_one_of_added_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_one_of_added_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyOneOfUpdatedCheck), d, osm, checker.INFO)

	require.Len(t, errs, 2)

	require.ElementsMatch(t, []checker.ApiChange{
		{
			Id:          "request-body-one-of-added",
			Text:        "added 'Rabbit' to the request body 'oneOf' list",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_one_of_added_revision.yaml",
			OperationId: "updatePets",
		},
		{
			Id:          "request-property-one-of-added",
			Text:        "added 'Breed3' to the '/oneOf[#/components/schemas/Dog]/breed' request property 'oneOf' list",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_one_of_added_revision.yaml",
			OperationId: "updatePets",
		}}, errs)
}

// CL: removing 'oneOf' schema from the request body or request body property
func TestRequestPropertyOneOfRemoved(t *testing.T) {
	s1, err := open("../data/checker/request_property_one_of_removed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_one_of_removed_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyOneOfUpdatedCheck), d, osm, checker.INFO)

	require.Len(t, errs, 2)

	require.ElementsMatch(t, []checker.ApiChange{
		{
			Id:          "request-body-one-of-removed",
			Text:        "removed 'Rabbit' from the request body 'oneOf' list",
			Comment:     "",
			Level:       checker.ERR,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_one_of_removed_revision.yaml",
			OperationId: "updatePets",
		},
		{
			Id:          "request-property-one-of-removed",
			Text:        "removed 'Breed3' from the '/oneOf[#/components/schemas/Dog]/breed' request property 'oneOf' list",
			Comment:     "",
			Level:       checker.ERR,
			Operation:   "POST",
			Path:        "/pets",
			Source:      "../data/checker/request_property_one_of_removed_revision.yaml",
			OperationId: "updatePets",
		}}, errs)
}
