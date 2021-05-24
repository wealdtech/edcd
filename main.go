package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"syscall"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	zerologger "github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	standardclaimdata "github.com/wealdtech/edcd/services/claimdata/standard"
	jsonrpcdaemon "github.com/wealdtech/edcd/services/daemon/jsonrpc"
	mockens "github.com/wealdtech/edcd/services/ens/mock"
	standardens "github.com/wealdtech/edcd/services/ens/standard"
	"github.com/wealdtech/edcd/services/metrics"
	nullmetrics "github.com/wealdtech/edcd/services/metrics/null"
	prometheusmetrics "github.com/wealdtech/edcd/services/metrics/prometheus"
	"github.com/wealdtech/edcd/util"
)

// ReleaseVersion is the release version for the code.
var ReleaseVersion = "0.1.0"

func main() {
	os.Exit(main2())
}

func main2() int {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := fetchConfig(); err != nil {
		zerologger.Error().Err(err).Msg("Failed to fetch configuration")
		return 1
	}

	if err := initLogging(); err != nil {
		log.Error().Err(err).Msg("Failed to initialise logging")
		return 1
	}

	// runCommands will not return if a command is run.
	runCommands(ctx)

	logModules()
	log.Info().Str("version", ReleaseVersion).Msg("Starting ENS domain claim daemon")

	if err := initProfiling(); err != nil {
		log.Error().Err(err).Msg("Failed to initialise profiling")
		return 1
	}

	runtime.GOMAXPROCS(runtime.NumCPU() * 8)

	log.Trace().Msg("Starting metrics service")
	monitor, err := startMonitor(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to start metrics service")
		return 1
	}
	if err := registerMetrics(ctx, monitor); err != nil {
		log.Error().Err(err).Msg("Failed to register metrics")
		return 1
	}
	setRelease(ctx, ReleaseVersion)
	setReady(ctx, false)

	if err := startServices(ctx, monitor); err != nil {
		log.Error().Err(err).Msg("Failed to initialise services")
		return 1
	}
	setReady(ctx, true)

	log.Info().Msg("All services operational")

	// Wait for signal.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	for {
		sig := <-sigCh
		if sig == syscall.SIGINT || sig == syscall.SIGTERM || sig == os.Interrupt || sig == os.Kill {
			cancel()
			break
		}
	}

	log.Info().Msg("Stopping ENS domain claim daemon")
	return 0
}

// fetchConfig fetches configuration from various sources.
func fetchConfig() error {
	pflag.String("base-dir", "", "base directory for configuration files")
	pflag.Bool("version", false, "show version and exit")
	pflag.String("log-level", "info", "minimum level of messsages to log")
	pflag.String("log-file", "", "redirect log output to a file")
	pflag.String("profile-address", "", "Address on which to run Go profile server")
	pflag.String("eth1client.address", "", "Address for Ethereum 1 node")
	pflag.String("jsonrpc.listen-address", "", "Listen address for JSON-RPC service")
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return errors.Wrap(err, "failed to bind pflags to viper")
	}

	if viper.GetString("base-dir") != "" {
		// User-defined base directory.
		viper.AddConfigPath(resolvePath(""))
		viper.SetConfigName("edcd")
	} else {
		// Home directory.
		home, err := homedir.Dir()
		if err != nil {
			return errors.Wrap(err, "failed to obtain home directory")
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".edcd")
	}

	// Environment settings.
	viper.SetEnvPrefix("EDCD")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return errors.Wrap(err, "failed to read configuration file")
		}
	}

	return nil
}

// initProfiling initialises the profiling server.
func initProfiling() error {
	profileAddress := viper.GetString("profile-address")
	if profileAddress != "" {
		go func() {
			log.Info().Str("profile_address", profileAddress).Msg("Starting profile server")
			runtime.SetMutexProfileFraction(1)
			if err := http.ListenAndServe(profileAddress, nil); err != nil {
				log.Warn().Str("profile_address", profileAddress).Err(err).Msg("Failed to run profile server")
			}
		}()
	}
	return nil
}

func startMonitor(ctx context.Context) (metrics.Service, error) {
	var monitor metrics.Service
	if viper.Get("metrics.prometheus.listen-address") != nil {
		var err error
		monitor, err = prometheusmetrics.New(ctx,
			prometheusmetrics.WithLogLevel(util.LogLevel("metrics.prometheus")),
			prometheusmetrics.WithAddress(viper.GetString("metrics.prometheus.listen-address")),
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to start prometheus metrics service")
		}
		log.Info().Str("listen_address", viper.GetString("metrics.prometheus.listen-address")).Msg("Started prometheus metrics service")
	} else {
		log.Debug().Msg("No metrics service supplied; monitor not starting")
		monitor = &nullmetrics.Service{}
	}
	return monitor, nil
}

func startServices(ctx context.Context, monitor metrics.Service) error {
	log.Trace().Msg("Starting ENS service")
	ens, err := standardens.New(ctx,
		standardens.WithLogLevel(util.LogLevel("ens")),
		standardens.WithMonitor(monitor),
		standardens.WithTimeout(viper.GetDuration("claimdata.timeout")),
		standardens.WithConnectionURL(viper.GetString("eth1client.address")),
	)
	if err != nil {
		return errors.Wrap(err, "failed to start ENS service")
	}
	fmt.Printf("Not using ENS %v\n", ens)

	log.Trace().Msg("Starting claim data service")
	claimData, err := standardclaimdata.New(ctx,
		standardclaimdata.WithLogLevel(util.LogLevel("claimdata")),
		standardclaimdata.WithMonitor(monitor),
		standardclaimdata.WithTimeout(viper.GetDuration("claimdata.timeout")),
		standardclaimdata.WithDomainControls(viper.GetStringMap("claimdata.domain-controls")),
		standardclaimdata.WithENS(mockens.New()),
	)
	if err != nil {
		return errors.Wrap(err, "failed to start claim data service")
	}

	log.Trace().Msg("Starting daemon service")
	_, err = jsonrpcdaemon.New(ctx,
		jsonrpcdaemon.WithLogLevel(util.LogLevel("jsonrpc")),
		jsonrpcdaemon.WithMonitor(monitor),
		jsonrpcdaemon.WithClaimData(claimData),
		jsonrpcdaemon.WithListenAddress(viper.GetString("jsonrpc.listen-address")),
	)
	if err != nil {
		return errors.Wrap(err, "failed to start ENS service")
	}

	return nil
}

func logModules() {
	buildInfo, ok := debug.ReadBuildInfo()
	if ok {
		log.Trace().Str("path", buildInfo.Path).Msg("Main package")
		for _, dep := range buildInfo.Deps {
			log := log.Trace()
			if dep.Replace == nil {
				log = log.Str("path", dep.Path).Str("version", dep.Version)
			} else {
				log = log.Str("path", dep.Replace.Path).Str("version", dep.Replace.Version)
			}
			log.Msg("Dependency")
		}
	}
}

// resolvePath resolves a potentially relative path to an absolute path.
func resolvePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	baseDir := viper.GetString("base-dir")
	if baseDir == "" {
		homeDir, err := homedir.Dir()
		if err != nil {
			log.Fatal().Err(err).Msg("Could not determine a home directory")
		}
		baseDir = homeDir
	}
	return filepath.Join(baseDir, path)
}

func runCommands(ctx context.Context) {
	if viper.GetBool("version") {
		fmt.Printf("%s\n", ReleaseVersion)
		os.Exit(0)
	}
}
