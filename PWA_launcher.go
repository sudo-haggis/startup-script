package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// Wait a few seconds for the desktop environment to fully load
	fmt.Println("Starting PWA launcher...")
	time.Sleep(5 * time.Second)

	// Desktop files to launch
	pwaFiles := []string{
		"Claude.desktop",
		"github.desktop",
		"notion.desktop",
		"portainer.desktop",
		"whatsapp-web.desktop",
	}

	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		return
	}
	
	// Desktop directory path
	desktopDir := filepath.Join(homeDir, "Desktop")
	
	// Launch each PWA
	for _, file := range pwaFiles {
		fullPath := filepath.Join(desktopDir, file)
		fmt.Printf("Processing %s...\n", file)
		
		// Extract and execute the command from the desktop file
		execCommand, err := extractExecFromDesktopFile(fullPath)
		if err != nil {
			fmt.Printf("Error extracting command from %s: %v\n", file, err)
			continue
		}
		
		fmt.Printf("Found command: %s\n", execCommand)
		
		// Parse the command and arguments
		cmdParts := parseCommand(execCommand)
		if len(cmdParts) == 0 {
			fmt.Printf("Error parsing command from %s\n", file)
			continue
		}
		
		// Execute the command
		cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
		err = cmd.Start()
		if err != nil {
			fmt.Printf("Error launching %s: %v\n", file, err)
		} else {
			fmt.Printf("Successfully launched %s\n", file)
		}
		
		// Small delay between launches to prevent overwhelming the system
		time.Sleep(1 * time.Second)
	}
	
	fmt.Println("All PWAs launched!")
}

// extractExecFromDesktopFile reads a desktop file and extracts the Exec= line
func extractExecFromDesktopFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Exec=") {
			// Remove the Exec= prefix
			execLine := strings.TrimPrefix(line, "Exec=")
			
			// Remove any field codes like %f, %u, etc.
			execLine = removeFieldCodes(execLine)
			
			return execLine, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("no Exec line found in desktop file")
}

// removeFieldCodes removes desktop entry field codes like %f, %u, etc.
func removeFieldCodes(cmd string) string {
	// List of field codes to remove
	fieldCodes := []string{"%f", "%F", "%u", "%U", "%d", "%D", "%n", "%N", "%i", "%c", "%k", "%v", "%m"}
	
	result := cmd
	for _, code := range fieldCodes {
		result = strings.Replace(result, code, "", -1)
	}
	
	// Trim any extra spaces
	return strings.TrimSpace(result)
}

// parseCommand splits a command string into command and arguments
func parseCommand(cmdStr string) []string {
	var result []string
	
	// Simple parsing for quoted arguments
	inQuote := false
	current := ""
	for _, char := range cmdStr {
		if char == '"' || char == '\'' {
			inQuote = !inQuote
			continue
		}
		
		if char == ' ' && !inQuote {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	
	if current != "" {
		result = append(result, current)
	}
	
	return result
}
