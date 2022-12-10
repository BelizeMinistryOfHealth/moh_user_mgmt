package main

import (
	"bz.moh.epi/users/internal/server"
	"context"
	log "github.com/sirupsen/logrus"
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

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	appConf := server.AppConf{
		ProjectID:       projectID,
		Port:            port,
		FirestoreApiKey: apiKey,
	}

	ctx := context.Background()
	deps := server.RegisterHandlers(ctx, appConf)
	server.NewServer(ctx, deps, port)

}
