package analysis

import (
	"codemap/backend/internal/models"
	"encoding/json"
	"fmt"
	"os/exec"
)

// Run executes the Node.js analysis tool and returns the parsed data.
func Run(toolsPath string, targetDir string) (*models.Analysis, error) {
	// The command and its directory are now configured externally.
	cmd := exec.Command("node", "main.js", targetDir)
	cmd.Dir = toolsPath

	fmt.Printf("🔧 ANALYSIS: Running in: %s\n", toolsPath)
	fmt.Printf("🔧 ANALYSIS: Target directory: %s\n", targetDir)
	fmt.Printf("🔧 ANALYSIS: Command: %v\n", cmd.Args)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("❌ ANALYSIS FAILED: %v\n", err)
		fmt.Printf("❌ OUTPUT: %s\n", string(output))
		return nil, fmt.Errorf("failed to run analysis tool : %w\nOutput : %s", err, string(output))
	}

	fmt.Printf("✅ ANALYSIS OUTPUT LENGTH: %d bytes\n", len(output))

	var analysisResult models.Analysis
	err = json.Unmarshal(output, &analysisResult)
	if err != nil {
		fmt.Printf("❌ JSON UNMARSHAL FAILED: %v\n", err)
		// Show first 500 characters of output for debugging
		if len(output) > 500 {
			fmt.Printf("❌ RAW OUTPUT (first 500 chars): %s...\n", string(output[:500]))
		} else {
			fmt.Printf("❌ RAW OUTPUT: %s\n", string(output))
		}
		return nil, fmt.Errorf("failed to unmarshal analysis result: %w", err)
	}

	fmt.Printf("✅ ANALYSIS SUCCESS: Found %d files\n", len(analysisResult.Files))
	return &analysisResult, nil
}
