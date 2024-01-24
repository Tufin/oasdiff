package delta

import (
	"github.com/tufin/oasdiff/diff"
)

const coefficient = 0.5

func Get(asymmetric bool, diffReport *diff.Diff) float64 {
	if diffReport.Empty() {
		return 0
	}

	return getEndpointsDelta(asymmetric, diffReport.EndpointsDiff)
}

func ratio(asymmetric bool, added int, deleted int, modifiedDelta float64, all int) float64 {
	if asymmetric {
		added = 0
	}

	return (float64(added+deleted) + modifiedDelta) / float64(all)
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
