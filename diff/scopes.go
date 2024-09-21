package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/utils"
)

func getScopesDiff(scopes1, scopes2 openapi3.StringMap) *StringMapDiff {
	diff := getStringMapDiff(utils.StringMap(scopes1), utils.StringMap(scopes2))

	if diff.Empty() {
		return nil
	}

	return diff
}
