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
)

// GetClaimDataArgs are the arguments for the GetClaimData method.
type GetClaimDataArgs struct {
	Domain string `json:"domain"`
}

// GetClaimDataArgs are the results for the GetClaimData method.
type GetClaimDataResults struct {
	Message   string
	Node      string
	Label     string
	NewOwner  string
	Signature string
}

func (s *Service) GetClaimData(r *http.Request, args *GetClaimDataArgs, results *GetClaimDataResults) error {
	ctx := context.Background()
	node, label, newOwner, signature, err := s.ens.GetClaimData(ctx, args.Domain)
	if err != nil {
		results.Message = err.Error()
	} else {
		results.Message = "Success"
		results.Node = fmt.Sprintf("%#x", node)
		results.Label = label
		results.NewOwner = fmt.Sprintf("%#x", newOwner)
		results.Signature = fmt.Sprintf("%#x", signature)
	}
	return nil
}
