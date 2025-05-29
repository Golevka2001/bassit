package assets

import (
	"os"
	"path/filepath"
	"testing"

	C "github.com/Golevka2001/bassit/constant"

	"github.com/stretchr/testify/assert"
)

func TestExtractTo(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "assets_test")
	assert.NoError(t, err, "Failed to create temp dir")
	defer os.RemoveAll(tempDir)

	// Test extraction
	err = ExtractTo(tempDir)
	assert.NoError(t, err, "ExtractTo failed")

	// Verify config file exists
	configPath := filepath.Join(tempDir, "config.yaml")
	_, err = os.Stat(configPath)
	assert.False(t, os.IsNotExist(err), "Config file not extracted: %s", configPath)

	// Verify audio file exists
	audioPath := filepath.Join(tempDir, "audio/note/C2.wav")
	_, err = os.Stat(audioPath)
	assert.False(t, os.IsNotExist(err), "Audio file not extracted: %s", audioPath)

	// Verify NoteSampleDir is set correctly
	expectedNoteSampleDir := filepath.Join(tempDir, "audio/note/")
	assert.Equal(t, expectedNoteSampleDir, C.NoteSampleDir, "NoteSampleDir not set correctly")

	// Verify RubberbandCommand is set (OS-specific)
	rbDstPrefix := filepath.Join(tempDir, "3rdparty/rubberband/")
	switch C.OS {
	case "windows":
		expectedCmd := filepath.Join(rbDstPrefix, "rubberband-r3.exe")
		assert.Equal(t, expectedCmd, C.RubberbandCommand, "RubberbandCommand not set correctly for Windows")
	case "darwin":
		expectedCmd := filepath.Join(rbDstPrefix, "rubberband-r3")
		assert.Equal(t, expectedCmd, C.RubberbandCommand, "RubberbandCommand not set correctly for Darwin")
	}
}

func TestExtractFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "extract_file_test")
	assert.NoError(t, err, "Failed to create temp dir")
	defer os.RemoveAll(tempDir)

	// Test extracting config file
	dst := filepath.Join(tempDir, "test_config.yaml")
	err = extractFile(cfgFile, dst, 0644)
	assert.NoError(t, err, "extractFile failed")

	// Verify file exists and has correct permissions
	info, err := os.Stat(dst)
	assert.NoError(t, err, "File not created")

	assert.Equal(t, os.FileMode(0644), info.Mode().Perm(), "File permissions incorrect")

	// Test extracting file that already exists (should not return error)
	err = extractFile(cfgFile, dst, 0644)
	assert.NoError(t, err, "extractFile should not fail when file already exists")
}

func TestExtractFileNonExistentSource(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "extract_file_error_test")
	assert.NoError(t, err, "Failed to create temp dir")
	defer os.RemoveAll(tempDir)

	dst := filepath.Join(tempDir, "nonexistent.txt")
	err = extractFile("nonexistent/file.txt", dst, 0644)
	assert.Error(t, err, "extractFile should fail with nonexistent source file")
}
