package iperf

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	api "github.com/BGrewell/go-iperf/api/go"
	"google.golang.org/grpc"
)

func NewController(port int) (controller *Controller, err error) {
	c := &Controller{
		Port:    port,
		clients: make(map[string]*Client),
		servers: make(map[string]*Server),
		clientLock: sync.Mutex{},
		serverLock: sync.Mutex{},
	}
	err = c.startListener()

	return c, err
}

// Controller is a helper in the go-iperf package that is designed to run on both the client and the server side. On the
// server side it listens for new gRPC connections, when a connection is made by a client the client can tell it to
// start a new iperf server instance. It will start a instance on an unused port and return the port number to the
// client. This allows the entire iperf setup and session to be performed from the client side.
type Controller struct {
	api.UnimplementedCommandServer
	Port int
	cmdClient api.CommandClient
	clientLock sync.Mutex
	serverLock sync.Mutex
	clients map[string]*Client
	servers map[string]*Server
}

// StartServer is the handler for the gRPC function StartServer()
func (c *Controller) GrpcRequestServer(context.Context, *api.StartServerRequest) (*api.StartServerResponse, error) {
	srv, err := c.NewServer()
	srv.SetOneOff(true)
	if err != nil {
		return nil, err
	}
	err = srv.Start()
	if err != nil {
		return nil, err
	}

	c.serverLock.Lock()
	c.servers[srv.Id] = srv
	c.serverLock.Unlock()
	
	reply := &api.StartServerResponse{
		Id:         srv.Id,
		ListenPort: int32(srv.Port()),
	}
	
	return reply, nil
}

// StartListener starts a command listener which is used to accept gRPC connections from another go-iperf controller
func (c *Controller) startListener() (err error) {

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "0.0.0.0", c.Port))
	if err != nil {
		return err
	}

	gs := grpc.NewServer()
	api.RegisterCommandServer(gs, c)

	go func() {
		err := gs.Serve(listener)
		if err != nil {
			log.Fatal(err)
		}
	}()

	time.Sleep(250 * time.Millisecond)
	return nil
}

// NewServer gets a new instance of an iperf server on a free port
func (c *Controller) NewServer() (server *Server, err error) {
	freePort, err := GetUnusedTcpPort()
	s := NewServer()
	s.SetPort(freePort)
	c.serverLock.Lock()
	c.servers[s.Id] = s
	c.serverLock.Unlock()
	return s, nil
}

// StopServer shuts down an iperf server and frees any actively used resources
func (c *Controller) StopServer(id string) (err error) {
	c.serverLock.Lock()
	delete(c.servers, id)
	c.serverLock.Unlock()
	return nil
}

// NewClient gets a new instance of an iperf client and also starts up a matched iperf server instance on the specified
// serverAddr. If it fails to connect to the gRPC interface of the controller on the remote side it will return an error
func (c *Controller) NewClient(serverAddr string) (client *Client, err error) {
	grpc, err := GetConnectedClient(serverAddr, c.Port)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
	defer cancel()
	reply, err := grpc.GrpcRequestServer(ctx, &api.StartServerRequest{})
	srvPort := int(reply.ListenPort)
	fmt.Printf("[!] server is listening on port %d\n", srvPort)
	
	cli := NewClient(serverAddr)
	cli.SetPort(srvPort)
	c.clientLock.Lock()
	c.clients[cli.Id] = cli
	c.clientLock.Unlock()

	return cli, nil
}

// StopClient will clean up the server side connection and shut down any actively used resources
func (c *Controller) StopClient(id string) (err error) {
	c.clientLock.Lock()
	delete(c.clients, id)
	c.clientLock.Unlock()
	return nil
}

func GetConnectedClient(addr string, port int) (client api.CommandClient, err error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", addr, port), grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second*2))
	if err != nil {
		return nil, err
	}
	client = api.NewCommandClient(conn)
	return client, nil
}

func GetUnusedTcpPort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}

	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil

}

//func GetUnusedUdpPort() (int, error) {
//	addr, err := net.ResolveUDPAddr("udp", "localhost:0")
//	if err != nil {
//		return 0, err
//	}
//
//	l, err := net.ListenUDP("udp", addr)
//	if err != nil {
//		return 0, err
//	}
//
//	defer l.Close()
//	return l.LocalAddr().(*net.UDPAddr).Port, nil
//}