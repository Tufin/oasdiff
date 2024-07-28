package internal

import (
	"github.com/tufin/oasdiff/load"
)

type FlattenFlags struct {
	CommonFlags

	spec *load.Source
}

func NewFlattenFlags() *FlattenFlags {
	return &FlattenFlags{
		CommonFlags: NewCommonFlags(),
	}
}
