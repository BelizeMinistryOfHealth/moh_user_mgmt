package api

import (
	"bz.moh.epi/users/internal/auth"
	"context"
	"testing"
)

func TestUserApi_CanRetrieveOwnUser(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	user := createTestUser(t, ctx, userStore, auth.BFLA, auth.PeerNavigatorRole)
	got, err := userApi.GetUser(ctx, GetUserRequest{
		ID:          user.ID,
		RequestedBy: user.Email,
	})
	if err != nil {
		t.Fatalf("error getting user: %v", err)
	}
	if got.ID != user.ID {
		t.Errorf("GetUser() want: %s, got: %s", user.ID, got.ID)
	}
}

func TestUserApi_NonAdminUser_CanNotRetrieveOtherUser(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	user := createTestUser(t, ctx, userStore, auth.BFLA, auth.PeerNavigatorRole)
	otherUser := createTestUser(t, ctx, userStore, auth.BFLA, auth.PeerNavigatorRole)
	_, err := userApi.GetUser(ctx, GetUserRequest{
		ID:          otherUser.ID,
		RequestedBy: user.Email,
	})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
func TestUserApi_BFLAAdminUser_CanRetrieveBFLAUser(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	user := createTestUser(t, ctx, userStore, auth.BFLA, auth.AdminRole)
	otherUser := createTestUser(t, ctx, userStore, auth.BFLA, auth.PeerNavigatorRole)
	_, err := userApi.GetUser(ctx, GetUserRequest{
		ID:          otherUser.ID,
		RequestedBy: user.Email,
	})
	if err != nil {
		t.Errorf("expected no error retrieving User, got %v", err)
	}
}

func TestUserApi_BFLAAdminUser_CanNotRetrieveUsersFromOtherOrgs(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	adminUser := createTestUser(t, ctx, userStore, auth.BFLA, auth.AdminRole)
	csoUser := createTestUser(t, ctx, userStore, auth.CSO, auth.PeerNavigatorRole)
	nacUser := createTestUser(t, ctx, userStore, auth.NAC, auth.AdminRole)
	mohwUser := createTestUser(t, ctx, userStore, auth.MOHW, auth.AdminRole)

	var testCases = []struct {
		name string
		user *auth.User
	}{
		{
			name: "Can Not Retrieve CSO User",
			user: csoUser,
		},
		{
			name: "Can Not Retrieve NAC User",
			user: nacUser,
		},
		{
			name: "Can Not Retrieve MOHW User",
			user: mohwUser,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := userApi.GetUser(ctx, GetUserRequest{
				ID:          tc.user.ID,
				RequestedBy: adminUser.Email,
			})
			if err == nil {
				t.Errorf("expected error, got nil")
			}
		})
	}
}
func TestUserApi_CSOAdminUser_CanNotRetrieveUsersFromOtherOrgs(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	adminUser := createTestUser(t, ctx, userStore, auth.CSO, auth.AdminRole)
	bflaUser := createTestUser(t, ctx, userStore, auth.BFLA, auth.PeerNavigatorRole)
	nacUser := createTestUser(t, ctx, userStore, auth.NAC, auth.AdminRole)
	mohwUser := createTestUser(t, ctx, userStore, auth.MOHW, auth.AdminRole)

	var testCases = []struct {
		name string
		user *auth.User
	}{
		{
			name: "Can Not Retrieve BFLA User",
			user: bflaUser,
		},
		{
			name: "Can Not Retrieve NAC User",
			user: nacUser,
		},
		{
			name: "Can Not Retrieve MOHW User",
			user: mohwUser,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := userApi.GetUser(ctx, GetUserRequest{
				ID:          tc.user.ID,
				RequestedBy: adminUser.Email,
			})
			if err == nil {
				t.Errorf("expected error, got nil")
			}
		})
	}
}
func TestUserApi_CSOAdminUser_CanRetrieveCSOUser(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	user := createTestUser(t, ctx, userStore, auth.CSO, auth.AdminRole)
	otherUser := createTestUser(t, ctx, userStore, auth.CSO, auth.PeerNavigatorRole)
	_, err := userApi.GetUser(ctx, GetUserRequest{
		ID:          otherUser.ID,
		RequestedBy: user.Email,
	})
	if err != nil {
		t.Errorf("expected no error retrieving User, got %v", err)
	}
}
func TestUserApi_MOHWdminUser_CanRetrieveUsersFromAllOrgs(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	adminUser := createTestUser(t, ctx, userStore, auth.MOHW, auth.AdminRole)
	csoUser := createTestUser(t, ctx, userStore, auth.CSO, auth.PeerNavigatorRole)
	nacUser := createTestUser(t, ctx, userStore, auth.NAC, auth.AdminRole)
	mohwUser := createTestUser(t, ctx, userStore, auth.MOHW, auth.AdminRole)
	bflaUser := createTestUser(t, ctx, userStore, auth.BFLA, auth.AdminRole)

	var testCases = []struct {
		name string
		user *auth.User
	}{
		{
			name: "Can Retrieve CSO User",
			user: csoUser,
		},
		{
			name: "Can Retrieve NAC User",
			user: nacUser,
		},
		{
			name: "Can Retrieve MOHW User",
			user: mohwUser,
		},
		{
			name: "Can Retrieve BFLA User",
			user: bflaUser,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := userApi.GetUser(ctx, GetUserRequest{
				ID:          tc.user.ID,
				RequestedBy: adminUser.Email,
			})
			if err != nil {
				t.Errorf("MOHW Admin User error retrieving user: %v", err)
			}
		})
	}
}
func TestUserApi_NACAdminUser_CanRetrieveUsersFromAllOrgs(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	adminUser := createTestUser(t, ctx, userStore, auth.NAC, auth.AdminRole)
	csoUser := createTestUser(t, ctx, userStore, auth.CSO, auth.PeerNavigatorRole)
	nacUser := createTestUser(t, ctx, userStore, auth.NAC, auth.AdminRole)
	bflaUser := createTestUser(t, ctx, userStore, auth.BFLA, auth.AdminRole)

	var testCases = []struct {
		name string
		user *auth.User
	}{
		{
			name: "Can Retrieve CSO User",
			user: csoUser,
		},
		{
			name: "Can Retrieve NAC User",
			user: nacUser,
		},
		{
			name: "Can Retrieve BFLA User",
			user: bflaUser,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := userApi.GetUser(ctx, GetUserRequest{
				ID:          tc.user.ID,
				RequestedBy: adminUser.Email,
			})
			if err != nil {
				t.Errorf("NAC Admin User error retrieving user: %v", err)
			}
		})
	}
}

func TestUserApi_NACAdminUser_CanNotRetrieveMOHWUser(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	adminUser := createTestUser(t, ctx, userStore, auth.NAC, auth.AdminRole)
	mohwUser := createTestUser(t, ctx, userStore, auth.MOHW, auth.AdminRole)
	_, err := userApi.GetUser(ctx, GetUserRequest{
		ID:          mohwUser.ID,
		RequestedBy: adminUser.Email,
	})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
