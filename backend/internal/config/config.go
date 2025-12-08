package config

import (
	"os"
)

type AppConfig struct {
	Port         string
	Neo4jURI     string
	Neo4jUser    string
	Neo4jPass    string
	ToolsPath    string
	TempUploads  string
	S3Bucket     string
	S3Region     string
	AWSAccessKey string
	AWSSecretKey string
	PostgresUrl  string
}

// getEnv reads an environment variable or returns a default value.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Load loads configuration from environment variables or uses defaults.
func Load() *AppConfig {
	return &AppConfig{
		Port:         getEnv("PORT", "8080"),
		Neo4jURI:     getEnv("NEO4J_URI", "path"),
		Neo4jUser:    getEnv("NEO4J_USERNAME", "neo4j"),
		Neo4jPass:    getEnv("NEO4J_PASSWORD", "your_neo4j_password"),
		ToolsPath:    getEnv("TOOLS_PATH", "../tools"),
		TempUploads:  getEnv("TEMP_UPLOADS", os.TempDir()),
		S3Bucket:     getEnv("S3_BUCKET", "your-bucket-name"),
		S3Region:     getEnv("S3_REGION", "your-region"),
		AWSAccessKey: getEnv("AWS_ACCESS_KEY", ""),
		AWSSecretKey: getEnv("AWS_SECRET_KEY", ""),
		PostgresUrl:  getEnv("POSTGRES_URL", ""),
	}
}
