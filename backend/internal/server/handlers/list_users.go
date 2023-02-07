package handlers

import (
	"bz.moh.epi/users/internal/auth"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// ListUsers returns a list of all users
func (s *UserCrudService) ListUsers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close() //nolint:errcheck
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	// Verify that user is admin
	token := r.Context().Value("authToken").(auth.JwtToken)
	if !token.Admin {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	users, err := s.UserApi.ListUsers(r.Context(), token.Email)
	if err != nil {
		log.WithError(err).Error("error retrieving users from the user store")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var response []UserResponse
	for i := range users {
		response = append(response, UserResponse{
			ID:        users[i].ID,
			FirstName: users[i].FirstName,
			LastName:  users[i].LastName,
			Email:     users[i].Email,
			Org:       users[i].Org.String(),
			Role:      users[i].Role.String(),
			Enabled:   users[i].Enabled,
		})
	}
	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.WithFields(log.Fields{
			"users":    users,
			"response": response,
		}).WithError(err).Error("error encoding users")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
