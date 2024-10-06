package main

import (
	"context"
	"log"

	"github.com/wal99d/blockr/node"
	"github.com/wal99d/blockr/proto"
	"google.golang.org/grpc"
)

func main() {
	// node := node.NewNode()
	//
	// go func() {
	//
	// 	for {
	//
	// 		time.Sleep(2 * time.Second)
	// 		makeTransaction()
	// 	}
	// }()
	//

	makeNode(":3000", []string{})
	makeNode(":4000", []string{":3000"})
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

	client, err := grpc.Dial(":3000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	c := proto.NewNodeClient(client)

	v := &proto.Version{

		Version:    "blockr-0.1",
		Height:     1,
		ListenAddr: ":4000",
	}

	_, err = c.Handshake(context.TODO(), v)
	if err != nil {
		log.Fatal(err)
	}

}
