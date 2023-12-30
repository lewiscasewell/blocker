package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/lewiscasewell/blocker/node"
	"github.com/lewiscasewell/blocker/proto"
	"google.golang.org/grpc"
)

func main() {
	node := node.NewNode()

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	proto.RegisterNodeServer(grpcServer, node)
	fmt.Println("node running on port 8080")

	go func() {
		for {
			time.Sleep(2 * time.Second)
			mainTransaction()
		}
	}()

	grpcServer.Serve(ln)

}

func mainTransaction() {
	client, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	c := proto.NewNodeClient(client)

	version := &proto.Version{
		Version: "blocker-0.0.1",
		Height:  1,
	}

	_, err = c.Handshake(context.TODO(), version)
	if err != nil {
		log.Fatal(err)
	}
}
