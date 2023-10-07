package internal

import (
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

type Flags interface {
	toConfig() *diff.Config

	getComposed() bool
	getBase() load.Source
	getRevision() load.Source
	getFlatten() bool
}
