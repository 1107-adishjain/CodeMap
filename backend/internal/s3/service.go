package s3

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Service struct {
	uploader *s3manager.Uploader
	bucket   string
}

func NewS3Service(region, accessKey, secretKey, bucket string) (*Service, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	uploader := s3manager.NewUploader(sess)

	return &Service{
		uploader: uploader,
		bucket:   bucket,
	}, nil
}

// UploadZipFile uploads a zip file to S3 and returns the S3 key
func (s *Service) UploadZipFile(zipData []byte, filename string) (string, error) {
	// Generate unique key with timestamp
	timestamp := time.Now().Format("20060102-150405")
	key := fmt.Sprintf("projects/%s-%s", timestamp, filename)

	_, err := s.uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(zipData),
		ContentType: aws.String("application/zip"),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload zip to S3: %w", err)
	}

	return key, nil
}

// UploadGitRepo clones a GitHub repo and uploads it as a zip to S3
// Returns: (s3Key, tempDirPath, error) - caller must clean up tempDirPath
func (s *Service) UploadGitRepo(repoURL string) (string, string, error) {
	// Create temp directory for cloning
	tempDir, err := os.MkdirTemp("", "git-clone-*")
	if err != nil {
		return "", "", fmt.Errorf("failed to create temp dir: %w", err)
	}
	// NOTE: No defer cleanup - caller must handle cleanup

	// Extract repo name from URL
	repoName := extractRepoName(repoURL)
	cloneDir := filepath.Join(tempDir, repoName)

	// Clone the repository
	cmd := exec.Command("git", "clone", repoURL, cloneDir)
	if err := cmd.Run(); err != nil {
		return "", "", fmt.Errorf("failed to clone repository: %w", err)
	}

	// Create zip file
	zipPath := filepath.Join(tempDir, repoName+".zip")
	if err := createZipFromDir(cloneDir, zipPath); err != nil {
		return "", "", fmt.Errorf("failed to create zip: %w", err)
	}

	// Read zip file
	zipData, err := os.ReadFile(zipPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to read zip file: %w", err)
	}

	// Upload to S3
	timestamp := time.Now().Format("20060102-150405")
	key := fmt.Sprintf("projects/%s-%s.zip", timestamp, repoName)

	_, err = s.uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(zipData),
		ContentType: aws.String("application/zip"),
	})

	if err != nil {
		return "", "", fmt.Errorf("failed to upload git repo to S3: %w", err)
	}

	// Return S3 key and clone directory path (caller must cleanup tempDir)
	return key, cloneDir, nil
}

// extractRepoName extracts repository name from GitHub URL
func extractRepoName(repoURL string) string {
	// Remove .git suffix if present
	repoURL = strings.TrimSuffix(repoURL, ".git")

	// Split by "/" and get the last part
	parts := strings.Split(repoURL, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return "unknown-repo"
}

// createZipFromDir creates a zip file from a directory
func createZipFromDir(sourceDir, zipPath string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip .git directory
		if strings.Contains(path, ".git") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		// Skip root directory
		if relPath == "." {
			return nil
		}

		// Create zip entry
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(relPath)

		if info.IsDir() {
			header.Name += "/"
			_, err := zipWriter.CreateHeader(header)
			return err
		}

		// Write file content
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})
}
