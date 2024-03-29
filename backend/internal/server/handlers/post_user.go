package handlers

import (
	"bz.moh.epi/users/internal/auth"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// PostUserRequest is the payload that is expected by the handler that creates a new user
type PostUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Org       string `json:"org"`
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
	var requestPayload PostUserRequest
	if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
		log.WithFields(log.Fields{
			"body": r.Body,
		}).WithError(err).Error("Error decoding request when creating a new user")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	role, err := auth.ToUserRole(requestPayload.Role)
	if err != nil {
		log.WithFields(log.Fields{
			"body": r.Body,
			"user": requestPayload,
		}).WithError(err).Error("Invalid User Role")
		http.Error(w, "User Role provided is not valid", http.StatusBadRequest)
		return
	}
	org, err := auth.ToOrg(requestPayload.Org)
	if err != nil {
		log.WithFields(log.Fields{
			"body": r.Body,
			"user": requestPayload,
		}).WithError(err).Error("Invalid Organization")
		http.Error(w, "Organization provided is not valid", http.StatusBadRequest)
		return
	}

	createRequest := auth.CreateUserRequest{
		FirstName: requestPayload.FirstName,
		LastName:  requestPayload.LastName,
		Email:     requestPayload.Email,
		Org:       org,
		Role:      role,
		CreatedBy: token.Email,
	}

	user, err := s.UserApi.CreateUser(r.Context(), createRequest)

	//user, err := s.UserStore.CreateUser(r.Context(), createRequest)
	if err != nil {
		log.WithFields(log.Fields{
			"body": r.Body,
		}).WithError(err).Error("Error creating a new user")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response := UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Org:       user.Org.String(),
		Role:      user.Role.String(),
		Enabled:   user.Enabled,
	}
	// Return user as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.WithFields(log.Fields{
			"user":     user,
			"response": response,
			"request":  r.Body,
		}).WithError(err).Error("Error encoding user as json")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
