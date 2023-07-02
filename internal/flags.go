package internal

import "github.com/tufin/oasdiff/diff"

type Flags interface {
	toConfig() *diff.Config

	getComposed() bool
	getBase() string
	getRevision() string
}
