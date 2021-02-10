package diff

type DiffSummary struct {
	Diff              bool `json:"diff"`
	AddedEndpoints    int  `json:"addedEndpoints"`
	DeletedEndpoints  int  `json:"deletedEndpoints"`
	ModifiedEndpoints int  `json:"modifiedEndpoints"`
}

func getDiffSummary(diffResult *DiffResult) *DiffSummary {

	return &DiffSummary{
		Diff:              !diffResult.empty(),
		AddedEndpoints:    len(diffResult.AddedEndpoints),
		DeletedEndpoints:  len(diffResult.DeletedEndpoints),
		ModifiedEndpoints: len(diffResult.ModifiedEndpoints),
	}
}
