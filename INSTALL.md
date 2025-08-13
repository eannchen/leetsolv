# LeetSolv Installation Guide

This guide explains how to install LeetSolv on your system using various methods.

## Prerequisites

- **Go 1.24.4 or later** (for building from source)
- **Git** (for cloning the repository)
- **Internet connection** (for downloading releases)

## Installation Methods

### Method 1: Using Installation Scripts (Recommended)

#### Linux/macOS
```bash
# Download and run the installation script
curl -fsSL https://raw.githubusercontent.com/eannchen/leetsolv/main/install.sh | bash

# Or download first, then run
wget https://raw.githubusercontent.com/eannchen/leetsolv/main/install.sh
chmod +x install.sh
./install.sh
```

#### Windows (Command Prompt)
```cmd
# Download and run the batch file
powershell -Command "Invoke-WebRequest -Uri 'https://raw.githubusercontent.com/eannchen/leetsolv/main/install.bat' -OutFile 'install.bat'"
install.bat
```

#### Windows (PowerShell)
```powershell
# Download and run the PowerShell script
Invoke-WebRequest -Uri 'https://raw.githubusercontent.com/eannchen/leetsolv/main/install.ps1' -OutFile 'install.ps1'
.\install.ps1

# If you get execution policy errors, run:
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### Method 2: Building from Source

#### Clone the Repository
```bash
git clone https://github.com/eannchen/leetsolv.git
cd leetsolv
```

#### Build the Application
```bash
# Build for current platform
make build

# Or build for all platforms
make build-all

# Or use Go directly
go build -o leetsolv
```

#### Install Locally
```bash
# Install to your Go bin directory
make install

# Or use Go directly
go install
```

### Method 3: Manual Download

1. Go to the [Releases page](https://github.com/eannchen/leetsolv/releases)
2. Download the appropriate binary for your platform
3. Make it executable (Linux/macOS): `chmod +x leetsolv-<platform>`
4. Move it to a directory in your PATH or run it directly

## Platform-Specific Instructions

### Linux
- **Ubuntu/Debian**: The installation script will install to `/usr/local/bin`
- **CentOS/RHEL**: Same as Ubuntu
- **Arch Linux**: Can use AUR or the installation script

### macOS
- **Intel Macs**: Use the `darwin-amd64` binary
- **Apple Silicon**: Use the `darwin-arm64` binary
- The installation script will install to `/usr/local/bin`

### Windows
- **Windows 10/11**: Use the `windows-amd64.exe` binary
- **Windows ARM**: Use the `windows-arm64.exe` binary
- The installation scripts will install to `%USERPROFILE%\AppData\Local\Programs\leetsolv`

## Verification

After installation, verify that LeetSolv is working:

```bash
# Check if the command is available
leetsolv --version

# Or run the version command
leetsolv version

# Check the help
leetsolv help
```

## Configuration

LeetSolv will create its configuration directory at:
- **Linux/macOS**: `~/.leetsolv/`
- **Windows**: `%USERPROFILE%\.leetsolv\`

Configuration files will be created automatically when you first run the application.

## Troubleshooting

### Common Issues

#### Permission Denied
```bash
# Make the binary executable
chmod +x leetsolv

# Or install with sudo (not recommended)
sudo ./install.sh
```

#### Command Not Found
```bash
# Check if the binary is in your PATH
which leetsolv

# Add the installation directory to PATH
export PATH="$PATH:/path/to/leetsolv"
```

#### Windows Execution Policy
```powershell
# If you get execution policy errors
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

#### Go Version Issues
```bash
# Check your Go version
go version

# Update Go if needed (should be 1.24.4 or later)
```

### Getting Help

If you encounter issues:

1. Check the [GitHub Issues](https://github.com/eannchen/leetsolv/issues)
2. Run `leetsolv help` to see available commands
3. Check the logs in the configuration directory

## Uninstallation

### Using Dedicated Uninstall Scripts (Recommended)
```bash
# Linux/macOS
curl -fsSL https://raw.githubusercontent.com/eannchen/leetsolv/main/uninstall.sh | bash

# Or download first, then run
wget https://raw.githubusercontent.com/eannchen/leetsolv/main/uninstall.sh
chmod +x uninstall.sh
./uninstall.sh

# Windows (Command Prompt)
# Download and run
powershell -Command "Invoke-WebRequest -Uri 'https://raw.githubusercontent.com/eannchen/leetsolv/main/uninstall.bat' -OutFile 'uninstall.bat'"
uninstall.bat

# Windows (PowerShell)
# Download and run
Invoke-WebRequest -Uri 'https://raw.githubusercontent.com/eannchen/leetsolv/main/uninstall.ps1' -OutFile 'uninstall.ps1'
.\uninstall.ps1
```

### Using Installation Scripts
```bash
# Linux/macOS
./install.sh --uninstall

# Windows (Command Prompt)
install.bat --uninstall

# Windows (PowerShell)
.\install.ps1 -Uninstall
```

### Manual Uninstallation
```bash
# Remove the binary
rm /usr/local/bin/leetsolv  # Linux/macOS
# or
rmdir /s /q "%USERPROFILE%\AppData\Local\Programs\leetsolv"  # Windows

# Remove configuration (optional)
rm -rf ~/.leetsolv  # Linux/macOS
# or
rmdir /s /q "%USERPROFILE%\.leetsolv"  # Windows
```

## Development Installation

For developers who want to work on LeetSolv:

```bash
# Clone the repository
git clone https://github.com/eannchen/leetsolv.git
cd leetsolv

# Install dependencies
go mod download

# Run tests
make test

# Build and run
make build
./leetsolv
```

## Support

- **Documentation**: [README.md](README.md)
- **Issues**: [GitHub Issues](https://github.com/eannchen/leetsolv/issues)
- **Discussions**: [GitHub Discussions](https://github.com/eannchen/leetsolv/discussions)

---

**Note**: Replace `eannchen` with your actual GitHub username in all URLs and commands above.
