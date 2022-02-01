package diff

// RequiredPropertiesDiff describes the changes between a pair of lists of required properties
type RequiredPropertiesDiff struct {
	StringsDiff
}

// Empty indicates whether a change was found in this element
func (diff *RequiredPropertiesDiff) Empty() bool {
	if diff == nil {
		return true
	}

	return diff.StringsDiff.Empty()
}

func (diff *RequiredPropertiesDiff) removeNonBreaking() {
	if diff.Empty() {
		return
	}

	if diff.StringsDiff.Empty() {
		return
	}

	diff.Deleted = nil
}

func getRequiredPropertiesDiff(config *Config, state *state, strings1, strings2 StringList) *RequiredPropertiesDiff {

	diff := getRequiredPropertiesDiffInternal(strings1, strings2)

	if config.BreakingOnly {
		diff.removeNonBreaking()
	}

	if diff.Empty() {
		return nil
	}

	return diff
}

func getRequiredPropertiesDiffInternal(strings1, strings2 StringList) *RequiredPropertiesDiff {
	if stringsDiff := getStringsDiff(strings1, strings2); stringsDiff != nil {
		return &RequiredPropertiesDiff{
			StringsDiff: *stringsDiff,
		}
	}
	return nil
}
