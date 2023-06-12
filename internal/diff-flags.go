package internal

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/utils"
	"golang.org/x/exp/slices"
)

type DiffFlags struct {
	base                     string
	revision                 string
	composed                 bool
	prefixBase               string
	prefixRevision           string
	stripPrefixBase          string
	stripPrefixRevision      string
	matchPath                string
	filterExtension          string
	summary                  bool
	format                   string
	failOnDiff               bool
	circularReferenceCounter int
	matchPathParams          bool
	excludeElements          []string
}

func (flags *DiffFlags) isFailOnDiff() bool {
	return flags.failOnDiff
}

func (flags *DiffFlags) isExcludeEndpoints() bool {
	return slices.Contains(flags.excludeElements, "endpoints")
}

func (flags *DiffFlags) toConfig() *diff.Config {
	config := diff.NewConfig()
	config.PathFilter = flags.matchPath
	config.FilterExtension = flags.filterExtension
	config.PathPrefixBase = flags.prefixBase
	config.PathPrefixRevision = flags.prefixRevision
	config.PathStripPrefixBase = flags.stripPrefixBase
	config.PathStripPrefixRevision = flags.stripPrefixRevision
	config.MatchPathParams = flags.matchPathParams
	config.SetExcludeElements(utils.StringList(flags.excludeElements).ToStringSet(), false, false, false)

	return config
}

func (flags *DiffFlags) validate() *ReturnError {
	if flags.base == "" {
		return getErrInvalidFlags(fmt.Errorf("please specify the \"-base\" flag=the path of the original OpenAPI spec in YAML or JSON format"))
	}
	if flags.revision == "" {
		return getErrInvalidFlags(fmt.Errorf("please specify the \"-revision\" flag=the path of the revised OpenAPI spec in YAML or JSON format"))
	}
	if !slices.Contains([]string{"yaml", "json", "text", "html"}, flags.format) {
		return getErrUnsupportedDiffFormat(flags.format)
	}
	if flags.format == "json" && !flags.isExcludeEndpoints() {
		return getErrInvalidFlags(fmt.Errorf("json format requires \"-exclude-elements endpoints\""))
	}
	if invalidElements := diff.ValidateExcludeElements(flags.excludeElements); len(invalidElements) > 0 {
		return getErrInvalidFlags(fmt.Errorf("invalid exclude-elements=%s", flags.excludeElements))
	}

	return nil
}
