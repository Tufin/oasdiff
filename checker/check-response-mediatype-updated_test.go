package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: adding a new media type to response
func TestAddNewMediaType(t *testing.T) {
	s1, err := open("../data/checker/add_new_media_type_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/add_new_media_type_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseMediaTypeUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "response-media-type-added",
		Text:        "added the media type 'application/xml' for the response with the status '200'",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/add_new_media_type_revision.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: removing a new media type to response
func TestDeleteNewMediaType(t *testing.T) {
	s1, err := open("../data/checker/add_new_media_type_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/add_new_media_type_base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseMediaTypeUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "response-media-type-removed",
		Text:        "removed the media type 'application/xml' for the response with the status '200'",
		Comment:     "",
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/add_new_media_type_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}
