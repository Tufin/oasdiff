package diff

import "github.com/tufin/oasdiff/utils"

type direction int

const (
	directionRequest direction = iota
	directionResponse
)

type state struct {
	visitedSchemasBase     utils.VisitedRefs
	visitedSchemasRevision utils.VisitedRefs
	cache                  directionalSchemaDiffCache
	direction              direction
}

func newState() *state {
	return &state{
		visitedSchemasBase:     utils.VisitedRefs{},
		visitedSchemasRevision: utils.VisitedRefs{},
		cache:                  newDirectionalSchemaDiffCache(),
		direction:              directionRequest,
	}
}

func (state *state) setDirection(direction direction) {
	state.direction = direction
}
