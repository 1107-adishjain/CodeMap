package main

import (
	"codemap/backend/internal/config"
	"codemap/backend/internal/database"
	"codemap/backend/internal/s3"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

type application struct {
	config *config.AppConfig
	db     *database.DB
	logger *log.Logger
	s3     *s3.Service
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables.")
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	cfg := config.Load()

	dbNeo4j, err := database.NewDB(cfg)
	if err != nil {
		logger.Fatalf("Could not connect to database: %v", err)
	} else {
		logger.Println("connected to database")
	}
	defer dbNeo4j.Close(context.Background())

	db, err := database.DBinit(cfg.PostgresUrl)
	if err !=nil {
		logger.Fatalf("Could not connect to Postgres database: %v", err)
	}else{
		logger.Println("connected to Postgres database")
	}
	defer database.DBclose(db)

	s3Service, err := s3.NewS3Service(cfg.S3Region, cfg.AWSAccessKey, cfg.AWSSecretKey, cfg.S3Bucket)
	if err != nil {
		logger.Fatalf("Could not initialize S3 service: %v", err)
	} else {
		logger.Println("S3 succesfully initalized")
	}

	app := &application{
		config: cfg,
		db:     dbNeo4j,
		logger: logger,
		s3:     s3Service,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      app.routes(db),
		ErrorLog:     log.New(logger.Writer(), "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		logger.Printf("Caught signal: %v. Shutting down server...", s)

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	logger.Printf("Starting server on %s", srv.Addr)

	err = srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		logger.Fatalf("Server error: %v", err)
	}

	err = <-shutdownError
	if err != nil {
		logger.Fatalf("Error during shutdown: %v", err)
	}

	logger.Println("Server stopped gracefully.")
}
