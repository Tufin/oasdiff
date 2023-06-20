package internal

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/tufin/oasdiff/utils"
)

// enumSliceValue is like stringSliceValue with allowed values
type enumSliceValue struct {
	value         *[]string
	allowedValues []string
	changed       bool
}

func newEnumSliceValue(allowedValues []string, val []string, p *[]string) *enumSliceValue {
	esv := new(enumSliceValue)
	esv.allowedValues = allowedValues
	esv.value = p
	*esv.value = val
	return esv
}

func readAsCSV(val string) ([]string, error) {
	if val == "" {
		return []string{}, nil
	}
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)

	records, err := csvReader.Read()
	if err != nil {
		return []string{}, err
	}

	return records, nil
}

func writeAsCSV(vals []string) (string, error) {
	b := &bytes.Buffer{}
	w := csv.NewWriter(b)
	err := w.Write(vals)
	if err != nil {
		return "", err
	}
	w.Flush()
	return strings.TrimSuffix(b.String(), "\n"), nil
}

func checkAllowedValues(values []string, allowed []string) error {
	if notAllowed := utils.StringList(values).ToStringSet().Minus(utils.StringList(allowed).ToStringSet()); !notAllowed.Empty() {
		verb := "are"
		if len(notAllowed) == 1 {
			verb = "is"
		}
		return fmt.Errorf("%s %s not one of the allowed values: %s", strings.Join(notAllowed.ToStringList(), ","), verb, strings.Join(allowed, ","))
	}
	return nil
}

func (s *enumSliceValue) Set(val string) error {
	v, err := readAsCSV(val)
	if err != nil {
		return err
	}

	if err := checkAllowedValues(v, s.allowedValues); err != nil {
		return err
	}

	if !s.changed {
		*s.value = v
	} else {
		*s.value = append(*s.value, v...)
	}
	s.changed = true
	return nil
}

func (s *enumSliceValue) Type() string {
	return "csv"
}

func (s *enumSliceValue) String() string {
	str, _ := writeAsCSV(*s.value)
	return str
}
