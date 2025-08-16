#!/bin/bash

# LeetSolv Installation Script
# This script installs LeetSolv CLI application on Linux and macOS

set -e
# Use set -u to exit if an unset variable is used.
set -u
# Use set -o pipefail to exit if a command in a pipe fails.
set -o pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO_OWNER="eannchen"
REPO_NAME="leetsolv"
BINARY_NAME="leetsolv"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="$HOME/.leetsolv"
BACKUP_DIR="$CONFIG_DIR/backup"

# Temp directory variable that will be cleaned up on exit.
TEMP_DIR=""

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
    # Exit with a non-zero status code on error.
    exit 1
}

# Cleanup function to be called on script exit.
cleanup() {
    if [ -n "$TEMP_DIR" ] && [ -d "$TEMP_DIR" ]; then
        rm -rf "$TEMP_DIR"
    fi
}

# Trap the EXIT signal to run the cleanup function.
trap cleanup EXIT

# Function to detect OS and architecture
detect_platform() {
    OS=""
    ARCH=""

    case "$(uname -s)" in
        Linux*)     OS="linux";;
        Darwin*)    OS="darwin";;
        *)          print_error "Unsupported operating system. This script only supports Linux and macOS.";;
    esac

    case "$(uname -m)" in
        x86_64)    ARCH="amd64";;
        aarch64)   ARCH="arm64";;
        arm64)     ARCH="arm64";;
        *)         print_error "Unsupported architecture: $(uname -m). This script only supports x86_64 (amd64) and aarch64/arm64.";;
    esac

    print_status "Detected platform: $OS-$ARCH"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to backup existing installation
backup_existing() {
    if command_exists "$BINARY_NAME"; then
        print_status "Backing up existing installation..."
        mkdir -p "$BACKUP_DIR"
        local backup_path="$BACKUP_DIR/$(date +%Y%m%d_%H%M%S)_$BINARY_NAME"
        # Use command -v instead of which for better portability.
        cp "$(command -v "$BINARY_NAME")" "$backup_path"
        print_success "Backup created at: $backup_path"
    fi
}

# Function to download latest release
download_release() {
    print_status "Checking for latest release..."

    local api_url="https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/latest"
    local latest_release=""

    # Prioritize jq for robust JSON parsing.
    if command_exists "jq"; then
        if command_exists "curl"; then
            latest_release=$(curl -s "$api_url" | jq -r .tag_name)
        elif command_exists "wget"; then
            latest_release=$(wget -qO- "$api_url" | jq -r .tag_name)
        fi
    else
        print_warning "jq not found. Falling back to grep/sed for version detection."
        # Fallback to grep/sed if jq is not available
        if command_exists "curl"; then
            latest_release=$(curl -s "$api_url" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
        elif command_exists "wget"; then
            latest_release=$(wget -qO- "$api_url" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
        else
            print_error "Neither curl nor wget found. Please install one of them."
        fi
    fi

    if [ -z "$latest_release" ] || [ "$latest_release" == "null" ]; then
        print_warning "Could not determine latest release from GitHub API. Using 'latest' tag."
        latest_release="latest"
    fi

    print_status "Latest release: $latest_release"
    local download_url="https://github.com/$REPO_OWNER/$REPO_NAME/releases/download/$latest_release/${BINARY_NAME}-${OS}-${ARCH}"
    print_status "Downloading from: $download_url"

    # Create temporary directory
    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"

    # Download binary
    # Added explicit failure checks for downloads.
    if command_exists "curl"; then
        if ! curl -f -L -o "$BINARY_NAME" "$download_url"; then
            print_error "Download failed using curl. Please check the URL and your connection."
        fi
    else
        if ! wget -O "$BINARY_NAME" "$download_url"; then
            print_error "Download failed using wget. Please check the URL and your connection."
        fi
    fi

    chmod +x "$BINARY_NAME"
    print_success "Download completed"
}

# Function to install binary
install_binary() {
    print_status "Installing binary to $INSTALL_DIR..."

    if [ -w "$INSTALL_DIR" ]; then
        mv "$BINARY_NAME" "$INSTALL_DIR/"
    else
        print_warning "No write permission to $INSTALL_DIR. Trying with sudo..."
        if command_exists "sudo"; then
            sudo mv "$BINARY_NAME" "$INSTALL_DIR/"
        else
            print_error "No sudo available and no write permission to $INSTALL_DIR. Please move the binary manually from '$TEMP_DIR/$BINARY_NAME' to a directory in your PATH."
        fi
    fi

    print_success "Binary installed to $INSTALL_DIR"
}

# Function to verify installation
verify_installation() {
    print_status "Verifying installation..."

    if command_exists "$BINARY_NAME"; then
        local install_path
        # Use command -v instead of which for better portability.
        install_path=$(command -v "$BINARY_NAME")
        local version
        version=$("$BINARY_NAME" version 2>/dev/null || echo "unknown")
        print_success "LeetSolv installed successfully!"
        print_status "Version: $version"
        print_status "Location: $install_path"
    else
        print_error "Installation verification failed. '$BINARY_NAME' not found in your PATH."
    fi
}

# Function to create configuration directory
setup_config() {
    print_status "Setting up configuration directory..."
    mkdir -p "$CONFIG_DIR"
    print_success "Configuration directory ready at $CONFIG_DIR"
}

# Function to show post-install instructions
show_post_install() {
    echo
    print_success "Installation completed successfully!"
    echo
    echo "Next steps:"
    echo "1. Run '$BINARY_NAME' to start the application."
    echo "2. Run '$BINARY_NAME help' to see available commands."
    echo "3. Configuration files are in '$CONFIG_DIR'."
    echo
    echo "To uninstall, run: 'sudo rm \"$INSTALL_DIR/$BINARY_NAME\"' or use the --uninstall flag with this script."
    echo
}

# Main installation function
main() {
    echo -e "${BLUE}╭───────────────────────────────────────────────────╮${NC}"
    echo -e "${BLUE}│                                                   │${NC}"
    echo -e "${BLUE}│    ░▒▓   LeetSolv — CLI SRS for LeetCode   ▓▒░    │${NC}"
    echo -e "${BLUE}│                                                   │${NC}"
    echo -e "${BLUE}│                Installation Script                │${NC}"
    echo -e "${BLUE}│                                                   │${NC}"
    echo -e "${BLUE}╰───────────────────────────────────────────────────╯${NC}"
    echo

    detect_platform
    backup_existing
    download_release
    install_binary
    setup_config
    verify_installation
    show_post_install
}

# Uninstall function with interactive cleanup.
uninstall() {
    print_status "Uninstalling LeetSolv..."
    local install_path
    # Use command -v instead of which for better portability.
    install_path=$(command -v "$BINARY_NAME" || echo "")

    if [ -n "$install_path" ]; then
        print_status "Removing binary from $install_path"
        if [ -w "$(dirname "$install_path")" ]; then
            rm -f "$install_path"
        else
            sudo rm -f "$install_path"
        fi
        print_success "LeetSolv binary uninstalled."
    else
        print_warning "LeetSolv binary not found in PATH."
    fi

    if [ -d "$CONFIG_DIR" ]; then
        echo
        print_warning "Configuration directory found at '$CONFIG_DIR'."
        read -p "Do you want to remove all configurations and backups? [y/N] " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            print_status "Removing configuration directory..."
            rm -rf "$CONFIG_DIR"
            print_success "Configuration directory removed."
        else
            print_status "Skipping removal of configuration directory."
        fi
    fi
    echo
    print_success "Uninstallation process finished."
}


# Handle command line arguments
case "${1:-}" in
    --help|-h)
        echo "Usage: $0 [COMMAND]"
        echo
        echo "Commands:"
        echo "  (no command)   Install or update LeetSolv"
        echo "  --uninstall    Uninstall LeetSolv"
        echo "  --help, -h     Show this help message"
        echo "  --version, -v  Show installer version"
        exit 0
        ;;
    --version|-v)
        echo "LeetSolv Installer v1.1.0"
        exit 0
        ;;
    --uninstall)
        uninstall
        exit 0
        ;;
    "")
        # No arguments, proceed with installation
        main
        ;;
    *)
        print_error "Unknown option: $1"
        ;;
esac