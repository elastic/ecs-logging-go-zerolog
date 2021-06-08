// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package ecszerolog

import (
	"io"
	"runtime"
	"time"

	"github.com/rs/zerolog"
	"go.elastic.co/ecszerolog/internal"
)

const (
	ecsVersion = "1.6.0"
	originKey  = "log.origin"
)

// Option indicates optional configuration that can be used when creating a logger.
type Option func(l zerolog.Logger) zerolog.Logger

// Level will set the minimum accepted level by the logger.
func Level(level zerolog.Level) Option {
	return func(l zerolog.Logger) zerolog.Logger {
		return l.Level(level)
	}
}

// origin is a hook to add the file name and line to an entry.
var origin zerolog.HookFunc = func(e *zerolog.Event, _ zerolog.Level, _ string) {
	_, file, line, ok := runtime.Caller(zerolog.CallerSkipFrameCount)
	if !ok {
		return
	}
	e.Dict(originKey, zerolog.Dict().
		Str("file.name", file).
		Int("file.line", line),
	)
}

// Origin is an option to add the log.origin attribute populated with file.name and file.line attributes in all entries.
func Origin() Option {
	return func(l zerolog.Logger) zerolog.Logger {
		return l.Hook(origin)
	}
}

// ErrorStack enables stack_trace output for supported errors.
//
// currently supports: github.com/pkg/errors
func ErrorStack() Option {
	return func(l zerolog.Logger) zerolog.Logger {
		return l.With().Stack().Logger()
	}
}

// New will configure zerolog global options and return an ECS compliant logger.
//
// All entries will have @timestamp and ecs.version attributes automatically added.
func New(w io.Writer, options ...Option) zerolog.Logger {
	zerolog.MessageFieldName = "message"
	zerolog.ErrorFieldName = "error.message"
	zerolog.ErrorStackMarshaler = internal.MarshallStack
	zerolog.ErrorStackFieldName = "error.stack_trace"
	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.999Z" // RFC3339 at millisecond resolution in zulu timezone
	zerolog.TimestampFieldName = "@timestamp"
	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	zerolog.LevelFieldName = "log.level"
	zerolog.CallerSkipFrameCount = 4

	l := zerolog.New(w).With().Timestamp().Str("ecs.version", ecsVersion).Logger()
	for _, option := range options {
		l = option(l)
	}
	return l
}
