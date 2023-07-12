package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: adding a new media type to request body
func TestRequestBodyMediaTypeAdded(t *testing.T) {
	s1, _ := open("../data/checker/request_body_media_type_updated_base.yaml")
	s2, err := open("../data/checker/request_body_media_type_updated_revision.yaml")
	require.Empty(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestBodyMediaTypeChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-body-media-type-added",
		Text:        "added the media type application/json to the request body",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/request_body_media_type_updated_revision.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: adding a new media type to request body
func TestRequestBodyMediaTypeRemoved(t *testing.T) {
	s1, _ := open("../data/checker/request_body_media_type_updated_revision.yaml")
	s2, err := open("../data/checker/request_body_media_type_updated_base.yaml")
	require.Empty(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestBodyMediaTypeChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-body-media-type-removed",
		Text:        "removed the media type application/json from the request body",
		Comment:     "",
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/request_body_media_type_updated_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}
