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
	err := a.checkUserPermissions(ctx, permissionRequest{
		Org:         user.Org,
		RequestedBy: user.CreatedBy,
	})
	if err != nil {
		return nil, fmt.Errorf("error checking user permissions: %w", err)
	}

	u, err := a.UserStore.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}
	return u, nil
}

type UpdateUserRequest struct {
	User      *auth.User
	UpdatedBy string
}

func (a *UserApi) UpdateUser(ctx context.Context, request UpdateUserRequest) error {
	err := a.checkUserPermissions(ctx, permissionRequest{
		Org:         request.User.Org,
		RequestedBy: request.UpdatedBy,
	})
	if err != nil {
		return fmt.Errorf("error checking user permissions: %w", err)
	}

	oldUserData, err := a.UserStore.GetUserByID(ctx, request.User.ID)
	if err != nil {
		return fmt.Errorf("error retrieving user to udpate: %w", err)
	}
	if request.User.Email == request.UpdatedBy && oldUserData.Role != request.User.Role {
		return fmt.Errorf("user cannot change their own role") //nolint: goerr113
	}

	if err := a.UserStore.UpdateUser(ctx, request.User, request.UpdatedBy); err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}

type permissionRequest struct {
	Org         auth.Org
	RequestedBy string
}

func (a *UserApi) checkUserPermissions(ctx context.Context, user permissionRequest) error {
	requestedBy, err := a.UserStore.GetUserByEmail(ctx, user.RequestedBy)
	if err != nil {
		return fmt.Errorf("error verifying person creating or editing user: %w", err)
	}

	if requestedBy.Role != auth.AdminRole {
		return fmt.Errorf("only admins can create users") //nolint: goerr113
	}
	if requestedBy.Org != user.Org && (requestedBy.Org != auth.MOHW && requestedBy.Org != auth.NAC) {
		return fmt.Errorf("only MOHW and NAC admins can create and edit users for other organizations") //nolint: goerr113
	}
	if requestedBy.Org == auth.NAC && user.Org == auth.MOHW {
		return fmt.Errorf("only MOHW admins can create and edit users for MOHW") //nolint: goerr113
	}
	return nil
}
