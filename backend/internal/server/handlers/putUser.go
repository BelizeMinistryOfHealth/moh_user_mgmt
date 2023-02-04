package handlers

import (
	"bz.moh.epi/users/internal/auth"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type PutUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Org       string `json:"org"`
	Role      string `json:"role"`
}

func (s *UserCrudService) PutUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close() //nolint:errcheck
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodPut {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Verify that the user is admin
	token := r.Context().Value("authToken").(auth.JwtToken)
	if !token.Admin {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	// Get Body
	var requestPayload PutUserRequest
	if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
		log.WithFields(log.Fields{
			"body": r.Body,
		}).WithError(err).Error("Error decoding request when updating a user")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// Get the ID param
	path := r.RequestURI
	paths := strings.Split(path, "/")
	id := paths[len(paths)-1]

	role, err := auth.ToUserRole(requestPayload.Role)
	if err != nil {
		log.WithFields(log.Fields{
			"body": r.Body,
			"user": requestPayload,
		}).WithError(err).Error("Invalid User Role")
		http.Error(w, "User Role provided is no valid", http.StatusBadRequest)
		return
	}

	var user = auth.User{
		ID:        id,
		FirstName: requestPayload.FirstName,
		LastName:  requestPayload.LastName,
		Email:     requestPayload.Email,
		Org:       requestPayload.Org,
		Role:      role,
	}
	if err := s.UserStore.UpdateUser(r.Context(), &user); err != nil {
		log.WithFields(log.Fields{
			"body": r.Body,
			"user": user,
		}).WithError(err).Error("Error updating user")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.WithFields(log.Fields{
			"user":    user,
			"request": r.Body,
		}).WithError(err).Error("Error encoding user as json")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

}
