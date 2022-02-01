package diff

type state struct {
	visitedSchemasBase     visitedRefs
	visitedSchemasRevision visitedRefs
	cache                  schemaDiffCache
}

func newState() *state {
	return &state{
		visitedSchemasBase:     visitedRefs{},
		visitedSchemasRevision: visitedRefs{},
		cache:                  schemaDiffCache{},
	}
}
