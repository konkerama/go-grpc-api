package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/konkerama/go-grpc-api/pkg/pb/orders/v1"
)

const (
	defaultName = "world"
	product     = "apples"
	quantity    = 10
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()

	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)
	orders := pb.NewOrdersClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greetings: %s", r.GetMessage())

	o, err := orders.CreateOrder(ctx, &pb.CreateOrderRequest{ProductName: product, Quantity: quantity})
	if err != nil {
		log.Fatalf("could not create order: %v", err)
	}
	log.Printf("Greetings: %s", o.GetOrderID())

}
