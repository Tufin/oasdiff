package diff

type DiffSummary struct {
	Diff bool `json:"diff"`
	PathDiffSummary
}

type PathDiffSummary struct {
	AddedEndpoints    int `json:"addedEndpoints"`
	DeletedEndpoints  int `json:"deletedEndpoints"`
	ModifiedEndpoints int `json:"modifiedEndpoints"`
}
