package grpc

import (
	"auth/internal"
	userv1 "auth/internal/api/grpc/gen/course/auth/user/v1"
	"errors"
	"fmt"
)

var pbRoles = map[internal.Role]userv1.Role{
	internal.RoleAdmin: userv1.Role_ROLE_ADMIN,
	internal.RoleUser:  userv1.Role_ROLE_USER,
}

type toProtobufMapper struct{}

func (a *toProtobufMapper) mapRole(category internal.Role) (userv1.Role, error) {
	pbRole, ok := pbRoles[category]
	if !ok {
		return userv1.Role_ROLE_INVALID, errors.New("category id is invalid")
	}

	return pbRole, nil
}

func (a *toProtobufMapper) mapUser(user *internal.User) (*userv1.User, error) {
	role, err := a.mapRole(user.Role)
	if err != nil {
		return nil, fmt.Errorf("can't map category id: %w", err)
	}

	return &userv1.User{
		Id:    int32(user.Id),
		Email: user.Email,
		Role:  role,
	}, nil
}
