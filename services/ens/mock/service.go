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

package mock

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

// Service is the mock ENS service.
type Service struct{}

// New creates a new mock ENS service.
func New() *Service {
	return &Service{}
}

// SignatureHash obtains the signature hash for a domain from its parent.
// This is a mock; it always returns the same hash.
func (s *Service) SignatureHash(ctx context.Context,
	name string,
	domain string,
	owner common.Address,
) (
	[]byte,
	error,
) {
	return []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x12, 0x13, 0x14}, nil
}
