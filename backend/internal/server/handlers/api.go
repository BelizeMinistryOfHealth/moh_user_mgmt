package handlers

import (
	"bz.moh.epi/users/internal"
	"bz.moh.epi/users/internal/api"
	"bz.moh.epi/users/internal/auth"
	"context"
	"github.com/gorilla/mux"
	"net/http"
)

// API registers all the server handlers
func API(ctx context.Context, app *internal.App) (*mux.Router, error) {
	r := mux.NewRouter()
	mids := NewChain(EnableCors(), VerifyToken(app.UserStore), JsonContentType())
	r.HandleFunc("/health", TestHandler)
	userCrudService := NewUserCrudService(app.UserStore, app.UserApi)
	r.HandleFunc("/user",
		mids.Then(userCrudService.PostUser)).
		Methods(http.MethodOptions, http.MethodPost)
	r.HandleFunc("/users",
		mids.Then(userCrudService.ListUsers)).
		Methods(http.MethodOptions, http.MethodGet)
	r.HandleFunc("/users/{id}",
		mids.Then(userCrudService.GetUserByID)).
		Methods(http.MethodOptions, http.MethodGet)
	r.HandleFunc("/users/{id}", mids.
		Then(userCrudService.PutUser)).
		Methods(http.MethodOptions, http.MethodPut)
	r.HandleFunc("/applications",
		mids.Then(userCrudService.ListApplications)).
		Methods(http.MethodOptions, http.MethodGet)
	return r, nil
}

type UserCrudService struct {
	UserStore *auth.UserStore
	UserApi   *api.UserApi
}

func NewUserCrudService(userStore *auth.UserStore, userApi *api.UserApi) *UserCrudService {
	return &UserCrudService{
		UserStore: userStore,
		UserApi:   userApi,
	}
}
