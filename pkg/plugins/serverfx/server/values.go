package server

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Descriptor struct {
	Server               interface{}
	GRPCRegistrar        interface{}
	GRPCGatewayRegistrar func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
}

type Config struct {
	Host     string `yaml:"host"`
	GrpcPort int    `yaml:"grpc_port"`
	HttpPort int    `yaml:"http_port"`
}
