package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// ExtensionsDiff describes the changes between a pair of sets of specification extensions: https://swagger.io/specification/#specification-extensions
type ExtensionsDiff InterfaceMapDiff

// Breaking indicates whether this element includes a breaking change
func (diff *ExtensionsDiff) Breaking() bool {
	return false
}

// Empty indicates whether a change was found in this element
func (diff *ExtensionsDiff) Empty() bool {
	return (*InterfaceMapDiff)(diff).Empty()
}

func getExtensionsDiff(config *Config, extensions1, extensions2 openapi3.ExtensionProps) *ExtensionsDiff {
	diff := getExtensionsDiffInternal(config, extensions1, extensions2)
	if diff.Empty() {
		return nil
	}

	if config.BreakingOnly && !diff.Breaking() {
		return nil
	}

	return (*ExtensionsDiff)(diff)
}

func getExtensionsDiffInternal(config *Config, extensions1, extensions2 openapi3.ExtensionProps) *InterfaceMapDiff {
	return getInterfaceMapDiff(config, false, extensions1.Extensions, extensions2.Extensions, config.IncludeExtensions)
}
