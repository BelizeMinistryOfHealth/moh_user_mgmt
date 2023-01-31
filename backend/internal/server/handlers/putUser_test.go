package handlers

import (
	"bytes"
	"bz.moh.epi/users/internal/auth"
	"bz.moh.epi/users/internal/db"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPutUser_FailsIfNotAdmin(t *testing.T) {
	ctx := context.Background()
	firebaseConfig := &firebase.Config{ProjectID: projectID}
	firestoreClient, err := db.NewFirestoreClient(ctx, firebaseConfig)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	userStore, _ := auth.NewStore(firestoreClient, apiKey)
	mids := NewChain(verifyTokenNonAdmin())
	userCrudService := NewUserCrudService(&userStore)
	createRequest := auth.CreateUserRequest{
		FirstName:        "Dan",
		LastName:         "Th",
		Email:            "dan@th.com",
		UserApplications: nil,
	}
	usr := createUser(ctx, t, userStore, createRequest)

	putBody := PutUserRequest{
		FirstName:    "Dan",
		LastName:     "Don",
		Email:        "dan@th.com",
		Applications: nil,
	}
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/users/%s", usr.ID), putRequestToJSON(putBody))
	mids.Then(userCrudService.PutUser)(res, req)

	if res.Code != 401 {
		t.Errorf("want: %d; got: %d", 401, res.Code)
	}
}
func TestPutUser_Success(t *testing.T) {
	ctx := context.Background()
	firebaseConfig := &firebase.Config{ProjectID: projectID}
	firestoreClient, err := db.NewFirestoreClient(ctx, firebaseConfig)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	userStore, _ := auth.NewStore(firestoreClient, apiKey)
	mids := NewChain(verifyTokenAdmin())
	userCrudService := NewUserCrudService(&userStore)
	createRequest := auth.CreateUserRequest{
		FirstName:        "Dan",
		LastName:         "Th",
		Email:            "dan@th.com",
		UserApplications: nil,
	}
	usr := createUser(ctx, t, userStore, createRequest)

	putBody := PutUserRequest{
		FirstName:    "Dan",
		LastName:     "Th",
		Email:        "dan@th.com",
		Applications: nil,
	}
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/users/%s", usr.ID), putRequestToJSON(putBody))
	mids.Then(userCrudService.PutUser)(res, req)

	if res.Code != 200 {
		t.Errorf("want: %d; got: %d", 200, res.Code)
	}
	wantUser, _ := userStore.GetUserByID(ctx, usr.ID)
	if wantUser.LastName != putBody.LastName {
		t.Errorf("Updating user failed want: %s; got: %s", wantUser.LastName, putBody.LastName)
	}
}
func TestPutUser_UpdatesApplications(t *testing.T) {
	ctx := context.Background()
	firebaseConfig := &firebase.Config{ProjectID: projectID}
	firestoreClient, err := db.NewFirestoreClient(ctx, firebaseConfig)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	userStore, _ := auth.NewStore(firestoreClient, apiKey)
	mids := NewChain(verifyTokenAdmin())
	userCrudService := NewUserCrudService(&userStore)
	createRequest := auth.CreateUserRequest{
		FirstName:        "Dan",
		LastName:         "Th",
		Email:            "dan@th.com",
		UserApplications: nil,
	}
	usr := createUser(ctx, t, userStore, createRequest)

	putBody := PutUserRequest{
		FirstName: "Dan",
		LastName:  "Th",
		Email:     "dan@th.com",
		Applications: []auth.UserApplication{
			{ApplicationID: "app1",
				Name:        "app",
				Permissions: []string{"read"},
			}},
	}
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/users/%s", usr.ID), putRequestToJSON(putBody))
	mids.Then(userCrudService.PutUser)(res, req)

	if res.Code != 200 {
		t.Errorf("want: %d; got: %d", 200, res.Code)
	}
	wantUser, _ := userStore.GetUserByID(ctx, usr.ID)
	if diff := cmp.Diff(wantUser.UserApplications, putBody.Applications); diff != "" {
		t.Errorf("Updating user applications failed want: %s; got: %s", wantUser.UserApplications, putBody.Applications)
	}
}

func putRequestToJSON(req PutUserRequest) io.Reader {
	b, _ := json.Marshal(req)
	return bytes.NewBuffer(b)
}
