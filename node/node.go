package node

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/lewiscasewell/blocker/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type Node struct {
	version    string
	listenAddr string
	logger     zap.SugaredLogger
	peerLock   sync.RWMutex
	peers      map[proto.NodeClient]*proto.Version
	proto.UnimplementedNodeServer
}

func NewNode() *Node {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.DisableStacktrace = true
	loggerConfig.DisableCaller = true
	loggerConfig.DisableStacktrace = true
	loggerConfig.DisableCaller = true

	logger, _ := loggerConfig.Build()

	return &Node{
		peers:   make(map[proto.NodeClient]*proto.Version),
		version: "blocker-0.0.1",
		logger:  *logger.Sugar(),
	}
}

func (n *Node) AddPeer(client proto.NodeClient, v *proto.Version) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()

	n.logger.Debugf("[%s] new peer connected (%s) - height (%d)\n", n.listenAddr, v.ListenAddr, v.Height)

	n.peers[client] = v
}

func (n *Node) RemovePeer(client proto.NodeClient) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()

	delete(n.peers, client)
}

func (n *Node) BootstrapNetwork(addrs []string) error {
	for _, addr := range addrs {
		c, err := makeNodeClient(addr)
		if err != nil {
			fmt.Println("makeNodeClient error", err)
			continue
		}

		version, err := c.Handshake(context.Background(), n.getVersion())
		if err != nil {
			fmt.Println("Handshake error", err)
			continue
		}

		n.AddPeer(c, version)
	}

	return nil

}

func (n *Node) Start(listenAddr string) error {
	n.listenAddr = listenAddr
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	proto.RegisterNodeServer(grpcServer, n)

	return grpcServer.Serve(ln)
}

func (n *Node) Handshake(ctx context.Context, v *proto.Version) (*proto.Version, error) {

	p, _ := peer.FromContext(ctx)

	c, err := makeNodeClient(v.ListenAddr)
	if err != nil {
		return nil, err
	}

	n.AddPeer(c, v)

	fmt.Printf("received handshake from %s: %+v\n", p.Addr, v)

	return n.getVersion(), nil
}

func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*proto.Ack, error) {
	peer, _ := peer.FromContext(ctx)
	fmt.Println("received transaction", peer)
	return &proto.Ack{}, nil
}

func (n *Node) getVersion() *proto.Version {
	return &proto.Version{
		Version:    "blocker-0.0.1",
		Height:     0,
		ListenAddr: n.listenAddr,
	}
}

func makeNodeClient(listenAddr string) (proto.NodeClient, error) {
	a := fmt.Sprintf("localhost%s", listenAddr)
	conn, err := grpc.Dial(a, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return proto.NewNodeClient(conn), nil
}
