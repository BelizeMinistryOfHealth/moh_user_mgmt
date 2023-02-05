package handlers

import (
	"bytes"
	"bz.moh.epi/users/internal/api"
	"bz.moh.epi/users/internal/auth"
	"context"
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPutUser_FailsIfNotAdmin(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := api.CreateUserApi(*userStore)
	nonAdminUser := createTestUser(t, *userStore, auth.BFLA, auth.PeerNavigatorRole)
	mids := NewChain(verifyUserToken(*nonAdminUser))
	userCrudService := NewUserCrudService(userStore, userApi)
	email := gofakeit.Email()
	createRequest := auth.CreateUserRequest{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     email,
		Org:       auth.BFLA,
		Role:      auth.PeerNavigatorRole,
	}
	usr, err := createUser(ctx, t, *userStore, createRequest)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	putBody := PutUserRequest{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     email,
		Org:       "BFLA",
		Role:      "SR",
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
	userStore := createUserStore(t, ctx)
	userApi := api.CreateUserApi(*userStore)
	adminUser := createTestUser(t, *userStore, auth.BFLA, auth.AdminRole)
	mids := NewChain(verifyUserToken(*adminUser))
	userCrudService := NewUserCrudService(userStore, userApi)
	email := gofakeit.Email()
	createRequest := auth.CreateUserRequest{
		FirstName: "Dan",
		LastName:  "Th",
		Email:     email,
		Org:       auth.BFLA,
		Role:      auth.SrRole,
	}
	usr, err := createUser(ctx, t, *userStore, createRequest)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	putBody := PutUserRequest{
		FirstName: "Dan",
		LastName:  "Th",
		Email:     email,
		Org:       "BFLA",
		Role:      "PeerNavigatorRole",
	}
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/users/%s", usr.ID), putRequestToJSON(putBody))
	mids.Then(userCrudService.PutUser)(res, req)

	if res.Code != 200 {
		t.Errorf("want: %d; got: %d", 200, res.Code)
	}
	wantUser, err := userStore.GetUserByID(ctx, usr.ID)
	if err != nil {
		t.Fatalf("failed to get user by id: %v", err)
	}
	if wantUser.LastName != putBody.LastName {
		t.Errorf("Updating user failed want: %s; got: %s", wantUser.LastName, putBody.LastName)
	}
}
func TestPutUser_UpdatesRoles(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := api.CreateUserApi(*userStore)
	adminUser := createTestUser(t, *userStore, auth.BFLA, auth.AdminRole)
	mids := NewChain(verifyUserToken(*adminUser))
	userCrudService := NewUserCrudService(userStore, userApi)
	email := gofakeit.Email()
	createRequest := auth.CreateUserRequest{
		FirstName: "Dan",
		LastName:  "Th",
		Email:     email,
		Org:       auth.BFLA,
		Role:      auth.PeerNavigatorRole,
	}
	usr, err := createUser(ctx, t, *userStore, createRequest)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	putBody := PutUserRequest{
		FirstName: "Dan",
		LastName:  "Th",
		Email:     email,
		Org:       "BFLA",
		Role:      "AdherenceCounselorRole",
	}
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/users/%s", usr.ID), putRequestToJSON(putBody))
	mids.Then(userCrudService.PutUser)(res, req)

	if res.Code != 200 {
		t.Errorf("want: %d; got: %d", 200, res.Code)
	}
	gotUser, _ := userStore.GetUserByID(ctx, usr.ID)
	if gotUser.Role.String() != putBody.Role {
		t.Errorf("Updating user failed got: %s; want: %s", gotUser.Role, putBody.Role)
	}
}

func putRequestToJSON(req PutUserRequest) io.Reader {
	b, _ := json.Marshal(req)
	return bytes.NewBuffer(b)
}
