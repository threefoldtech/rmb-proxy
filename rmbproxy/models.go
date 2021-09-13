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
	Retqueue string `json:"retqueue"`
}

// App is the main app objects
type App struct {
	resolver TwinResolver
}

type TwinExplorerResolver struct {
	client *substrate.Substrate
}

type twinClient struct {
	dstIP string
}

type TwinResolver interface {
	Resolve(twinID int) (TwinClient, error)
}

type TwinClient interface {
	SubmitMessage(msg bytes.Buffer) (string, error)
	GetResult(msgId MessageIdentifier) (string, error)
}
