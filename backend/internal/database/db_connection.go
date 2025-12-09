package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// DB struct holds SQL and other DB handles
type DB struct {
	SQL    *sql.DB
	Driver neo4j.DriverWithContext
}

func DBinit(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func DBclose(db *sql.DB) error {
	return db.Close()
}
