package gotransport

import (
	"gocloud.dev/server"
)

type OptFunc func(*server.Options)

func New(opts ...OptFunc) *server.Server {
	o := &server.Options{
		RequestLogger:         nil,
		HealthChecks:          nil,
		TraceExporter:         nil,
		DefaultSamplingPolicy: nil,
		Driver:                nil,
	}
	for _, f := range opts {
		f(o)
	}
	return server.New(o)
}
