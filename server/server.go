package main

import (
	chat "awesomeProject/publish"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"time"
)

var (
	port = flag.Int("port", 50051, "The server port")
	addr = flag.String("addr", "localhost:50052", "the address to connect to")
	name = flag.String("name", "Ping", "Name to greet")
)

type server struct {
	chat.UnimplementedChatServiceServer
}

func (s *server) SayHello(ctx context.Context, in *chat.Message) (*chat.Message, error) {
	log.Printf("Received: %v", in.GetBody())
	go sendPing()
	return &chat.Message{Body: "Server 1 response "}, nil
}

func sendPing() {
	time.Sleep(1 * time.Second)
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := chat.NewChatServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &chat.Message{Body: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetBody())
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	chat.RegisterChatServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
