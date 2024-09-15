package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/wal99d/blockr/node"
	"github.com/wal99d/blockr/proto"
	"google.golang.org/grpc"
)

func main() {
	node := node.NewNode()
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}
	proto.RegisterNodeServer(grpcServer, node)
	fmt.Println("node running on port:", "3000")
	go func() {

		for {

			time.Sleep(2 * time.Second)
			makeTransaction()
		}
	}()
	grpcServer.Serve(ln)

}

func makeTransaction() {

	client, err := grpc.Dial(":3000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	c := proto.NewNodeClient(client)

	v := &proto.Version{

		Version: "blockr-0.1",
		Height:  1,
	}

	_, err = c.Handshake(context.TODO(), v)
	if err != nil {
		log.Fatal(err)
	}

}
