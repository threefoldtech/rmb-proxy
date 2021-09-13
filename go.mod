module github.com/threefoldtech/rmb_proxy_server

go 1.16

require (
	github.com/gorilla/mux v1.8.0
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.23.0
	github.com/threefoldtech/zos v0.5.2
)

replace github.com/threefoldtech/zos => ../zos
