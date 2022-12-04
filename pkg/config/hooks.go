// Copyright © 2022 Krishna Iyer Easwaran
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
