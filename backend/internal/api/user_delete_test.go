package api

import (
	"bz.moh.epi/users/internal/auth"
	"context"
	"testing"
)

func TestUserApi_MOHW_AdminUsersCanDeleteAnyUser(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	adminUser := createTestUser(t, ctx, userStore, auth.MOHW, auth.AdminRole)
	mohwUser := createTestUser(t, ctx, userStore, auth.MOHW, auth.AdminRole)
	nacUser := createTestUser(t, ctx, userStore, auth.NAC, auth.AdminRole)
	bflaUser := createTestUser(t, ctx, userStore, auth.BFLA, auth.PeerNavigatorRole)
	csoUser := createTestUser(t, ctx, userStore, auth.GOJOVEN, auth.SrRole)

	var testCases = []struct {
		name      string
		inputUser *auth.User
	}{
		{
			name:      "MOHW Admin can delete MOHW User",
			inputUser: mohwUser,
		},
		{
			name:      "MOHW Admin can delete NAC User",
			inputUser: nacUser,
		},
		{
			name:      "MOHW Admin can delete BFLA User",
			inputUser: bflaUser,
		},
		{
			name:      "MOHW Admin can delete CSO User",
			inputUser: csoUser,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := userApi.DeleteUser(ctx, DeleteUserRequest{
				DeletedBy: adminUser.Email,
				ID:        tc.inputUser.ID,
			})
			if err != nil {
				t.Errorf("Unexpected error deleting user: %v", err)
			}
		})
	}
}

func TestUserApi_NAC_AdminUserCanDeleteUsers(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	nacAdminUser := createTestUser(t, ctx, userStore, auth.NAC, auth.AdminRole)
	nacUser := createTestUser(t, ctx, userStore, auth.NAC, auth.AdminRole)
	bflaUser := createTestUser(t, ctx, userStore, auth.BFLA, auth.PeerNavigatorRole)
	csoUser := createTestUser(t, ctx, userStore, auth.GOJOVEN, auth.SrRole)

	var testCases = []struct {
		name      string
		inputUser *auth.User
	}{
		{
			name:      "NAC Admin can delete NAC User",
			inputUser: nacUser,
		},
		{
			name:      "NAC Admin can delete BFLA User",
			inputUser: bflaUser,
		},
		{
			name:      "NAC Admin can delete CSO User",
			inputUser: csoUser,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := userApi.DeleteUser(ctx, DeleteUserRequest{
				DeletedBy: nacAdminUser.Email,
				ID:        tc.inputUser.ID,
			})
			if err != nil {
				t.Errorf("Unexpected error deleting user: %v", err)
			}
		})
	}
}

func TestUserApi_NAC_AdminUserCanNotDeleteMOHWUser(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	nacAdminuser := createTestUser(t, ctx, userStore, auth.NAC, auth.AdminRole)
	mohwUser := createTestUser(t, ctx, userStore, auth.MOHW, auth.AdminRole)
	err := userApi.DeleteUser(ctx, DeleteUserRequest{
		DeletedBy: nacAdminuser.Email,
		ID:        mohwUser.ID,
	})

	if err == nil {
		t.Errorf("NAC should not be able to delete MOHW user")
	}
}

func TestUserApi_BFLA_AdminUserCanDeleteBFLAUsers(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	bflaAdminUser := createTestUser(t, ctx, userStore, auth.BFLA, auth.AdminRole)
	bflaTestUser1 := createTestUser(t, ctx, userStore, auth.BFLA, auth.PeerNavigatorRole)
	bflaTestUser2 := createTestUser(t, ctx, userStore, auth.BFLA, auth.AdherenceCounselorRole)
	bflaTestUser3 := createTestUser(t, ctx, userStore, auth.BFLA, auth.AdminRole)

	var testCases = []struct {
		name      string
		inputUser *auth.User
	}{
		{
			name:      "BFLA Admin can delete BFLA Peer Navigator",
			inputUser: bflaTestUser1,
		},
		{
			name:      "BFLA Admin can delete BFLA Adherence Counselor",
			inputUser: bflaTestUser2,
		},
		{
			name:      "BFLA Admin can delete BFLA AdminRole",
			inputUser: bflaTestUser3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := userApi.DeleteUser(ctx, DeleteUserRequest{
				DeletedBy: bflaAdminUser.Email,
				ID:        tc.inputUser.ID,
			})
			if err != nil {
				t.Errorf("Unexpected error deleting user: %v", err)
			}
		})
	}
}

func TestUserApi_BFLA_AdminUserCanNotDeleteUsersFromOtherOrgs(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	bflaAdminUser := createTestUser(t, ctx, userStore, auth.BFLA, auth.AdminRole)
	mohwUser := createTestUser(t, ctx, userStore, auth.MOHW, auth.AdminRole)
	nacUser := createTestUser(t, ctx, userStore, auth.NAC, auth.AdminRole)
	csoUser := createTestUser(t, ctx, userStore, auth.GOJOVEN, auth.SrRole)

	var testCases = []struct {
		name      string
		inputUser *auth.User
	}{
		{
			name:      "BFLA Admin can not delete MOHW User",
			inputUser: mohwUser,
		},
		{
			name:      "BFLA Admin can not delete NAC User",
			inputUser: nacUser,
		},
		{
			name:      "BFLA Admin can not delete CSO User",
			inputUser: csoUser,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := userApi.DeleteUser(ctx, DeleteUserRequest{
				DeletedBy: bflaAdminUser.Email,
				ID:        tc.inputUser.ID,
			})
			if err == nil {
				t.Errorf("Expected error deleting user, got nil")
			}
		})
	}
}

func TestUserApi_NonAdminUsers_CanNotDeleteUsers(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	bflaPeerNavigatorUser := createTestUser(t, ctx, userStore, auth.BFLA, auth.PeerNavigatorRole)
	bflaAdherenceCounselorUser := createTestUser(t, ctx, userStore, auth.BFLA, auth.AdherenceCounselorRole)
	mohwUser := createTestUser(t, ctx, userStore, auth.MOHW, auth.SrRole)
	nacUser := createTestUser(t, ctx, userStore, auth.NAC, auth.SrRole)

	var testCases = []struct {
		name          string
		requestedUser *auth.User
		inputUser     *auth.User
	}{
		{
			name:          "BFLA Peer Navigator can not delete BFLA Adherence Counselor",
			requestedUser: bflaAdherenceCounselorUser,
			inputUser:     bflaAdherenceCounselorUser,
		},
		{
			name:          "BFLA Adherence Counselor can not delete BFLA Peer Navigator",
			requestedUser: bflaPeerNavigatorUser,
			inputUser:     bflaPeerNavigatorUser,
		},
		{
			name:          "Non Admin MOHW can not delete BFLA Peer Navigator",
			inputUser:     bflaPeerNavigatorUser,
			requestedUser: mohwUser,
		},
		{
			name:          "Non Admin MOHW can not delete BFLA Adherence Counselor",
			inputUser:     bflaAdherenceCounselorUser,
			requestedUser: mohwUser,
		},
		{
			name:          "Non Admin NAC can not delete BFLA Peer Navigator",
			inputUser:     bflaPeerNavigatorUser,
			requestedUser: nacUser,
		},
		{
			name:          "Non Admin NAC can not delete BFLA Adherence Counselor",
			inputUser:     bflaAdherenceCounselorUser,
			requestedUser: nacUser,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := userApi.DeleteUser(ctx, DeleteUserRequest{
				DeletedBy: tc.requestedUser.Email,
				ID:        tc.inputUser.ID,
			})
			if err == nil {
				t.Errorf("Expected error deleting user, got nil")
			}
		})
	}

}

func TestUserApi_WhenDeletingAnInexistentUser(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	adminUser := createTestUser(t, ctx, userStore, auth.MOHW, auth.AdminRole)

	err := userApi.DeleteUser(ctx, DeleteUserRequest{
		DeletedBy: adminUser.Email,
		ID:        "inexistent",
	})
	if err == nil {
		t.Errorf("Expected error deleting user, got nil")
	}
}
