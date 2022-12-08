package main

import (
	"bz.moh.epi/users/internal/server"
	"os"
	"strconv"
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
	p := os.Getenv("PORT")
	port := 8080
	if p != "" {
		pint, err := strconv.Atoi(p)
		if err == nil {
			port = pint
		}
	}

	appConf := server.AppConf{
		ProjectID:       projectID,
		Port:            port,
		FirestoreApiKey: apiKey,
	}

	server.NewServer(appConf)

}
