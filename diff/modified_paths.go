package diff

// ModifiedPaths is a map of paths to their respective diffs
type ModifiedPaths map[string]*PathDiff

func (modifiedPaths ModifiedPaths) addPathDiff(config *Config, state *state, path1 string, pathItemPair *pathItemPair) error {

	diff, err := getPathDiff(config, state, pathItemPair)
	if err != nil {
		return err
	}

	if !diff.Empty() {
		modifiedPaths[path1] = diff
	}

	return nil
}
