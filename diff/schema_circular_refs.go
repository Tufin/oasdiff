package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/utils"
)

type circularRefStatus int

const (
	circularRefStatusNone circularRefStatus = iota
	circularRefStatusDiff
	circularRefStatusNoDiff
)

func getCircularRefsDiff(visited1, visited2 utils.VisitedRefs, schema1, schema2 *openapi3.SchemaRef) circularRefStatus {

	if schema1 == nil || schema2 == nil ||
		schema1.Value == nil || schema2.Value == nil {
		return circularRefStatusNone
	}

	circular1 := visited1.IsVisited(schema1.Ref)
	circular2 := visited2.IsVisited(schema2.Ref)

	// neither are circular
	if !circular1 && !circular2 {
		return circularRefStatusNone
	}

	// one ref is circular but the other isn't
	if circular1 != circular2 {
		return circularRefStatusDiff
	}

	// now we know that both refs are circular

	// if they don't reference the same schema name, we consider them to be different
	if schema1.Ref != schema2.Ref {
		return circularRefStatusDiff
	}

	// otherwise, consider them equal
	return circularRefStatusNoDiff
}
