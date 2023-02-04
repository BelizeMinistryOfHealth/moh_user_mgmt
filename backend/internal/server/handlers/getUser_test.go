package handlers

import (
	"bz.moh.epi/users/internal/auth"
	"bz.moh.epi/users/internal/db"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createUser(ctx context.Context, t *testing.T, userStore auth.UserStore, user auth.CreateUserRequest) (*auth.User, error) {
	u, err := userStore.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	if u == nil {
		return nil, fmt.Errorf("want: a non nil user, got nil") //nolint: goerr113
	}

	t.Cleanup(func() {
		err = userStore.DeleteUserByID(ctx, u.ID)
		if err != nil {
			t.Errorf("cleaning up user failed: %v", err)
		}
	})
	return u, nil
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
		FirstName: "Roberto",
		LastName:  "Guerra",
		Email:     gofakeit.Email(),
		Org:       "BFLA",
		Role:      auth.PeerNavigatorRole,
		CreatedBy: "some@mail.com",
	}
	wantUser, err := createUser(ctx, t, userStore, userRequest)
	if err != nil {
		t.Errorf("error creating user: %v", err)
	}
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
