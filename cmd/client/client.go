package main

import (
	"context"
	"fmt"

	"github.com/v-starostin/goph-keeper/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(":9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	client := pb.NewAuthClient(conn)
	ctx := context.Background()
	res, err := client.Authenticate(ctx, &pb.AuthenticateRequest{
		Username: "alice",
		Password: "qwerty",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("res.AccessToken", res.GetAccessToken())
	fmt.Println("res.RefreshToken", res.GetRefreshToken())

	res2, err := client.Refresh(ctx, &pb.RefreshRequest{
		AccessToken:  "alice_access_token",
		RefreshToken: "alice_refresh_token",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("res2.Refresh", res2.GetRefreshToken())
	fmt.Println("res2.Access", res2.GetAccessToken())
}
