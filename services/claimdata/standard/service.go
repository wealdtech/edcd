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
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	zerologger "github.com/rs/zerolog/log"
	"github.com/wealdtech/edcd/services/ens"
)

// Service is the ENS service.
type Service struct {
	timeout        time.Duration
	domainControls map[string]*domainControl
	ens            ens.Service
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
	log = zerologger.With().Str("service", "claimdata").Str("impl", "standard").Logger()
	if parameters.logLevel != log.GetLevel() {
		log = log.Level(parameters.logLevel)
	}

	if err := registerMetrics(ctx, parameters.monitor); err != nil {
		return nil, errors.New("failed to register metrics")
	}

	// Parse domain controls
	domainControls, err := parseDomainControls(parameters.domainControls)
	if err != nil {
		return nil, errors.Wrap(err, "invalid domain controls")
	}

	s := &Service{
		timeout:        parameters.timeout,
		domainControls: domainControls,
		ens:            parameters.ens,
	}

	return s, nil
}
