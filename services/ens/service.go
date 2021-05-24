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

package ens

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

// Service defines the ENS service.
type Service interface {
	// SignatureHash obtains the signature hash for a domain from its parent.
	SignatureHash(ctx context.Context,
		name string,
		domain string,
		owner common.Address,
	) (
		[]byte,
		error,
	)

	//// DomainOwner gets the owner for a domain.
	//DomainOwner(ctx context.Context,
	//domain string,
	//)(
	//
	//)
	//	// GetClaimData gets the claim data for a domain.
	//	GetClaimData(ctx context.Context,
	//		domain string,
	//	) (
	//		nameHash [32]byte,
	//		label string,
	//		owner common.Address,
	//		signature []byte,
	//		err error,
	//	)
}
