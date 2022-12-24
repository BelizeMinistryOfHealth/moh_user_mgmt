package handlers

import (
	"bz.moh.epi/users/internal/auth"
	"bz.moh.epi/users/internal/db"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go/v4"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApplicationsApi_CanListApplications(t *testing.T) {
	ctx := context.Background()
	firebaseConfig := &firebase.Config{ProjectID: projectID}
	firestoreClient, err := db.NewFirestoreClient(ctx, firebaseConfig)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	userStore, _ := auth.NewStore(firestoreClient, apiKey)
	mids := NewChain(verifyTokenAdmin())
	userCrudService := NewUserCrudService(&userStore)
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/applications", nil)
	mids.Then(userCrudService.ListApplications)(res, req)
	if res.Code != 200 {
		t.Errorf("want: %d; got: %d", 200, res.Code)
	}
	var apps []auth.UserApplication
	if err := json.NewDecoder(res.Body).Decode(&apps); err != nil {
		t.Errorf("could not decode the applications")
	}
	if len(apps) == 0 {
		t.Errorf("wanted non-empty user applications list to be returned")
	}
}
