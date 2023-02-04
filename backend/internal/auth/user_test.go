package auth

import (
	"bz.moh.epi/users/internal/db"
	"context"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"os"
	"testing"
)

var projectID = os.Getenv("PROJECT_ID") //nolint: gochecknoglobals
var apiKey = os.Getenv("API_KEY")       //nolint: gochecknoglobals

func createUser(ctx context.Context, t *testing.T, userStore UserStore, user CreateUserRequest) (*User, error) {
	u, err := userStore.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	if u == nil {
		return nil, fmt.Errorf("want: a non nil user, got nil") //nolint: goerr113
	}

	t.Cleanup(func() {
		err = userStore.DeleteUserByID(ctx, u.ID)
		if err != nil {
			t.Errorf("cleaning up user failed: %v", err)
		}
	})
	return u, nil
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
	firebaseConfig := &firebase.Config{ProjectID: projectID}
	firestoreClient, err := db.NewFirestoreClient(ctx, firebaseConfig)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}

	// Create User Store
	user := CreateUserRequest{
		FirstName: "Roberto",
		LastName:  "Guerra",
		Email:     "uris77@gmail.com",
		Role:      PeerNavigatorRole,
		Org:       "BFLA",
		CreatedBy: "some@mail.com",
	}
	userStore, _ := NewStore(firestoreClient, apiKey)

	createUser(ctx, t, userStore, user) //nolint: errcheck

}

func TestUserStore_UpdateUser(t *testing.T) {
	ctx := context.Background()
	firebaseConfig := &firebase.Config{ProjectID: projectID}
	firestoreClient, err := db.NewFirestoreClient(ctx, firebaseConfig)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	user := CreateUserRequest{
		FirstName: "Roberto",
		LastName:  "Guerra",
		Email:     "uris77@gmail.com",
		Org:       "BFLA",
		Role:      PeerNavigatorRole,
		CreatedBy: "some@mailcom",
	}
	userStore, err := NewStore(firestoreClient, apiKey)
	if err != nil {
		t.Fatalf("failed to create user store: %v", err)
	}

	testUser, err := createUser(ctx, t, userStore, user)
	if err != nil {
		t.Fatalf("error creating user: %v", err)
	}

	var permissionsTcs = []struct {
		name  string
		input UserRole
		want  string
	}{
		{"Should add role",
			PeerNavigatorRole,
			"PeerNavigatorRole",
		},
		{
			"Should change role",
			AdherenceCounselorRole,
			"AdherenceCounselorRole",
		},
	}

	for _, tt := range permissionsTcs {
		t.Run(tt.name, func(t *testing.T) {
			testUser.Role = tt.input
			err := userStore.UpdateUser(ctx, testUser)
			if err != nil {
				t.Errorf("unexpected error updating user: %v", err)
			}
			if tt.want != testUser.Role.String() {
				t.Errorf("UpdateUser mismatch want: %s got: %s:", tt.want, testUser.Role)
			}
		})
	}
}

func TestUserStore_GetUserByID(t *testing.T) {

	ctx := context.Background()
	firebaseConfig := &firebase.Config{ProjectID: projectID}
	firestoreClient, err := db.NewFirestoreClient(ctx, firebaseConfig)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	user := CreateUserRequest{
		FirstName: "Roberto",
		LastName:  "Guerra",
		Email:     "uris77@gmail.com",
		Org:       "BFLA",
		Role:      PeerNavigatorRole,
		CreatedBy: "some@mail.com",
	}
	userStore, err := NewStore(firestoreClient, apiKey)
	if err != nil {
		t.Fatalf("failed to create user store: %v", err)
	}

	testUser, err := createUser(ctx, t, userStore, user)
	if err != nil {
		t.Fatalf("error creating user: %v", err)
	}
	retrievedUser, err := userStore.GetUserByID(ctx, testUser.ID)
	if err != nil {
		t.Fatalf("error retrieving user by id: %v", err)
	}
	if retrievedUser.ID != testUser.ID {
		t.Errorf("GetUserByID failed, want: %s, got: %s", testUser.ID, retrievedUser.ID)
	}
}

func TestUserStore_CreateToken(t *testing.T) {
	ctx := context.Background()
	firebaseConfig := &firebase.Config{ProjectID: projectID}
	firestoreClient, err := db.NewFirestoreClient(ctx, firebaseConfig)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	user := CreateUserRequest{
		FirstName: "Roberto",
		LastName:  "Guerra",
		Email:     "uris77@gmail.com",
		Org:       "BFLA",
		Role:      PeerNavigatorRole,
		CreatedBy: "some@mail.com",
	}
	userStore, err := NewStore(firestoreClient, apiKey)
	if err != nil {
		t.Fatalf("failed to create user store: %v", err)
	}

	testUser, err := createUser(ctx, t, userStore, user)
	if err != nil {
		t.Fatalf("error creating user: %v", err)
	}

	token, err := userStore.CreateToken(ctx, testUser.ID)
	if err != nil {
		t.Errorf("error creating token for user: %v", err)
	}

	if token == "" {
		t.Errorf("wanted non-empty token, got empty string")
	}
}
