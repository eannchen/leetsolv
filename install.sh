#!/bin/bash

# LeetSolv Installation Script
# This script installs LeetSolv CLI application

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO_OWNER="eannchen"  # Change this to your GitHub username
REPO_NAME="leetsolv"
BINARY_NAME="leetsolv"
INSTALL_DIR="/usr/local/bin"
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

# Function to detect OS and architecture
detect_platform() {
    OS=""
    ARCH=""

    case "$(uname -s)" in
        Linux*)     OS="linux";;
        Darwin*)    OS="darwin";;
        CYGWIN*)   OS="windows";;
        MINGW*)    OS="windows";;
        MSYS*)     OS="windows";;
        *)         OS="unknown";;
    esac

    case "$(uname -m)" in
        x86_64)    ARCH="amd64";;
        aarch64)   ARCH="arm64";;
        arm64)     ARCH="arm64";;
        *)         ARCH="unknown";;
    esac

    if [ "$OS" = "unknown" ] || [ "$ARCH" = "unknown" ]; then
        print_error "Unsupported platform: $OS-$ARCH"
        exit 1
    fi

    print_status "Detected platform: $OS-$ARCH"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check if running as root
check_root() {
    if [ "$EUID" -eq 0 ]; then
        print_warning "Running as root. This is not recommended for security reasons."
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
}

# Function to backup existing installation
backup_existing() {
    if command_exists "$BINARY_NAME"; then
        print_status "Backing up existing installation..."
        mkdir -p "$BACKUP_DIR"
        local backup_path="$BACKUP_DIR/$(date +%Y%m%d_%H%M%S)_$BINARY_NAME"
        cp "$(which $BINARY_NAME)" "$backup_path"
        print_success "Backup created at: $backup_path"
    fi
}

# Function to download latest release
download_release() {
    print_status "Checking for latest release..."

    # Try to get latest release from GitHub API
    local latest_release
    if command_exists "curl"; then
        latest_release=$(curl -s "https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    elif command_exists "wget"; then
        latest_release=$(wget -qO- "https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    else
        print_error "Neither curl nor wget found. Please install one of them."
        exit 1
    fi

    if [ -z "$latest_release" ]; then
        print_warning "Could not determine latest release. Using 'latest' tag."
        latest_release="latest"
    fi

    print_status "Latest release: $latest_release"

    # Download URL
    local download_url="https://github.com/$REPO_OWNER/$REPO_NAME/releases/download/$latest_release/${BINARY_NAME}-${OS}-${ARCH}"
    if [ "$OS" = "windows" ]; then
        download_url="${download_url}.exe"
    fi

    print_status "Downloading from: $download_url"

    # Create temporary directory
    local temp_dir=$(mktemp -d)
    cd "$temp_dir"

    # Download binary
    if command_exists "curl"; then
        curl -L -o "$BINARY_NAME" "$download_url"
    else
        wget -O "$BINARY_NAME" "$download_url"
    fi

    # Make executable
    chmod +x "$BINARY_NAME"

    print_success "Download completed"
}

# Function to install binary
install_binary() {
    print_status "Installing binary..."

    # Check if we have write permissions to install directory
    if [ ! -w "$INSTALL_DIR" ]; then
        print_warning "No write permission to $INSTALL_DIR. Trying to use sudo..."
        if command_exists "sudo"; then
            sudo mv "$BINARY_NAME" "$INSTALL_DIR/"
        else
            print_error "No sudo available and no write permission to $INSTALL_DIR"
            print_status "You can manually move the binary to a directory in your PATH"
            print_status "Binary location: $(pwd)/$BINARY_NAME"
            exit 1
        fi
    else
        mv "$BINARY_NAME" "$INSTALL_DIR/"
    fi

    print_success "Binary installed to $INSTALL_DIR"
}

# Function to verify installation
verify_installation() {
    print_status "Verifying installation..."

    if command_exists "$BINARY_NAME"; then
        local version=$("$BINARY_NAME" version 2>/dev/null || echo "unknown")
        print_success "LeetSolv installed successfully!"
        print_status "Version: $version"
        print_status "Location: $(which $BINARY_NAME)"
        print_status "You can now run: $BINARY_NAME"
    else
        print_error "Installation verification failed"
        exit 1
    fi
}

# Function to create configuration directory
setup_config() {
    print_status "Setting up configuration directory..."
    mkdir -p "$HOME/.leetsolv"
    print_success "Configuration directory created at $HOME/.leetsolv"
}

# Function to show post-install instructions
show_post_install() {
    echo
    print_success "Installation completed successfully!"
    echo
    echo "Next steps:"
    echo "1. Run '$BINARY_NAME' to start the application"
    echo "2. Run '$BINARY_NAME help' to see available commands"
    echo "3. Configuration files will be created in $HOME/.leetsolv/"
    echo
    echo "To uninstall, run: sudo rm $INSTALL_DIR/$BINARY_NAME"
    echo
}

# Main installation function
main() {
    echo -e "${BLUE}╭───────────────────────────────────────────────────╮${NC}"
    echo -e "${BLUE}│                                                   │${NC}"
    echo -e "${BLUE}│    ░▒▓   LeetSolv — CLI SRS for LeetCode   ▓▒░    │${NC}"
    echo -e "${BLUE}│                                                   │${NC}"
    echo -e "${BLUE}│                Installation Script               │${NC}"
    echo -e "${BLUE}│                                                   │${NC}"
    echo -e "${BLUE}╰───────────────────────────────────────────────────╯${NC}"
    echo

    # Check prerequisites
    check_root
    detect_platform

    # Installation steps
    backup_existing
    download_release
    install_binary
    setup_config
    verify_installation
    show_post_install
}

# Handle command line arguments
case "${1:-}" in
    --help|-h)
        echo "Usage: $0 [OPTIONS]"
        echo "Options:"
        echo "  --help, -h     Show this help message"
        echo "  --version, -v  Show version information"
        echo "  --uninstall    Uninstall LeetSolv"
        exit 0
        ;;
    --version|-v)
        echo "LeetSolv Installer v1.0.0"
        exit 0
        ;;
    --uninstall)
        print_status "Uninstalling LeetSolv..."
        if command_exists "$BINARY_NAME"; then
            if [ -w "$INSTALL_DIR" ]; then
                rm -f "$INSTALL_DIR/$BINARY_NAME"
            else
                sudo rm -f "$INSTALL_DIR/$BINARY_NAME"
            fi
            print_success "LeetSolv uninstalled successfully"
        else
            print_warning "LeetSolv not found in PATH"
        fi
        exit 0
        ;;
    "")
        # No arguments, proceed with installation
        main
        ;;
    *)
        print_error "Unknown option: $1"
        echo "Use --help for usage information"
        exit 1
        ;;
esac
