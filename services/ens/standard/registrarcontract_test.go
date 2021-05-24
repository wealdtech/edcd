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

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/wealdtech/edcd/services/ens/standard"
	nullmetrics "github.com/wealdtech/edcd/services/metrics/null"
)

func TestSignatureHash(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name   string
		domain string
		owner  common.Address
		err    string
		res    []byte
	}{
		{
			name:   "sub.wealdtech.eth",
			domain: "wealdtech.eth",
			owner:  common.HexToAddress("000102030405060708090a0b0c0d0e0f10111213"),
		},
		{
			name:   "sub.wealdtech.com",
			domain: "wealdtech.com",
			owner:  common.HexToAddress("000102030405060708090a0b0c0d0e0f10111213"),
			err:    "no registrar for wealdtech.com",
		},
	}

	monitor := nullmetrics.New()
	s, err := standard.New(ctx,
		standard.WithLogLevel(zerolog.Disabled),
		standard.WithMonitor(monitor),
		standard.WithTimeout(10*time.Second),
		standard.WithConnectionURL("localhost:8545/"),
	)
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := s.SignatureHash(ctx, test.name, test.domain, test.owner)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.res, res)
			}
		})
	}
}
