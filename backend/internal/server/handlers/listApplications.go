package handlers

import (
	"bz.moh.epi/users/internal/auth"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// ListApplications returns a list of applications if the user is an Admin.
func (s *UserCrudService) ListApplications(w http.ResponseWriter, r *http.Request) {
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
	apps, err := s.UserStore.ListApplications(r.Context())
	if err != nil {
		log.WithError(err).Error("error retrieving users from the database")
		http.Error(w, "error retrieving users from the database", http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(apps); err != nil {
		log.WithFields(log.Fields{
			"applications": apps,
		}).WithError(err).Error("error encoding applications")
		http.Error(w, "error encoding applications", http.StatusInternalServerError)
		return
	}
}
