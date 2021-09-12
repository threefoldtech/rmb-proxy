package rmbClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/threefoldtech/zos/pkg/substrate"
)

func submitURL(timeIP string) string {
	return fmt.Sprintf("http://%s:8051/zbus-cmd", timeIP)
}

func resultURL(timeIP string) string {
	return fmt.Sprintf("http://%s:8051/zbus-result", timeIP)
}

func NewTwinResolver(substrateURL string) (TwinResolver, error) {
	client, err := substrate.NewSubstrate(substrateURL)
	if err != nil {
		return nil, err
	}

	return &TwinExplorerResolver{
		client: client,
	}, nil
}

func (r TwinExplorerResolver) Resolve(timeID int) (TwinClient, error) {
	log.Debug().Int("twin", timeID).Msg("resolving twin")

	twin, err := r.client.GetTwin(uint32(timeID))
	if err != nil {
		return nil, err
	}
	log.Debug().Str("ip", twin.IP).Msg("resolved twin ip")

	return &twinClient{
		dstIP: twin.IP,
	}, nil
}

func (c *twinClient) readError(r io.Reader) string {
	var body struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r).Decode(&body); err != nil {
		return fmt.Sprintf("failed to read response body: %s", err)
	}

	return body.Message
}

func (c *twinClient) SubmitMessage(msg Message) error {
	var buffer bytes.Buffer
	if err := json.NewEncoder(&buffer).Encode(msg); err != nil {
		return err
	}
	resp, err := http.Post(submitURL(c.dstIP), "application/json", &buffer)
	// check on response for non-communication errors?
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to Submit the message: %s (%s)", resp.Status, c.readError(resp.Body))
	}

	return nil
}

func (c *twinClient) GetMessage(msg Message) (string, error) {
	var buffer bytes.Buffer
	if err := json.NewEncoder(&buffer).Encode(msg); err != nil {
		return "", err
	}
	resp, err := http.Post(resultURL(c.dstIP), "application/json", &buffer)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// body
		return "", fmt.Errorf("failed to send remote: %s (%s)", resp.Status, c.readError(resp.Body))
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	newStr := buf.String()

	return newStr, err
}
