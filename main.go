package main

import (
	"context"
	"log"

	"github.com/lewiscasewell/blocker/node"
	"github.com/lewiscasewell/blocker/proto"
	"google.golang.org/grpc"
)

func main() {
	makeNode(":4000", []string{})
	makeNode(":4001", []string{":4000"})

	select {}
}

func makeNode(listenAddr string, bootstrapNodes []string) *node.Node {
	n := node.NewNode()

	go n.Start(listenAddr)
	if len(bootstrapNodes) > 0 {
		if err := n.BootstrapNetwork(bootstrapNodes); err != nil {
			log.Fatal(err)
		}
	}

	return n

}

func makeTransaction() {
	client, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	c := proto.NewNodeClient(client)

	version := &proto.Version{
		Version:    "blocker-0.0.1",
		Height:     1,
		ListenAddr: ":4000",
	}

	_, err = c.Handshake(context.TODO(), version)
	if err != nil {
		log.Fatal(err)
	}
}
