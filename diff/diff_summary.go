package diff

type DiffSummary struct {
	Diff bool `json:"diff"`
	PathDiffSummary
	SchemaDiffSummary
}

type PathDiffSummary struct {
	AddedEndpoints    int `json:"addedEndpoints"`
	DeletedEndpoints  int `json:"deletedEndpoints"`
	ModifiedEndpoints int `json:"modifiedEndpoints"`
}

type SchemaDiffSummary struct {
	AddedSchemas   int `json:"addedSchemas"`
	DeletedSchemas int `json:"deletedSchemas"`
}
