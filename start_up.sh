#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}PWA Launcher Installation${NC}"
echo "--------------------------------"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed!${NC}"
    echo "Please install Go with: sudo apt install golang"
    exit 1
fi

# Create directory for our app
echo -e "${BLUE}Creating application directory...${NC}"
mkdir -p ~/go-pwa-launcher

# Create Go source file
echo -e "${BLUE}Creating Go source file...${NC}"
cat > ~/go-pwa-launcher/launcher.go << 'EOF'
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
		fmt.Printf("Launching %s...\n", file)
		
		// Using xdg-open to launch the desktop file
		cmd := exec.Command("xdg-open", fullPath)
		err := cmd.Start()
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
EOF

# Compile the Go program
echo -e "${BLUE}Compiling Go application...${NC}"
cd ~/go-pwa-launcher
go build -o pwa-launcher launcher.go

# Check if compilation was successful
if [ ! -f "./pwa-launcher" ]; then
    echo -e "${RED}Error: Compilation failed!${NC}"
    exit 1
fi

echo -e "${GREEN}Compilation successful!${NC}"

# Create autostart entry
echo -e "${BLUE}Creating autostart entry...${NC}"
mkdir -p ~/.config/autostart
cat > ~/.config/autostart/go-pwa-launcher.desktop << EOF
[Desktop Entry]
Type=Application
Name=Go PWA Launcher
Exec=$HOME/go-pwa-launcher/pwa-launcher
Terminal=false
X-GNOME-Autostart-enabled=true
Comment=Launches PWA applications at login using Go
EOF

echo -e "${GREEN}Installation complete!${NC}"
echo -e "The PWA launcher will run automatically at login."
echo -e "You can test it manually by running: ${BLUE}~/go-pwa-launcher/pwa-launcher${NC}"
