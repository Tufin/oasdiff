package internal

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

// enumValue is like stringValue with allowed values
type enumValue struct {
	value         *string
	allowedValues []string
}

func newEnumValue(allowedValues []string, val string, p *string) *enumValue {
	result := new(enumValue)
	result.allowedValues = allowedValues
	result.value = p
	*result.value = val
	return result
}

// String is used both by fmt.Print and by Cobra in help text
func (v *enumValue) String() string {
	return string(*v.value)
}

// Set must have pointer receiver so it doesn't change the value of a copy
func (v *enumValue) Set(s string) error {
	if slices.Contains(v.allowedValues, s) {
		*v.value = s
		return nil
	}
	return fmt.Errorf("must be %s", listOf(v.allowedValues))
}

func listOf(options []string) string {
	l := len(options)
	return strings.Join(options[:l-1], ", ") + ", or " + options[l-1]
}

// Type is only used in help text
func (f *enumValue) Type() string {
	return "string"
}
