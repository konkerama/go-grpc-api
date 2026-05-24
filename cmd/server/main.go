package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/konkerama/go-grpc-api/internal/db"

	"github.com/konkerama/go-grpc-api/internal/orders"
	pb "github.com/konkerama/go-grpc-api/pkg/pb/orders/v1"
	"github.com/lmittmann/tint"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// The gRPC interface requires your server to implement every single function defined in your proto file.
// If you add a new RPC method to your proto tomorrow but haven't written the Go code for it yet, your project will crash on compile.

// Embedding UnimplementedGreeterServer acts as a safety net.
// It provides "catch-all" placeholder methods for everything in your proto.
// It allows your server to compile immediately, returning an "Unimplemented" error to clients for any methods you haven't explicitly coded yet.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(_ context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

// type orders struct {
// 	pb.UnimplementedOrdersServer
// }
//
// func (s *orders) CreateOrder(_ context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderReply, error) {
// 	log.Printf("Received order for %v of type %v", in.GetQuantity(), in.GetProductName())
// 	return &pb.CreateOrderReply{OrderID: "test-id"}, nil
// }

type Closer func(context.Context) error

func main() {
	flag.Parse()

	var handler slog.Handler

	if os.Getenv("APP_ENV") == "production" {
		// Production stays structured JSON
		handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		})
	} else {
		// Local dev gets nice terminal colors via tint
		handler = tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelInfo,
			TimeFormat: time.Kitchen, // "3:04PM" or use "15:04:05" for 24h
			AddSource:  true,         // Keeps your file source lines, but tidier
		})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	// init db

	dbConfig := db.NewDBConfig()

	if err := db.RunMigrations(dbConfig.PGPool); err != nil {
		slog.Error("failed to run database migrations", "error", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		slog.Error("failed to listen", "error", err)
	}

	ordersModule := orders.Wire(dbConfig.PGPool)
	logger.Info("orders module dependencies wired successfully")

	s := grpc.NewServer()

	pb.RegisterGreeterServer(s, &server{})

	pb.RegisterOrdersServer(s, ordersModule.Controller)

	// gRPC Reflection is a built-in extension to the gRPC protocol that allows a server to describe its own API schema to anyone who asks.
	reflection.Register(s)

	slog.Info("server listening", "addr", lis.Addr())
	if err := s.Serve(lis); err != nil {
		slog.Error("failed to serve", "error", err)
	}

	// Block until we receive an OS signal to terminate
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	logger.Info("shutting down gRPC server gracefully...")
	s.GracefulStop()
	logger.Info("application exited cleanly")

}
