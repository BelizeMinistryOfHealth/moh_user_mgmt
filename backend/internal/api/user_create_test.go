package api

import (
	"bz.moh.epi/users/internal/auth"
	"bz.moh.epi/users/internal/db"
	"context"
	firebase "firebase.google.com/go/v4"
	"github.com/brianvoe/gofakeit/v6"
	"os"
	"testing"
)

var projectID = os.Getenv("PROJECT_ID") //nolint: gochecknoglobals
var apiKey = os.Getenv("API_KEY")       //nolint: gochecknoglobals

func TestCreateUserApi_NAC_AdminsCanCreateUsers(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	nacUser := createAdminUser(t, ctx, userStore, "NAC")

	var testCases = []struct {
		name  string
		input auth.CreateUserRequest
	}{
		{"NAC Admins can create users at NAC",
			auth.CreateUserRequest{
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
				Email:     gofakeit.Email(),
				Org:       "NAC",
				Role:      auth.AdminRole,
				CreatedBy: nacUser.Email,
			},
		},
		{
			"NAC Admins can create users at BFLA",
			auth.CreateUserRequest{
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
				Email:     gofakeit.Email(),
				Org:       "BFLA",
				Role:      auth.AdminRole,
				CreatedBy: nacUser.Email,
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userApi.CreateUser(ctx, tt.input)
			if err != nil {
				t.Errorf("CreateUser() error = %v", err)
			}
			if got == nil {
				t.Errorf("CreateUser() got = %v", got)
			}
		})
	}
}

func TestCreateUserApi_MOHW_AdminsCanCreateUsers(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	nacUser := createAdminUser(t, ctx, userStore, "NAC")

	var testCases = []struct {
		name  string
		input auth.CreateUserRequest
	}{
		{"MOHW Admins can create users at MOHW",
			auth.CreateUserRequest{
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
				Email:     gofakeit.Email(),
				Org:       "MOHW",
				Role:      auth.AdminRole,
				CreatedBy: nacUser.Email,
			},
		},
		{
			"MOHW Admins can create users at BFLA",
			auth.CreateUserRequest{
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
				Email:     gofakeit.Email(),
				Org:       "BFLA",
				Role:      auth.AdminRole,
				CreatedBy: nacUser.Email,
			},
		},
		{
			"MOHW Admins can create users at NAC",
			auth.CreateUserRequest{
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
				Email:     gofakeit.Email(),
				Org:       "NAC",
				Role:      auth.AdminRole,
				CreatedBy: nacUser.Email,
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userApi.CreateUser(ctx, tt.input)
			if err != nil {
				t.Errorf("CreateUser() error = %v", err)
			}
			if got == nil {
				t.Errorf("CreateUser() got = %v", got)
			}
		})
	}
}
func createUserStore(t *testing.T, ctx context.Context) auth.UserStore {
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
	userStore, _ := auth.NewStore(firestoreClient, apiKey)
	return userStore
}

func createAdminUser(t *testing.T, ctx context.Context, userStore auth.UserStore, org string) *auth.User {
	user := auth.CreateUserRequest{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Org:       org,
		Role:      auth.AdminRole,
		CreatedBy: gofakeit.Email(),
	}
	u, err := userStore.CreateUser(ctx, user)
	if err != nil {
		t.Fatalf("error creating user: %v", err)
	}
	return u
}
