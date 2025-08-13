#!/bin/bash

# LeetSolv Uninstall Script
# This script removes LeetSolv CLI application

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BINARY_NAME="leetsolv"
INSTALL_DIR="/usr/local/bin"
USER_INSTALL_DIR="$HOME/.local/bin"
CONFIG_DIR="$HOME/.leetsolv"
BACKUP_DIR="$HOME/.leetsolv/backup"

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
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to find LeetSolv installation
find_installation() {
    local locations=(
        "$INSTALL_DIR/$BINARY_NAME"
        "$USER_INSTALL_DIR/$BINARY_NAME"
        "$(which $BINARY_NAME 2>/dev/null)"
    )

    for location in "${locations[@]}"; do
        if [ -n "$location" ] && [ -f "$location" ]; then
            echo "$location"
            return 0
        fi
    done

    return 1
}

# Function to backup before uninstalling
backup_binary() {
    local binary_path="$1"
    if [ -n "$binary_path" ] && [ -f "$binary_path" ]; then
        print_status "Creating backup before uninstallation..."
        mkdir -p "$BACKUP_DIR"
        local backup_path="$BACKUP_DIR/$(date +%Y%m%d_%H%M%S)_$BINARY_NAME"
        cp "$binary_path" "$backup_path"
        print_success "Backup created at: $backup_path"
    fi
}

# Function to remove binary
remove_binary() {
    local binary_path="$1"
    if [ -n "$binary_path" ] && [ -f "$binary_path" ]; then
        print_status "Removing binary: $binary_path"

        # Check if we have write permissions
        if [ -w "$(dirname "$binary_path")" ]; then
            rm -f "$binary_path"
            print_success "Binary removed successfully"
        else
            print_warning "No write permission. Trying with sudo..."
            if command_exists "sudo"; then
                sudo rm -f "$binary_path"
                print_success "Binary removed successfully with sudo"
            else
                print_error "No sudo available and no write permission"
                print_status "Please remove manually: $binary_path"
                return 1
            fi
        fi
    fi
}

# Function to remove configuration
remove_config() {
    if [ -d "$CONFIG_DIR" ]; then
        print_status "Removing configuration directory: $CONFIG_DIR"

        # Ask user if they want to keep config
        read -p "Do you want to keep your configuration files? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            print_status "Configuration directory kept at: $CONFIG_DIR"
        else
            rm -rf "$CONFIG_DIR"
            print_success "Configuration directory removed"
        fi
    else
        print_status "No configuration directory found"
    fi
}

# Function to remove from PATH
remove_from_path() {
    local binary_path="$1"
    if [ -n "$binary_path" ]; then
        local install_dir="$(dirname "$binary_path")"

        # Check common shell configuration files
        local shell_configs=(
            "$HOME/.bashrc"
            "$HOME/.bash_profile"
            "$HOME/.zshrc"
            "$HOME/.profile"
        )

        for config in "${shell_configs[@]}"; do
            if [ -f "$config" ] && grep -q "$install_dir" "$config"; then
                print_status "Found PATH entry in: $config"
                read -p "Do you want to remove the PATH entry from $config? (y/N): " -n 1 -r
                echo
                if [[ $REPLY =~ ^[Yy]$ ]]; then
                    # Remove the PATH line (simple approach)
                    sed -i.bak "/$install_dir/d" "$config"
                    print_success "PATH entry removed from $config"
                    print_warning "Backup created as $config.bak"
                fi
            fi
        done
    fi
}

# Function to show uninstall summary
show_summary() {
    echo
    print_success "Uninstallation completed!"
    echo
    echo "Summary of actions:"
    echo "✓ Binary removed"
    echo "✓ Configuration handled"
    echo "✓ PATH entries checked"
    echo
    echo "If you want to reinstall later, run:"
    echo "curl -fsSL https://raw.githubusercontent.com/eannchen/leetsolv/main/install.sh | bash"
    echo
    echo "Backup files are stored in: $BACKUP_DIR"
}

# Function to check if LeetSolv is still running
check_running() {
    if pgrep -f "$BINARY_NAME" >/dev/null 2>&1; then
        print_warning "LeetSolv is currently running. Please close it before uninstalling."
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
}

# Main uninstall function
main() {
    echo -e "${BLUE}╭───────────────────────────────────────────────────╮${NC}"
    echo -e "${BLUE}│                                                   │${NC}"
    echo -e "${BLUE}│    ░▒▓   LeetSolv — Uninstall Script        ▓▒░    │${NC}"
    echo -e "${BLUE}│                                                   │${NC}"
    echo -e "${BLUE}╰───────────────────────────────────────────────────╯${NC}"
    echo

    # Check if LeetSolv is running
    check_running

    # Find installation
    print_status "Searching for LeetSolv installation..."
    local binary_path=$(find_installation)

    if [ -z "$binary_path" ]; then
        print_warning "LeetSolv not found in common locations"
        print_status "Checking if it's available in PATH..."

        if command_exists "$BINARY_NAME"; then
            binary_path=$(which "$BINARY_NAME")
            print_status "Found in PATH: $binary_path"
        else
            print_error "LeetSolv not found. It may already be uninstalled."
            exit 1
        fi
    fi

    # Confirm uninstallation
    echo
    print_warning "About to uninstall LeetSolv from: $binary_path"
    read -p "Are you sure you want to continue? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_status "Uninstallation cancelled"
        exit 0
    fi

    # Perform uninstallation
    backup_binary "$binary_path"
    remove_binary "$binary_path"
    remove_from_path "$binary_path"
    remove_config

    show_summary
}

# Handle command line arguments
case "${1:-}" in
    --help|-h)
        echo "Usage: $0 [OPTIONS]"
        echo "Options:"
        echo "  --help, -h     Show this help message"
        echo "  --force, -f    Skip confirmation prompts"
        echo "  --config-only  Remove only configuration files"
        exit 0
        ;;
    --force|-f)
        # Skip confirmations (for automation)
        FORCE=true
        main
        ;;
    --config-only)
        # Remove only configuration
        remove_config
        print_success "Configuration cleanup completed"
        exit 0
        ;;
    "")
        # No arguments, proceed with uninstallation
        main
        ;;
    *)
        print_error "Unknown option: $1"
        echo "Use --help for usage information"
        exit 1
        ;;
esac
