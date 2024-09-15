package node

import (
	"context"
	"fmt"

	"github.com/wal99d/blockr/proto"
	"google.golang.org/grpc/peer"
)

type Node struct {
	version string
	proto.UnimplementedNodeServer
}

func NewNode() *Node {
	return &Node{

		version: "blockr-0.1",
	}
}

func (n *Node) Handshake(ctx context.Context, version *proto.Version) (*proto.Version, error) {

	v := &proto.Version{

		Version: n.version,
		Height:  100,
	}

	peer, _ := peer.FromContext(ctx)
	fmt.Printf("received version from %s: %+v\n", version, peer.Addr)
	return v, nil
}

func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*proto.Ack, error) {

	peer, _ := peer.FromContext(ctx)
	fmt.Println("received tx from:", peer)
	return &proto.Ack{}, nil
}
