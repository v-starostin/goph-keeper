package main

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/v-starostin/goph-keeper/internal/handler"
	"github.com/v-starostin/goph-keeper/pkg/pb"
)

func main() {
	l, err := net.Listen("tcp", ":9090")
	if err != nil {
		fmt.Println(err)
		return
	}

	authHandler := handler.New(nil)

	server := grpc.NewServer()
	pb.RegisterAuthServer(server, authHandler)

	if err := server.Serve(l); err != nil {
		fmt.Println(err)
		return
	}
}
