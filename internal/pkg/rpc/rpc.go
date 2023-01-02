package rpc

import (
	"net"
	"net/rpc"

	"github.com/martijnkorbee/gobaboon/internal/pkg/server"
)

type RPCServer struct {
	// defaults to 4004
	Port string
}

func NewRPCServer(port string) *RPCServer {
	if port == "" {
		port = "4004"
	}

	return &RPCServer{
		Port: port,
	}
}

func (r *RPCServer) Run() error {

	err := rpc.Register(new(RPCServer))
	if err != nil {
		return err
	}

	listen, err := net.Listen("tcp", "127.0.0.1:"+r.Port)
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
