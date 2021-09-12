package rmbClient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func errorReply(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	fmt.Fprintf(w, "{\"status\": \"error\", \"message\": \"%s\"}", message)
}

func (a *App) sendMessage(w http.ResponseWriter, r *http.Request) {
	twinIDString := mux.Vars(r)["node_id"]
	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		errorReply(w, http.StatusBadRequest, "couldn't parse json")
		return
	}

	twinID, err := strconv.Atoi(twinIDString)
	if err != nil {
		errorReply(w, http.StatusBadRequest, "Invalid twinId")
		return
	}

	c, err := a.resolver.Resolve(twinID)
	if err != nil {
		log.Error().Err(err).Msg("failed to create mux server")
	}

	err = c.SubmitMessage(msg)
	if err != nil {
		errorReply(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msg.Retqueue))
}

func (a *App) getMessage(w http.ResponseWriter, r *http.Request) {
	twinIDString := mux.Vars(r)["node_id"]
	retueue := mux.Vars(r)["retqueue"]
	var msg Message
	msg.Retqueue = retueue

	twinID, err := strconv.Atoi(twinIDString)
	if err != nil {
		errorReply(w, http.StatusBadRequest, "Invalid twinId")
		return
	}

	c, err := a.resolver.Resolve(twinID)
	if err != nil {
		log.Error().Err(err).Msg("failed to create mux server")
	}

	data, err := c.GetMessage(msg)
	if err != nil {
		errorReply(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data))
}

// Setup is the server and do initial configurations
func Setup(router *mux.Router, substrate string, hostAddress string) error {
	c := cache.New(10*time.Minute, 15*time.Minute)
	resolver, err := NewTwinResolver(substrate)
	if err != nil {
		return errors.Wrap(err, "couldn't get a client to explorer resolver")
	}

	a := &App{
		resolver: resolver,
		lruCache: c,
	}
	log.Info().Str("listening on", hostAddress).Msg("Server started ...")
	router.HandleFunc("/twin/{node_id:[0-9]+}", a.sendMessage)
	router.HandleFunc("/twin/{node_id:[0-9]+}/{retqueue}", a.getMessage)

	return nil
}
