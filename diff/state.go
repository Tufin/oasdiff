package diff

type direction int

const (
	directionRequest direction = iota
	directionResponse
)

type state struct {
	visitedSchemasBase     visitedRefs
	visitedSchemasRevision visitedRefs
	cache                  directionalSchemaDiffCache
	direction              direction
}

func newState() *state {
	return &state{
		visitedSchemasBase:     visitedRefs{},
		visitedSchemasRevision: visitedRefs{},
		cache:                  newDirectionalSchemaDiffCache(),
		direction:              directionRequest,
	}
}

func (state *state) setDirection(direction direction) {
	state.direction = direction
}
