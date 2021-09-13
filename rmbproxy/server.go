package rmbproxy

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func errorReply(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	fmt.Fprintf(w, "{\"status\": \"error\", \"message\": \"%s\"}", message)
}

func (a *App) NewTwinClient(twinID int) (TwinClient, error) {
	log.Debug().Int("twin", twinID).Msg("resolving twin")

	twin, err := a.resolver.client.GetTwin(uint32(twinID))
	if err != nil {
		return nil, err
	}
	log.Debug().Str("ip", twin.IP).Msg("resolved twin ip")

	return &twinClient{
		dstIP: twin.IP,
	}, nil
}

func (a *App) sendMessage(w http.ResponseWriter, r *http.Request) {
	twinIDString := mux.Vars(r)["twin_id"]

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(r.Body)

	twinID, err := strconv.Atoi(twinIDString)
	if err != nil {
		errorReply(w, http.StatusBadRequest, "Invalid twinId")
		return
	}

	c, err := a.NewTwinClient(twinID)
	if err != nil {
		log.Error().Err(err).Msg("failed to create mux server")
	}

	data, err := c.SubmitMessage(*buffer)
	if err != nil {
		errorReply(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data))
}

func (a *App) getResult(w http.ResponseWriter, r *http.Request) {
	twinIDString := mux.Vars(r)["twin_id"]
	retqueue := mux.Vars(r)["retqueue"]

	reqBody := MessageIdentifier{
		Retqueue: retqueue,
	}

	twinID, err := strconv.Atoi(twinIDString)
	if err != nil {
		errorReply(w, http.StatusBadRequest, "Invalid twinId")
		return
	}

	c, err := a.NewTwinClient(twinID)
	if err != nil {
		log.Error().Err(err).Msg("failed to create mux server")
	}

	data, err := c.GetResult(reqBody)
	if err != nil {
		errorReply(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data))
}

func CreateServer(substrate string, address string) (*http.Server, error) {
	log.Info().Msg("Creating server")
	router := mux.NewRouter().StrictSlash(true)

	resolver, err := NewTwinResolver(substrate)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get a client to explorer resolver")
	}

	a := &App{
		resolver: *resolver,
	}

	router.HandleFunc("/twin/{twin_id:[0-9]+}", a.sendMessage)
	router.HandleFunc("/twin/{twin_id:[0-9]+}/{retqueue}", a.getResult)

	return &http.Server{
		Handler: router,
		Addr:    address,
	}, nil
}
