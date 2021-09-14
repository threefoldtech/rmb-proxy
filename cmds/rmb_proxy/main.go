package main

import (
	"flag"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/threefoldtech/rmb_proxy_server/pkg/logging"
	"github.com/threefoldtech/rmb_proxy_server/pkg/rmbproxy"
)

func main() {
	f := rmbproxy.Flags{}
	flag.StringVar(&f.Debug, "log-level", "info", "log level [debug|info|warn|error|fatal|panic]")
	flag.StringVar(&f.Substrate, "substrate", "wss://explorer.devnet.grid.tf/ws", "substrate url")
	flag.StringVar(&f.Address, "address", ":8080", "explorer running ip address")
	flag.Parse()

	logging.SetupLogging(f.Debug)

	if err := app(f); err != nil {
		log.Fatal().Msg(err.Error())
		if err == http.ErrServerClosed {
			log.Info().Msg("server stopped gracefully")
		} else {
			log.Error().Err(err).Msg("server stopped unexpectedly")
		}
	}

}

func app(f rmbproxy.Flags) error {
	s, err := rmbproxy.CreateServer(f.Substrate, f.Address)
	if err != nil {
		return errors.Wrap(err, "failed to create server")
	}

	log.Info().Str("listening on", f.Address).Msg("Server started ...")
	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
