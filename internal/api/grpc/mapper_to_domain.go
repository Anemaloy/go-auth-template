package grpc

import (
	"auth/internal"
	userv1 "auth/internal/api/grpc/gen/course/auth/user/v1"
	"errors"
)

var domainRoles = map[userv1.Role]internal.Role{
	userv1.Role_ROLE_USER:  internal.RoleUser,
	userv1.Role_ROLE_ADMIN: internal.RoleAdmin,
}

type toDomainMapper struct{}

func (a *toDomainMapper) mapRole(pbRole userv1.Role) (internal.Role, error) {
	role, ok := domainRoles[pbRole]
	if !ok {
		return 0, errors.New("role id is invalid")
	}

	return role, nil
}
