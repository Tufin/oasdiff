package diff

type DiffResult struct {
	PathsDiff *PathsDiff `json:"pathsDiff,omitempty"`
}

func (diffResult *DiffResult) empty() bool {
	return diffResult.PathsDiff == nil || diffResult.PathsDiff.empty()
}

func (diffResult *DiffResult) getSummary() *DiffSummary {

	result := DiffSummary{
		Diff: !diffResult.empty(),
	}

	if diffResult.PathsDiff != nil {
		result.PathsDiffSummary = diffResult.PathsDiff.getSummary()
	}

	return &result
}

func newDiffResult() *DiffResult {
	return &DiffResult{}
}

func (diffResult *DiffResult) FilterByRegex(filter string) {
	if diffResult.PathsDiff != nil {
		diffResult.PathsDiff.filterByRegex(filter)
	}
}

func (diffResult *DiffResult) GetSummary() *DiffSummary {
	return getDiffSummary(diffResult)
}
