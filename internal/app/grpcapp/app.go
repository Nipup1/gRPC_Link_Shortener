package grpcapp

import (
	"fmt"
	"go/link_shortener/internal/grpc_handlers"
	"go/link_shortener/internal/service/link"
	"net"

	"google.golang.org/grpc"
)

type App struct {
	gRPCServer *grpc.Server
	port int
}

func New(port int, linkShortener *link.Link) *App{
	gRPCServer := grpc.NewServer()

	grpc_handlers.Register(gRPCServer, linkShortener)

	return &App{
		gRPCServer: gRPCServer,
		port: port,
	}
}

func (a *App) MustRun(){
	if err := a.Run(); err != nil{
		panic(err)
	}
}

func (a *App) Run() error{
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil{
		return err
	}

	fmt.Printf("Server is listening on port: %d", a.port)

	if err := a.gRPCServer.Serve(l); err != nil{
		return err
	}

	return nil
}

func (a *App) Stop(){
	a.gRPCServer.GracefulStop()
}