package node

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/wal99d/blockr/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type Node struct {
	peerLock sync.RWMutex
	peers    map[proto.NodeClient]bool
	version  string
	proto.UnimplementedNodeServer
}

func NewNode() *Node {
	return &Node{
		peers:   make(map[proto.NodeClient]bool),
		version: "blockr-0.1",
	}
}

func (n *Node) addPeer(c proto.NodeClient) {

	n.peerLock.Lock()
	defer n.peerLock.Unlock()
	n.peers[c] = true
}

func (n *Node) deletePeer(c proto.NodeClient) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()
	delete(n.peers, c)
}

func (n *Node) Start(listenAddr string) error {
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	proto.RegisterNodeServer(grpcServer, n)

	fmt.Println("node running on port:", listenAddr)
	return grpcServer.Serve(ln)
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

func makeNodeClient(listenAddr string) (proto.NodeClient, error) {

	conn, err := grpc.NewClient(listenAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return proto.NewNodeClient(conn), nil
}
