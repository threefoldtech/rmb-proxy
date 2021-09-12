package rmbproxy

import (
	"github.com/patrickmn/go-cache"
	"github.com/threefoldtech/zos/pkg/substrate"
)

// DefaultExplorerURL is the default explorer graphql url
const DefaultExplorerURL string = "https://explorer.devnet.grid.tf/graphql/"

// App is the main app objects
type App struct {
	resolver TwinResolver
	lruCache *cache.Cache
}

type Message struct {
	Version    int    `json:"ver"`
	ID         string `json:"uid"`
	Command    string `json:"cmd"`
	Expiration int64  `json:"exp"`
	Retry      int    `json:"try"`
	Data       string `json:"dat"`
	TwinSrc    int    `json:"src"`
	TwinDst    []int  `json:"dst"`
	Retqueue   string `json:"ret"`
	Schema     string `json:"shm"`
	Epoch      int64  `json:"now"`
	Err        string `json:"err"`
}

type TwinExplorerResolver struct {
	client *substrate.Substrate
}

type twinClient struct {
	dstIP string
}

type TwinResolver interface {
	Resolve(timeID int) (TwinClient, error)
}

type TwinClient interface {
	SubmitMessage(msg Message) error
	GetMessage(msg Message) (string, error)
}
