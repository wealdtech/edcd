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
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestParseDomainControls(t *testing.T) {
	tests := []struct {
		name     string
		dcs      map[string]interface{}
		expected map[string]*domainControl
		err      string
	}{
		{
			name:     "Nil",
			expected: map[string]*domainControl{},
		},
		{
			name: "NotControl",
			dcs: map[string]interface{}{
				"wealdtech.eth": "not control",
			},
			err: "invalid configuration for wealdtech.eth",
		},
		{
			name: "OwnerAddressMissing",
			dcs: map[string]interface{}{
				"wealdtech.eth": map[string]interface{}{},
			},
			err: "owner-address missing for wealdtech.eth",
		},
		{
			name: "OwnerAddressInvalid",
			dcs: map[string]interface{}{
				"wealdtech.eth": map[string]interface{}{
					"owner-address": "invalid",
				},
			},
			err: "owner-address invalid for wealdtech.eth: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name: "OwnerAddressShort",
			dcs: map[string]interface{}{
				"wealdtech.eth": map[string]interface{}{
					"owner-address": "0x000102030405060708090a0b0c0d0e0f",
				},
			},
			err: "incorrect owner-address length for wealdtech.eth",
		},
		{
			name: "PassphraseMissing",
			dcs: map[string]interface{}{
				"wealdtech.eth": map[string]interface{}{
					"owner-address": "0x000102030405060708090a0b0c0d0e0f10111213",
				},
			},
			err: "passphrase missing for wealdtech.eth",
		},
		{
			name: "Good",
			dcs: map[string]interface{}{
				"wealdtech.eth": map[string]interface{}{
					"owner-address": "0x000102030405060708090a0b0c0d0e0f10111213",
					"passphrase":    "a secret",
				},
			},
			expected: map[string]*domainControl{
				"wealdtech.eth": {
					Domain:     "wealdtech.eth",
					Owner:      common.HexToAddress("000102030405060708090a0b0c0d0e0f10111213"),
					Passphrase: "a secret",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := parseDomainControls(test.dcs)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expected, res)
			}
		})
	}
}
