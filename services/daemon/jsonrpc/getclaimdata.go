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

package jsonrpc

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// GetClaimDataArgs are the arguments for the GetClaimData method.
type GetClaimDataArgs struct {
	Domain string `json:"domain"`
}

// GetClaimDataArgs are the results for the GetClaimData method.
type GetClaimDataResults struct {
	Message   string `json:"message,omitempty"`
	Node      string `json:"node,omitempty"`
	Label     string `json:"label,omitempty"`
	NewOwner  string `json:"newowner,omitempty"`
	Signature string `json:"siganture,omitempty"`
}

func (s *Service) GetClaimData(r *http.Request, args *GetClaimDataArgs, results *GetClaimDataResults) error {
	if args == nil {
		return errors.New("no arguments supplied")
	}

	ctx := context.Background()
	log.Trace().Str("domain", args.Domain).Msg("GetClaimData called")

	node, label, newOwner, signature, err := s.claimData.GetClaimData(ctx, args.Domain)
	if err != nil {
		log.Trace().Err(err).Msg("GetClaimData failed")
		return err
	}

	results.Message = "Success"
	results.Node = fmt.Sprintf("%#x", node)
	results.Label = label
	results.NewOwner = fmt.Sprintf("%#x", newOwner)
	results.Signature = fmt.Sprintf("%#x", signature)
	log.Trace().
		Str("nodehash", results.Node).
		Str("label", results.Label).
		Str("new_owner", results.NewOwner).
		Str("signature", results.Signature).
		Msg("GetClaimData succeeded")

	return nil
}
