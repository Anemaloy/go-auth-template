package grpc

import (
	"auth/internal"
	userv1 "auth/internal/api/grpc/gen/course/auth/user/v1"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	storage          internal.UserStorage
	toDomainMapper   *toDomainMapper
	toProtobufMapper *toProtobufMapper
	logger           *zap.Logger
	userv1.UnimplementedUserAPIServer
}

func NewUserServer(
	storage internal.UserStorage,
	logger *zap.Logger,
) *UserServer {
	return &UserServer{
		storage:          storage,
		logger:           logger,
		toDomainMapper:   &toDomainMapper{},
		toProtobufMapper: &toProtobufMapper{},
	}
}

func (a *UserServer) Create(_ context.Context, req *userv1.CreateRequest) (*userv1.CreateResponse, error) {
	if req.GetPassword() != req.GetPasswordConfirm() {
		return nil, status.Error(codes.InvalidArgument, "password not match")
	}

	role, err := a.toDomainMapper.mapRole(req.GetRole())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("can't map category id: %s", err))
	}

	user, err := a.storage.Create(req.GetName(), req.GetEmail(), req.GetPassword(), role)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("can't create user: %s", err))
	}

	return &userv1.CreateResponse{Id: int64(user.Id)}, nil
}

func (a *UserServer) Update(_ context.Context, req *userv1.UpdateRequest) (*userv1.UpdateResponse, error) {
	user, err := a.storage.Get(internal.UserId(req.GetId()))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("can't get user: %s", err))
	}
	if user == nil {
		return nil, status.Error(codes.NotFound, "user doesn't exist")
	}

	_, err = a.storage.Update(internal.UserId(req.GetId()), req.GetName(), req.GetEmail())
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("can't update user: %s", err))
	}

	return &userv1.UpdateResponse{}, nil
}

func (a *UserServer) Get(_ context.Context, req *userv1.GetRequest) (*userv1.GetResponse, error) {
	id := req.GetId()
	user, err := a.storage.Get(internal.UserId(int(id)))
	if err != nil {
		if errors.Is(err, errors.New("user not found")) {
			return nil, status.Error(codes.NotFound, fmt.Sprintf("user not found: %s", err))
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("can't get user: %s", err))
	}

	pbUser, err := a.toProtobufMapper.mapUser(user)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("can't map user: %s", err))
	}

	return &userv1.GetResponse{User: pbUser}, nil
}

func (a *UserServer) Delete(_ context.Context, req *userv1.DeleteRequest) (*userv1.DeleteResponse, error) {
	err := a.storage.Delete(internal.UserId(req.GetId()))
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("can't delete user: %s", err))
	}

	return &userv1.DeleteResponse{}, nil
}
