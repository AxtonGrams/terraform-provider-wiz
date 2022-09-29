package utils

import (
	"encoding/json"
	"fmt"
)

// PrettyPrint prints a struct in formatted json
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

// ConvertListToString converts schema.TypeList to a slice of strings
func ConvertListToString(input []interface{}) []string {
	strings := make([]string, 0)
	for _, b := range input {
		strings = append(strings, b.(string))
	}
	return strings
}

// ConvertBoolToPointer converts a bool to a pointer to bool
func ConvertBoolToPointer(in bool) *bool {
	t := new(bool)
	*t = in
	return t
}

// SliceOfStringToMDUList converts a slice of string to an ordered markdown list
func SliceOfStringToMDUList(input []string) string {
	var output string
	output = fmt.Sprintf("\n")
	for _, a := range input {
		output = output + fmt.Sprintf("        - %s\n", a)
	}
	return output
}
