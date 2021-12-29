package diff

// StringsDiff describes the changes between a pair of lists of strings
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

// Empty indicates whether a change was found in this element
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

	s1 := strings1.toStringSet()
	s2 := strings2.toStringSet()

	result.Added = s2.minus(s1).toStringList()
	result.Deleted = s1.minus(s2).toStringList()

	return result
}
