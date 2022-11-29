package db

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

// FirestoreClient is a client for communicating with Firestore.
type FirestoreClient struct {
	Client      *firestore.Client
	AuthClient  *auth.Client
	AdminClient *firebase.App
	projectId   string
}

func (c *FirestoreClient) Close() error {
	return c.Client.Close()
}

// NewFirestoreClient creates a new Firestore connection.
func NewFirestoreClient(ctx context.Context, config *firebase.Config) (*FirestoreClient, error) {
	app, err := firebase.NewApp(ctx, config)
	if err != nil {
		return nil, err
	}
	authClient, _ := app.Auth(ctx)

	c, err := firestore.NewClient(ctx, config.ProjectID)
	if err != nil {
		return nil, err
	}
	client := &FirestoreClient{
		Client:      c,
		AuthClient:  authClient,
		AdminClient: app,
		projectId:   config.ProjectID,
	}
	return client, nil
}
