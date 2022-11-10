package diff

// StringList is a list of string values
type StringList []string

func (list StringList) ToStringSet() StringSet {
	result := make(StringSet, len(list))

	for _, s := range list {
		result[s] = struct{}{}
	}

	return result
}

type DiffT interface {
	*ExampleDiff |
		*ServerDiff |
		*ParameterDiff |
		*VariableDiff |
		*SchemaDiff |
		*ResponseDiff |
		*MediaTypeDiff |
		*HeaderDiff |
		SecurityScopesDiff |
		*StringsDiff
}

func ToStringList[diff DiffT](m map[string]diff) StringList {
	keys := make(StringList, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

// Len implements the sort.Interface interface
func (list StringList) Len() int {
	return len(list)
}

// Less implements the sort.Interface interface
func (list StringList) Less(i, j int) bool {
	return list[i] < list[j]
}

// Swap implements the sort.Interface interface
func (list StringList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}
