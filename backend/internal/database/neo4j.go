package database

import (
	"codemap/backend/internal/helper"
	"codemap/backend/internal/models"
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Query executes a read-only Cypher query and returns the results as a slice of maps, which is ready to be converted to JSON.
func (db *DB) Query(ctx context.Context, cypher string, params map[string]any) ([]map[string]any, error) {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		res, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}

		records, err := res.Collect(ctx)
		if err != nil {
			return nil, err
		}

		var results []map[string]any
		for _, record := range records {
			results = append(results, record.AsMap())
		}

		return results, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed during query execution: %w", err)
	}

	return result.([]map[string]any), nil
}

// ImportAnalysis imports the entire analysis result into Neo4j within a single transaction.
func (db *DB) ImportAnalysis(ctx context.Context, analysisData *models.Analysis) error {
	session := db.Driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	// DEBUGGING: Print analysis data info
	fmt.Printf("🔍 IMPORTING ANALYSIS: %d files found\n", len(analysisData.Files))
	for i, file := range analysisData.Files {
		fmt.Printf("  File %d: %s (%s) - %d classes, %d functions\n",
			i+1, file.Path, file.Language, len(file.Classes), len(file.Functions))
		if file.Error != "" {
			fmt.Printf("    ERROR: %s\n", file.Error)
		}
	}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// First, clear the existing database to ensure a fresh import
		fmt.Println("🗑️  Clearing existing database...")
		if _, err := tx.Run(ctx, "MATCH (n) DETACH DELETE n", nil); err != nil {
			return nil, err
		}
		fmt.Println("✅ Database cleared successfully")

		// Create all nodes
		for _, file := range analysisData.Files {
			if err := helper.CreateNodesForFile(ctx, tx, file); err != nil {
				return nil, fmt.Errorf("failed to create nodes for file %s: %w", file.Path, err)
			}
		}

		// Create all relationships
		for _, file := range analysisData.Files {
			if err := helper.CreateRelationshipsForFile(ctx, tx, file); err != nil {
				return nil, fmt.Errorf("failed to create relationships for file %s: %w", file.Path, err)
			}
		}
		return nil, nil
	})

	if err != nil {
		return fmt.Errorf("failed to execute import transaction: %w", err)
	}

	fmt.Println("Successfully imported analysis into Neo4j.")
	return nil
}
