package main

type Config struct {
	HttpPublicListen string `long:"http-public-listen" description:"Listening host:port for public http-server" env:"HTTP_PUBLIC_LISTEN" required:"true"`
	GrpcListen       string `long:"grpc-listen" description:"Listening host:port for grpc-server" env:"GRPC_LISTEN" required:"true"`
}
