package main

import (
	"archive/zip"
	"codemap/backend/internal/analysis"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// healthCheckHandler is a simple handler to confirm the API is running.
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "available",
		"environment": "development",
		"version":     "1.0.0",
	}
	app.writeJSON(w, http.StatusOK, data)
}

// uploadHandler handles the file upload and analysis process.
func (app *application) uploadHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(500 << 20); err != nil { // 500MB max
		app.errorResponse(w, r, http.StatusBadRequest, "Could not parse multipart form.")
		return
	}
	file, handler, err := r.FormFile("codebase")
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Could not retrieve the file from form.")
		return
	}
	defer file.Close()
	// Use configured path for temporary uploads
	tempDir, err := os.MkdirTemp(app.config.TempUploads, "codemap-upload-*")
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, "Could not create temp directory.")
		return
	}
	defer os.RemoveAll(tempDir)
	zipPath := filepath.Join(tempDir, handler.Filename)
	tempFile, err := os.Create(zipPath)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, "Could not create temp file.")
		return
	}
	// Stream file to local temp file first (better for large files)
	_, err = io.Copy(tempFile, file)
	tempFile.Close() // Close file before reading again
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, "Could not save uploaded file.")
		return
	}

	// Read file for S3 upload (now from disk, not memory)
	zipData, err := os.ReadFile(zipPath)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, "Could not read temp file.")
		return
	}

	// Upload to S3
	s3Key, err := app.s3.UploadZipFile(zipData, handler.Filename)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, fmt.Sprintf("Failed to upload to S3: %v", err))
		return
	}

	app.logger.Printf("📤 Uploaded %s to S3: %s (Size: %d bytes)", handler.Filename, s3Key, len(zipData))
	unzipDest := filepath.Join(tempDir, "unzipped")
	if err := unzip(zipPath, unzipDest); err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, fmt.Sprintf("Failed to unzip file: %v", err))
		return
	}
	// Pass the configured tools path and the unzipped directory to the runner
	analysisResult, err := analysis.Run(app.config.ToolsPath, unzipDest)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	// Import the result into Neo4j
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute) // 5-minute timeout for import
	defer cancel()
	if err := app.db.ImportAnalysis(ctx, analysisResult); err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, fmt.Sprintf("Failed to import data to Neo4j: %v", err))
		return
	}
	// Send a success response
	app.writeJSON(w, http.StatusAccepted, map[string]string{
		"message": "Upload successful. Codebase has been analyzed and imported.",
		"s3_key":  s3Key,
	})
}

// githubHandler handles GitHub repository URL submission for analysis.
func (app *application) githubHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		RepoURL string `json:"repo_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if payload.RepoURL == "" {
		app.errorResponse(w, r, http.StatusBadRequest, "Repository URL is required")
		return
	}

	// Validate GitHub URL format
	if !strings.Contains(payload.RepoURL, "github.com") {
		app.errorResponse(w, r, http.StatusBadRequest, "Please provide a valid GitHub repository URL")
		return
	}

	// Upload GitHub repo to S3 and get the local temp directory path
	s3Key, tempDir, err := app.s3.UploadGitRepo(payload.RepoURL)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, fmt.Sprintf("Failed to clone and upload repository: %v", err))
		return
	}

	// Clean up temp directory after we're done
	defer os.RemoveAll(tempDir)

	app.logger.Printf("📤 Uploaded GitHub repo to S3: %s", s3Key)

	// Run analysis on the cloned repository
	analysisResult, err := analysis.Run(app.config.ToolsPath, tempDir)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, fmt.Sprintf("Analysis failed: %v", err))
		return
	}

	// Import analysis results into Neo4j
	if err := app.db.ImportAnalysis(r.Context(), analysisResult); err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, fmt.Sprintf("Failed to import analysis to Neo4j: %v", err))
		return
	}

	app.writeJSON(w, http.StatusAccepted, map[string]string{
		"message":        "GitHub repository analyzed and imported to Neo4j successfully.",
		"s3_key":         s3Key,
		"repo_url":       payload.RepoURL,
		"files_analyzed": fmt.Sprintf("%d", len(analysisResult.Files)),
	})
}

// queryHandler accepts a POST request with a Cypher query and returns the result.
func (app *application) queryHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Query  string         `json:"query"`
		Params map[string]any `json:"params"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// The handler now calls the database method, keeping the logic separate.
	results, err := app.db.Query(ctx, payload.Query, payload.Params)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, fmt.Sprintf("Failed to execute query: %v", err))
		return
	}

	// Send the results back as JSON.
	app.writeJSON(w, http.StatusOK, results)
}

// --- HELPER METHODS ---

// writeJSON is a helper for sending JSON responses.
func (app *application) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		app.logError(nil, err)
	}
}

// errorResponse is a helper for sending consistent error messages.
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message string) {
	errData := map[string]string{"error": message}
	app.writeJSON(w, status, errData)
}

// logError is a helper for logging errors.
func (app *application) logError(r *http.Request, err error) {
	app.logger.Println(err)
}

// unzip function to extract a zip archive.
func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()
	os.MkdirAll(dest, 0755)
	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
