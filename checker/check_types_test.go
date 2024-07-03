package checker

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
)

func breaking(t *testing.T, typeDiff *diff.StringsDiff, formatDiff *diff.ValueDiff, isJson bool, revisionTypes *openapi3.Types) {
	t.Helper()

	schemaDiff := &diff.SchemaDiff{
		Revision: &openapi3.Schema{
			Type: revisionTypes,
		},
	}

	mediaType := ""
	if isJson {
		mediaType = "application/json"
	}

	require.True(t, breakingTypeFormatChangedInRequestProperty(typeDiff, formatDiff, mediaType, schemaDiff), "breakingTypeFormatChangedInRequestProperty failed")
}

func notBreaking(t *testing.T, typeDiff *diff.StringsDiff, formatDiff *diff.ValueDiff, isJson bool, revisionTypes *openapi3.Types) {

	schemaDiff := &diff.SchemaDiff{
		Revision: &openapi3.Schema{
			Type: revisionTypes,
		},
	}

	mediaType := ""
	if isJson {
		mediaType = "application/json"
	}
	require.False(t, breakingTypeFormatChangedInRequestProperty(typeDiff, formatDiff, mediaType, schemaDiff), "breakingTypeFormatChangedInRequestProperty failed")
}

func TestStringtoInt(t *testing.T) {
	typeDiff := &diff.StringsDiff{
		Deleted: []string{"string"},
		Added:   []string{"integer"},
	}

	var formatDiff *diff.ValueDiff

	revisionType := &openapi3.Types{
		"integer",
	}

	breaking(t, typeDiff, formatDiff, false, revisionType)
}

func TestIntToString(t *testing.T) {
	typeDiff := &diff.StringsDiff{
		Deleted: []string{"integer"},
		Added:   []string{"string"},
	}

	var formatDiff *diff.ValueDiff

	revisionType := &openapi3.Types{
		"string",
	}
	notBreaking(t, typeDiff, formatDiff, false, revisionType)
}

func TestTypeDeleted(t *testing.T) {
	typeDiff := &diff.StringsDiff{
		Deleted: []string{"integer"},
		Added:   nil,
	}

	var formatDiff *diff.ValueDiff

	revisionType := &openapi3.Types{}

	notBreaking(t, typeDiff, formatDiff, false, revisionType)
}

func TestIntToStringJson(t *testing.T) {
	typeDiff := &diff.StringsDiff{
		Deleted: []string{"integer"},
		Added:   []string{"string"},
	}

	var formatDiff *diff.ValueDiff

	revisionType := &openapi3.Types{
		"string",
	}
	breaking(t, typeDiff, formatDiff, true, revisionType)
}

func TestIntToNumber(t *testing.T) {
	typeDiff := &diff.StringsDiff{
		Deleted: []string{"integer"},
		Added:   []string{"number"},
	}

	var formatDiff *diff.ValueDiff

	revisionType := &openapi3.Types{
		"number",
	}
	notBreaking(t, typeDiff, formatDiff, false, revisionType)
}

func TestUnchanged(t *testing.T) {
	var typeDiff *diff.StringsDiff
	var formatDiff *diff.ValueDiff

	revisionType := &openapi3.Types{
		"integer",
	}
	notBreaking(t, typeDiff, formatDiff, false, revisionType)
}

func TestFormatAdded(t *testing.T) {
	var typeDiff *diff.StringsDiff
	var formatDiff *diff.ValueDiff = &diff.ValueDiff{
		From: nil,
		To:   "int64",
	}

	revisionType := &openapi3.Types{
		"string",
	}
	breaking(t, typeDiff, formatDiff, false, revisionType)
}
