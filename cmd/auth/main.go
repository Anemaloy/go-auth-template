package main

import (
	"auth/internal"
	grpc_internal "auth/internal/api/grpc"
	userv1 "auth/internal/api/grpc/gen/course/auth/user/v1"
	"auth/internal/inmemory"
	"context"
	"fmt"
	"github.com/jessevdk/go-flags"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sync"
)

func main() {
	var cfg Config
	parser := flags.NewParser(&cfg, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		log.Fatal("Failed to parse config.", err)
	}

	logger := zap.Must(zap.NewDevelopment())

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	userStorage := inmemory.NewUserStorage()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := startGRPCServer(ctx, cfg.GrpcListen, userStorage, logger)
		if err != nil {
			logger.Error("can't start gRPC server or server return error while working", zap.Error(err))
		}
	}()

	go func() {
		<-ctx.Done()
	}()
	wg.Wait()
}

// startGRPCServer запускает gRPC сервер
func startGRPCServer(
	ctx context.Context,
	listen string,
	userStorage internal.UserStorage,
	logger *zap.Logger,
) error {
	lis, err := net.Listen("tcp", listen)
	if err != nil {
		return fmt.Errorf("failed to listen GRPC server: %w", err)
	}

	s := grpc.NewServer()
	userv1.RegisterUserAPIServer(s, grpc_internal.NewUserServer(userStorage, logger))
	reflection.Register(s)

	logger.Info("gRPC started", zap.String("address", listen))

	go func() {
		<-ctx.Done()
		s.GracefulStop()
	}()

	return s.Serve(lis)
}
