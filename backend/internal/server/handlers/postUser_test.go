package handlers

import (
	"bytes"
	"bz.moh.epi/users/internal/api"
	"bz.moh.epi/users/internal/auth"
	"bz.moh.epi/users/internal/db"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go/v4"
	"github.com/brianvoe/gofakeit/v6"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var apiKey = os.Getenv("API_KEY")       //nolint: gochecknoglobals
var projectID = os.Getenv("PROJECT_ID") //nolint: gochecknoglobals

func TestPostUser_FailsIfNotAdmin(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := api.CreateUserApi(*userStore)
	nonAdminUser := createTestUser(t, *userStore, auth.BFLA, auth.PeerNavigatorRole)
	mids := NewChain(verifyUserToken(*nonAdminUser))
	userCrudService := NewUserCrudService(userStore, userApi)
	postBody := PostUserRequest{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Role:      auth.PeerNavigatorRole.String(),
		Org:       auth.BFLA.String(),
	}
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/user", userToJSON(postBody))
	mids.Then(userCrudService.PostUser)(res, req)
	if res.Code != 401 {
		t.Errorf("want: %d; got: %d", 401, res.Code)
	}
}

func TestPostUser_Success(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := api.CreateUserApi(*userStore)
	adminUser := createTestUser(t, *userStore, auth.BFLA, auth.AdminRole)
	mids := NewChain(verifyUserToken(*adminUser))
	userCrudService := NewUserCrudService(userStore, userApi)
	postBody := PostUserRequest{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Role:      auth.PeerNavigatorRole.String(),
		Org:       auth.BFLA.String(),
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

func createUserStore(t *testing.T, ctx context.Context) *auth.UserStore {
	firebaseConfig := &firebase.Config{ProjectID: projectID}
	firestoreClient, err := db.NewFirestoreClient(ctx, firebaseConfig)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	store, err := auth.NewStore(firestoreClient, apiKey)
	if err != nil {
		t.Fatalf("failed to create user store: %v", err)
	}
	return &store

}
func createTestUser(t *testing.T, userStore auth.UserStore, org auth.Org, role auth.UserRole) *auth.User {
	ctx := context.Background()
	user, err := userStore.CreateUser(ctx, auth.CreateUserRequest{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Org:       org,
		Role:      role,
		CreatedBy: gofakeit.Email(),
	})
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	t.Cleanup(func() {
		if user != nil {
			userStore.DeleteUserByID(ctx, user.ID) //nolint:errcheck,gosec
		}
	})
	return user
}

func verifyUserToken(user auth.User) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			t := auth.JwtToken{
				Email: user.Email,
				Admin: auth.IsAdmin(user),
				Org:   user.Org,
				Role:  user.Role,
			}
			ctx := context.WithValue(r.Context(), "authToken", t) //nolint: staticcheck
			f(w, r.WithContext(ctx))
		}
	}
}
func userToJSON(u PostUserRequest) io.Reader {
	r, _ := json.Marshal(u)
	return bytes.NewReader(r)
}
