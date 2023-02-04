package handlers

import (
	"bytes"
	"bz.moh.epi/users/internal/auth"
	"bz.moh.epi/users/internal/db"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go/v4"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var apiKey = os.Getenv("API_KEY")       //nolint: gochecknoglobals
var projectID = os.Getenv("PROJECT_ID") //nolint: gochecknoglobals

func verifyTokenNonAdmin() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			t := auth.JwtToken{
				Email: "me@mail.com",
				Admin: false,
				Org:   "BFLA",
				Role:  auth.PeerNavigatorRole,
			}
			ctx := context.WithValue(r.Context(), "authToken", t) //nolint: staticcheck
			f(w, r.WithContext(ctx))
		}
	}
}
func TestPostUser_FailsIfNotAdmin(t *testing.T) {
	ctx := context.Background()
	firebaseConfig := &firebase.Config{ProjectID: projectID}
	firestoreClient, err := db.NewFirestoreClient(ctx, firebaseConfig)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	userStore, _ := auth.NewStore(firestoreClient, apiKey)
	mids := NewChain(verifyTokenNonAdmin())
	userCrudService := NewUserCrudService(&userStore)
	postBody := PostUserRequest{
		FirstName: "Dan",
		LastName:  "Th",
		Email:     "dan@th.com",
		Role:      "PeerNavigatorRole",
		Org:       "BFLA",
	}
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/user", userToJSON(postBody))
	mids.Then(userCrudService.PostUser)(res, req)
	if res.Code != 401 {
		t.Errorf("want: %d; got: %d", 401, res.Code)
	}
}

func verifyTokenAdmin() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			t := auth.JwtToken{
				Email: "me@mail.com",
				Admin: true,
				Org:   "BFLA",
				Role:  auth.SrRole,
			}
			ctx := context.WithValue(r.Context(), "authToken", t) //nolint: staticcheck
			f(w, r.WithContext(ctx))
		}
	}
}
func TestPostUser_Success(t *testing.T) {
	ctx := context.Background()
	firebaseConfig := &firebase.Config{ProjectID: projectID}
	firestoreClient, err := db.NewFirestoreClient(ctx, firebaseConfig)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	userStore, _ := auth.NewStore(firestoreClient, apiKey)
	mids := NewChain(verifyTokenAdmin())
	userCrudService := NewUserCrudService(&userStore)
	postBody := PostUserRequest{
		FirstName: "Dan",
		LastName:  "Th",
		Email:     "dan@th.com",
		Role:      "PeerNavigatorRole",
		Org:       "BFLA",
	}
	t.Cleanup(func() {
		userCrudService.UserStore.DeleteUserByEmail(ctx, postBody.Email) //nolint:errcheck,gosec
	})
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/user", userToJSON(postBody))
	mids.Then(userCrudService.PostUser)(res, req)
	if res.Code != 200 {
		t.Errorf("want: %d; got: %d", 200, res.Code)
	}
}
func userToJSON(u PostUserRequest) io.Reader {
	r, _ := json.Marshal(u)
	return bytes.NewReader(r)
}
