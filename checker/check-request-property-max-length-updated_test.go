package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: increasing max length of request body
func TestRequestBodyMaxLengthDecreasedCheck(t *testing.T) {
	s1, err := open("../data/checker/request_body_max_length_decreased_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_body_max_length_decreased_base.yaml")
	require.NoError(t, err)

	maxLength := uint64(50)
	newMaxLength := uint64(100)
	s1.Spec.Paths["/pets"].Post.RequestBody.Value.Content["application/json"].Schema.Value.MaxLength = &maxLength
	s2.Spec.Paths["/pets"].Post.RequestBody.Value.Content["application/json"].Schema.Value.MaxLength = &newMaxLength

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyMaxLengthUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestBodyMaxLengthIncreasedId,
		Args:        []any{maxLength, newMaxLength},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/pets",
		Source:      "../data/checker/request_body_max_length_decreased_base.yaml",
		OperationId: "addPet",
	}, errs[0])
}

// CL: decreasing max length of request body
func TestRequestBodyMaxLengthIncreasedCheck(t *testing.T) {
	s1, err := open("../data/checker/request_body_max_length_decreased_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_body_max_length_decreased_base.yaml")
	require.NoError(t, err)

	maxLength := uint64(100)
	newMaxLength := uint64(50)
	s1.Spec.Paths["/pets"].Post.RequestBody.Value.Content["application/json"].Schema.Value.MaxLength = &maxLength
	s2.Spec.Paths["/pets"].Post.RequestBody.Value.Content["application/json"].Schema.Value.MaxLength = &newMaxLength

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyMaxLengthUpdatedCheck), d, osm, checker.ERR)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestBodyMaxLengthDecreasedId,
		Args:        []any{newMaxLength},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/pets",
		Source:      "../data/checker/request_body_max_length_decreased_base.yaml",
		OperationId: "addPet",
	}, errs[0])
}

// CL: decreasing max length of request property
func TestRequestPropertyMaxLengthDecreasedCheck(t *testing.T) {
	s1, err := open("../data/checker/request_body_max_length_decreased_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_body_max_length_decreased_base.yaml")
	require.NoError(t, err)

	maxLength := uint64(100)
	newMaxLength := uint64(50)
	s1.Spec.Paths["/pets"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties["description"].Value.MaxLength = &maxLength
	s2.Spec.Paths["/pets"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties["description"].Value.MaxLength = &newMaxLength
	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyMaxLengthUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestPropertyMaxLengthDecreasedId,
		Args:        []any{"description", newMaxLength},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/pets",
		Source:      "../data/checker/request_body_max_length_decreased_base.yaml",
		OperationId: "addPet",
	}, errs[0])
}

// CL: increasing max length of request property
func TestRequestPropertyMaxLengthIncreasedCheck(t *testing.T) {
	s1, err := open("../data/checker/request_body_max_length_decreased_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_body_max_length_decreased_base.yaml")
	require.NoError(t, err)

	maxLength := uint64(50)
	newMaxLength := uint64(100)
	s1.Spec.Paths["/pets"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties["description"].Value.MaxLength = &maxLength
	s2.Spec.Paths["/pets"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties["description"].Value.MaxLength = &newMaxLength
	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyMaxLengthUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestPropertyMaxLengthIncreasedId,
		Args:        []any{"description", maxLength, newMaxLength},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/pets",
		Source:      "../data/checker/request_body_max_length_decreased_base.yaml",
		OperationId: "addPet",
	}, errs[0])
}
