package handlers

import (
	"bz.moh.epi/users/internal"
	"context"
	"github.com/gorilla/mux"
)

// API registers all the server handlers
func API(ctx context.Context, app *internal.App) (*mux.Router, error) {
	r := mux.NewRouter()
	r.HandleFunc("/health", TestHandler)
	return r, nil
}
