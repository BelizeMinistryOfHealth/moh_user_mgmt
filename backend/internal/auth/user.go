package auth

import (
	"bytes"
	"bz.moh.epi/users/internal/db"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"errors"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
	"net/http"
	"os"
)

const (
	passwordReset   = "PASSWORD_RESET"
	applicationName = "hiv_surveys"
)

// User represents a user account in the system
type User struct {
	ID               string            `json:"id" firestore:"id"`
	FirstName        string            `json:"firstName" firestore:"firstName"`
	LastName         string            `json:"lastName" firestore:"lastName"`
	Email            string            `json:"email" firestore:"email"`
	UserApplications []UserApplication `json:"userApplications" firestore:"userApplications"`
}

type CreateUserRequest struct {
	FirstName        string            `json:"firstName" firestore:"firstName"`
	LastName         string            `json:"lastName" firestore:"lastName"`
	Email            string            `json:"email" firestore:"email"`
	UserApplications []UserApplication `json:"userApplications" firestore:"userApplications"`
}

// IsAdmin indicates if a given slice of permissions has admin access
func IsAdmin(permissions []string) bool {
	isAdmin := false
	for idx := range permissions {
		if permissions[idx] == "admin" {
			isAdmin = true
		}
	}
	return isAdmin
}

// UserService describes the contract that any service that will be used for persisting and querying users
// should satisfy.
type UserService interface {
	CreateUser(ctx context.Context, user User) (*User, error)
}

// UserStore represents a store used for persisting and retrieving users from a data store.
type UserStore struct {
	db          *db.FirestoreClient
	collection  string
	adminClient *firebase.App
	authClient  *auth.Client
	AuthClient  *auth.Client
	apiKey      string
}

// NewStore creates a new store that provides ways for creating and mutating a user.
func NewStore(db *db.FirestoreClient, apiKey string) (UserStore, error) {
	return UserStore{
		db:          db,
		collection:  "epi_users",
		authClient:  db.AuthClient,
		adminClient: db.AdminClient,
		apiKey:      apiKey,
	}, nil
}

// CreateUser creates an auth user and also creates a record in the user collection.
// It will send a password reset email so the user, with a link that allows them to set a password
// for their account.
func (s *UserStore) CreateUser(ctx context.Context, user CreateUserRequest) (*User, error) {
	u := (&auth.UserToCreate{}).Email(user.Email).DisplayName(fmt.Sprintf("%s %s", user.FirstName, user.LastName)).Disabled(false)
	userRecord, err := s.authClient.CreateUser(ctx, u)
	if err != nil {
		return nil, AuthError{
			Reason: "failed creating auth user",
			Inner:  err,
		}
	}

	userPayload := &User{
		ID:               userRecord.UID,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		Email:            user.Email,
		UserApplications: user.UserApplications,
	}

	//Create userRecord in firestore
	_, err = s.db.Client.Collection(s.collection).Doc(userRecord.UID).Set(ctx, userPayload)
	if err != nil {
		return nil, AuthError{
			Reason: "failed inserting user to collection",
			Inner:  err,
		}
	}
	if err := s.SendPasswordResetEmail(user.Email); err != nil {
		return nil, AuthError{
			Reason: "failed sending password reset email",
			Inner:  err,
		}
	}
	return userPayload, nil
}

// SendPasswordResetEmail sends an email to the user with a link to reset their password.
func (s *UserStore) SendPasswordResetEmail(email string) error {
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:sendOobCode?key=%s", s.apiKey)
	emulatorHost := os.Getenv("FIRESTORE_AUTH_EMULATOR_HOST")
	if len(emulatorHost) > 0 {
		url = fmt.Sprintf("https://%s/identitytoolkit.googleapis.com/v1/accounts:sendOobCode?key=%s", emulatorHost, s.apiKey)
	}
	reqBody, _ := json.Marshal(map[string]string{
		"requestType": passwordReset,
		"email":       email,
	})
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody)) //nolint: gosec
	if err != nil {
		return AuthError{
			Reason: "failed posting PASSWORD_RESET",
			Inner:  err,
		}
	}
	defer resp.Body.Close() //nolint:errcheck
	var body map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return AuthError{
			Reason: "failed decoding PASSWORD_RESET response",
			Inner:  err,
		}
	}
	log.WithFields(log.Fields{
		"body": body,
	}).Info("result")
	return nil
}

// DeleteUser deletes a user from the firestore collection and also from Firestore Auth.
func (s *UserStore) DeleteUser(ctx context.Context, user User) error {
	_, err := s.authClient.DeleteUsers(ctx, []string{user.ID})
	if err != nil {
		return AuthError{
			Reason: "failed deleting auth user",
			Inner:  err,
		}
	}
	if _, err := s.db.Client.Collection(s.collection).Doc(user.ID).Delete(ctx); err != nil {
		return AuthError{
			Reason: "failed deleting user from collection",
			Inner:  err,
		}
	}
	return nil
}

// GetUserByEmail gets a user's record from firebase.
func (s *UserStore) GetUserByEmail(ctx context.Context, email string) (User, error) {
	iter := s.db.Client.Collection(s.collection).Where("email", "==", email).Limit(1).Documents(ctx)
	var user User
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return User{}, AuthError{
				Reason: "failed to fetch user by email",
				Inner:  err,
			}
		}
		if err := doc.DataTo(&user); err != nil {
			return User{}, AuthError{
				Reason: "failed to convert user data",
				Inner:  err,
			}
		}
	}
	u, _ := s.authClient.GetUser(ctx, user.ID)
	log.WithFields(log.Fields{
		"user": *u,
	}).Info("user from auth client")
	return user, nil
}

// GetUserByID retrieves a user with the matching id. It returns nil if no user is found.
func (s *UserStore) GetUserByID(ctx context.Context, ID string) (*User, error) {
	iter := s.db.Client.Collection(s.collection).Where("id", "==", ID).Limit(1).Documents(ctx)
	var user User
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, AuthError{
				Reason: fmt.Sprintf("failed to retrieve user with ID %s", ID),
				Inner:  err,
			}
		}
		if err := doc.DataTo(&user); err != nil {
			return nil, AuthError{
				Reason: fmt.Sprintf("error converting result to User struct: %v", doc),
				Inner:  err,
			}
		}
	}
	return &user, nil
}

// UpdateUser updates a user's permissions.
func (s *UserStore) UpdateUser(ctx context.Context, user *User) error {
	if _, err := s.db.Client.Collection(s.collection).Doc(user.ID).Update(ctx, []firestore.Update{
		{
			Path:  "userApplications",
			Value: user.UserApplications,
		},
		{
			Path:  "firstName",
			Value: user.FirstName,
		},
		{
			Path:  "lastName",
			Value: user.LastName,
		},
	}); err != nil {
		return AuthError{
			Reason: "failed updating user collection",
			Inner:  err,
		}
	}
	_, err := s.authClient.UpdateUser(ctx, user.ID, (&auth.UserToUpdate{}).DisplayName(fmt.Sprintf("%s %s", user.FirstName, user.LastName)))
	if err != nil {
		return AuthError{
			Reason: "failed to update user's name",
			Inner:  err,
		}
	}
	return nil
}

// ListUsers will return a list of all users stored in the data store.
func (s *UserStore) ListUsers(ctx context.Context) ([]User, error) {
	iter := s.db.Client.Collection(s.collection).Documents(ctx)
	var users []User
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, UserError{
				Reason: "failed iterating over user collection",
				Inner:  err,
			}
		}
		var u User
		if err := doc.DataTo(&u); err != nil {
			return nil, UserError{
				Reason: "failed to transform user data",
				Inner:  err,
			}
		}
		users = append(users, u)
	}
	return users, nil
}

// VerifyToken will verify that a token is valid. It will return an error if it is not valid.
// It will return a JwtToken if it is valid. This will hold information that can be used to quickly
// check user information without needing to query the database again for it.
func (s *UserStore) VerifyToken(ctx context.Context, t string) (JwtToken, error) {
	token, err := s.authClient.VerifyIDToken(ctx, t)
	if err != nil {
		return JwtToken{}, AuthError{
			Reason: "failed to verify ID Token",
			Inner:  err,
		}
	}

	claims := token.Claims
	email := claims["email"]
	user, err := s.GetUserByEmail(ctx, email.(string))
	if err != nil {
		return JwtToken{
			Email:       email.(string),
			Permissions: nil,
		}, AuthError{Inner: err, Reason: "could not find user record with provided email"}
	}
	applications := user.UserApplications
	if applications == nil {
		return JwtToken{
			Email:       email.(string),
			Admin:       false,
			Permissions: nil,
		}, nil
	}
	applicationPermissions := findApplication(user.UserApplications, applicationName)
	var permissions []string
	if applicationPermissions != nil && len(applicationPermissions.Permissions) > 0 {
		permissions = append(permissions, applicationPermissions.Permissions...)
	}
	return JwtToken{
		Email:       email.(string),
		Admin:       IsAdmin(permissions),
		Permissions: permissions,
	}, nil
}

func findApplication(apps []UserApplication, appName string) *UserApplication {
	for i := 0; i < len(apps); i++ {
		if apps[i].Name == appName {
			return &apps[i]
		}
	}
	return nil
}

// CreateToken creates a token for a user that matches the ID provided.
func (s *UserStore) CreateToken(ctx context.Context, ID string) (string, error) {
	token, err := s.authClient.CustomToken(ctx, ID)
	if err != nil {
		return "", fmt.Errorf("error creating custom token: %w", err)
	}

	return token, nil
}

// JwtToken represents a JWT token that has been verified
type JwtToken struct {
	Email       string
	Admin       bool
	Permissions []string
}

// DeleteUserByID deletes a user from the data store
func (s *UserStore) DeleteUserByID(ctx context.Context, id string) error {
	_, err := s.authClient.DeleteUsers(ctx, []string{id})
	if err != nil {
		return AuthError{
			Reason: "failed deleting auth user",
			Inner:  err,
		}
	}
	if _, err := s.db.Client.Collection(s.collection).Doc(id).Delete(ctx); err != nil {
		return AuthError{
			Reason: "failed deleting user from collection",
			Inner:  err,
		}
	}
	return nil
}

// DeleteUserByEmail deletes a user with the matching email. Useful primarily for tests where
// we do not know the user's ID.
func (s *UserStore) DeleteUserByEmail(ctx context.Context, email string) error {
	user, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("error fetching user by email when deleting user by email%w", err)
	}
	err = s.authClient.DeleteUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("DeleteUserByEmail: error deleting user with email %s: %w", email, err)
	}
	if _, err := s.db.Client.Collection(s.collection).Doc(user.ID).Delete(ctx); err != nil {
		return fmt.Errorf("DeleteUserByEmail: error deleting user with email %s from collection: %w", user.Email, err)
	}
	return nil
}
