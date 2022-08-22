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

package logger

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	ctx := context.Background()

	logger, err := New(ctx, false)
	assert.Nil(t, err)
	assert.NotNil(t, logger)
	defer logger.Clean()

	logger.Info("This is an info message")
	logger.Error("This is an error message")
	logger.Warn("This is a warning message")

	// With debug.
	logger, err = New(ctx, true)
	assert.Nil(t, err)
	assert.NotNil(t, logger)
	defer logger.Clean()
	logger.Debug("This is a debug message")

	// With Attributes
	logger.WithField("test", "value").Info("This is an info message with a field")
	logger.WithFields(Fields(
		"test-key", "test-value",
		"test-other-key", 1,
	)).Info("This is an info message with multiple field")
}
