package handler

import (
	"context"

	"github.com/v-starostin/goph-keeper/internal/pb"
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

func (a *Auth) Authenticate(context.Context, *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	//return nil, status.Errorf(codes.Unimplemented, "method Authenticate not implemented")
	return &pb.AuthenticateResponse{
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
	}, nil
}
func (a *Auth) Refresh(context.Context, *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	return &pb.RefreshResponse{
		AccessToken:  "access_refresh_token",
		RefreshToken: "refresh_refresh_token",
	}, nil
}
