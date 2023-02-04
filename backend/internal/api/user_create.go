package api

import (
	"bz.moh.epi/users/internal/auth"
	"context"
	"fmt"
)

type UserApi struct {
	UserStore auth.UserStore
}

func CreateUserApi(userStore auth.UserStore) *UserApi {
	return &UserApi{
		UserStore: userStore,
	}
}

// CreateUser creates a new user. Only admins are allowed to create users. Only admins can create users for their
// respective organizations. Admin users for NAC and MOHW can create users for any organization.
func (a *UserApi) CreateUser(ctx context.Context, user auth.CreateUserRequest) (*auth.User, error) {
	createdBy, err := a.UserStore.GetUserByEmail(ctx, user.CreatedBy)
	if err != nil {
		return nil, fmt.Errorf("error verifying person creating new user: %w", err)
	}

	if createdBy.Role != auth.AdminRole {
		return nil, fmt.Errorf("only admins can create users") //nolint: goerr113
	}
	if createdBy.Org != user.Org && (createdBy.Org != auth.MOHW && createdBy.Org != auth.NAC) {
		return nil, fmt.Errorf("only MOHW and NAC admins can create users for other organizations") //nolint: goerr113
	}

	if createdBy.Org == auth.NAC && user.Org == auth.MOHW {
		return nil, fmt.Errorf("only MOHW admins can create users for MOHW") //nolint: goerr113
	}

	u, err := a.UserStore.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}
	return u, nil
}
