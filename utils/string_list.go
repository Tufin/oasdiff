package utils

import "strings"

type StringList []string

func (stringList *StringList) String() string {
	return strings.Join(*stringList, ", ")
}

func (stringList *StringList) Set(s string) error {
	*stringList = strings.Split(s, ",")
	return nil
}
