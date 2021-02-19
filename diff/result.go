package diff

type Diff struct {
	PathDiff   *PathDiff             `json:"endpoints,omitempty"`
	SchemaDiff *SchemaCollectionDiff `json:"schemas,omitempty"`
}

func (diff *Diff) empty() bool {
	return diff.PathDiff == nil &&
		diff.SchemaDiff == nil
}

func (diff *Diff) getSummary() *Summary {

	result := Summary{
		Diff: !diff.empty(),
	}

	if diff.PathDiff != nil {
		result.PathSummary = diff.PathDiff.getSummary()
	}

	if diff.SchemaDiff != nil {
		result.SchemaSummary = diff.SchemaDiff.getSummary()
	}

	return &result
}

func newDiff() *Diff {
	return &Diff{}
}

// FilterByRegex filters diff endpoints by regex
func (diff *Diff) FilterByRegex(filter string) {
	if diff.PathDiff != nil {
		diff.PathDiff.filterByRegex(filter)
	}
}
