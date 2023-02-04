package handlers

import (
	"bz.moh.epi/users/internal/api"
	"bz.moh.epi/users/internal/auth"
	"bz.moh.epi/users/internal/db"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
)

func TestUserCrudService_ListUsers_NonAdminNotAllowed(t *testing.T) {
	ctx := context.Background()
	firebaseConfig := &firebase.Config{ProjectID: projectID}
	firestoreClient, err := db.NewFirestoreClient(ctx, firebaseConfig)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	userStore, _ := auth.NewStore(firestoreClient, apiKey)
	userApi := api.CreateUserApi(userStore)
	mids := NewChain(verifyTokenNonAdmin())
	userCrudService := NewUserCrudService(&userStore, userApi)
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	mids.Then(userCrudService.ListUsers)(res, req)
	if res.Code != 401 {
		t.Errorf("want: %d; got: %d", 401, res.Code)
	}
}

func TestUserCrudService_ListUsers_AdminUserCanListUsers(t *testing.T) {
	ctx := context.Background()
	firebaseConfig := &firebase.Config{ProjectID: projectID}
	firestoreClient, err := db.NewFirestoreClient(ctx, firebaseConfig)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	userStore, _ := auth.NewStore(firestoreClient, apiKey)
	userApi := api.CreateUserApi(userStore)
	mids := NewChain(verifyTokenAdmin())
	userCrudService := NewUserCrudService(&userStore, userApi)
	want := createMultipleUsers(ctx, userStore)
	t.Cleanup(func() {
		for i := range want {
			userStore.DeleteUserByID(ctx, want[i].ID) //nolint:errcheck,gosec
		}
	})
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	mids.Then(userCrudService.ListUsers)(res, req)
	if res.Code != 200 {
		t.Errorf("Status Code | want: %d; got: %d", 200, res.Code)
	}
	var got []auth.User
	if err = json.Unmarshal(res.Body.Bytes(), &got); err != nil {
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
	//if len(got) != len(want) {
	//	t.Fatalf("Unexpected number of users want: %d; got: %d", len(want), len(got))
	//}
	//if diff := cmp.Diff(want, got); diff != "" {
	//	t.Errorf("want: %v;\n got: %v", want, got)
	//}
}

func createMultipleUsers(ctx context.Context, s auth.UserStore) []auth.User {
	var users []auth.User
	for i := 0; i < 5; i++ {
		req := auth.CreateUserRequest{
			FirstName: fmt.Sprintf("first'Name%d", i),
			LastName:  "LastName",
			Email:     fmt.Sprintf("%d@mail.com", i),
		}
		user, err := s.CreateUser(ctx, req)
		if err == nil {
			users = append(users, *user)
		}
	}

	return users
}
