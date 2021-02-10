package diff

type DiffSummary struct {
	Diff              bool `json:"diff"`
	DeletedEndpoints  int  `json:"deletedEndpoints"`
	ModifiedEndpoints int  `json:"modifiedEndpoints"`
}

func getDiffSummary(diffResult *DiffResult) *DiffSummary {

	return &DiffSummary{
		Diff:              !diffResult.empty(),
		DeletedEndpoints:  len(diffResult.DeletedEndpoints),
		ModifiedEndpoints: len(diffResult.ModifiedEndpoints),
	}
}
