package database

import (
	"codemap/backend/internal/config"
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func NewDB(cfg *config.AppConfig) (*DB, error) {
	driver, err := neo4j.NewDriverWithContext(
		cfg.Neo4jURI,
		neo4j.BasicAuth(cfg.Neo4jUser, cfg.Neo4jPass, ""),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create neo4j driver: %w", err)
	}

	// Verify the connection to the database
	err = driver.VerifyConnectivity(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not verify neo4j connectivity: %w", err)
	}

	fmt.Println("Successfully connected to Neo4j.")
	return &DB{Driver: driver}, nil
}

// Close gracefully closes the database driver.
func (db *DB) Close(ctx context.Context) {
	db.Driver.Close(ctx)
}
