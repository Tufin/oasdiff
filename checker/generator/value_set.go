package generator

import (
	"fmt"
	"io"
)

type valueSets []valueSet

func (vs valueSets) generate(out io.Writer) {
	for _, v := range vs {
		v.generate(out)
	}
}

type valueSet struct {
	adjective  string
	hierarchy  []string
	attributed []bool
	nouns      []string
	actions    []string
}

func (v valueSet) generate(out io.Writer) {
	for _, noun := range v.nouns {
		for _, action := range v.actions {
			fmt.Fprintln(out, fmt.Sprintf("%s: %s", generateId(v.hierarchy, noun, action), generateMessage(v.hierarchy, v.attributed, noun, v.adjective, action)))
		}
	}
}
