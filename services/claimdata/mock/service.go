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
	"errors"

	"github.com/ethereum/go-ethereum/common"
)

// Service is a mock claim data service.
type Service struct{}

// New creates a new claim data service.
func New() *Service {
	return &Service{}
}

// GetClaimData gets the claim data for a domain.
func (s *Service) GetClaimData(ctx context.Context,
	domain string,
) ([32]byte, string, common.Address, []byte, error) {
	return [32]byte{}, "", common.Address{}, nil, errors.New("mock")
}