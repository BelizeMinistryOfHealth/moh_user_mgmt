package auth

import (
	"bytes"
	"bz.moh.epi/users/internal/db"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
	"net/http"
	"os"
)

type linkType string

const (
	passwordReset   linkType = "PASSWORD_RESET"
	applicationName          = "hiv_surveys"
)

type Application struct {
	ID          string   `json:"id" firestore:"id"`
	Name        string   `json:"name" firestore:"name"`
	Permissions []string `json:"permissions" firestore:"permissions"`
}

type UserApplication struct {
	ApplicationID string   `json:"id" firestore:"id"`
	Name          string   `json:"name" firestore:"name"`
	Permissions   []string `json:"permissions" firestore:"permissions"`
}

type User struct {
	ID              string          `json:"id" firestore:"id"`
	FirstName       string          `json:"firstName" firestore:"firstName"`
	LastName        string          `json:"lastName" firestore:"lastName"`
	Email           string          `json:"email" firestore:"email"`
	UserApplication UserApplication `json:"userApplication" firestore:"userApplication"`
}

func IsAdmin(permissions []string) bool {
	isAdmin := false
	for idx := range permissions {
		if permissions[idx] == "admin:write" {
			isAdmin = true
		}
	}
	return isAdmin
}

type UserStore struct {
	db          *db.FirestoreClient
	collection  string
	adminClient *firebase.App
	authClient  *auth.Client
	apiKey      string
}

// NewStore creates a new store that provides ways for creating and mutating a user.
func NewStore(ctx context.Context, db *db.FirestoreClient, config *firebase.Config, apiKey string) (UserStore, error) {
	app, err := firebase.NewApp(ctx, config)
	if err != nil {
		return UserStore{}, err
	}
	authClient, err := app.Auth(ctx)
	if err != nil {
		return UserStore{}, err
	}

	return UserStore{
		db:          db,
		collection:  "epi_users",
		authClient:  authClient,
		adminClient: app,
		apiKey:      apiKey,
	}, nil
}

// CreateUser creates an auth user and also creates a record in the user collection.
// It will send a password reset email so the user, with a link that allows them to set a password
// for their account.
func (s *UserStore) CreateUser(ctx context.Context, user User) (*User, error) {
	u := (&auth.UserToCreate{}).Email(user.Email).DisplayName(fmt.Sprintf("%s %s", user.FirstName, user.LastName)).Disabled(false)
	userRecord, err := s.authClient.CreateUser(ctx, u)
	if err != nil {
		return nil, AuthError{
			Reason: "failed creating auth user",
			Inner:  err,
		}
	}

	userRecord.CustomClaims = map[string]interface{}{"applications": map[string]interface{}{user.UserApplication.Name: user.UserApplication.Permissions}}
	ur, err := s.authClient.UpdateUser(ctx, userRecord.UID, (&auth.UserToUpdate{}).CustomClaims(map[string]interface{}{"applications": userRecord.CustomClaims["applications"]}))
	if err != nil {
		return nil, AuthError{
			Reason: "failed updating auth user claims",
			Inner:  err,
		}
	}
	log.WithFields(log.Fields{
		"userRecord": *ur,
	}).Info("updated user with claims")
	//Create userRecord in firestore
	user.ID = userRecord.UID
	_, err = s.db.Client.Collection(s.collection).Doc(user.ID).Set(ctx, user)
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
	return &user, nil
}

// SendPasswordResetEmail sends an email to the user with a link to reset their password.
func (s *UserStore) SendPasswordResetEmail(email string) error {
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:sendOobCode?key=%s", s.apiKey)
	emulatorHost := os.Getenv("FIRESTORE_AUTH_EMULATOR_HOST")
	if len(emulatorHost) > 0 {
		url = fmt.Sprintf("https://%s/identitytoolkit.googleapis.com/v1/accounts:sendOobCode?key=%s", emulatorHost, s.apiKey)
	}
	reqBody, _ := json.Marshal(map[string]string{
		"requestType": string(passwordReset),
		"email":       email,
	})
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return AuthError{
			Reason: "failed posting PASSWORD_RESET",
			Inner:  err,
		}
	}
	defer resp.Body.Close()
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
		if err == iterator.Done {
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

// UpdateUser updates a user's permissions.
func (s *UserStore) UpdateUser(ctx context.Context, user User) error {
	userToUpdate := (&auth.UserToUpdate{}).
		CustomClaims(map[string]interface{}{"permissions": user.UserApplication.Permissions})
	if _, err := s.authClient.UpdateUser(ctx, user.ID, userToUpdate); err != nil {
		return AuthError{
			Reason: "failed to update auth user",
			Inner:  err,
		}
	}
	if _, err := s.db.Client.Collection(s.collection).Doc(user.ID).Update(ctx, []firestore.Update{
		{
			Path:  "permissions",
			Value: user.UserApplication.Permissions,
		},
	}); err != nil {
		return AuthError{
			Reason: "failed updating user collection",
			Inner:  err,
		}
	}
	return nil
}

func (s *UserStore) ListUsers(ctx context.Context) ([]User, error) {
	iter := s.db.Client.Collection(s.collection).Documents(ctx)
	var users []User
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
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
	_, err = s.GetUserByEmail(ctx, email.(string))
	if err != nil {
		return JwtToken{
			Email:       email.(string),
			Permissions: nil,
		}, AuthError{Inner: err, Reason: "could not find user record with provided email"}
	}
	applications := claims["applications"]
	applicationPermissions := applications.(map[string]interface{})[applicationName]
	var permissions []string
	if applicationPermissions != nil && len(applicationPermissions.([]string)) > 0 {
		for idx := range applicationPermissions.([]interface{}) {
			permissions = append(permissions, applicationPermissions.([]interface{})[idx].(string))
		}
	}
	return JwtToken{
		Email:       email.(string),
		Permissions: permissions,
	}, nil
}

type JwtToken struct {
	Email       string
	Permissions []string
}

func (s *UserStore) DeleteUserById(ctx context.Context, id string) error {
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
