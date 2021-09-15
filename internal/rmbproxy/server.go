package rmbproxy

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	// swagger configuration
	_ "github.com/threefoldtech/rmb_proxy_server/docs"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger"
)

func errorReply(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	fmt.Fprintf(w, "{\"status\": \"error\", \"message\": \"%s\"}", message)
}

// NewTwinClient : create new TwinClient
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

// sendMessage godoc
// @Summary submit the message
// @Description submit the message
// @Tags Result
// @Accept  json
// @Produce  json
// @Param msg body Message true "rmb.Message"
// @Param twin_id path int true "twin id"
// @Success 200 {object} MessageIdentifier
// @Router /twin/{twin_id} [post]
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
		log.Error().Err(err).Msg("failed to create TwinClient")
		errorReply(w, http.StatusInternalServerError, "failed to create TwinClient")
		return
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

// getResult godoc
// @Summary Get the message result
// @Description Get the message result
// @Tags Result
// @Accept  json
// @Produce  json
// @Param twin_id path int true "twin id"
// @Param retqueue path string true "message retqueue"
// @Success 200 {object} string ""
// @Router /twin/{twin_id}/{retqueue} [get]
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

// CreateServer : Create the app server
// @title RMB proxy API
// @version 1.0
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email soberkoder@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
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
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	return &http.Server{
		Handler: router,
		Addr:    address,
	}, nil
}
