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
	nacUser := createAdminUser(t, ctx, userStore, auth.NAC)

	var testCases = []struct {
		name  string
		input auth.CreateUserRequest
	}{
		{"NAC Admins can create users at NAC",
			auth.CreateUserRequest{
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
				Email:     gofakeit.Email(),
				Org:       auth.NAC,
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
				Org:       auth.BFLA,
				Role:      auth.AdminRole,
				CreatedBy: nacUser.Email,
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userApi.CreateUser(ctx, tt.input)
			t.Cleanup(func() {
				if got != nil {
					_ = userStore.DeleteUserByID(ctx, got.ID) //nolint: errcheck,gosec
				}
			})
			if err != nil {
				t.Errorf("CreateUser() error = %v", err)
			}
			if got == nil {
				t.Errorf("CreateUser() got = %v", got)
			}
		})
	}
}

func TestCreateUserApi_NAC_AdminsCanNotCreateMOHWUsers(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	nacUser := createAdminUser(t, ctx, userStore, auth.NAC)
	createRequest := auth.CreateUserRequest{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Org:       auth.MOHW,
		Role:      auth.AdminRole,
		CreatedBy: nacUser.Email,
	}
	_, err := userApi.CreateUser(ctx, createRequest)
	if err == nil {
		t.Errorf("CreateUser() should return an error = %v", err)
	}
}

func TestCreateUserApi_RequestorMustExistInDatabase(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	createRequest := auth.CreateUserRequest{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Org:       auth.MOHW,
		Role:      auth.AdminRole,
		CreatedBy: gofakeit.Email(),
	}
	_, err := userApi.CreateUser(ctx, createRequest)
	if err == nil {
		t.Errorf("CreateUser() should return an error because requestor is not in database ")
	}
}

func TestCreateUserApi_NonAdminRolesCanNotCreateUsers(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	nonAdminUser := createNonAdminUser(t, ctx, userStore, auth.MOHW)
	createRequest := auth.CreateUserRequest{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Org:       auth.MOHW,
		Role:      auth.SrRole,
		CreatedBy: nonAdminUser.Email,
	}
	_, err := userApi.CreateUser(ctx, createRequest)
	if err == nil {
		t.Errorf("CreateUser() should return an error because non admin can not create users")
	}

}

func TestCreateUserApi_MOHW_AdminsCanCreateUsers(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	mohwUser := createAdminUser(t, ctx, userStore, auth.MOHW)

	var testCases = []struct {
		name  string
		input auth.CreateUserRequest
	}{
		{"MOHW Admins can create users at MOHW",
			auth.CreateUserRequest{
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
				Email:     gofakeit.Email(),
				Org:       auth.MOHW,
				Role:      auth.AdminRole,
				CreatedBy: mohwUser.Email,
			},
		},
		{
			"MOHW Admins can create users at BFLA",
			auth.CreateUserRequest{
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
				Email:     gofakeit.Email(),
				Org:       auth.BFLA,
				Role:      auth.AdminRole,
				CreatedBy: mohwUser.Email,
			},
		},
		{
			"MOHW Admins can create users at NAC",
			auth.CreateUserRequest{
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
				Email:     gofakeit.Email(),
				Org:       auth.NAC,
				Role:      auth.AdminRole,
				CreatedBy: mohwUser.Email,
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userApi.CreateUser(ctx, tt.input)
			t.Cleanup(func() {
				if got != nil {
					userStore.DeleteUserByID(ctx, got.ID) //nolint:errcheck,gosec
				}
			})
			if err != nil {
				t.Errorf("CreateUser() error = %v", err)
			}
			if got == nil {
				t.Errorf("CreateUser() got = %v", got)
			}
		})
	}
}

func TestCreateUserApi_BFLA_AdminsCanCreateUsers(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	bflaUser := createAdminUser(t, ctx, userStore, auth.BFLA)

	var testCases = []struct {
		name  string
		input auth.CreateUserRequest
	}{
		{"BFLA Admins can create Peer Navigators at BFLA",
			auth.CreateUserRequest{
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
				Email:     gofakeit.Email(),
				Org:       auth.BFLA,
				Role:      auth.PeerNavigatorRole,
				CreatedBy: bflaUser.Email,
			},
		},
		{
			"BFLA Admins can create Adherence Counselors at BFLA",
			auth.CreateUserRequest{
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
				Email:     gofakeit.Email(),
				Org:       auth.BFLA,
				Role:      auth.AdherenceCounselorRole,
				CreatedBy: bflaUser.Email,
			},
		},
		{
			"BFLA Admins can create SRs at BFLA",
			auth.CreateUserRequest{
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
				Email:     gofakeit.Email(),
				Org:       auth.BFLA,
				Role:      auth.SrRole,
				CreatedBy: bflaUser.Email,
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userApi.CreateUser(ctx, tt.input)
			t.Cleanup(func() {
				if got != nil {
					userStore.DeleteUserByID(ctx, got.ID) //nolint:errcheck,gosec
				}

			})
			if err != nil {
				t.Errorf("CreateUser() error = %v", err)
			}
			if got == nil {
				t.Errorf("CreateUser() got = %v", got)
			}
		})
	}
}

func TestCreateUserApi_BFLA_CanNotCreateMOHWAndNACUsers(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	bflaUser := createAdminUser(t, ctx, userStore, auth.BFLA)
	var testCases = []struct {
		name  string
		input auth.CreateUserRequest
	}{

		{
			"BFLA Admins can not create MOHW Admins",
			auth.CreateUserRequest{
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
				Email:     gofakeit.Email(),
				Org:       auth.MOHW,
				Role:      auth.AdminRole,
				CreatedBy: bflaUser.Email,
			},
		},
		{
			"BFLA Admins can not create GOJOVEN Admins",
			auth.CreateUserRequest{
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
				Email:     gofakeit.Email(),
				Org:       auth.GOJOVEN,
				Role:      auth.AdminRole,
				CreatedBy: bflaUser.Email,
			},
		},
		{
			"BFLA Admins can not create NAC Admins",
			auth.CreateUserRequest{
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
				Email:     gofakeit.Email(),
				Org:       auth.NAC,
				Role:      auth.AdminRole,
				CreatedBy: bflaUser.Email,
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := userApi.CreateUser(ctx, tt.input)
			if err == nil {
				t.Errorf("CreateUser() expected an error")
			}
		})
	}
}
func TestCreateUserApi_GOJOVEN_CanNotCreateMOHWAndNACUsers(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	csoUser := createAdminUser(t, ctx, userStore, auth.GOJOVEN)
	var testCases = []struct {
		name  string
		input auth.CreateUserRequest
	}{

		{
			"GOJOVEN Admins can not create MOHW Admins",
			auth.CreateUserRequest{
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
				Email:     gofakeit.Email(),
				Org:       auth.MOHW,
				Role:      auth.AdminRole,
				CreatedBy: csoUser.Email,
			},
		},
		{
			"GOJOVEN Admins can not create BFLA Admins",
			auth.CreateUserRequest{
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
				Email:     gofakeit.Email(),
				Org:       auth.BFLA,
				Role:      auth.AdminRole,
				CreatedBy: csoUser.Email,
			},
		},
		{
			"GOJOVEN Admins can not create NAC Admins",
			auth.CreateUserRequest{
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
				Email:     gofakeit.Email(),
				Org:       auth.NAC,
				Role:      auth.AdminRole,
				CreatedBy: csoUser.Email,
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := userApi.CreateUser(ctx, tt.input)
			if err == nil {
				t.Errorf("CreateUser() expected an error")
			}
		})
	}
}
func TestCreateUserApi_GOJOVEN_AdminsCanCreateUsers(t *testing.T) {
	ctx := context.Background()
	userStore := createUserStore(t, ctx)
	userApi := CreateUserApi(userStore)
	csoUser := createAdminUser(t, ctx, userStore, auth.GOJOVEN)

	var testCases = []struct {
		name  string
		input auth.CreateUserRequest
	}{
		{"GOJOVEN Admins can create Peer Navigators at GOJOVEN",
			auth.CreateUserRequest{
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
				Email:     gofakeit.Email(),
				Org:       auth.GOJOVEN,
				Role:      auth.PeerNavigatorRole,
				CreatedBy: csoUser.Email,
			},
		},
		{
			"GOJOVEN Admins can create Adherence Counselors at GOJOVEN",
			auth.CreateUserRequest{
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
				Email:     gofakeit.Email(),
				Org:       auth.GOJOVEN,
				Role:      auth.AdherenceCounselorRole,
				CreatedBy: csoUser.Email,
			},
		},
		{
			"GOJOVEN Admins can create SRs at GOJOVEN",
			auth.CreateUserRequest{
				FirstName: gofakeit.FirstName(),
				LastName:  gofakeit.LastName(),
				Email:     gofakeit.Email(),
				Org:       auth.GOJOVEN,
				Role:      auth.SrRole,
				CreatedBy: csoUser.Email,
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userApi.CreateUser(ctx, tt.input)
			t.Cleanup(func() {
				if got != nil {
					userStore.DeleteUserByID(ctx, got.ID) // nolint:errcheck,gosec
				}
			})
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

func createAdminUser(t *testing.T, ctx context.Context, userStore auth.UserStore, org auth.Org) *auth.User {
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

	t.Cleanup(func() {
		userStore.DeleteUserByEmail(ctx, u.Email) //nolint:errcheck,gosec
	})
	return u
}
func createNonAdminUser(t *testing.T, ctx context.Context, userStore auth.UserStore, org auth.Org) *auth.User {
	user := auth.CreateUserRequest{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Org:       org,
		Role:      auth.PeerNavigatorRole,
		CreatedBy: gofakeit.Email(),
	}
	u, err := userStore.CreateUser(ctx, user)
	if err != nil {
		t.Fatalf("error creating user: %v", err)
	}

	t.Cleanup(func() {
		userStore.DeleteUserByEmail(ctx, u.Email) //nolint:errcheck,gosec
	})
	return u
}

func createTestUser(t *testing.T, ctx context.Context, userStore auth.UserStore, org auth.Org, role auth.UserRole) *auth.User {
	user := auth.CreateUserRequest{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Org:       org,
		Role:      role,
		CreatedBy: gofakeit.Email(),
	}
	u, err := userStore.CreateUser(ctx, user)
	if err != nil {
		t.Fatalf("error creating user: %v", err)
	}

	t.Cleanup(func() {
		userStore.DeleteUserByID(ctx, u.ID) //nolint:errcheck,gosec
	})
	return u
}
