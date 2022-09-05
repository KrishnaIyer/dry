// Copyright © 2022 Krishna Iyer Easwaran
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package assertions is an opinionated wrapper around the testify/assert package.
package assertions

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	IsNil    = "nil"
	IsNotNil = "not-nil"
)

type Condition int

// Assertion is an assertion.
type Assertion struct {
	t *testing.T
}

// New creates a new test assertion.
func New(t *testing.T) *Assertion {
	return &Assertion{
		t: t,
	}
}

// Assert tests a condition.
// If the condition is not met, the test will error and fail in-place.
func (a *Assertion) Assert(condition string, actual interface{}, expected ...interface{}) bool {
	switch condition {
	case IsNil:
		if !assert.Nil(a.t, actual) {
			log.Fatalf("Expected nil, got %v", actual)
			return false
		}
	case IsNotNil:
		if !assert.NotNil(a.t, actual) {
			log.Fatal("Expected value to not be nil, but it was!")
			return false
		}
	default:
		return false
	}
	return true
}
