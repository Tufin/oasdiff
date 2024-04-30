package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

// CL: adding a new media type to request body
func TestRequestBodyMediaTypeAdded(t *testing.T) {
	s1, err := open("../data/checker/request_body_media_type_updated_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_body_media_type_updated_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestBodyMediaTypeChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestBodyMediaTypeAddedId,
		Args:        []any{"application/json"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/request_body_media_type_updated_revision.yaml"),
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: removing media type from request body
func TestRequestBodyMediaTypeRemoved(t *testing.T) {
	s1, err := open("../data/checker/request_body_media_type_updated_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_body_media_type_updated_base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestBodyMediaTypeChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestBodyMediaTypeRemovedId,
		Args:        []any{"application/json"},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/request_body_media_type_updated_base.yaml"),
		OperationId: "createOneGroup",
	}, errs[0])
}
