package rpc

import (
	"fmt"
	"github.com/martijnkorbee/gobaboon/pkg/server"
	"net"
	"net/rpc"
)

type RPCServer struct {
	// same as server host
	Host string
	// defaults to 4004
	Port string
}

func NewRPCServer(host, port string) *RPCServer {
	if port == "" {
		port = "4004"
	}

	return &RPCServer{
		Host: host,
		Port: port,
	}
}

func (r *RPCServer) Run() error {

	err := rpc.Register(new(RPCServer))
	if err != nil {
		return err
	}

	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%s", r.Host, r.Port))
	if err != nil {
		return err
	}

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}

func (r *RPCServer) SetMaintenanceMode(on bool, resp *string) error {
	if on {
		server.MaintenanceMode = true
		*resp = "Server in maintenance mode."
	}

	if !on {
		server.MaintenanceMode = false
		*resp = "Server is live."
	}

	return nil
}
