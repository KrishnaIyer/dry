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

// Package logger is a highly opinionated wrapper around the Zap logger
package logger

import (
	"context"
	"reflect"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var reflectTypeToZapFieldType = map[reflect.Kind]zapcore.FieldType{
	reflect.String: zapcore.StringType,
	reflect.Int64:  zapcore.Int64Type,
	reflect.Int:    zapcore.Int64Type,
}

type loggerKeyType string

var loggerKey loggerKeyType = "logger"

// Logger wraps zap.Logger.
type Logger struct {
	ctx    context.Context
	logger *zap.Logger
	fields []zap.Field
}

// Options is the logger options.
type Options struct {
}

// Field represents a logger field.
type Field struct {
	Key   string
	Value interface{}
}

// Fields accepts the values as a string:interface pair.
// If the pairs don't match, they will be skipped
func Fields(a ...any) (fields []zapcore.Field) {
	for i := 0; i < len(a)-1; i = i + 2 {
		if str, ok := a[i].(string); ok {
			fields = append(fields, field(str, a[i+1]))
		}
	}
	return
}

// New creates a new logger. Make sure to call defer logger.Clean() after calling this.
// Log messages are JSON Encoded and timestamps are RFC3339 encoded.
func New(ctx context.Context, debug bool) (*Logger, error) {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	cfg := zapcore.EncoderConfig{
		EncodeTime:    zapcore.RFC3339TimeEncoder,
		MessageKey:    "msg",
		TimeKey:       "ts",
		CallerKey:     "caller",
		StacktraceKey: "trace",
	}
	if !debug {
		config.DisableStacktrace = true
		config.DisableCaller = true
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		cfg.CallerKey = ""
	}

	zap.RegisterEncoder("custom", func(ec zapcore.EncoderConfig) (zapcore.Encoder, error) {
		return zapcore.NewJSONEncoder(cfg), nil
	})
	config.Encoding = "custom"
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	return &Logger{
		ctx:    ctx,
		logger: logger,
		fields: make([]zap.Field, 0),
	}, nil
}

// Clean cleans up the log states. Make sure to call this after creating a new logger.
func (l *Logger) Clean() {
	l.logger.Sync()
}

// NewContextWithLogger returns a new context with a logger and panics if it doesn't match the interface.
func NewContextWithLogger(parentCtx context.Context, logger *Logger) context.Context {
	if logger == nil {
		panic("Nil Logger")
	}
	return context.WithValue(parentCtx, loggerKey, logger)
}

// NewLoggerFromContext retrieves a logger from a context and panics if there isn't one.
func NewLoggerFromContext(ctx context.Context) *Logger {
	val := ctx.Value(loggerKey)
	logger, ok := val.(*Logger)
	if !ok {
		panic("No logger in context")
	}
	return logger
}

// Debug logs a Debug level message.
func (l *Logger) Debug(msg string) {
	l.logger.With(l.fields...).Debug(msg)
}

// Info logs a Info level message.
func (l *Logger) Info(msg string) {
	l.logger.With(l.fields...).Info(msg)
}

// Warn logs a Warning level message.
func (l *Logger) Warn(msg string) {
	l.logger.With(l.fields...).Warn(msg)
}

// Error logs a Error level message.
func (l *Logger) Error(msg string) {
	l.logger.With(l.fields...).Error(msg)
}

// Fatal logs a Fatal message.
func (l *Logger) Fatal(msg string) {
	l.logger.With(l.fields...).Fatal(msg)
}

func field(key string, val interface{}) (field zapcore.Field) {
	kind := reflect.TypeOf(val).Kind()
	zapType := reflectTypeToZapFieldType[kind]
	switch kind {
	case reflect.String:
		field = zap.Field{Key: key, String: val.(string), Type: zapType}
	case reflect.Int:
		field = zap.Field{Key: key, Integer: int64(val.(int)), Type: zapType}
	case reflect.Int64:
		field = zap.Field{Key: key, Integer: int64(val.(int)), Type: zapType}
	case reflect.Float64:
		field = zap.Field{Key: key, Interface: float64(val.(float64)), Type: zapType}
	default:
		// Skip this since we don't know the type
	}
	return
}

// WithField returns a new logger with the provided field.
func (l *Logger) WithField(key string, val interface{}) *Logger {
	return &Logger{
		ctx:    l.ctx,
		logger: l.logger,
		fields: append(l.fields, field(key, val)),
	}
}

// WithFields returns a logger with the providedfields.
func (l *Logger) WithFields(f []zapcore.Field) *Logger {
	return &Logger{
		ctx:    l.ctx,
		logger: l.logger,
		fields: append(l.fields, f...),
	}
}
