package delta

import (
	"github.com/tufin/oasdiff/diff"
)

func getParametersDelta(asymmetric bool, d *diff.ParametersDiffByLocation) float64 {
	if d.Empty() {
		return 0
	}

	added := d.Added.Len()
	deleted := d.Deleted.Len()
	modified := d.Modified.Len()
	unchanged := d.Unchanged.Len()
	all := added + deleted + modified + unchanged

	// TODO: drill down into modified
	modifiedDelta := coefficient * float64(modified)

	return ratio(asymmetric, added, deleted, modifiedDelta, all)
}
