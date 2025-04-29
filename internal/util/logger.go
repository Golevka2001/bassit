// bassit/internal/util/logger.go
package util

import (
    "fmt"
    "log"
    "os"
)

var logger *log.Logger

func init() {
    // Initialize the logger
    logger = log.New(os.Stderr, "[BASSIT] ", log.LstdFlags)
}

// LogError logs an error message
func LogError(message string, err error) {
    logger.Printf("ERROR: %s: %v\n", message, err)
}

// LogInfo logs an informational message
func LogInfo(message string) {
    logger.Printf("INFO: %s\n", message)
}

// LogDebug logs a debug message
func LogDebug(message string) {
    if os.Getenv("BASS_DEBUG") == "1" {
        logger.Printf("DEBUG: %s\n", message)
    }
}

// Confirm asks for user confirmation
func Confirm(prompt string) bool {
    fmt.Printf("%s [y/N]: ", prompt)
    var response string
    fmt.Scanln(&response)
    return response == "y" || response == "Y"
}
