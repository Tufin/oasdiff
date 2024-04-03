package diff

// ExtensionsDiff describes the changes between a pair of sets of specification extensions: https://swagger.io/specification/#specification-extensions
type ExtensionsDiff InterfaceMapDiff

// Empty indicates whether a change was found in this element
func (diff *ExtensionsDiff) Empty() bool {
	return (*InterfaceMapDiff)(diff).Empty()
}

func getExtensionsDiff(config *Config, state *state, extensions1, extensions2 map[string]interface{}) (*ExtensionsDiff, error) {
	diff, err := getExtensionsDiffInternal(config, state, extensions1, extensions2)
	if err != nil {
		return nil, err
	}

	if diff.Empty() {
		return nil, nil
	}

	return (*ExtensionsDiff)(diff), nil
}

func getExtensionsDiffInternal(config *Config, state *state, extensions1, extensions2 map[string]interface{}) (*InterfaceMapDiff, error) {
	return getInterfaceMapDiff(extensions1, extensions2, config.IncludeExtensions)
}
