package api

import (
	"bz.moh.epi/users/internal/auth"
	"context"
	"sort"
	"testing"
)

func TestUserApi_MOHW_AdminsCanListAllUsers(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	mohwAdminUser := createAdminUser(t, ctx, userStore, auth.MOHW)
	nacTestUser1 := createTestUser(t, ctx, userStore, auth.NAC, auth.AdherenceCounselorRole)
	nacTestUser2 := createTestUser(t, ctx, userStore, auth.NAC, auth.AdminRole)
	mohwTestUser1 := createTestUser(t, ctx, userStore, auth.MOHW, auth.AdminRole)
	mohwTestUser2 := createTestUser(t, ctx, userStore, auth.MOHW, auth.AdminRole)
	bflaTestUser1 := createTestUser(t, ctx, userStore, auth.BFLA, auth.AdminRole)
	bflaTestUser2 := createTestUser(t, ctx, userStore, auth.BFLA, auth.AdherenceCounselorRole)
	csoTestUser1 := createTestUser(t, ctx, userStore, auth.CSO, auth.SrRole)

	returnedUsers, err := userApi.ListUsers(ctx, mohwAdminUser.Email)
	if err != nil {
		t.Fatalf("Error listing users: %v", err)
	}

	// Sort the returned users by email
	sort.Slice(returnedUsers, func(i, j int) bool {
		return returnedUsers[i].Email < returnedUsers[j].Email
	})

	// Sort the expected users by email
	expectedUsers := []auth.User{
		*mohwAdminUser,
		*mohwTestUser1,
		*mohwTestUser2,
		*nacTestUser1,
		*nacTestUser2,
		*bflaTestUser1,
		*bflaTestUser2,
		*csoTestUser1,
	}
	// Sort the expected users by email
	sort.Slice(expectedUsers, func(i, j int) bool {
		return expectedUsers[i].Email < expectedUsers[j].Email
	})

	if len(returnedUsers) != len(expectedUsers) {
		t.Errorf("Number of users | want: %d; got: %d", len(expectedUsers), len(returnedUsers))
	}
}
func TestUserApi_NAC_AdminsCanListAllUsersExceptMOHW(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	nacAdminUser := createAdminUser(t, ctx, userStore, auth.NAC)
	nacTestUser1 := createTestUser(t, ctx, userStore, auth.NAC, auth.AdherenceCounselorRole)
	nacTestUser2 := createTestUser(t, ctx, userStore, auth.NAC, auth.AdminRole)
	_ = createTestUser(t, ctx, userStore, auth.MOHW, auth.AdminRole)
	_ = createTestUser(t, ctx, userStore, auth.MOHW, auth.AdminRole)
	bflaTestUser1 := createTestUser(t, ctx, userStore, auth.BFLA, auth.AdminRole)
	bflaTestUser2 := createTestUser(t, ctx, userStore, auth.BFLA, auth.AdherenceCounselorRole)
	csoTestUser1 := createTestUser(t, ctx, userStore, auth.CSO, auth.SrRole)

	returnedUsers, err := userApi.ListUsers(ctx, nacAdminUser.Email)
	if err != nil {
		t.Fatalf("Error listing users: %v", err)
	}

	// Sort the returned users by email
	sort.Slice(returnedUsers, func(i, j int) bool {
		return returnedUsers[i].Email < returnedUsers[j].Email
	})

	// Sort the expected users by email
	expectedUsers := []auth.User{
		*nacAdminUser,
		*nacTestUser1,
		*nacTestUser2,
		*bflaTestUser1,
		*bflaTestUser2,
		*csoTestUser1,
	}
	// Sort the expected users by email
	sort.Slice(expectedUsers, func(i, j int) bool {
		return expectedUsers[i].Email < expectedUsers[j].Email
	})

	if len(returnedUsers) != len(expectedUsers) {
		t.Errorf("Number of users | want: %d; got: %d", len(expectedUsers), len(returnedUsers))
	}
}
func TestUserApi_BFLA_AdminsCanOnlyListBFLAUsers(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	bflaAdminUser := createAdminUser(t, ctx, userStore, auth.BFLA)
	_ = createTestUser(t, ctx, userStore, auth.NAC, auth.AdherenceCounselorRole)
	_ = createTestUser(t, ctx, userStore, auth.NAC, auth.AdminRole)
	_ = createTestUser(t, ctx, userStore, auth.MOHW, auth.AdminRole)
	_ = createTestUser(t, ctx, userStore, auth.MOHW, auth.AdminRole)
	bflaTestUser1 := createTestUser(t, ctx, userStore, auth.BFLA, auth.AdminRole)
	bflaTestUser2 := createTestUser(t, ctx, userStore, auth.BFLA, auth.AdherenceCounselorRole)
	_ = createTestUser(t, ctx, userStore, auth.CSO, auth.SrRole)

	returnedUsers, err := userApi.ListUsers(ctx, bflaAdminUser.Email)
	if err != nil {
		t.Fatalf("Error listing users: %v", err)
	}

	// Sort the returned users by email
	sort.Slice(returnedUsers, func(i, j int) bool {
		return returnedUsers[i].Email < returnedUsers[j].Email
	})

	// Sort the expected users by email
	expectedUsers := []auth.User{
		*bflaAdminUser,
		*bflaTestUser1,
		*bflaTestUser2,
	}
	// Sort the expected users by email
	sort.Slice(expectedUsers, func(i, j int) bool {
		return expectedUsers[i].Email < expectedUsers[j].Email
	})

	if len(returnedUsers) != len(expectedUsers) {
		t.Errorf("Number of users | want: %d; got: %d", len(expectedUsers), len(returnedUsers))
	}
}
func TestUserApi_CSO_AdminsCanOnlyListCSOUsers(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	csoAdminUser := createAdminUser(t, ctx, userStore, auth.CSO)
	_ = createTestUser(t, ctx, userStore, auth.NAC, auth.AdherenceCounselorRole)
	_ = createTestUser(t, ctx, userStore, auth.NAC, auth.AdminRole)
	_ = createTestUser(t, ctx, userStore, auth.MOHW, auth.AdminRole)
	_ = createTestUser(t, ctx, userStore, auth.MOHW, auth.AdminRole)
	_ = createTestUser(t, ctx, userStore, auth.BFLA, auth.AdminRole)
	_ = createTestUser(t, ctx, userStore, auth.BFLA, auth.AdherenceCounselorRole)
	csoTestUser := createTestUser(t, ctx, userStore, auth.CSO, auth.SrRole)

	returnedUsers, err := userApi.ListUsers(ctx, csoAdminUser.Email)
	if err != nil {
		t.Fatalf("Error listing users: %v", err)
	}

	// Sort the returned users by email
	sort.Slice(returnedUsers, func(i, j int) bool {
		return returnedUsers[i].Email < returnedUsers[j].Email
	})

	// Sort the expected users by email
	expectedUsers := []auth.User{
		*csoAdminUser,
		*csoTestUser,
	}
	// Sort the expected users by email
	sort.Slice(expectedUsers, func(i, j int) bool {
		return expectedUsers[i].Email < expectedUsers[j].Email
	})

	if len(returnedUsers) != len(expectedUsers) {
		t.Errorf("Number of users | want: %d; got: %d", len(expectedUsers), len(returnedUsers))
	}
}
