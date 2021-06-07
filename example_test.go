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

package ecszerolog_test

import (
	"os"
	"time"

	"go.elastic.co/ecszerolog"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ExampleLogger() {
	logger := ecszerolog.New(os.Stdout, ecszerolog.Timestamp())
	log.Logger = logger

	zerolog.TimestampFunc = func() time.Time { return time.Unix(0, 0).UTC() } // set time on entry to a known value

	log.Info().Str("foo", "bar").Msg("hello")
	// Output:
	// {"log.level":"info","ecs.version":"1.6.0","foo":"bar","@timestamp":"1970-01-01T00:00:00Z","message":"hello"}
}

func ExampleLoggerError() {
	logger := ecszerolog.New(os.Stdout)
	log.Logger = logger

	err := errors.New("something bad happened")
	log.Error().Err(err).Msg("An error has occured")
	// Output:
	// {"log.level":"error","ecs.version":"1.6.0","error.message":"something bad happened","message":"An error has occured"}
}
