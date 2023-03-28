package utils

import (
	"sort"
	"strings"
)

// StringList is a list of string values
type StringList []string

func (stringList *StringList) String() string {
	return strings.Join(*stringList, ", ")
}

func (stringList *StringList) Set(s string) error {
	*stringList = strings.Split(s, ",")
	return nil
}

func (stringList *StringList) Contains(s string) bool {
	for _, item := range *stringList {
		if s == item {
			return true
		}
	}
	return false
}

func (stringList *StringList) Minus(other StringList) StringList {
	return stringList.ToStringSet().Minus(other.ToStringSet()).ToStringList()
}

func (list StringList) ToStringSet() StringSet {
	result := make(StringSet, len(list))

	for _, s := range list {
		result[s] = struct{}{}
	}

	return result
}

func (list StringList) Sort() StringList {
	sort.Sort(list)
	return list
}

// Len implements the sort.Interface interface
func (list StringList) Len() int {
	return len(list)
}

// Less implements the sort.Interface interface
func (list StringList) Less(i, j int) bool {
	return list[i] < list[j]
}

// Swap implements the sort.Interface interface
func (list StringList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}
