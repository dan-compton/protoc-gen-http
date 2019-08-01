package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/flowup/protoc-gen-http/examples/hello"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"sync"
)

var (
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:9090", "gRPC server endpoint")
)

func run() error {
	server := &hello.GreeterService{}

	s := grpc.NewServer()
	hello.RegisterGreeterServer(s, server)

	reflection.Register(s)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	fmt.Println("Starting")

	go func() {
		lis, err := net.Listen("tcp", ":9090")
		if err != nil {
			panic(err)
		}

		fmt.Println("Starting gRPC server")
		if err := s.Serve(lis); err != nil {
			fmt.Errorf("grpc: failed to serve: %v", err)
		}

		wg.Done()
	}()

	go func() {
		// Register gRPC server endpoint
		// Note: Make sure the gRPC server is running properly and accessible
		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{grpc.WithInsecure()}
		err := hello.RegisterGreeterHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
		if err != nil {
			panic(err)
		}

		// Start HTTP server (and proxy calls to gRPC server endpoint)
		err = http.ListenAndServe(":8080", mux)
		if err != nil {
			panic(err)
		}

		wg.Done()
	}()

	wg.Wait()
	return nil
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		fmt.Errorf("%+v", err)
	}
}
