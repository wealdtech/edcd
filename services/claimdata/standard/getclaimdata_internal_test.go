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
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	mockens "github.com/wealdtech/edcd/services/ens/mock"
)

func TestManagedDomain(t *testing.T) {
	tests := []struct {
		name   string
		fqdn   string
		domain string
		err    string
	}{
		{
			name: "Nil",
			err:  "domain not allowed",
		},
		{
			name: "KnownTLD",
			fqdn: "com",
			err:  "domain not allowed",
		},
		{
			name: "UnknownTLD",
			fqdn: "net",
			err:  "domain not allowed",
		},
		{
			name:   "KnownDomain",
			fqdn:   "example.com",
			domain: "com",
		},
		{
			name:   "KnownSubdomain",
			fqdn:   "foo.example.com",
			domain: "example.com",
		},
		{
			name: "UnknownDomain",
			fqdn: "example.net",
			err:  "domain not supported",
		},
		{
			name:   "KnownSubdomain2",
			fqdn:   "foo.example.net",
			domain: "example.net",
		},
	}

	ctx := context.Background()
	dcs := map[string]interface{}{
		"com": map[string]interface{}{
			"owner-address": "0x000102030405060708090a0b0c0d0e0f10111213",
			"passphrase":    "a secret",
		},
		"example.com": map[string]interface{}{
			"owner-address": "0x0102030405060708090a0b0c0d0e0f1011121314",
			"passphrase":    "a secret",
		},
		"example.net": map[string]interface{}{
			"owner-address": "0x02030405060708090a0b0c0d0e0f101112131415",
			"passphrase":    "a secret",
		},
	}
	s, err := New(ctx,
		WithDomainControls(dcs),
		WithENS(mockens.New()),
	)
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, _, err := s.managedDomain(ctx, test.fqdn)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.domain, res.Domain)
			}
		})
	}
}

func TestOwnerForDomain(t *testing.T) {
	tests := []struct {
		name    string
		domain  string
		address common.Address
		err     string
	}{
		{
			name:    "Present",
			domain:  "wealdtech.eth.link",
			address: common.HexToAddress("a34C6BCAe6F46ac6470443CCea67d937f6060c7E"),
		},
		{
			name:   "NoTXTRecords",
			domain: "a.com",
			err:    "no owner found for domain a.com",
		},
		{
			name:   "NoARecord",
			domain: "example.com",
			err:    "no owner found for domain example.com",
		},
	}

	ctx := context.Background()
	dcs := map[string]interface{}{
		"com": map[string]interface{}{
			"owner-address": "0x000102030405060708090a0b0c0d0e0f10111213",
			"passphrase":    "a secret",
		},
	}
	s, err := New(ctx,
		WithDomainControls(dcs),
		WithENS(mockens.New()),
	)
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := s.ownerForDomain(ctx, test.domain)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.address, res)
			}
		})
	}
}

func TestNormalizeDomain(t *testing.T) {
	tests := []struct {
		name       string
		domain     string
		normalized string
	}{
		{
			name:       "Nil",
			domain:     "",
			normalized: "",
		},
		{
			name:       "Root",
			domain:     ".",
			normalized: ".",
		},
		{
			name:       "TLDLeading",
			domain:     ".com",
			normalized: "com",
		},
		{
			name:       "TLDTrailing",
			domain:     "com.",
			normalized: "com",
		},
		{
			name:       "TLDLeadingAndTrailing",
			domain:     ".com.",
			normalized: "com",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := normalizeDomain(test.domain)
			require.Equal(t, test.normalized, res)
		})
	}
}
