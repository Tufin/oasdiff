package diff

type direction int

const (
	directionRequest direction = iota
	directionResponse
)

type state struct {
	visitedSchemasBase     visitedRefs
	visitedSchemasRevision visitedRefs
	cache                  schemaDiffCache
	direction              direction
}

func newState() *state {
	return &state{
		visitedSchemasBase:     visitedRefs{},
		visitedSchemasRevision: visitedRefs{},
		cache:                  schemaDiffCache{},
		direction:              directionRequest,
	}
}

func (state *state) setDirection(direction direction) {
	state.direction = direction
}
