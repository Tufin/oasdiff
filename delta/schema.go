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

func modifiedLeafDelta(asymmetric bool, modified float64) float64 {
	if asymmetric {
		return modified / 2
	}

	return modified
}

func boolToFloat64(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}
