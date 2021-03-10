package diff

// StringMapDiff is the diff between two string maps
type StringMapDiff struct {
	Added    StringList   `json:"added,omitempty" yaml:"added,omitempty"`
	Deleted  StringList   `json:"deleted,omitempty" yaml:"deleted,omitempty"`
	Modified ModifiedKeys `json:"modified,omitempty" yaml:"modified,omitempty"`
}

// ModifiedKeys maps keys to their respective diffs
type ModifiedKeys map[string]*ValueDiff

func newStringMapDiffDiff() *StringMapDiff {
	return &StringMapDiff{
		Added:    StringList{},
		Deleted:  StringList{},
		Modified: ModifiedKeys{},
	}
}

// Empty return true if there is no diff
func (diff *StringMapDiff) Empty() bool {
	if diff == nil {
		return true
	}

	return len(diff.Added) == 0 &&
		len(diff.Deleted) == 0 &&
		len(diff.Modified) == 0
}

func getStringMapDiff(strings1, strings2 StringMap) *StringMapDiff {
	diff := getStringMapDiffInternal(strings1, strings2)
	if diff.Empty() {
		return nil
	}
	return diff
}

func getStringMapDiffInternal(strings1, strings2 StringMap) *StringMapDiff {
	result := newStringMapDiffDiff()

	for k1, v1 := range strings1 {
		if v2, ok := strings2[k1]; ok {
			if v1 != v2 {
				result.Modified[k1] = getValueDiff(v1, v2)
			}
		} else {
			result.Deleted = append(result.Deleted, k1)
		}
	}

	for k2 := range strings2 {
		if _, ok := strings1[k2]; !ok {
			result.Added = append(result.Added, k2)
		}
	}

	return result
}
