package main

import (
	"fmt"
)

func ToStringArray(in []interface{}) []string {
	result := make([]string, len(in))

	for i, val := range in {
		result[i] = fmt.Sprintf("%s", val)
	}

	return result
}

func ToInterfaceArray(in []string) []interface{} {
	result := make([]interface{}, len(in))

	for i, val := range in {
		result[i] = interface{}(val)
	}

	return result
}
