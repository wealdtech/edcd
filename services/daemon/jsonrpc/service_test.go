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

package jsonrpc_test

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/wealdtech/edcd/services/daemon/jsonrpc"
	mockens "github.com/wealdtech/edcd/services/ens/mock"
	nullmetrics "github.com/wealdtech/edcd/services/metrics/null"
)

func TestService(t *testing.T) {
	ctx := context.Background()

	ens := mockens.New()
	monitor := nullmetrics.New()

	tests := []struct {
		name   string
		params []jsonrpc.Parameter
		err    string
	}{
		{
			name: "MonitorMissing",
			params: []jsonrpc.Parameter{
				jsonrpc.WithLogLevel(zerolog.Disabled),
				jsonrpc.WithMonitor(nil),
				jsonrpc.WithListenAddress(":14732"),
				jsonrpc.WithENS(ens),
			},
			err: "problem with parameters: no monitor specified",
		},
		{
			name: "ListenAddressMissing",
			params: []jsonrpc.Parameter{
				jsonrpc.WithLogLevel(zerolog.Disabled),
				jsonrpc.WithMonitor(monitor),
				jsonrpc.WithENS(ens),
			},
			err: "problem with parameters: no listen address specified",
		},
		{
			name: "ENSMissing",
			params: []jsonrpc.Parameter{
				jsonrpc.WithLogLevel(zerolog.Disabled),
				jsonrpc.WithMonitor(monitor),
				jsonrpc.WithListenAddress(":14732"),
			},
			err: "problem with parameters: no ENS service specified",
		},
		{
			name: "Good",
			params: []jsonrpc.Parameter{
				jsonrpc.WithLogLevel(zerolog.Disabled),
				jsonrpc.WithMonitor(monitor),
				jsonrpc.WithListenAddress(":14732"),
				jsonrpc.WithENS(ens),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := jsonrpc.New(ctx, test.params...)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
