package handlers

import (
	"bz.moh.epi/users/internal/auth"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// PostUserRequest is the payload that is expected by the handler that creates a new user
type PostUserRequest struct {
	FirstName    string                 `json:"firstName"`
	LastName     string                 `json:"lastName"`
	Email        string                 `json:"email"`
	Applications []auth.UserApplication `json:"applications"`
}

// PostUser is the handler for creating a new user
func (s *UserCrudService) PostUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close() //nolint:errcheck
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	// Verify that user is admin
	token := r.Context().Value("authToken").(auth.JwtToken)
	if !token.Admin {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	// Get the body
	var requestPayload auth.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
		log.WithFields(log.Fields{
			"body": r.Body,
		}).WithError(err).Error("Error decoding request when creating a new user")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	user, err := s.UserStore.CreateUser(r.Context(), requestPayload)
	if err != nil {
		log.WithFields(log.Fields{
			"body": r.Body,
		}).WithError(err).Error("Error creating a new user")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Return user as json
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.WithFields(log.Fields{
			"user":    user,
			"request": r.Body,
		}).WithError(err).Error("Error encoding user as json")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
