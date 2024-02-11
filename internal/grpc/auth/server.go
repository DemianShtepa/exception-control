package auth

import (
	"context"
	"errors"
	"github.com/DemianShtepa/exception-control/internal/services/auth"
	"github.com/DemianShtepa/exception-control/protos/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Register(ctx context.Context, email, password string) error
	Login(ctx context.Context, email, password string) (string, error)
}

type serverApi struct {
	gen.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	gen.RegisterAuthServer(gRPC, &serverApi{auth: auth})
}

func (s *serverApi) Register(c context.Context, r *gen.RegisterRequest) (*gen.RegisterResponse, error) {
	email := r.GetEmail()
	if email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	password := r.GetPassword()
	if password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "password is required")
	}

	err := s.auth.Register(c, email, password)
	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &gen.RegisterResponse{}, nil
}

func (s *serverApi) Login(c context.Context, r *gen.LoginRequest) (*gen.LoginResponse, error) {
	email := r.GetEmail()
	if email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	password := r.GetPassword()
	if password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "password is required")
	}

	token, err := s.auth.Login(c, email, password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &gen.LoginResponse{
		Token: token,
	}, nil
}
