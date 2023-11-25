package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: changing request parameter type to enum
func TestRequestParameterBecameEnum(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_became_enum_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_became_enum_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterBecameEnumCheck), d, osm, checker.ERR)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestParameterBecameEnumId,
		Text:        "the 'path' request parameter 'groupId' was restricted to a list of enum values",
		Args:        []any{"path", "groupId"},
		Comment:     "",
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/request_parameter_became_enum_revision.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}
