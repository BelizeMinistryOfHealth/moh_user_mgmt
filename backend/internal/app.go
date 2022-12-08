package internal

import (
	"bz.moh.epi/users/internal/auth"
	"bz.moh.epi/users/internal/db"
)

// App represents the components of the application that are needed internally for
// us to communicate with different infrastructural elements.
type App struct {
	Firestore       *db.FirestoreClient
	UserStore       *auth.UserStore
	ProjectID       string
	FirestoreApiKey string
}
