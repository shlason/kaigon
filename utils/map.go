package utils

import (
	"reflect"
)

func FilterNilMap(m *map[string]interface{}) map[string]interface{} {
	var s []string

	for k, v := range *m {
		if reflect.ValueOf(v).Kind() == reflect.Ptr {
			if reflect.ValueOf(v).IsNil() {
				s = append(s, k)
			}
		} else {
			if v == nil {
				s = append(s, k)
			}
		}
	}

	for _, v := range s {
		delete(*m, v)
	}

	return *m
}
