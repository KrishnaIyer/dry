package config

import (
	"fmt"
	"reflect"
	"strings"
)

// stringSliceToStringMapHookFunc is a hook for mapstructure that decodes []string to map[string]string.
func stringSliceToStringMapHookFunc(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	if f.Kind() != reflect.Slice || f.Elem().Kind() != reflect.String ||
		t.Kind() != reflect.Map || t.Elem().Kind() != reflect.String {
		return data, nil
	}
	sl := data.([]string)
	m := make(map[string]string, len(sl))
	for _, s := range sl {
		p := strings.SplitN(s, "=", 2)
		if len(p) != 2 {
			return nil, fmt.Errorf("Invalid format for input %v", s)
		}
		m[p[0]] = p[1]
	}

	return m, nil
}
