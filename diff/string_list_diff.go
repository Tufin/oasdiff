package diff

// StringsDiff is the diff between two lists of strings
type StringsDiff struct {
	Added   StringList `json:"added,omitempty"`
	Deleted StringList `json:"deleted,omitempty"`
}

func newStringsDiff() *StringsDiff {
	return &StringsDiff{
		Added:   StringList{},
		Deleted: StringList{},
	}
}

func (stringsDiff *StringsDiff) empty() bool {
	return len(stringsDiff.Added) == 0 &&
		len(stringsDiff.Deleted) == 0
}

func getStringsDiff(strings1, strings2 StringList) *StringsDiff {
	result := newStringsDiff()

	s1 := stringListToSet(strings1)
	s2 := stringListToSet(strings2)

	result.Added = s2.minus(s1).toStringList()
	result.Deleted = s1.minus(s2).toStringList()

	if result.empty() {
		return nil
	}

	return result
}
