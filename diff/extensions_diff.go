package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// ExtensionsDiff describes the changes between a pair of sets of specification extensions: https://swagger.io/specification/#specification-extensions
type ExtensionsDiff InterfaceMapDiff

func (diff *ExtensionsDiff) Empty() bool {
	return (*InterfaceMapDiff)(diff).Empty()
}

func getExtensionsDiff(config *Config, extensions1, extensions2 openapi3.ExtensionProps) *ExtensionsDiff {
	diff := getExtensionsDiffInternal(config, extensions1, extensions2)
	if diff.Empty() {
		return nil
	}
	return (*ExtensionsDiff)(diff)
}

func getExtensionsDiffInternal(config *Config, extensions1, extensions2 openapi3.ExtensionProps) *InterfaceMapDiff {
	return getInterfaceMapDiff(extensions1.Extensions, extensions2.Extensions, config.IncludeExtensions)
}
