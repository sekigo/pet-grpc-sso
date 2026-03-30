package auth

import (
	"context"

	kir_sso_v1 "github.com/sekigo/pet-grpc/gen/go/sso"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email string, password string, AppId int) (token string, err error)

	RegisterNewUser(ctx context.Context, email string, password string) (userID int64, err error)

	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type serverAPI struct {
	kir_sso_v1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	kir_sso_v1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})

}

const (
	emptyValue = 0
)

func (s *serverAPI) Login(ctx context.Context, req *kir_sso_v1.LoginRequest) (*kir_sso_v1.LoginResponse, error) {
	// panic("impl")

	if err := LoginValidation(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.Email, req.Password, int(req.AppId))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &kir_sso_v1.LoginResponse{Token: token}, nil

}

func (s *serverAPI) Register(
	ctx context.Context,
	req *kir_sso_v1.RegisterRequest,
) (*kir_sso_v1.RegisterResponse, error) {

	if err := RegisterValidation(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal er")
	}

	return &kir_sso_v1.RegisterResponse{UserId: userID}, nil

}

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	req *kir_sso_v1.IsAdminRequest,
) (*kir_sso_v1.IsAdminResponse, error) {
	if err := IsAdminValidation(req); err != nil {
		return nil, err
	}

	flag, err := s.auth.IsAdmin(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &kir_sso_v1.IsAdminResponse{IsAdmin: flag}, nil
}

func LoginValidation(req *kir_sso_v1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	return nil
}

func RegisterValidation(req *kir_sso_v1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

func IsAdminValidation(req *kir_sso_v1.IsAdminRequest) error {

	if req.UserId == emptyValue {
		return status.Error(codes.Internal, "user_id is required")
	}

	return nil
}
