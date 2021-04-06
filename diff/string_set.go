package diff

// StringSet is a set of string values
type StringSet map[string]struct{}

func (stringSet StringSet) toStringList() StringList {
	result := make(StringList, len(stringSet))
	i := 0
	for s := range stringSet {
		result[i] = s
		i++
	}
	return result
}

func (stringSet StringSet) add(s string) {
	stringSet[s] = struct{}{}
}

func (stringSet StringSet) contains(s string) bool {
	_, ok := stringSet[s]
	return ok
}

func (stringSet StringSet) minus(other StringSet) StringSet {
	result := StringSet{}

	for s := range stringSet {
		if !other.contains(s) {
			result.add(s)
		}
	}

	return result
}

func (stringSet StringSet) equals(other StringSet) bool {
	return stringSet.minus(other).Empty() &&
		other.minus(stringSet).Empty()
}

// Empty indicates whether a change was found in this element
func (stringSet StringSet) Empty() bool {
	return len(stringSet) == 0
}
