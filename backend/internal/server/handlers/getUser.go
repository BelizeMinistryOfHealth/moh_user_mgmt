package handlers

import (
	"bz.moh.epi/users/internal/auth"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func (s *UserCrudService) GetUserByID(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close() //nolint:errcheck
	if r.Method == http.MethodOptions {
		return
	}
	// Verify that user is admin
	token := r.Context().Value("authToken").(auth.JwtToken)
	if !token.Admin {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	// Get the ID param
	path := r.RequestURI
	paths := strings.Split(path, "/")
	id := paths[len(paths)-1]

	user, err := s.UserStore.GetUserByID(r.Context(), id)
	if err != nil {
		log.WithError(err).Error("error retrieving user by ID")
		http.Error(w, "error retrieving user from the database", http.StatusInternalServerError)
		return
	}
	log.WithFields(log.Fields{
		"id":      id,
		"user":    user,
		"request": r.RequestURI,
	}).Info("Retrieved user")
	if err = json.NewEncoder(w).Encode(user); err != nil {
		log.WithFields(log.Fields{
			"user": user,
		}).WithError(err).Error("error encoding user")
		http.Error(w, "error encoding user", http.StatusInternalServerError)
		return
	}
}
