package database

import (
	"github.com/google/uuid"
	"time"
)

func (db *DB) CreateProject(userID, name, s3Key string) (string, error) {
	projectID := uuid.New().String()
	_, err := db.SQL.Exec(
		"INSERT INTO projects (id, user_id, name, s3_key, status, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
		projectID, userID, name, s3Key, "pending", time.Now(),
	)
	return projectID, err
}

func (db *DB) UpdateProjectStatus(projectID, status string) error {
	_, err := db.SQL.Exec("UPDATE projects SET status = $1 WHERE id = $2", status, projectID)
	return err
}