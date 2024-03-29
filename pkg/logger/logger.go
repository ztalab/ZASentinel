// Copyright 2022-present The Ztalab Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logger

import (
	"context"
	"io"

	"github.com/sirupsen/logrus"
)

// Define key
const (
	TraceIDKey = "trace_id"
	UserIDKey  = "user_id"
	TagKey     = "tag"
	VersionKey = "version"
	StackKey   = "stack"
)

var version string

type Logger = logrus.Logger

type Entry = logrus.Entry

type Hook = logrus.Hook

func StandardLogger() *Logger {
	return logrus.StandardLogger()
}

func SetLevel(level int) {
	logrus.SetLevel(logrus.Level(level))
}

func SetFormatter(format string) {
	switch format {
	case "json":
		logrus.SetFormatter(new(logrus.JSONFormatter))
	default:
		logrus.SetFormatter(new(logrus.TextFormatter))
	}
}

func SetOutput(out io.Writer) {
	logrus.SetOutput(out)
}

func SetVersion(v string) {
	version = v
}

func AddHook(hook Hook) {
	logrus.AddHook(hook)
}

type (
	traceIDKey struct{}
	userIDKey  struct{}
	tagKey     struct{}
	stackKey   struct{}
)

func NewTraceIDContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceID)
}

func FromTraceIDContext(ctx context.Context) string {
	v := ctx.Value(traceIDKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func NewUserIDContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

func FromUserIDContext(ctx context.Context) string {
	v := ctx.Value(userIDKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func NewTagContext(ctx context.Context, tag string) context.Context {
	return context.WithValue(ctx, tagKey{}, tag)
}

func FromTagContext(ctx context.Context) string {
	v := ctx.Value(tagKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func NewStackContext(ctx context.Context, stack error) context.Context {
	return context.WithValue(ctx, stackKey{}, stack)
}

func FromStackContext(ctx context.Context) error {
	v := ctx.Value(stackKey{})
	if v != nil {
		if s, ok := v.(error); ok {
			return s
		}
	}
	return nil
}

func WithErrorStack(ctx context.Context, err error) *Entry {
	if ctx == nil {
		ctx = context.Background()
	}
	return WithContext(NewStackContext(ctx, err))
}

func WithContext(ctx context.Context) *Entry {
	if ctx == nil {
		ctx = context.Background()
	}

	return logrus.WithContext(ctx)
}

// Define logrus alias
var (
	Tracef = logrus.Tracef
	Debugf = logrus.Debugf
	Infof  = logrus.Infof
	Warnf  = logrus.Warnf
	Errorf = logrus.Errorf
	Fatalf = logrus.Fatalf
	Panicf = logrus.Panicf
	Printf = logrus.Printf
)
