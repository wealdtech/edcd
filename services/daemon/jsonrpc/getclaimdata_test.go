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
	"net/http"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	mockclaimdata "github.com/wealdtech/edcd/services/claimdata/mock"
	"github.com/wealdtech/edcd/services/daemon/jsonrpc"
	nullmetrics "github.com/wealdtech/edcd/services/metrics/null"
)

func TestGetClaimData(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		args *jsonrpc.GetClaimDataArgs
		err  string
	}{
		{
			name: "Nil",
			err:  "no arguments supplied",
		},
		{
			name: "DomainMissing",
			args: &jsonrpc.GetClaimDataArgs{},
			err:  "no domain supplied",
		},
		{
			name: "Data",
			args: &jsonrpc.GetClaimDataArgs{
				Domain: "test.com",
			},
		},
	}

	s, err := jsonrpc.New(ctx,
		jsonrpc.WithLogLevel(zerolog.Disabled),
		jsonrpc.WithMonitor(nullmetrics.New()),
		jsonrpc.WithListenAddress(":14732"),
		jsonrpc.WithClaimData(mockclaimdata.New()),
	)
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := &jsonrpc.GetClaimDataResults{}
			r := &http.Request{}
			err := s.GetClaimData(r, test.args, res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
