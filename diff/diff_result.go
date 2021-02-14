package diff

type DiffResult struct {
	PathDiff *PathDiff `json:"pathDiff,omitempty"`
}

func (diffResult *DiffResult) empty() bool {
	return diffResult.PathDiff == nil || diffResult.PathDiff.empty()
}

func (diffResult *DiffResult) getSummary() *DiffSummary {

	result := DiffSummary{
		Diff: !diffResult.empty(),
	}

	if diffResult.PathDiff != nil {
		result.PathDiffSummary = diffResult.PathDiff.getSummary()
	}

	return &result
}

func newDiffResult() *DiffResult {
	return &DiffResult{}
}

// FilterByRegex filters diff endpoints by regex
func (diffResult *DiffResult) FilterByRegex(filter string) {
	if diffResult.PathDiff != nil {
		diffResult.PathDiff.filterByRegex(filter)
	}
}
