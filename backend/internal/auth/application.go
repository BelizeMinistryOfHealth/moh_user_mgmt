package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

const (
	applicationCollection = "applications"
)

// UserApplication describes what applications a user has access to.
type UserApplication struct {
	ApplicationID string   `json:"id" firestore:"id"`
	Name          string   `json:"name" firestore:"name"`
	Permissions   []string `json:"permissions" firestore:"permissions"`
}

type CreateApplicationRequest struct {
	Name        string   `json:"name" firestore:"name"`
	Permissions []string `json:"permissions" firestore:"permissions"`
}

func (s *UserStore) CreateApplication(ctx context.Context, req CreateApplicationRequest) (*UserApplication, error) {
	ID := uuid.New().String()
	userApplication := UserApplication{
		ApplicationID: ID,
		Name:          req.Name,
		Permissions:   req.Permissions,
	}
	_, err := s.db.Client.Collection(applicationCollection).Doc(ID).Create(ctx, userApplication)
	if err != nil {
		return nil, fmt.Errorf("error persisting new user application: %w", err)
	}
	return &userApplication, nil
}

// ListApplications returns all the user applications in the database.
func (s *UserStore) ListApplications(ctx context.Context) ([]UserApplication, error) {
	iter := s.db.Client.Collection(applicationCollection).Documents(ctx)
	var apps []UserApplication
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return apps, fmt.Errorf("error retrieving user applications: %w", err)
		}
		var app UserApplication
		if err := doc.DataTo(&app); err != nil {
			return apps, fmt.Errorf("error translating raw user application to go: %w", err)
		}
		apps = append(apps, app)
	}
	return apps, nil
}

func (s *UserStore) DeleteApplicationByID(ctx context.Context, ID string) error {
	if _, err := s.db.Client.Collection(applicationCollection).Doc(ID).Delete(ctx); err != nil {
		return fmt.Errorf("error deleting application with id %s : %w ", ID, err)
	}
	return nil
}
