package node

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/wal99d/blockr/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type Node struct {
	peerLock   sync.RWMutex
	peers      map[proto.NodeClient]*proto.Version
	version    string
	listenAddr string
	logger     *zap.SugaredLogger
	proto.UnimplementedNodeServer
}

func NewNode() *Node {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.TimeKey = ""
	logger, _ := loggerConfig.Build()
	return &Node{
		peers:   make(map[proto.NodeClient]*proto.Version),
		version: "blockr-0.1",
		logger:  logger.Sugar(),
	}
}

func (n *Node) addPeer(c proto.NodeClient, v *proto.Version) {

	n.peerLock.Lock()
	defer n.peerLock.Unlock()
	n.logger.Debugw("new peer connected", "address", v.ListenAddr, "height", v.Height)
	n.peers[c] = v
}

func (n *Node) deletePeer(c proto.NodeClient) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()
	delete(n.peers, c)
}

func (n *Node) Start(listenAddr string) error {
	n.listenAddr = listenAddr
	var (
		opts       = []grpc.ServerOption{}
		grpcServer = grpc.NewServer(opts...)
	)
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	proto.RegisterNodeServer(grpcServer, n)

	fmt.Println("node running on port:", listenAddr)
	return grpcServer.Serve(ln)
}

func (n *Node) Handshake(ctx context.Context, version *proto.Version) (*proto.Version, error) {

	c, err := makeNodeClient(version.ListenAddr)
	if err != nil {
		return nil, err
	}
	n.addPeer(c, version)
	return n.getVersion(), nil
}

func (n *Node) BootstrapNetwork(addrs []string) error {

	for _, addr := range addrs {
		c, err := makeNodeClient(addr)
		if err != nil {
			return err
		}
		v, err := c.Handshake(context.Background(), n.getVersion())
		if err != nil {
			fmt.Println("handshake error:", err)
			continue
		}
		n.addPeer(c, v)
	}
	return nil
}

func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*proto.Ack, error) {

	peer, _ := peer.FromContext(ctx)
	fmt.Println("received tx from:", peer)
	return &proto.Ack{}, nil
}

func makeNodeClient(listenAddr string) (proto.NodeClient, error) {

	conn, err := grpc.Dial(listenAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return proto.NewNodeClient(conn), nil
}

func (n *Node) getVersion() *proto.Version {
	return &proto.Version{
		Version:    "blockr-0.1",
		Height:     1,
		ListenAddr: n.listenAddr,
	}
}
