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
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"go.elastic.co/ecszerolog/internal/spec"
)

type testOut struct {
	m map[string]interface{}
}

func (o *testOut) Write(p []byte) (int, error) {
	return len(p), json.Unmarshal(p, &o.m)
}

func (o *testOut) validate(t *testing.T, keys ...string) {
	// Fields `log.origin.file.line` and `log.origin.file.name` are logged as
	// a map under key log.origin. By using the zap jsonEncoder this cannot
	// be changed. Remove the nested and add a dotted version of the fields
	if caller, ok := o.m[originKey].(map[string]interface{}); ok {
		for name, val := range caller {
			o.m[fmt.Sprintf("%s.%s", originKey, name)] = val
			delete(o.m, name)
		}
	}

	for _, s := range keys {
		require.Contains(t, o.m, s)
	}

	for name, field := range spec.V1.Fields {
		val, ok := o.m[name]
		if field.Required { // all required fields must be present in the log line
			require.True(t, ok, "%s is required", name)
			require.NotNil(t, val)
		}
		if !ok { // custom field not defined in spec
			continue
		}
		if field.Type != "" { // the defined type must be met
			var ok bool
			switch field.Type {
			case "string":
				_, ok = val.(string)
			case "datetime":
				var s string
				s, ok = val.(string)
				if _, err := time.Parse("2006-01-02T15:04:05.000Z0700", s); err == nil {
					ok = true
				}
			case "integer":
				// json.Unmarshal unmarshals into float64 instead of int
				if _, floatOK := val.(float64); floatOK {
					_, err := strconv.ParseInt(fmt.Sprintf("%v", val), 10, 64)
					if err == nil {
						ok = true
					}
				}
			default:
				panic(fmt.Errorf("unhandled type %s from specification for field %s", field.Type, name))
			}
			require.True(t, ok, fmt.Sprintf("%s: %v", name, val))
		}
	}
}

func TestNew(t *testing.T) {
	to := &testOut{}
	logger := New(to)

	logger.Print("hello world")
	to.validate(t)
}

func TestLoggerLevel(t *testing.T) {
	to := &testOut{}
	logger := New(to, Level(zerolog.InfoLevel))

	if logger.GetLevel() != zerolog.InfoLevel {
		t.Errorf("Expected InfoLevel, got %v", logger.GetLevel())
	}
	logger.Error().Err(fmt.Errorf("oh no")).Msg("an error")
	to.validate(t, "error.message")
}

func TestLoggerErrorStack(t *testing.T) {
	to := &testOut{}
	logger := New(to, ErrorStack())

	err := errors.New("something bad happened")
	logger.Error().Err(err).Msg("An error has occured")
	to.validate(t, "error.message", "error.stack_trace")
}

func TestLoggerOrigin(t *testing.T) {
	to := &testOut{}
	logger := New(to, Origin())

	logger.Info().Msg("hello world")
	to.validate(t, "log.origin.file.name", "log.origin.file.line")
}
