package handlers

import (
	"bz.moh.epi/users/internal/api"
	"bz.moh.epi/users/internal/auth"
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteUser(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := api.CreateUserApi(*userStore)
	adminUser := createTestUser(t, *userStore, auth.MOHW, auth.AdminRole)
	userRequest := auth.CreateUserRequest{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Org:       auth.NAC,
		Role:      auth.AdherenceCounselorRole,
	}
	testUser, err := createUser(ctx, t, *userStore, userRequest)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	mids := NewChain(verifyUserToken(*adminUser))
	userCrudService := NewUserCrudService(userStore, userApi)
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%s", testUser.ID), nil)
	mids.Then(userCrudService.DeleteUser)(res, req)
	if res.Code != 200 {
		t.Errorf("want: %d; got: %d", 200, res.Code)
	}
}
