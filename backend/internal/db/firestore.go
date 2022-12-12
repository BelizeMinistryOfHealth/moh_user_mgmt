package db

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"fmt"
)

// FirestoreClient is a client for communicating with Firestore.
type FirestoreClient struct {
	Client      *firestore.Client
	AuthClient  *auth.Client
	AdminClient *firebase.App
	projectId   string
}

func (c *FirestoreClient) Close() error {
	if err := c.Client.Close(); err != nil {
		return fmt.Errorf("error closing firestore connection: %w", err)
	}
	return nil
}

// NewFirestoreClient creates a new Firestore connection.
func NewFirestoreClient(ctx context.Context, config *firebase.Config) (*FirestoreClient, error) {
	app, err := firebase.NewApp(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error creating new firebase app: %w", err)
	}
	authClient, _ := app.Auth(ctx)

	c, err := firestore.NewClient(ctx, config.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("error creating new firestore client; %w", err)
	}
	client := &FirestoreClient{
		Client:      c,
		AuthClient:  authClient,
		AdminClient: app,
		projectId:   config.ProjectID,
	}
	return client, nil
}
