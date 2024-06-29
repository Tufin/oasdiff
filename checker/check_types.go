package checker

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/utils"
)

func breakingTypeFormatChangedInResponseProperty(typeDiff *diff.StringsDiff, formatDiff *diff.ValueDiff, mediaType string, schemaDiff *diff.SchemaDiff) bool {

	if typeDiff != nil {
		return !isTypeContained(typeDiff.Deleted, typeDiff.Added, mediaType)
	}

	if formatDiff != nil {
		return !isFormatContained(schemaDiff.Revision.Type, formatDiff.From, formatDiff.To)
	}

	return false
}

func breakingTypeFormatChangedInRequestProperty(typeDiff *diff.StringsDiff, formatDiff *diff.ValueDiff, mediaType string, schemaDiff *diff.SchemaDiff) bool {

	if typeDiff != nil {
		return !isTypeContained(typeDiff.Added, typeDiff.Deleted, mediaType)
	}

	if formatDiff != nil {
		return !isFormatContained(schemaDiff.Revision.Type, formatDiff.To, formatDiff.From)
	}

	return false
}

func breakingTypeFormatChangedInRequestParam(typeDiff *diff.StringsDiff, formatDiff *diff.ValueDiff, schemaDiff *diff.SchemaDiff) bool {
	if typeDiff != nil {
		return !isTypeContained(typeDiff.Added, typeDiff.Deleted, "")
	}

	if formatDiff != nil {
		return !isFormatContained(schemaDiff.Revision.Type, formatDiff.To, formatDiff.From)
	}

	return false
}

// isTypeContained checks if type2 is contained in type1
func isTypeContained(type1, type2 utils.StringList, mediaType string) bool {

	if type1.Is("number") && type2.Is("integer") {
		return true
	}

	// anything can be changed to string, unless it's json or xml
	return (type1.Empty() || type1.Is("string")) &&
		!isJsonMediaType(mediaType) && mediaType != "application/xml"
}

// isFormatContained checks if format2 is contained in format1
func isFormatContained(revisionType *openapi3.Types, format1, format2 interface{}) bool {

	if revisionType == nil || len(*revisionType) > 1 {
		return false
	}

	switch revisionType.Slice()[0] {
	case "number":
		return format1 == "double" && format2 == "float"
	case "integer":
		return (format1 == "int64" && format2 == "int32") ||
			(format1 == "bigint" && format2 == "int32") ||
			(format1 == "bigint" && format2 == "int64")
	case "string":
		return (format1 == "date-time" && format2 == "date" ||
			format1 == "date-time" && format2 == "time")
	}

	return false
}

func isJsonMediaType(mediaType string) bool {
	return mediaType == "application/json" ||
		(strings.HasPrefix(mediaType, "application/vnd.") && strings.HasSuffix(mediaType, "+json"))
}
