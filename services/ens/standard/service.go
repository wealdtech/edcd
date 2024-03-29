// Copyright © 2021 Weald Technology Trading.
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

package standard

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	zerologger "github.com/rs/zerolog/log"
)

// Service is the ENS service.
type Service struct {
	base    *url.URL
	timeout time.Duration
}

// module-wide log.
var log zerolog.Logger

// New creates a new claim data service.
func New(ctx context.Context, params ...Parameter) (*Service, error) {
	parameters, err := parseAndCheckParameters(params...)
	if err != nil {
		return nil, errors.Wrap(err, "problem with parameters")
	}

	// Set logging.
	log = zerologger.With().Str("service", "ens").Str("impl", "standard").Logger()
	if parameters.logLevel != log.GetLevel() {
		log = log.Level(parameters.logLevel)
	}

	if err := registerMetrics(ctx, parameters.monitor); err != nil {
		return nil, errors.New("failed to register metrics")
	}

	// Connect to Ethereum 1.
	connectionURL := parameters.connectionURL
	if !strings.HasPrefix(connectionURL, "http") {
		connectionURL = fmt.Sprintf("http://%s", parameters.connectionURL)
	}
	base, err := url.Parse(connectionURL)
	if err != nil {
		return nil, errors.Wrap(err, "invalid URL")
	}

	s := &Service{
		base:    base,
		timeout: parameters.timeout,
	}

	return s, nil
}
