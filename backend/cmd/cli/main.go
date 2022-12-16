package main

import (
	"bz.moh.epi/users/internal/auth"
	"bz.moh.epi/users/internal/db"
	"context"
	firebase "firebase.google.com/go/v4"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	projectID := os.Getenv("PROJECT_ID")
	apiKey := os.Getenv("API_KEY")

	if projectID == "" {
		panic("Provide a PROJECT_ID environment variable")
	}
	if apiKey == "" {
		panic("Provide an API_KEY environment variable")
	}

	firstName := os.Getenv("FIRST_NAME")
	if firstName == "" {
		panic("Provide a FIRST_NAME environment variable")
	}
	lastName := os.Getenv("LAST_NAME")
	if lastName == "" {
		panic("Provide a FIRST_NAME environment variable")
	}
	email := os.Getenv("EMAIL")
	if email == "" {
		panic("Provide an EMAIL environment variable")
	}
	ctx := context.Background()
	firebaseConfig := &firebase.Config{ProjectID: projectID}
	firestoreClient, err := db.NewFirestoreClient(ctx, firebaseConfig)
	if err != nil {
		log.Errorf("failed to create firestore client: %v", err)
		os.Exit(-1)
	}
	userStore, _ := auth.NewStore(firestoreClient, apiKey)
	userRequest := auth.CreateUserRequest{
		FirstName:        firstName,
		LastName:         lastName,
		Email:            email,
		UserApplications: []auth.UserApplication{{ApplicationID: "1", Name: "hiv_survey", Permissions: []string{"admin"}}},
	}
	_, err = userStore.CreateUser(ctx, userRequest)
	if err != nil {
		log.Errorf("Error creating user: %v", err)
	}
	os.Exit(1)
}
