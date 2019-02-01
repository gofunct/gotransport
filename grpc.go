package gotransport

import (
	"github.com/gofunct/gotransport/grpc/api"
	"github.com/gofunct/gotransport/grpc/engine"
)

func ServeGrpc(servers ...api.Server) error {
	s := engine.New(
		engine.WithDefaultLogger(),
		engine.WithServers(
			servers...,
		),
	)
	return s.Serve()
}
