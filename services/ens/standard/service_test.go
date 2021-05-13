// Copyright © 2021 Weald Technology Limited.
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

package standard_test

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/wealdtech/edcd/services/ens/standard"
	nullmetrics "github.com/wealdtech/edcd/services/metrics/null"
)

func TestService(t *testing.T) {
	ctx := context.Background()

	monitor := nullmetrics.New()
	domainControls := map[string]interface{}{
		"wealdtech.eth": map[string]interface{}{
			"owner-address": "0x388Ea662EF2c223eC0B047D41Bf3c0f362142ad5",
			"passphrase":    "a secret",
		},
	}

	tests := []struct {
		name   string
		params []standard.Parameter
		err    string
	}{
		{
			name: "MonitorMissing",
			params: []standard.Parameter{
				standard.WithLogLevel(zerolog.Disabled),
				standard.WithMonitor(nil),
				standard.WithTimeout(10 * time.Second),
				standard.WithConnectionURL("http://localhost:8545/"),
				standard.WithDomainControls(domainControls),
			},
			err: "problem with parameters: no monitor specified",
		},
		{
			name: "TimeoutZero",
			params: []standard.Parameter{
				standard.WithLogLevel(zerolog.Disabled),
				standard.WithMonitor(monitor),
				standard.WithTimeout(0),
				standard.WithConnectionURL("http://localhost:8545/"),
				standard.WithDomainControls(domainControls),
			},
			err: "problem with parameters: no timeout specified",
		},
		{
			name: "ConnectionURLMissing",
			params: []standard.Parameter{
				standard.WithLogLevel(zerolog.Disabled),
				standard.WithMonitor(monitor),
				standard.WithTimeout(10 * time.Second),
				standard.WithDomainControls(domainControls),
			},
			err: "problem with parameters: no connection URL specified",
		},
		{
			name: "ConnectionURLBad",
			params: []standard.Parameter{
				standard.WithLogLevel(zerolog.Disabled),
				standard.WithMonitor(monitor),
				standard.WithTimeout(10 * time.Second),
				standard.WithConnectionURL("http://\x01:8545/"),
				standard.WithDomainControls(domainControls),
			},
			err: "invalid URL: parse \"http://\\x01:8545/\": net/url: invalid control character in URL",
		},
		{
			name: "DomainControlsMissing",
			params: []standard.Parameter{
				standard.WithLogLevel(zerolog.Disabled),
				standard.WithMonitor(monitor),
				standard.WithTimeout(10 * time.Second),
				standard.WithConnectionURL("http://localhost:8545/"),
			},
			err: "problem with parameters: no domain controls specified",
		},
		{
			name: "DomainControlsBad",
			params: []standard.Parameter{
				standard.WithLogLevel(zerolog.Disabled),
				standard.WithMonitor(monitor),
				standard.WithTimeout(10 * time.Second),
				standard.WithConnectionURL("http://localhost:8545/"),
				standard.WithDomainControls(map[string]interface{}{
					"wealdtech.eth": "bad",
				}),
			},
			err: "invalid domain controls: invalid configuration for wealdtech.eth",
		},
		{
			name: "Good",
			params: []standard.Parameter{
				standard.WithLogLevel(zerolog.Disabled),
				standard.WithMonitor(monitor),
				standard.WithTimeout(10 * time.Second),
				standard.WithConnectionURL("localhost:8545/"),
				standard.WithDomainControls(domainControls),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := standard.New(ctx, test.params...)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
