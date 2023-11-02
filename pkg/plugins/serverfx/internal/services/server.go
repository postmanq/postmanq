package services

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/postmanq/postmanq/pkg/logfx/log"
	"github.com/postmanq/postmanq/pkg/plugins/serverfx/server"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"reflect"
)

type unionServer struct {
	ctx           context.Context
	logger        log.Logger
	cfg           server.Config
	grpcServer    *grpc.Server
	gatewayServer *http.Server
	mux           *runtime.ServeMux
	opts          []grpc.DialOption
	errors        chan error
}

func (s *unionServer) Register(descriptor server.Descriptor) error {
	reflect.ValueOf(descriptor.GRPCRegistrar).Call([]reflect.Value{
		reflect.ValueOf(s.grpcServer),
		reflect.ValueOf(descriptor.Server),
	})
	if descriptor.GRPCGatewayRegistrar != nil {
		return descriptor.GRPCGatewayRegistrar(s.ctx, s.mux, fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.GrpcPort), s.opts)
	}

	return nil
}

func (s *unionServer) Start() error {
	if s.cfg.GrpcPort > 0 {
		go func(logger log.Logger) {
			netAddress := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.GrpcPort)

			logger.Infof("start server at %s", netAddress)
			socket, err := net.Listen("tcp", netAddress)
			if err != nil {
				s.errors <- err
				return
			}

			s.errors <- s.grpcServer.Serve(socket)
		}(s.logger.Named("grpc_server"))
	}

	if s.cfg.HttpPort > 0 {
		go func(logger log.Logger) {
			netAddress := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.HttpPort)

			s.gatewayServer = &http.Server{
				Addr:    netAddress,
				Handler: s.mux,
			}

			logger.Infof("start gateway at %s", netAddress)

			s.errors <- s.gatewayServer.ListenAndServe()
		}(s.logger.Named("http_server"))
	}

	return <-s.errors
}
