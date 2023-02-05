package api

import (
	"bz.moh.epi/users/internal/auth"
	"context"
	"testing"
)

func TestUserApi_NAC_AdminsCanUpdateUsers(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	nacAdminUser := createAdminUser(t, ctx, userStore, auth.NAC)
	nacTestUser := createTestUser(t, ctx, userStore, auth.NAC, auth.AdherenceCounselorRole)
	bflaTestUser := createTestUser(t, ctx, userStore, auth.BFLA, auth.AdherenceCounselorRole)
	csoTestUser := createTestUser(t, ctx, userStore, auth.CSO, auth.AdherenceCounselorRole)

	var testCases = []struct {
		name  string
		input UpdateUserRequest
	}{
		{
			name: "NAC Admin can update users for NAC",
			input: UpdateUserRequest{
				User: &auth.User{
					ID:        nacTestUser.ID,
					FirstName: nacTestUser.FirstName,
					LastName:  nacTestUser.LastName,
					Email:     nacTestUser.Email,
					Org:       nacTestUser.Org,
					Role:      auth.PeerNavigatorRole,
				},
				UpdatedBy: nacAdminUser.Email,
			},
		},
		{
			name: "NAC Admin can update users for BFLA",
			input: UpdateUserRequest{
				User: &auth.User{
					ID:        bflaTestUser.ID,
					FirstName: bflaTestUser.FirstName,
					LastName:  bflaTestUser.LastName,
					Email:     bflaTestUser.Email,
					Org:       bflaTestUser.Org,
					Role:      auth.SrRole,
				},
				UpdatedBy: nacAdminUser.Email,
			},
		},
		{
			name: "NAC Admin can update users for CSO",
			input: UpdateUserRequest{
				User: &auth.User{
					ID:        csoTestUser.ID,
					FirstName: csoTestUser.FirstName,
					LastName:  csoTestUser.LastName,
					Email:     csoTestUser.Email,
					Org:       csoTestUser.Org,
					Role:      auth.SrRole,
				},
				UpdatedBy: nacAdminUser.Email,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if err := userApi.UpdateUser(ctx, tt.input); err != nil {
				t.Errorf("error updating user: %v", err)
			}
			got, err := userStore.GetUserByID(ctx, tt.input.User.ID)
			if err != nil {
				t.Errorf("error getting updated user: %v", err)
			}
			if got.Role != tt.input.User.Role {
				t.Errorf("want: %v, got: %v", tt.input.User.Role, got.Role)
			}
		})
	}
}

func TestUserApi_NAC_CanNotUpdateMOHWUsers(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	nacAdminUser := createAdminUser(t, ctx, userStore, auth.NAC)
	mohwTestUser := createTestUser(t, ctx, userStore, auth.MOHW, auth.AdherenceCounselorRole)

	var testCases = []struct {
		name  string
		input UpdateUserRequest
	}{
		{
			name: "NAC Admin can not update users for MOHW",
			input: UpdateUserRequest{
				User: &auth.User{
					ID:        mohwTestUser.ID,
					FirstName: mohwTestUser.FirstName,
					LastName:  mohwTestUser.LastName,
					Email:     mohwTestUser.Email,
					Org:       mohwTestUser.Org,
					Role:      auth.PeerNavigatorRole,
				},
				UpdatedBy: nacAdminUser.Email,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if err := userApi.UpdateUser(ctx, tt.input); err == nil {
				t.Errorf("expected error updating user")
			}
		})
	}
}

func TestUserApi_MOHW_CanUpdateAllUsers(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	mohwAdminUser := createAdminUser(t, ctx, userStore, auth.MOHW)
	nacTestUser := createTestUser(t, ctx, userStore, auth.NAC, auth.AdherenceCounselorRole)
	bflaTestUser := createTestUser(t, ctx, userStore, auth.BFLA, auth.AdherenceCounselorRole)
	csoTestUser := createTestUser(t, ctx, userStore, auth.CSO, auth.AdherenceCounselorRole)
	mohwTestUser := createTestUser(t, ctx, userStore, auth.MOHW, auth.AdherenceCounselorRole)

	var testCases = []struct {
		name  string
		input UpdateUserRequest
	}{
		{
			name: "MOHW Admin can update users for NAC",
			input: UpdateUserRequest{
				User: &auth.User{
					ID:        nacTestUser.ID,
					FirstName: nacTestUser.FirstName,
					LastName:  nacTestUser.LastName,
					Email:     nacTestUser.Email,
					Org:       nacTestUser.Org,
					Role:      auth.PeerNavigatorRole,
				},
				UpdatedBy: mohwAdminUser.Email,
			},
		},
		{
			name: "MOHW Admin can update users for BFLA",
			input: UpdateUserRequest{
				User: &auth.User{
					ID:        bflaTestUser.ID,
					FirstName: bflaTestUser.FirstName,
					LastName:  bflaTestUser.LastName,
					Email:     bflaTestUser.Email,
					Org:       bflaTestUser.Org,
					Role:      auth.SrRole,
				},
				UpdatedBy: mohwAdminUser.Email,
			},
		},
		{
			name: "MOHW Admin can update users for CSO",
			input: UpdateUserRequest{
				User: &auth.User{
					ID:        csoTestUser.ID,
					FirstName: csoTestUser.FirstName,
					LastName:  csoTestUser.LastName,
					Email:     csoTestUser.Email,
					Org:       csoTestUser.Org,
					Role:      auth.SrRole,
				},
				UpdatedBy: mohwAdminUser.Email,
			},
		},
		{
			name: "MOHW Admin can update users for MOHW",
			input: UpdateUserRequest{
				User: &auth.User{
					ID:        mohwTestUser.ID,
					FirstName: mohwTestUser.FirstName,
					LastName:  mohwTestUser.LastName,
					Email:     mohwTestUser.Email,
					Org:       mohwTestUser.Org,
					Role:      auth.PeerNavigatorRole,
				},
				UpdatedBy: mohwAdminUser.Email,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if err := userApi.UpdateUser(ctx, tt.input); err != nil {
				t.Errorf("error updating user: %v", err)
			}
			got, err := userStore.GetUserByID(ctx, tt.input.User.ID)
			if err != nil {
				t.Errorf("error getting updated user: %v", err)
			}
			if got.Role != tt.input.User.Role {
				t.Errorf("want: %v, got: %v", tt.input.User.Role, got.Role)
			}
		})
	}
}

func TestUserApi_BFLA_AdminsCanUpdateBFLAUsers(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	bflaAdminUser := createAdminUser(t, ctx, userStore, auth.BFLA)
	bflaTestUser := createTestUser(t, ctx, userStore, auth.BFLA, auth.AdherenceCounselorRole)

	var testCases = []struct {
		name  string
		input UpdateUserRequest
	}{
		{
			name: "BFLA Admin can update users for BFLA",
			input: UpdateUserRequest{
				User: &auth.User{
					ID:        bflaTestUser.ID,
					FirstName: bflaTestUser.FirstName,
					LastName:  bflaTestUser.LastName,
					Email:     bflaTestUser.Email,
					Org:       bflaTestUser.Org,
					Role:      auth.PeerNavigatorRole,
				},
				UpdatedBy: bflaAdminUser.Email,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if err := userApi.UpdateUser(ctx, tt.input); err != nil {
				t.Errorf("error updating user: %v", err)
			}
			got, err := userStore.GetUserByID(ctx, tt.input.User.ID)
			if err != nil {
				t.Errorf("error getting updated user: %v", err)
			}
			if got.Role != tt.input.User.Role {
				t.Errorf("want: %v, got: %v", tt.input.User.Role, got.Role)
			}
		})
	}
}

func TestUserApi_CSO_Admin_CanUpdateCSOUsers(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	csoAdminUser := createAdminUser(t, ctx, userStore, auth.CSO)
	csoTestUser := createTestUser(t, ctx, userStore, auth.CSO, auth.AdherenceCounselorRole)

	var testCases = []struct {
		name  string
		input UpdateUserRequest
	}{
		{
			name: "CSO Admin can update users for CSO",
			input: UpdateUserRequest{
				User: &auth.User{
					ID:        csoTestUser.ID,
					FirstName: csoTestUser.FirstName,
					LastName:  csoTestUser.LastName,
					Email:     csoTestUser.Email,
					Org:       csoTestUser.Org,
					Role:      auth.PeerNavigatorRole,
				},
				UpdatedBy: csoAdminUser.Email,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if err := userApi.UpdateUser(ctx, tt.input); err != nil {
				t.Errorf("error updating user: %v", err)
			}
			got, err := userStore.GetUserByID(ctx, tt.input.User.ID)
			if err != nil {
				t.Errorf("error getting updated user: %v", err)
			}
			if got.Role != tt.input.User.Role {
				t.Errorf("want: %v, got: %v", tt.input.User.Role, got.Role)
			}
		})
	}
}

func TestUserApi_BFLA_AdminsCanNotUpdateUsersInOtherOrgs(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	bflaAdminUser := createAdminUser(t, ctx, userStore, auth.BFLA)
	mohwTestUser := createTestUser(t, ctx, userStore, auth.MOHW, auth.AdherenceCounselorRole)
	nacTestUser := createTestUser(t, ctx, userStore, auth.NAC, auth.AdherenceCounselorRole)
	csoTestUser := createTestUser(t, ctx, userStore, auth.CSO, auth.AdherenceCounselorRole)

	var testCases = []struct {
		name  string
		input UpdateUserRequest
	}{
		{
			name: "BFLA Admin can not update users for MOHW",
			input: UpdateUserRequest{
				User: &auth.User{
					ID:        mohwTestUser.ID,
					FirstName: mohwTestUser.FirstName,
					LastName:  mohwTestUser.LastName,
					Email:     mohwTestUser.Email,
					Org:       mohwTestUser.Org,
					Role:      auth.PeerNavigatorRole,
				},
				UpdatedBy: bflaAdminUser.Email,
			},
		},
		{
			name: "BFLA Admin can not update users for NAC",
			input: UpdateUserRequest{
				User: &auth.User{
					ID:        nacTestUser.ID,
					FirstName: nacTestUser.FirstName,
					LastName:  nacTestUser.LastName,
					Email:     nacTestUser.Email,
					Org:       nacTestUser.Org,
					Role:      auth.PeerNavigatorRole,
				},
				UpdatedBy: bflaAdminUser.Email,
			},
		},
		{
			name: "BFLA Admin can not update users for CSO",
			input: UpdateUserRequest{
				User: &auth.User{
					ID:        csoTestUser.ID,
					FirstName: csoTestUser.FirstName,
					LastName:  csoTestUser.LastName,
					Email:     csoTestUser.Email,
					Org:       csoTestUser.Org,
					Role:      auth.PeerNavigatorRole,
				},
				UpdatedBy: bflaAdminUser.Email,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if err := userApi.UpdateUser(ctx, tt.input); err == nil {
				t.Errorf("expected error updating user")
			}
		})
	}
}

func TestUserApi_CSO_AdminsCanNotUpdateUsersInOtherOrgs(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	csoAdminUser := createAdminUser(t, ctx, userStore, auth.CSO)
	mohwTestUser := createTestUser(t, ctx, userStore, auth.MOHW, auth.AdherenceCounselorRole)
	nacTestUser := createTestUser(t, ctx, userStore, auth.NAC, auth.AdherenceCounselorRole)
	bflaTestUser := createTestUser(t, ctx, userStore, auth.BFLA, auth.AdherenceCounselorRole)

	var testCases = []struct {
		name  string
		input UpdateUserRequest
	}{
		{
			name: "CSO Admin can not update users for MOHW",
			input: UpdateUserRequest{
				User: &auth.User{
					ID:        mohwTestUser.ID,
					FirstName: mohwTestUser.FirstName,
					LastName:  mohwTestUser.LastName,
					Email:     mohwTestUser.Email,
					Org:       mohwTestUser.Org,
					Role:      auth.PeerNavigatorRole,
				},
				UpdatedBy: csoAdminUser.Email,
			},
		},
		{
			name: "CSO Admin can not update users for NAC",
			input: UpdateUserRequest{
				User: &auth.User{
					ID:        nacTestUser.ID,
					FirstName: nacTestUser.FirstName,
					LastName:  nacTestUser.LastName,
					Email:     nacTestUser.Email,
					Org:       nacTestUser.Org,
					Role:      auth.PeerNavigatorRole,
				},
				UpdatedBy: csoAdminUser.Email,
			},
		},
		{
			name: "CSO Admin can not update users for BFLA",
			input: UpdateUserRequest{
				User: &auth.User{
					ID:        bflaTestUser.ID,
					FirstName: bflaTestUser.FirstName,
					LastName:  bflaTestUser.LastName,
					Email:     bflaTestUser.Email,
					Org:       bflaTestUser.Org,
					Role:      auth.PeerNavigatorRole,
				},
				UpdatedBy: csoAdminUser.Email,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if err := userApi.UpdateUser(ctx, tt.input); err == nil {
				t.Errorf("expected error updating user")
			}
		})
	}
}

func TestUserApi_AdminUsersCanNotUpdateTheirRole(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	mohwAdminUser := createAdminUser(t, ctx, userStore, auth.MOHW)
	nacAdminUser := createAdminUser(t, ctx, userStore, auth.NAC)
	bflaAdminUser := createAdminUser(t, ctx, userStore, auth.BFLA)
	csoAdminUser := createAdminUser(t, ctx, userStore, auth.CSO)

	var testCase = []struct {
		name  string
		input UpdateUserRequest
	}{
		{
			name: "MOHW Admin can not update their role",
			input: UpdateUserRequest{
				User: &auth.User{
					ID:        mohwAdminUser.ID,
					FirstName: mohwAdminUser.FirstName,
					LastName:  mohwAdminUser.LastName,
					Email:     mohwAdminUser.Email,
					Org:       mohwAdminUser.Org,
					Role:      auth.PeerNavigatorRole,
				},
				UpdatedBy: mohwAdminUser.Email,
			},
		},
		{
			name: "NAC Admin can not update their role",
			input: UpdateUserRequest{
				User: &auth.User{
					ID:        nacAdminUser.ID,
					FirstName: nacAdminUser.FirstName,
					LastName:  nacAdminUser.LastName,
					Email:     nacAdminUser.Email,
					Org:       nacAdminUser.Org,
					Role:      auth.PeerNavigatorRole,
				},
				UpdatedBy: nacAdminUser.Email,
			},
		},
		{
			name: "BFLA Admin can not update their role",
			input: UpdateUserRequest{
				User: &auth.User{
					ID:        bflaAdminUser.ID,
					FirstName: bflaAdminUser.FirstName,
					LastName:  bflaAdminUser.LastName,
					Email:     bflaAdminUser.Email,
					Org:       bflaAdminUser.Org,
					Role:      auth.PeerNavigatorRole,
				},
				UpdatedBy: bflaAdminUser.Email,
			},
		},
		{
			name: "CSO Admin can not update their role",
			input: UpdateUserRequest{
				User: &auth.User{
					ID:        csoAdminUser.ID,
					FirstName: csoAdminUser.FirstName,
					LastName:  csoAdminUser.LastName,
					Email:     csoAdminUser.Email,
					Org:       csoAdminUser.Org,
					Role:      auth.PeerNavigatorRole,
				},
				UpdatedBy: csoAdminUser.Email,
			},
		},
	}

	for _, tt := range testCase {
		t.Run(tt.name, func(t *testing.T) {
			if err := userApi.UpdateUser(ctx, tt.input); err == nil {
				t.Errorf("expected error updating user")
			}
		})
	}
}
