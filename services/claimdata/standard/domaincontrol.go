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
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

// domainControl contains information about control of a domain.
type domainControl struct {
	Domain     string
	Owner      common.Address
	Passphrase string
}

func parseDomainControls(dcs map[string]interface{}) (map[string]*domainControl, error) {
	domainControls := make(map[string]*domainControl)

	for domain, dc := range dcs {
		control, isControl := dc.(map[string]interface{})
		if !isControl {
			return nil, fmt.Errorf("invalid configuration for %s", domain)
		}

		ownerAddress, exists := control["owner-address"].(string)
		if !exists {
			return nil, fmt.Errorf("owner-address missing for %s", domain)
		}
		owner, err := hex.DecodeString(strings.TrimPrefix(ownerAddress, "0x"))
		if err != nil {
			return nil, errors.Wrapf(err, "owner-address invalid for %s", domain)
		}
		if len(owner) != 20 {
			return nil, fmt.Errorf("incorrect owner-address length for %s", domain)
		}
		address := common.BytesToAddress(owner)

		passphrase, exists := control["passphrase"].(string)
		if !exists {
			return nil, fmt.Errorf("passphrase missing for %s", domain)
		}

		domainControls[domain] = &domainControl{
			Domain:     domain,
			Owner:      address,
			Passphrase: passphrase,
		}
	}

	return domainControls, nil
}
