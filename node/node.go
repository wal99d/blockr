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
	// handle the logic where we decide to accept or drop
	// the incoming conn
	n.peers[c] = v
	// connect to all peers in the received list of peers.
	if len(v.PeerList) > 0 {

		go n.bootstrapNetwork(v.PeerList)
	}
	n.logger.Debugw("new peer successfully connected",
		"we", n.listenAddr,
		"remoteNode", v.ListenAddr,
		"height", v.Height)
}

func (n *Node) deletePeer(c proto.NodeClient) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()
	delete(n.peers, c)
}

func (n *Node) Start(listenAddr string, bootstrapNodes []string) error {
	n.listenAddr = listenAddr
	var (
		opts       = []grpc.ServerOption{}
		grpcServer = grpc.NewServer(opts...)
	)
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		n.logger.Errorw("error", err)
		return err
	}
	proto.RegisterNodeServer(grpcServer, n)

	n.logger.Infow("node started..", "port", n.listenAddr)
	//bootstrap the network with a list of already known nodes
	// in the network
	if len(bootstrapNodes) > 0 {
		go n.bootstrapNetwork(bootstrapNodes)
	}
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

func (n *Node) bootstrapNetwork(addrs []string) error {

	for _, addr := range addrs {
		if !n.canConnectWith(addr) {

			continue
		}
		n.logger.Debugw("dailing remote node", "we", n.listenAddr, "remote", addr)
		c, v, err := n.dialRemoteNode(addr)
		if err != nil {
			return err
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
		PeerList:   n.getPeerList(),
	}
}

func (n *Node) getPeerList() []string {

	n.peerLock.RLock()
	defer n.peerLock.RUnlock()

	peers := []string{}
	for _, version := range n.peers {

		peers = append(peers, version.ListenAddr)
	}
	return peers
}

func (n *Node) dialRemoteNode(addr string) (proto.NodeClient, *proto.Version, error) {

	c, err := makeNodeClient(addr)
	if err != nil {
		return nil, nil, err
	}
	v, err := c.Handshake(context.Background(), n.getVersion())
	if err != nil {
		return nil, nil, err
	}
	return c, v, nil
}

func (n *Node) canConnectWith(addr string) bool {
	if n.listenAddr == addr {
		return false
	}
	connectedPeers := n.getPeerList()

	for _, connectedAddr := range connectedPeers {

		if addr == connectedAddr {
			return false
		}
	}
	return true
}
