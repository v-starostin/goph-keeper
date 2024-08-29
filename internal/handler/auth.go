package handler

import (
	"context"

	"github.com/v-starostin/goph-keeper/pkg/pb"
)

type AuthService interface {
	Authenticate()
	Refresh()
}

type Auth struct {
	pb.UnimplementedAuthServer
	service AuthService
}

func New(s AuthService) *Auth {
	return &Auth{service: s}
}

func (a *Auth) Authenticate(ctx context.Context, in *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	//return nil, status.Errorf(codes.Unimplemented, "method Authenticate not implemented")
	username := in.GetUsername()

	return &pb.AuthenticateResponse{
		AccessToken:  username + "'s " + "access_token",
		RefreshToken: username + "'s " + "refresh_token",
	}, nil
}
func (a *Auth) Refresh(ctx context.Context, in *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	return &pb.RefreshResponse{
		AccessToken:  "access_refresh_token",
		RefreshToken: "refresh_refresh_token",
	}, nil
}
