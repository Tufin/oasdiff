package diff

type DiffSummary struct {
	Diff              bool `json:"diff"`
	MissingEndpoints  int  `json:"missingEndpoints"`
	ModifiedEndpoints int  `json:"modifiedEndpoints"`
}

func getDiffSummary(diffResult *DiffResult) *DiffSummary {

	return &DiffSummary{
		Diff:              !diffResult.empty(),
		MissingEndpoints:  len(diffResult.MissingEndpoints),
		ModifiedEndpoints: len(diffResult.ModifiedEndpoints),
	}
}
