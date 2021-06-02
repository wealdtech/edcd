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

package jsonrpc

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	zerologger "github.com/rs/zerolog/log"
	"github.com/wealdtech/edcd/services/claimdata"
	"github.com/wealdtech/edcd/services/daemon/jsonrpc/codecs/mapping"
)

// Service is the JSON-RPC daemon service.
type Service struct {
	srv       *http.Server
	claimData claimdata.Service
}

// module-wide log.
var log zerolog.Logger

// New creates a new JSON-RPC daemon service.
func New(ctx context.Context, params ...Parameter) (*Service, error) {
	parameters, err := parseAndCheckParameters(params...)
	if err != nil {
		return nil, errors.Wrap(err, "problem with parameters")
	}

	// Set logging.
	log = zerologger.With().Str("service", "daemon").Str("impl", "jsonrpc").Logger()
	if parameters.logLevel != log.GetLevel() {
		log = log.Level(parameters.logLevel)
	}

	if err := registerMetrics(ctx, parameters.monitor); err != nil {
		return nil, errors.New("failed to register metrics")
	}

	rpcServer := rpc.NewServer()

	mappingCodec := mapping.New(ctx)
	rpcServer.RegisterCodec(mappingCodec, "application/json")

	s := &Service{
		claimData: parameters.claimData,
	}

	if err := rpcServer.RegisterService(s, "ENSService"); err != nil {
		return nil, errors.Wrap(err, "Failed to register ENS service")
	}
	mappingCodec.Add("ens_getclaimdata", "ENSService.GetClaimData")

	router := mux.NewRouter()
	router.Handle("/", rpcServer)

	s.srv = &http.Server{
		Addr:    parameters.listenAddress,
		Handler: router,
	}

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		for {
			sig := <-sigCh
			if sig == syscall.SIGINT || sig == syscall.SIGTERM || sig == os.Interrupt || sig == os.Kill {
				if err := s.srv.Shutdown(ctx); err != nil {
					log.Warn().Err(err).Msg("Failed to shutdown service")
				}
				break
			}
		}
	}()

	go func() {
		log.Trace().Str("listen_address", parameters.listenAddress).Msg("Starting daemon")
		if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server shut down unexpectedly")
		}
	}()

	return s, nil
}
