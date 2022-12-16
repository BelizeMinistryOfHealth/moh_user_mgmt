package handlers

import (
	"bz.moh.epi/users/internal/auth"
	"context"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

// EnableCors enables CORS
func EnableCors() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Referer, Connection")
			f(w, r)
		}
	}
}

// VerifyToken is a middleware that checks the authorization header for every a handler and
// adds information related to the user to the context, which is passed down to the handler.
func VerifyToken(store *auth.UserStore) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				f(w, r)
				return
			}
			token := strings.Split(r.Header.Get("Authorization"), " ")
			if len(token) != 2 {
				log.WithFields(log.Fields{
					"token": token,
				}).Error("token has invalid format")
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			jwtToken, err := store.VerifyToken(r.Context(), token[1])
			if err != nil {
				log.WithFields(log.Fields{
					"token": token,
				}).WithError(err).Error("error verifying token")
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "authToken", jwtToken) //nolint: staticcheck
			f(w, r.WithContext(ctx))
		}
	}
}

// JsonContentType sets the response content type to application/json
func JsonContentType() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			f(w, r)
		}
	}
}
