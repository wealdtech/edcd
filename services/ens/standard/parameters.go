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
	"errors"
	"time"

	"github.com/rs/zerolog"
	"github.com/wealdtech/edcd/services/metrics"
	nullmetrics "github.com/wealdtech/edcd/services/metrics/null"
)

type parameters struct {
	logLevel       zerolog.Level
	monitor        metrics.Service
	timeout        time.Duration
	connectionURL  string
	domainControls map[string]interface{}
}

// Parameter is the interface for service parameters.
type Parameter interface {
	apply(*parameters)
}

type parameterFunc func(*parameters)

func (f parameterFunc) apply(p *parameters) {
	f(p)
}

// WithLogLevel sets the log level for the module.
func WithLogLevel(logLevel zerolog.Level) Parameter {
	return parameterFunc(func(p *parameters) {
		p.logLevel = logLevel
	})
}

// WithMonitor sets the monitor for the module.
func WithMonitor(monitor metrics.Service) Parameter {
	return parameterFunc(func(p *parameters) {
		p.monitor = monitor
	})
}

// WithTimeout sets the timeout for requests for this module.
func WithTimeout(timeout time.Duration) Parameter {
	return parameterFunc(func(p *parameters) {
		p.timeout = timeout
	})
}

// WithConnectionURL sets the Ethereum 1 connection URL service for this module.
func WithConnectionURL(url string) Parameter {
	return parameterFunc(func(p *parameters) {
		p.connectionURL = url
	})
}

// WithDomainControls sets the domain controls for this module.
func WithDomainControls(controls map[string]interface{}) Parameter {
	return parameterFunc(func(p *parameters) {
		p.domainControls = controls
	})
}

// parseAndCheckParameters parses and checks parameters to ensure that mandatory parameters are present and correct.
func parseAndCheckParameters(params ...Parameter) (*parameters, error) {
	parameters := parameters{
		logLevel: zerolog.GlobalLevel(),
		monitor:  nullmetrics.New(),
		timeout:  30 * time.Second,
	}
	for _, p := range params {
		if params != nil {
			p.apply(&parameters)
		}
	}

	if parameters.timeout == 0 {
		return nil, errors.New("no timeout specified")
	}
	if parameters.monitor == nil {
		return nil, errors.New("no monitor specified")
	}
	if parameters.connectionURL == "" {
		return nil, errors.New("no connection URL specified")
	}
	if parameters.domainControls == nil {
		return nil, errors.New("no domain controls specified")
	}

	return &parameters, nil
}
