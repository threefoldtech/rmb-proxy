package rmbproxy

import (
	"bytes"

	"github.com/threefoldtech/zos/pkg/substrate"
)

type flags struct {
	debug     string
	substrate string
	address   string
}
type MessageIdentifier struct {
	ID       string `json:"id"`
	Retqueue string `json:"retqueue"`
}

type Resp struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

// App is the main app objects
type App struct {
	resolver TwinExplorerResolver
}

type TwinExplorerResolver struct {
	client *substrate.Substrate
}

type twinClient struct {
	dstIP string
}

type TwinClient interface {
	SubmitMessage(msg bytes.Buffer) (string, error)
	GetResult(msgId MessageIdentifier) (string, error)
}
