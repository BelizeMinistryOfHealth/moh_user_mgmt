package server

import (
	"bz.moh.epi/users/internal"
	"bz.moh.epi/users/internal/auth"
	"bz.moh.epi/users/internal/db"
	"bz.moh.epi/users/internal/server/handlers"
	"context"
	"errors"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// AppConf represents the configuration required to boot the server
type AppConf struct {
	ProjectID       string
	Port            int
	FirestoreApiKey string
}

// Deps represents the dependencies that the server requires
type Deps struct {
	router          *mux.Router
	firestoreClient *db.FirestoreClient
}

// RegisterHandlers creates a router with all its handlers
// It will instantiate all the db connections required for the handlers
// to work.
func RegisterHandlers(ctx context.Context, cnf AppConf) Deps {
	firebaseConfig := &firebase.Config{ProjectID: cnf.ProjectID}
	firestoreClient, err := db.NewFirestoreClient(ctx, firebaseConfig)
	if err != nil {
		log.Errorf("failed to create firestore client: %v", err)
		os.Exit(-1)
	}
	userStore, _ := auth.NewStore(firestoreClient, cnf.FirestoreApiKey)

	app := &internal.App{
		Firestore:       firestoreClient,
		UserStore:       &userStore,
		ProjectID:       cnf.ProjectID,
		FirestoreApiKey: cnf.FirestoreApiKey,
	}

	router, err := handlers.API(ctx, app)
	if err != nil {
		log.Errorf("failed to initiate handlers: %v", err)
		panic("Could not initiate handlers")
	}
	return Deps{
		router:          router,
		firestoreClient: firestoreClient,
	}
}

// NewServer starts a http server
func NewServer(ctx context.Context, deps Deps, port int) {
	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 60,
		Handler:      deps.router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v\n", err)
		}
	}()
	log.Println("Server Started")
	<-done
	log.Println("Server Stopped")
	wait := time.Duration(30)
	// Create a deadline to wait for.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), wait)
	defer func() {
		// Close other resources here
		if err := deps.firestoreClient.Close(); err != nil {
			log.Errorf("Failed to close firestore connection: %v", err)
		}
		cancel()
	}()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server Shotdown failed: %v", err)
	}
	log.Println("Server Exited Properly")
	os.Exit(0)
}
