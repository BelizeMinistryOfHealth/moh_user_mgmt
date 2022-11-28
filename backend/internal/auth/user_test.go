package auth

import (
	"bz.moh.epi/users/internal/db"
	"context"
	firebase "firebase.google.com/go/v4"
	"os"
	"testing"
)

var projectID = os.Getenv("PROJECT_ID")
var apiKey = os.Getenv("API_KEY")

func createUser(t *testing.T, ctx context.Context, userStore UserStore, user User) {
	u, err := userStore.CreateUser(ctx, user)
	if err != nil {
		t.Errorf("error creating user: %v", err)
	}

	if u == nil {
		t.Errorf("want: a non nil user, got nil")
	}

	t.Cleanup(func() {
		err = userStore.DeleteUserById(ctx, u.ID)
		if err != nil {
			t.Errorf("cleaning up user failed: %v", err)
		}
	})
}

func TestUserStore_CreateUser(t *testing.T) {
	ctx := context.Background()
	if len(projectID) == 0 {
		t.Fatalf("PROJECT_ID environment variable is missing")
	}
	if len(apiKey) == 0 {
		t.Fatalf("API_KEY environment variable is missing")
	}

	emulatorHost := os.Getenv("FIRESTORE_EMULATOR_HOST")
	if len(emulatorHost) == 0 {
		t.Fatalf("FIRESTORE_EMULATOR_HOST environment variable is missing")
	}
	emulatorAuthHost := os.Getenv("FIREBASE_AUTH_EMULATOR_HOST")
	if len(emulatorAuthHost) == 0 {
		t.Fatalf("FIREBASE_AUTH_EMULATOR_HOST environment variable is missing")
	}

	firestoreClient, err := db.NewFirestoreClient(ctx, projectID)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}

	// Create User Store
	user := User{
		ID:        "1",
		FirstName: "Roberto",
		LastName:  "Guerra",
		Email:     "uris77@gmail.com",
		UserApplication: UserApplication{
			ApplicationID: "b4718c64-5b4f-4649-ab1a-e8cb5c887a92",
			Name:          "hiv_surveys",
			Permissions:   []string{},
		}}
	userStore, err := NewStore(ctx, firestoreClient, &firebase.Config{
		ProjectID: projectID,
	}, apiKey)

	createUser(t, ctx, userStore, user)

}
