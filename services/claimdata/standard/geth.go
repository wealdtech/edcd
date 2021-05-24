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
	"bytes"
	"context"
	"fmt"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/node"
	"github.com/pkg/errors"
)

// gethAccount obtains the Geth account for the given address.
func (s *Service) gethAccount(ctx context.Context, address common.Address) (accounts.Wallet, accounts.Account, error) {
	dir := node.DefaultDataDir()

	ks := keystore.NewKeyStore(filepath.Join(dir, "keystore"), keystore.StandardScryptN, keystore.StandardScryptP)
	am := accounts.NewManager(&accounts.Config{InsecureUnlockAllowed: false}, ks)

	for _, wallet := range am.Wallets() {
		for _, account := range wallet.Accounts() {
			if bytes.Equal(address.Bytes(), account.Address.Bytes()) {
				return wallet, account, nil
			}
		}
	}

	return nil, accounts.Account{}, errors.New("not found")
}

func (s *Service) sign(ctx context.Context, account accounts.Account, tx *types.Transaction) error {
	ks := keystore.NewKeyStore("/path/to/keystore", keystore.StandardScryptN, keystore.StandardScryptP)

	if err := ks.Unlock(account, "Signer password"); err != nil {
		return err
	}
	signature, err := ks.SignHash(account, tx.Hash().Bytes())
	if err != nil {
		return err
	}
	fmt.Printf("Signature is %x\n", signature)

	return nil

}
