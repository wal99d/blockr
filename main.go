package main

import (
	"context"
	"log"
	"time"

	"github.com/wal99d/blockr/node"
	"github.com/wal99d/blockr/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	makeNode(":3000", []string{})
	time.Sleep(2 * time.Second)
	makeNode(":4000", []string{":3000"})
	time.Sleep(4 * time.Second)
	makeNode(":6000", []string{":4000"})
	// makeTransaction()
	select {}
}

func makeNode(listenAddr string, bootstrapNodes []string) *node.Node {
	n := node.NewNode()
	go n.Start(listenAddr, bootstrapNodes)
	return n
}

func makeTransaction() {

	client, err := grpc.NewClient(":3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	c := proto.NewNodeClient(client)

	v := &proto.Version{

		Version:    "blockr-0.1",
		Height:     1,
		ListenAddr: ":8000",
	}

	_, err = c.Handshake(context.TODO(), v)
	if err != nil {
		log.Fatal(err)
	}

}
