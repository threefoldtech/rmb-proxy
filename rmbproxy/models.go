package rmbproxy

import (
	"bytes"

	"github.com/threefoldtech/zos/pkg/substrate"
)

type MessageIdentifier struct {
	ID       string `json:"id"`
	Retqueue string `json:"retqueue"`
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
