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
	"fmt"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	ens "github.com/wealdtech/go-ens/v3"
)

var registrarABI = `[{"inputs":[{"internalType":"bytes32","name":"node","type":"bytes32"},{"internalType":"address","name":"owner","type":"address"}],"name":"getSignatureHash","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"}]`

// SignatureHash obtains the signature hash for a domain from its parent.
func (s *Service) SignatureHash(ctx context.Context,
	name string,
	domain string,
	owner common.Address,
) (
	[]byte,
	error,
) {
	log := log.With().Str("domain", domain).Logger()

	nameHash, err := ens.NameHash(name)
	if err != nil {
		return nil, err
	}
	log.Trace().Str("name_hash", fmt.Sprintf("%#x", nameHash)).Msg("Calculated name hash")

	abi, err := abi.JSON(strings.NewReader(registrarABI))
	if err != nil {
		return nil, err
	}

	data, err := abi.Pack("getSignatureHash", nameHash, owner)
	if err != nil {
		return nil, err
	}

	backend, err := ethclient.Dial(s.base.String())
	if err != nil {
		return nil, err
	}

	registrarAddress, err := ens.RegistrarContractAddress(backend, domain)
	if err != nil {
		return nil, err
	}
	log.Trace().Str("address", fmt.Sprintf("%#x", registrarAddress)).Msg("Obtained registrar address")

	// TODO can from be 0?
	msg := ethereum.CallMsg{From: owner, To: &registrarAddress, Data: data}
	res, err := backend.CallContract(ctx, msg, nil)
	if err != nil {
		return nil, err
	}

	// TODO need to unpack this?
	fmt.Printf("Res is %#x (%d)\n", res, len(res))

	return res, nil
}
