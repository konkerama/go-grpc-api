package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/konkerama/go-grpc-api/pkg/pb/orders/v1"
	"google.golang.org/grpc"
)

const (
	DB_URL = "postgres://postgres:postgres@localhost:5432/postgres"
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

type orders struct {
	pb.UnimplementedOrdersServer
}

func (s *orders) CreateOrder(_ context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderReply, error) {
	log.Printf("Received order for %v of type %v", in.GetQuantity(), in.GetProductName())
	return &pb.CreateOrderReply{OrderID: "test-id"}, nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterGreeterServer(s, &server{})
	pb.RegisterOrdersServer(s, &orders{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
