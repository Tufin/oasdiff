package diff

type DiffSummary struct {
	Diff bool `json:"diff"`
	PathsDiffSummary
}

type PathsDiffSummary struct {
	AddedEndpoints    int `json:"addedEndpoints"`
	DeletedEndpoints  int `json:"deletedEndpoints"`
	ModifiedEndpoints int `json:"modifiedEndpoints"`
}

func getDiffSummary(diffResult *DiffResult) *DiffSummary {

	return diffResult.getSummary()
}
