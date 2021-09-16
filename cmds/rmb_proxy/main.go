package main

import (
	"crypto/tls"
	"flag"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/threefoldtech/rmb_proxy_server/internal/rmbproxy"
	logging "github.com/threefoldtech/rmb_proxy_server/pkg"
)

const (
	CertDefaultCacheDir = "/tmp/certs"
)

func main() {
	f := rmbproxy.Flags{}
	flag.StringVar(&f.Debug, "log-level", "info", "log level [debug|info|warn|error|fatal|panic]")
	flag.StringVar(&f.Substrate, "substrate", "wss://explorer.devnet.grid.tf/ws", "substrate url")
	flag.StringVar(&f.Address, "address", ":443", "explorer running ip address")
	flag.StringVar(&f.Domain, "domain", "", "domain on which the server will be served")
	flag.StringVar(&f.TLSEmail, "email", "", "tmail address to generate certificate with")
	flag.StringVar(&f.CA, "ca", "https://acme-staging-v02.api.letsencrypt.org/directory", "certificate authority used to generate certificate")
	flag.StringVar(&f.CertCacheDir, "cert-cache-dir", CertDefaultCacheDir, "path to store generated certs in")
	flag.Parse()
	if f.Domain == "" {
		log.Fatal().Err(errors.New("domain is required"))
	}
	if f.TLSEmail == "" {
		log.Fatal().Err(errors.New("email is required"))
	}
	logging.SetupLogging(f.Debug)

	if err := app(f); err != nil {
		log.Fatal().Msg(err.Error())
	}
}

func app(f rmbproxy.Flags) error {
	s, err := rmbproxy.CreateServer(f.Substrate, f.Address)
	if err != nil {
		return errors.Wrap(err, "failed to create server")
	}
	config := rmbproxy.CertificateConfig{
		Domain:   f.Domain,
		Email:    f.TLSEmail,
		CA:       f.CA,
		CacheDir: f.CertCacheDir,
	}
	cm := rmbproxy.NewCertificateManager(config)
	go func() {
		if err := cm.ListenForChallenges(); err != nil {
			log.Error().Err(err).Msg("error occurred when listening for challenges")
		}
	}()
	kpr, err := rmbproxy.NewKeypairReloader(cm)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initiate key reloader")
	}
	s.TLSConfig = &tls.Config{
		GetCertificate: kpr.GetCertificateFunc(),
	}

	log.Info().Str("listening on", f.Address).Msg("Server started ...")
	if err := s.ListenAndServeTLS("", ""); err != nil {
		if err == http.ErrServerClosed {
			log.Info().Msg("server stopped gracefully")
		} else {
			log.Error().Err(err).Msg("server stopped unexpectedly")
		}
	}
	return nil
}
