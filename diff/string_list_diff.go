package diff

// StringsDiff is the diff between two lists of strings
type StringsDiff struct {
	Added   StringList `json:"added,omitempty" yaml:"added,omitempty"`
	Deleted StringList `json:"deleted,omitempty" yaml:"deleted,omitempty"`
}

func newStringsDiff() *StringsDiff {
	return &StringsDiff{
		Added:   StringList{},
		Deleted: StringList{},
	}
}

// Empty return true if there is no diff
func (stringsDiff *StringsDiff) Empty() bool {
	if stringsDiff == nil {
		return true
	}

	return len(stringsDiff.Added) == 0 &&
		len(stringsDiff.Deleted) == 0
}

func getStringsDiff(strings1, strings2 StringList) *StringsDiff {
	diff := getStringsDiffInternal(strings1, strings2)
	if diff.Empty() {
		return nil
	}
	return diff
}

func getStringsDiffInternal(strings1, strings2 StringList) *StringsDiff {
	result := newStringsDiff()

	s1 := stringListToSet(strings1)
	s2 := stringListToSet(strings2)

	result.Added = s2.minus(s1).toStringList()
	result.Deleted = s1.minus(s2).toStringList()

	return result
}
