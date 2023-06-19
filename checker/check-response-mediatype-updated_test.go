package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: adding a new media type to response is detected
func TestAddNewMediaType(t *testing.T) {
	s1, err := open("../data/checker/add_new_media_type_base.yaml")
	require.Empty(t, err)
	s2, err := open("../data/checker/add_new_media_type_revision.yaml")
	require.Empty(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseMediaTypeUpdated), d, osm, checker.INFO)
	require.NotEmpty(t, errs)
	require.Equal(t, checker.BackwardCompatibilityErrors{
		{
			Id:          "response-media-type-added",
			Text:        "added the media type 'application/xml' for the response with the status '200'",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/api/v1.0/groups",
			Source:      "../data/checker/add_new_media_type_revision.yaml",
			OperationId: "createOneGroup",
		}}, errs)
}

// CL: Removing a new media type to response is detected
func TestDeleteNewMediaType(t *testing.T) {
	s1, err := open("../data/checker/add_new_media_type_revision.yaml")
	require.Empty(t, err)
	s2, err := open("../data/checker/add_new_media_type_base.yaml")
	require.Empty(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseMediaTypeUpdated), d, osm, checker.INFO)
	require.NotEmpty(t, errs)
	require.Equal(t, checker.BackwardCompatibilityErrors{
		{
			Id:          "response-media-type-removed",
			Text:        "removed the media type 'application/xml' for the response with the status '200'",
			Comment:     "",
			Level:       checker.ERR,
			Operation:   "POST",
			Path:        "/api/v1.0/groups",
			Source:      "../data/checker/add_new_media_type_base.yaml",
			OperationId: "createOneGroup",
		}}, errs)
}
