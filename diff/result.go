package diff

type Result struct {
	PathDiff   *PathDiff             `json:"endpoints,omitempty"`
	SchemaDiff *SchemaCollectionDiff `json:"schemas,omitempty"`
}

func (diffResult *Result) empty() bool {
	return diffResult.PathDiff == nil &&
		diffResult.SchemaDiff == nil
}

func (diffResult *Result) getSummary() *Summary {

	result := Summary{
		Diff: !diffResult.empty(),
	}

	if diffResult.PathDiff != nil {
		result.PathSummary = diffResult.PathDiff.getSummary()
	}

	if diffResult.SchemaDiff != nil {
		result.SchemaSummary = diffResult.SchemaDiff.getSummary()
	}

	return &result
}

func newResult() *Result {
	return &Result{}
}

// FilterByRegex filters diff endpoints by regex
func (diffResult *Result) FilterByRegex(filter string) {
	if diffResult.PathDiff != nil {
		diffResult.PathDiff.filterByRegex(filter)
	}
}
