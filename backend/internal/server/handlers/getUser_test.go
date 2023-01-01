package handlers

import (
	"bz.moh.epi/users/internal/auth"
	"bz.moh.epi/users/internal/db"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createUser(ctx context.Context, t *testing.T, userStore auth.UserStore, user auth.CreateUserRequest) *auth.User {
	u, err := userStore.CreateUser(ctx, user)
	if err != nil {
		t.Errorf("error creating user: %v", err)
	}

	if u == nil {
		t.Errorf("want: a non nil user, got nil")
	}

	t.Cleanup(func() {
		err = userStore.DeleteUserByID(ctx, u.ID)
		if err != nil {
			t.Errorf("cleaning up user failed: %v", err)
		}
	})
	return u
}
func TestGetUser(t *testing.T) {
	ctx := context.Background()
	firebaseConfig := &firebase.Config{ProjectID: projectID}
	firestoreClient, err := db.NewFirestoreClient(ctx, firebaseConfig)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	userStore, _ := auth.NewStore(firestoreClient, apiKey)
	userRequest := auth.CreateUserRequest{
		FirstName:        "Roberto",
		LastName:         "Guerra",
		Email:            "uris77@gmail.com",
		UserApplications: []auth.UserApplication{},
	}
	wantUser := createUser(ctx, t, userStore, userRequest)
	mids := NewChain(verifyTokenAdmin())
	userCrudService := NewUserCrudService(&userStore)
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s", wantUser.ID), nil)
	mids.Then(userCrudService.GetUserByID)(res, req)
	if res.Code != 200 {
		t.Errorf("GET users/:id want: 200; got: %d", res.Code)
	}

	var got auth.User
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Errorf("could not decode user: %v", err)
	}

	if got.ID != wantUser.ID {
		t.Errorf("GET users/:id want: %s; got: %s", wantUser.ID, got.ID)
	}

}
