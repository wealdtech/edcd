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

	"github.com/ethereum/go-ethereum/common"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
	"github.com/wealdtech/go-ens/v3"
)

// GetClaimData gets the claim data for a domain.
func (s *Service) GetClaimData(ctx context.Context,
	domain string,
) (
	[32]byte,
	string,
	common.Address,
	[]byte,
	error,
) {
	log := log.With().Str("domain", domain).Logger()

	domainControl, label, err := s.managedDomain(ctx, domain)
	if err != nil {
		return [32]byte{}, "", common.Address{}, nil, err
	}
	log.Trace().Str("parent_domain", domainControl.Domain).Str("parent_owner", fmt.Sprintf("%#x", domainControl.Owner)).Msg("Obtained parent domain")

	nameHash, err := ens.NameHash(domainControl.Domain)
	if err != nil {
		return [32]byte{}, "", common.Address{}, nil, err
	}

	owner, err := s.ownerForDomain(ctx, domainControl.Domain)
	if err != nil {
		return [32]byte{}, "", common.Address{}, nil, err
	}
	log.Trace().Str("owner", fmt.Sprintf("%#x", owner)).Msg("Obtained domain owner")

	hash, err := s.ens.SignatureHash(ctx, domain, domainControl.Domain, owner)
	if err != nil {
		return [32]byte{}, "", common.Address{}, nil, err
	}
	log.Trace().Str("hash", fmt.Sprintf("%#x", hash)).Msg("Obtained signature hash")

	// TODO sign the hash.
	var sig []byte
	log.Trace().Str("signature", fmt.Sprintf("%#x", sig)).Msg("Signed hash")

	return nameHash, label, owner, sig, nil
}

// managedDomain finds the managed domain given a fully-qualified domain name.
func (s *Service) managedDomain(ctx context.Context, fqdn string) (*domainControl, string, error) {
	domain := normalizeDomain(fqdn)
	separatorIndex := strings.Index(domain, ".")
	if separatorIndex == -1 {
		return nil, "", errors.New("domain not allowed")
	}
	label := domain[:separatorIndex]
	domain = domain[separatorIndex+1:]
	domainControl, exists := s.domainControls[domain]
	if !exists {
		return nil, "", errors.New("domain not supported")
	}
	return domainControl, label, nil
}

// normalizeDomain normalizes an input domain.
func normalizeDomain(domain string) string {
	// If this is the root we keep it.
	if domain == "." {
		return domain
	}

	// Trim leading and trailing periods.
	domain = strings.TrimPrefix(domain, ".")
	domain = strings.TrimSuffix(domain, ".")

	return domain
}

func (s *Service) ownerForDomain(ctx context.Context, domain string) (common.Address, error) {
	c := new(dns.Client)
	c.Net = "udp"
	//				c.DialTimeout = *timeoutDial
	//	c.ReadTimeout = *timeoutRead
	//	c.WriteTimeout = *timeoutWrite
	m := &dns.Msg{
		MsgHdr: dns.MsgHdr{
			Authoritative:     false,
			AuthenticatedData: false,
			CheckingDisabled:  false,
			RecursionDesired:  true,
			Opcode:            dns.OpcodeQuery,
		},
		Question: make([]dns.Question, 1),
	}
	m.Rcode = dns.RcodeSuccess

	//	o := &dns.OPT{
	//		Hdr: dns.RR_Header{
	//			Name:   ".",
	//			Rrtype: dns.TypeOPT,
	//		},
	//	}
	//	o.SetDo()
	//	o.SetUDPSize(dns.DefaultMsgSize)
	//	m.Extra = append(m.Extra, o)

	m.Question[0] = dns.Question{Name: dns.Fqdn(domain), Qtype: dns.TypeTXT, Qclass: dns.ClassINET}
	m.Id = dns.Id()

	r, _, err := c.Exchange(m, "127.0.0.53:53")
	if err != nil {
		return common.Address{}, err
	}
	if r.Id != m.Id {
		return common.Address{}, errors.New("query ID mismatch")
	}

	for _, rr := range r.Answer {
		txtRR, isTxtRR := rr.(*dns.TXT)
		if !isTxtRR {
			continue
		}
		for _, txt := range txtRR.Txt {
			// Only continue if this is an address record.
			if !strings.HasPrefix(txt, "a=0x") {
				continue
			}
			address := common.HexToAddress(txt[4:])
			if address.String() == "0x0000000000000000000000000000000000000000" {
				return common.Address{}, fmt.Errorf("invalid record %s", txt)
			}
			return address, nil
		}
	}

	return common.Address{}, fmt.Errorf("no owner found for domain %s", domain)
}
