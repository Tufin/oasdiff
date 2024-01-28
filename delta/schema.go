package delta

import (
	"github.com/tufin/oasdiff/diff"
)

func getSchemaDelta(asymmetric bool, d *diff.SchemaDiff) float64 {
	if d.Empty() {
		return 0
	}

	// consider additional fields of schema
	typeDelta := modifiedLeafDelta(asymmetric, boolToFloat64(!d.TypeDiff.Empty()))

	return typeDelta
}
