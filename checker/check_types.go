package checker

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/utils"
)

func breakingTypeFormatChangedInResponseProperty(typeDiff *diff.StringsDiff, formatDiff *diff.ValueDiff, mediaType string, schemaDiff *diff.SchemaDiff) bool {

	if typeDiff != nil {
		return !isTypeContained(typeDiff.Deleted, typeDiff.Added, isStronglyTyped(mediaType))
	}

	if formatDiff != nil {
		return !isFormatContained(schemaDiff.Revision.Type, formatDiff.From, formatDiff.To)
	}

	return false
}

func breakingTypeFormatChangedInRequestProperty(typeDiff *diff.StringsDiff, formatDiff *diff.ValueDiff, mediaType string, schemaDiff *diff.SchemaDiff) bool {

	if typeDiff != nil {
		return !isTypeContained(typeDiff.Added, typeDiff.Deleted, isStronglyTyped(mediaType))
	}

	if formatDiff != nil {
		return !isFormatContained(schemaDiff.Revision.Type, formatDiff.To, formatDiff.From)
	}

	return false
}

func breakingTypeFormatChangedInRequestParam(typeDiff *diff.StringsDiff, formatDiff *diff.ValueDiff, schemaDiff *diff.SchemaDiff) bool {

	if typeDiff != nil {
		return !isTypeContained(typeDiff.Added, typeDiff.Deleted, false)
	}

	if formatDiff != nil {
		return !isFormatContained(schemaDiff.Revision.Type, formatDiff.To, formatDiff.From)
	}

	return false
}

// isTypeContained checks if type2 is contained in type1
// note that we don't support multiple types currenty
func isTypeContained(to, from utils.StringList, stronglyTyped bool) bool {

	if to.Is("number") && from.Is("integer") {
		return true
	}

	// anything can be changed to string, unless it's "strongly typed"
	if !stronglyTyped {
		return to.Empty() || to.Is("string")
	}

	return false
}

// isStronglyTyped checks if the media type is strongly typed, for example:
// in text format, all numbers can also be interpreted as strings (1 can be a number or a string)
// but in json, a number (1) is not the same as a string ("1")
func isStronglyTyped(mediaType string) bool {
	return isJsonMediaType(mediaType)
}

func isJsonMediaType(mediaType string) bool {
	return mediaType == "application/json" ||
		(strings.HasPrefix(mediaType, "application/vnd.") && strings.HasSuffix(mediaType, "+json"))
}

// isFormatContained checks if from is contained in to
func isFormatContained(revisionType *openapi3.Types, to, from interface{}) bool {

	if revisionType == nil || len(*revisionType) > 1 {
		return false
	}

	// we don't support multiple types currenty, so just take the first one
	switch getSingleType(revisionType) {
	case "number":
		return to == "double" && from == "float"
	case "integer":
		return (to == "int64" && from == "int32") ||
			(to == "bigint" && from == "int32") ||
			(to == "bigint" && from == "int64")
	case "string":
		return (to == "date-time" && from == "date" ||
			to == "date-time" && from == "time")
	}

	return false
}

func getSingleType(types *openapi3.Types) string {
	if types == nil || len(*types) == 0 {
		return ""
	}

	return (*types)[0]
}

func getBaseType(schemaDiff *diff.SchemaDiff) utils.StringList {
	return schemaDiff.Base.Type.Slice()
}

func getRevisionType(schemaDiff *diff.SchemaDiff) utils.StringList {
	return schemaDiff.Revision.Type.Slice()
}

func getBaseFormat(schemaDiff *diff.SchemaDiff) string {
	return schemaDiff.Base.Format
}

func getRevisionFormat(schemaDiff *diff.SchemaDiff) string {
	return schemaDiff.Revision.Format
}
