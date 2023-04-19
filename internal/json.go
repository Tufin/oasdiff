package internal

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func PrintJSON(output interface{}) error {
	if reflect.ValueOf(output).IsNil() {
		return nil
	}

	bytes, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", bytes)
	return nil
}
