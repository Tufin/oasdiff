package diff

// StringList is a list of string values
type StringList []string

func (list StringList) toStringSet() StringSet {
	result := make(StringSet, len(list))

	for _, s := range list {
		result[s] = struct{}{}
	}

	return result
}
