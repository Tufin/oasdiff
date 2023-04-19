package internal

import (
	"fmt"
	"reflect"

	"gopkg.in/yaml.v3"
)

func PrintYAML(output interface{}) error {
	if reflect.ValueOf(output).IsNil() {
		return nil
	}

	bytes, err := yaml.Marshal(output)
	if err != nil {
		return err
	}
	fmt.Printf("%s", bytes)
	return nil
}
