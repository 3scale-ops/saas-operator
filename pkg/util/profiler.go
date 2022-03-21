package util

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/go-logr/logr"
	"github.com/kelseyhightower/envconfig"
)

type ProfilerConfig struct {
	EnableProfiler      bool   `envconfig:"PROFILER_ENABLE" default:"false"`
	ProfilerBindAddress string `envconfig:"PROFILER_BIND_ADDRESS" default:"0.0.0.0"`
	ProfilerBinbPort    string `envconfig:"PROFILER_BIND_ORT" default:"6060"`
	Log                 logr.Logger
}

func (cfg ProfilerConfig) Setup() error {

	if err := envconfig.Process("profiler", &cfg); err != nil {
		cfg.Log.Error(err, "unable to fetch profile env config")
		return err
	}

	if !cfg.EnableProfiler {
		cfg.Log.Info("profiler not enabled")
		return nil
	}

	return cfg.enable()

}

func (cfg ProfilerConfig) enable() error {
	cfg.Log.Info("profiler enabled")
	go func() {
		pprofListenerAddr := fmt.Sprintf("%s:%s", cfg.ProfilerBindAddress, cfg.ProfilerBinbPort)
		cfg.Log.Info("starting pprof server", "pprofListenerAddr", pprofListenerAddr)
		if err := http.ListenAndServe(pprofListenerAddr, nil); err != nil {
			cfg.Log.Error(err, "unable to start the pprof server")
			return
		}
		cfg.Log.Info("pprof server successfully started",
			"pprofUrl", fmt.Sprintf("http://%s/debug/pprof/", pprofListenerAddr),
		)
	}()
	return nil
}
