package auth

import (
	"bz.moh.epi/users/internal/db"
	"context"
	firebase "firebase.google.com/go/v4"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"os"
	"testing"
)

var projectID = os.Getenv("PROJECT_ID") //nolint: gochecknoglobals
var apiKey = os.Getenv("API_KEY")       //nolint: gochecknoglobals

func createUser(ctx context.Context, t *testing.T, userStore UserStore, user CreateUserRequest) *User {
	u, err := userStore.CreateUser(ctx, user)
	if err != nil {
		t.Errorf("error creating user: %v", err)
	}

	if u == nil {
		t.Errorf("want: a non nil user, got nil")
	}

	t.Cleanup(func() {
		err = userStore.DeleteUserByID(ctx, u.ID)
		if err != nil {
			t.Errorf("cleaning up user failed: %v", err)
		}
	})
	return u
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
		UserApplications: []UserApplication{{
			ApplicationID: "b4718c64-5b4f-4649-ab1a-e8cb5c887a92",
			Name:          "hiv_surveys",
			Permissions:   []string{},
		}}}
	userStore, _ := NewStore(firestoreClient, apiKey)

	createUser(ctx, t, userStore, user)

}

func TestUserStore_UpdateUser(t *testing.T) {
	ctx := context.Background()
	firebaseConfig := &firebase.Config{ProjectID: projectID}
	firestoreClient, err := db.NewFirestoreClient(ctx, firebaseConfig)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	ID := uuid.New().String()
	user := CreateUserRequest{
		FirstName:        "Roberto",
		LastName:         "Guerra",
		Email:            "uris77@gmail.com",
		UserApplications: []UserApplication{},
	}
	userStore, err := NewStore(firestoreClient, apiKey)
	if err != nil {
		t.Fatalf("failed to create user store: %v", err)
	}

	testUser := createUser(ctx, t, userStore, user)

	var permissionsTcs = []struct {
		name  string
		input UserApplication
		want  []UserApplication
	}{
		{"Should add permissions",
			UserApplication{
				ApplicationID: ID,
				Name:          "application 1",
				Permissions:   []string{"view"},
			},
			[]UserApplication{{
				ApplicationID: ID,
				Name:          "application 1",
				Permissions:   []string{"view"},
			}},
		},
		{
			"Should be able to add more than one permissions to the same application",
			UserApplication{
				ApplicationID: ID,
				Name:          "application 1",
				Permissions:   []string{"view", "edit"},
			},
			[]UserApplication{{
				ApplicationID: ID,
				Name:          "application 1",
				Permissions:   []string{"view", "edit"},
			}},
		},
		{
			"Should be able to add multiple applications",
			UserApplication{
				ApplicationID: ID,
				Name:          "application 1",
				Permissions:   []string{"view", "edit"},
			},
			[]UserApplication{{
				ApplicationID: ID,
				Name:          "application 1",
				Permissions:   []string{"view", "edit"},
			}},
		},
	}

	for _, tt := range permissionsTcs {
		t.Run(tt.name, func(t *testing.T) {
			testUser.UserApplications = []UserApplication{tt.input}
			err := userStore.UpdateUser(ctx, testUser)
			if err != nil {
				t.Errorf("unexpected error updating user: %v", err)
			}
			if diff := cmp.Diff(tt.want, testUser.UserApplications); diff != "" {
				t.Errorf("UpdateUser mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUserStore_UpdateUser_AddMultipleApplications(t *testing.T) {
	ctx := context.Background()
	firebaseConfig := &firebase.Config{ProjectID: projectID}
	firestoreClient, err := db.NewFirestoreClient(ctx, firebaseConfig)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	ID := uuid.New().String()
	user := CreateUserRequest{
		FirstName:        "Roberto",
		LastName:         "Guerra",
		Email:            "uris77@gmail.com",
		UserApplications: []UserApplication{},
	}
	userStore, err := NewStore(firestoreClient, apiKey)
	if err != nil {
		t.Fatalf("failed to create user store: %v", err)
	}

	app2ID := uuid.New().String()

	testUser := createUser(ctx, t, userStore, user)
	want := []UserApplication{
		{
			ApplicationID: ID,
			Name:          "application 1",
			Permissions:   []string{"view", "edit"},
		},
		{
			ApplicationID: app2ID,
			Name:          "application 2",
			Permissions:   []string{"view", "edit"},
		},
	}
	testUser.UserApplications = []UserApplication{
		{
			ApplicationID: ID,
			Name:          "application 1",
			Permissions:   []string{"view", "edit"},
		},
		{
			ApplicationID: app2ID,
			Name:          "application 2",
			Permissions:   []string{"view", "edit"},
		},
	}

	err = userStore.UpdateUser(ctx, testUser)
	if err != nil {
		t.Errorf("failed to update user: %v", err)
	}

	if diff := cmp.Diff(want, testUser.UserApplications); diff != "" {
		t.Errorf("UpdateUser mismatch (-want +got)\n%s", diff)
	}
}

func TestUserStore_UpdateUser_UpdatesNames(t *testing.T) {
	ctx := context.Background()
	firebaseConfig := &firebase.Config{ProjectID: projectID}
	firestoreClient, err := db.NewFirestoreClient(ctx, firebaseConfig)
	if err != nil {
		t.Fatalf("failed to create firestore client: %v", err)
	}
	user := CreateUserRequest{
		FirstName:        "Roberto",
		LastName:         "Guerra",
		Email:            "uris77@gmail.com",
		UserApplications: []UserApplication{},
	}
	userStore, err := NewStore(firestoreClient, apiKey)
	if err != nil {
		t.Fatalf("failed to create user store: %v", err)
	}

	testUser := createUser(ctx, t, userStore, user)

	var testCases = []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Should update first name",
			input: "UpdateName",
			want:  "UpdateName",
		},
		{
			name:  "Should update last name",
			input: "UpdatedLastName",
			want:  "UpdatedLastName",
		},
	}

	for idx, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if idx == 0 {
				testUser.FirstName = tt.input
				err := userStore.UpdateUser(ctx, testUser)
				if err != nil {
					t.Errorf("error updating user first name: %v", err)
				}

				u, _ := userStore.GetUserByEmail(ctx, testUser.Email)
				if u.FirstName != tt.want {
					t.Errorf("UpdateUser first name | want: %s got: %s", tt.want, u.FirstName)
				}

			}
			if idx == 1 {
				testUser.LastName = tt.input
				err := userStore.UpdateUser(ctx, testUser)
				if err != nil {
					t.Errorf("error updating user first name: %v", err)
				}
				u, _ := userStore.GetUserByEmail(ctx, testUser.Email)
				if u.LastName != tt.want {
					t.Errorf("UpdateUser last name | want: %s got: %s", tt.want, u.LastName)
				}

			}
		})
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
		FirstName:        "Roberto",
		LastName:         "Guerra",
		Email:            "uris77@gmail.com",
		UserApplications: []UserApplication{},
	}
	userStore, err := NewStore(firestoreClient, apiKey)
	if err != nil {
		t.Fatalf("failed to create user store: %v", err)
	}

	testUser := createUser(ctx, t, userStore, user)

	token, err := userStore.CreateToken(ctx, testUser.ID)
	if err != nil {
		t.Errorf("error creating token for user: %v", err)
	}

	if token == "" {
		t.Errorf("wanted non-empty token, got empty string")
	}
}

func createApplication(ctx context.Context, t *testing.T, userStore UserStore, request CreateApplicationRequest) *UserApplication {
	app, err := userStore.CreateApplication(ctx, request)
	if err != nil {
		t.Fatalf("error creating user application: %v", err)
	}
	t.Cleanup(func() {
		err = userStore.DeleteApplicationByID(ctx, app.ApplicationID)
		if err != nil {
			t.Errorf("cleaning up user applicatino (%s) failed: %v", app.ApplicationID, err)
		}
	})
	return app
}
func TestUserStore_ListPermissions(t *testing.T) {
	ctx := context.Background()
	firebaseConfig := &firebase.Config{ProjectID: projectID}
	firestoreClient, err := db.NewFirestoreClient(ctx, firebaseConfig)
	if err != nil {
		t.Fatalf("error creating firestore client: %v", err)
	}
	userStore, err := NewStore(firestoreClient, apiKey)
	if err != nil {
		t.Fatalf("error creating user store: %v", err)
	}

	// Create a list of applications
	createApplication(ctx, t, userStore,
		CreateApplicationRequest{ //nolint: errcheck,gosec
			Name:        "application 1",
			Permissions: []string{"admin"},
		})
	createApplication(ctx, t, userStore, CreateApplicationRequest{ //nolint:errcheck,gosec
		Name:        "application 2",
		Permissions: []string{},
	})
	applications, err := userStore.ListApplications(ctx)
	if err != nil {
		t.Fatalf("error retrieving applications: %v", err)
	}
	// It includes the default one created when the docker container is built
	if len(applications) != 3 {
		t.Errorf("wanted 2 applications but got %d", len(applications))
	}
}
