package handler

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/v-starostin/goph-keeper/pkg/pb"
)

type AuthService interface {
	Register(ctx context.Context, username, password string) error
	Authenticate(ctx context.Context, username, password string) (string, string, error)
	Refresh(ctx context.Context, access, refresh string) (string, string, error)
}

type Auth struct {
	pb.UnimplementedAuthServer
	service AuthService
}

func New(s AuthService) *Auth {
	return &Auth{service: s}
}

func (a *Auth) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	u := in.GetUsername()
	p := in.GetPassword()

	var errs []error
	if u == "" {
		errs = append(errs, errors.New("username is empty"))
	}
	if p == "" {
		errs = append(errs, errors.New("password is empty"))
	}
	if len(errs) > 0 {
		return nil, status.Error(codes.InvalidArgument, errors.Join(errs...).Error())
	}

	if err := a.service.Register(ctx, u, p); err != nil {
		log.Println("Registration error", err)
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code.Name() == "unique_violation" {
				return nil, status.Error(codes.AlreadyExists, "User already exists")
			}
		}

		return nil, status.Error(codes.Internal, "Internal server error")
	}

	accessToken, refreshToken, err := a.service.Authenticate(ctx, u, p)
	if err != nil {
		log.Println("Authentication error", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &pb.RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *Auth) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	u := in.GetUsername()
	p := in.GetPassword()

	var errs []error
	if u == "" {
		errs = append(errs, errors.New("username is empty"))
	}
	if p == "" {
		errs = append(errs, errors.New("password is empty"))
	}
	if len(errs) > 0 {
		return nil, status.Error(codes.InvalidArgument, errors.Join(errs...).Error())
	}

	accessToken, refreshToken, err := a.service.Authenticate(ctx, u, p)
	if err != nil {
		log.Println("Authentication error", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.Unauthenticated, "Invalid credentials")
		}

		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &pb.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *Auth) Refresh(ctx context.Context, in *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	access := in.GetAccessToken()
	refresh := in.GetRefreshToken()

	var errs []error
	if access == "" {
		errs = append(errs, errors.New("access token is empty"))
	}
	if refresh == "" {
		errs = append(errs, errors.New("refresh token is empty"))
	}
	if len(errs) > 0 {
		return nil, status.Error(codes.InvalidArgument, errors.Join(errs...).Error())
	}

	newAccess, newRefresh, err := a.service.Refresh(ctx, access, refresh)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &pb.RefreshResponse{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
	}, nil
}
