package main

import (
	"flag"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/threefoldtech/rmb_proxy_server/rmbproxy"
	"github.com/threefoldtech/rmb_proxy_server/tools/logging"
)

type flags struct {
	debug     string
	substrate string
	address   string
}

func main() {
	f := flags{}
	flag.StringVar(&f.debug, "log-level", "info", "log level [debug|info|warn|error|fatal|panic]")
	flag.StringVar(&f.substrate, "substrate", "wss://explorer.devnet.grid.tf/ws", "substrate url")
	flag.StringVar(&f.address, "address", ":8080", "explorer running ip address")
	flag.Parse()

	logging.SetupLogging(f.debug)

	if err := app(f); err != nil {
		log.Fatal().Msg(err.Error())
		if err == http.ErrServerClosed {
			log.Info().Msg("server stopped gracefully")
		} else {
			log.Error().Err(err).Msg("server stopped unexpectedly")
		}
	}

}

func app(f flags) error {
	s, err := rmbproxy.CreateServer(f.substrate, f.address)
	if err != nil {
		return errors.Wrap(err, "failed to create server")
	}

	log.Info().Str("listening on", f.address).Msg("Server started ...")
	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
