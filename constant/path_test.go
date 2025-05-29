package constant

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetBaseDir(t *testing.T) {
	// Save original baseDir to restore after tests
	originalBaseDir := baseDir

	t.Run("empty path creates default base directory", func(t *testing.T) {
		// Reset baseDir
		baseDir = ""

		// Create a temporary directory for testing
		tempDir := t.TempDir()
		testBaseDir := filepath.Join(tempDir, ".config", "bassit")

		// Mock the BaseDir value
		BaseDir = testBaseDir

		err := SetBaseDir("")
		assert.NoError(t, err)
		// Check if directory was created
		_, err = os.Stat(testBaseDir)
		assert.NoError(t, err, "Expected directory to be created at %s", testBaseDir)
	})

	t.Run("valid absolute path", func(t *testing.T) {
		tempDir := t.TempDir()

		err := SetBaseDir(tempDir)
		assert.NoError(t, err)
		assert.Equal(t, tempDir, baseDir)
	})

	t.Run("valid relative path", func(t *testing.T) {
		tempDir := t.TempDir()
		relativePath := "."

		// Change to temp directory
		oldWd, _ := os.Getwd()
		os.Chdir(tempDir)
		defer os.Chdir(oldWd)

		err := SetBaseDir(relativePath)
		assert.NoError(t, err)

		expectedPath, _ := filepath.Abs(relativePath)
		assert.Equal(t, expectedPath, baseDir)
	})

	t.Run("non-existent path", func(t *testing.T) {
		nonExistentPath := "/path/that/does/not/exist"

		err := SetBaseDir(nonExistentPath)
		assert.Error(t, err, "Expected error for non-existent path")
		if err != nil { // Ensure err is not nil before checking its message
			assert.True(t, strings.HasPrefix(err.Error(), "path does not exist"), "Expected path error for non-existent path, got: %v", err)
		}
	})

	// Restore original baseDir
	t.Cleanup(func() {
		baseDir = originalBaseDir
	})
}
