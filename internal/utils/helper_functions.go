package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
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

// Missing returns the elements in a that are missing from b
func Missing(a, b []string) []string {
	type void struct{}

	// create map with length of the 'a' slice
	ma := make(map[string]void, len(a))
	diffs := []string{}
	// Convert first slice to map with empty struct (0 bytes)
	for _, ka := range a {
		ma[ka] = void{}
	}
	// find missing values in a
	for _, kb := range b {
		if _, ok := ma[kb]; !ok {
			diffs = append(diffs, kb)
		}
	}
	return diffs
}

// Unique returns the unique values in a slice of strings
func Unique(s []string) []string {
	inResult := make(map[string]bool)
	var result []string
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return result
}

// ConvertSliceToGenericArray returns a generic array from a slice of strings
func ConvertSliceToGenericArray(s []string) []interface{} {
	var output = make([]interface{}, 0)
	for _, b := range s {
		output = append(output, b)
	}
	return output
}

// RemoveNullAndEmptyValues removes null and empty values from a nested map limited in depth traversal
func RemoveNullAndEmptyValues(m map[string]interface{}, depth int) {
	if depth == 0 {
		return
	}
	for k, v := range m {
		if v == nil || (reflect.TypeOf(v).Kind() == reflect.String && v.(string) == "") {
			delete(m, k)
		} else if childMap, ok := v.(map[string]interface{}); ok {
			RemoveNullAndEmptyValues(childMap, depth-1)
		} else if childSlice, ok := v.([]interface{}); ok {
			for _, child := range childSlice {
				if childMap, ok := child.(map[string]interface{}); ok {
					RemoveNullAndEmptyValues(childMap, depth-1)
				}
			}
		}
	}
}
