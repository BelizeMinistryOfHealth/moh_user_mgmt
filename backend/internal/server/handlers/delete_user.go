package handlers

import (
	"bz.moh.epi/users/internal/api"
	"bz.moh.epi/users/internal/auth"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func (s *UserCrudService) DeleteUser(w http.ResponseWriter, r *http.Request) {
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

	err := s.UserApi.DeleteUser(r.Context(), api.DeleteUserRequest{
		ID:        id,
		DeletedBy: token.Email,
	})

	if err != nil {
		log.WithFields(log.Fields{
			"id":          id,
			"requestedBy": token.Email,
			"role":        token.Role,
			"org":         token.Org,
		}).WithError(err).Error("error deleting user by ID")
		http.Error(w, "error deleting user from the database", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("{\"message\": \"User deleted successfully\"}")) //nolint:errcheck,gosec
}
