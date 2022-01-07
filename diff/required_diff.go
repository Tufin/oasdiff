package diff

// RequiredPropertiesDiff describes the changes between a pair of lists of Required Properties
type RequiredPropertiesDiff struct {
	*StringsDiff
}

func (diff *RequiredPropertiesDiff) removeNonBreaking() {
	if diff.Empty() {
		return
	}

	diff.Deleted = nil
}

func getRequiredPropertiesDiff(config *Config, strings1, strings2 StringList) *RequiredPropertiesDiff {
	diff := &RequiredPropertiesDiff{
		StringsDiff: getStringsDiff(strings1, strings2),
	}

	if config.BreakingOnly {
		diff.removeNonBreaking()
	}

	if diff.Empty() {
		return nil
	}

	return diff
}
