package handlers

import (
	"bz.moh.epi/users/internal/api"
	"bz.moh.epi/users/internal/auth"
	"context"
	"encoding/json"
	"github.com/brianvoe/gofakeit/v6"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
)

func TestUserCrudService_ListUsers_NonAdminNotAllowed(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := api.CreateUserApi(*userStore)
	nonAdminUser := createTestUser(t, *userStore, auth.GOJOVEN, auth.SrRole)
	mids := NewChain(verifyUserToken(*nonAdminUser))
	userCrudService := NewUserCrudService(userStore, userApi)
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	mids.Then(userCrudService.ListUsers)(res, req)
	if res.Code != 401 {
		t.Errorf("want: %d; got: %d", 401, res.Code)
	}
}

func TestUserCrudService_ListUsers_AdminUserCanListUsers(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := api.CreateUserApi(*userStore)
	adminUser := createTestUser(t, *userStore, auth.GOJOVEN, auth.AdminRole)
	mids := NewChain(verifyUserToken(*adminUser))
	userCrudService := NewUserCrudService(userStore, userApi)
	// This user should not be included the results
	createTestUser(t, *userStore, auth.BFLA, auth.SrRole)
	createdUsers := createMultipleUsers(ctx, *userStore)
	t.Cleanup(func() {
		for i := range createdUsers {
			userStore.DeleteUserByID(ctx, createdUsers[i].ID) //nolint:errcheck,gosec
		}
	})
	// Add the admin user to the list of users that we expect to be returned.
	var want []auth.User
	want = append(want, createdUsers...)
	want = append(want, *adminUser)
	//t.Logf("want: %v ; size: %d createdUsers: %d", want, len(want), len(createdUsers))
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	mids.Then(userCrudService.ListUsers)(res, req)
	if res.Code != 200 {
		t.Errorf("Status Code | want: %d; got: %d", 200, res.Code)
	}
	var got []UserResponse
	if err := json.Unmarshal(res.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unarmshal users (%v): %v", res.Body, err)
	}
	// Sort the slices so that they can be compared. Otherwise, even if the values are the same,
	// the slices won't be equal.
	sort.Slice(got, func(i, j int) bool {
		return got[i].Email < got[j].Email
	})
	sort.Slice(want, func(i, j int) bool {
		return want[i].Email < got[j].Email
	})
	var nonGOJOVENUsers []UserResponse
	for i := range got {
		if got[i].Org != "GoJoven" {
			nonGOJOVENUsers = append(nonGOJOVENUsers, got[i])
		}
	}
	if len(nonGOJOVENUsers) > 0 {
		t.Fatalf("Non GoJoven users returned: %v", nonGOJOVENUsers)
	}
}

func createMultipleUsers(ctx context.Context, s auth.UserStore) []auth.User {
	var users []auth.User
	for i := 0; i < 5; i++ {
		req := auth.CreateUserRequest{
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
			Email:     gofakeit.Email(),
			Org:       auth.GOJOVEN,
			Role:      auth.SrRole,
		}
		user, err := s.CreateUser(ctx, req)
		if err == nil {
			users = append(users, *user)
		}
	}

	return users
}
