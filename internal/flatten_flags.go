package internal

import "github.com/tufin/oasdiff/load"

type FlattenFlags struct {
	source                   load.Source
	format                   string
	circularReferenceCounter int
}
