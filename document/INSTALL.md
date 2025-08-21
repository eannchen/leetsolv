# LeetSolv Installation Guide

This guide walks you through installing LeetSolv on Linux, macOS, and Windows. Choose the installation method that best fits your platform and preferences.

## Installation Methods

### Method 1: Using Installation Scripts (Linux/macOS)

```bash
# Download and run the installation script
curl -fsSL https://raw.githubusercontent.com/eannchen/leetsolv/main/install.sh | bash

# Or download first, then run
wget https://raw.githubusercontent.com/eannchen/leetsolv/main/install.sh
chmod +x install.sh
./install.sh
```

#### What the Install Script Does

The `install.sh` script automates the entire installation process:

1. **Platform Detection**: Automatically detects your operating system and CPU architecture
2. **Latest Release**: Fetches the latest version from GitHub releases
3. **Binary Download**: Downloads the correct binary for your platform
4. **Installation**: Installs the binary to `/usr/local/bin` (may require sudo)
5. **Configuration Setup**: Creates the configuration directory at `~/.leetsolv`
6. **Backup**: Backs up any existing installation before updating
7. **Verification**: Tests that the installation was successful

#### Binary Selection Logic

The script uses `uname` commands to determine the correct binary:

| Operating System   | Processor     | uname -s Output | uname -m Output | Final Binary Name     |
| ------------------ | ------------- | --------------- | --------------- | --------------------- |
| Ubuntu/Debian/etc. | Intel / AMD   | Linux           | x86_64          | leetsolv-linux-amd64  |
| Ubuntu/Debian/etc. | ARM           | Linux           | aarch64         | leetsolv-linux-arm64  |
| macOS              | Intel         | Darwin          | x86_64          | leetsolv-darwin-amd64 |
| macOS              | Apple Silicon | Darwin          | arm64           | leetsolv-darwin-arm64 |

#### Script Options

```bash
# Show help
./install.sh --help

# Show installer version
./install.sh --version

# Uninstall LeetSolv
./install.sh --uninstall
```

#### What Gets Created

- **Binary**: `/usr/local/bin/leetsolv`
- **Config Directory**: `~/.leetsolv/`
- **Backup Directory**: `~/.leetsolv/backup/` (for existing installations)

### Method 2: Manual Download (All Platforms)

1. Go to the [Releases page](https://github.com/eannchen/leetsolv/releases)
2. Download the appropriate binary for your platform:
   - **Linux**: `leetsolv-linux-amd64` or `leetsolv-linux-arm64`
   - **macOS**: `leetsolv-darwin-amd64` or `leetsolv-darwin-arm64`
   - **Windows**: `leetsolv-windows-amd64.exe` or `leetsolv-windows-arm64.exe`

#### Linux/macOS Setup
```bash
# Make it executable
chmod +x leetsolv-<platform>

# Rename and move to PATH (optional)
sudo mv leetsolv-<platform> /usr/local/bin/leetsolv

# Or run directly from current directory
./leetsolv-<platform>
```

#### Windows Setup
```cmd
# Rename the binary (optional)
ren leetsolv-windows-amd64.exe leetsolv.exe

# Move to a directory in PATH, for example:
move leetsolv.exe C:\Windows\System32\

# Or add current directory to PATH temporarily
set PATH=%PATH%;%CD%

# Run the application
leetsolv.exe
```

**Windows PATH Setup (Permanent):**
1. Create a folder like `C:\Program Files\LeetSolv\`
2. Move `leetsolv.exe` to this folder
3. Add `C:\Program Files\LeetSolv\` to your system PATH:
   - Press `Win + R`, type `sysdm.cpl`, press Enter
   - Go to "Advanced" tab → "Environment Variables"
   - Under "System Variables", find and edit "Path"
   - Add the new path and click OK

### Method 3: Building from Source

#### Prerequisites

- **Git**
- **Go 1.25.0 or later**


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
```

#### Install Locally
```bash
# Install to your Go bin directory
make install
```

## Configuration

Read [README.md](../README.md) for more information.

## Uninstallation

### Using the install script (Linux/macOS)
```bash
# Download and run with uninstall flag
curl -fsSL https://raw.githubusercontent.com/eannchen/leetsolv/main/install.sh | bash -s -- --uninstall

# Or if you have the script locally
./install.sh --uninstall
```

### Manual uninstallation

#### Linux/macOS
```bash
# Remove the binary
sudo rm /usr/local/bin/leetsolv

# Remove configuration directory (optional)
rm -rf ~/.leetsolv
```

#### Windows
```cmd
# Remove the binary (if installed in System32)
del C:\Windows\System32\leetsolv.exe

# Or remove from custom location
del "C:\Program Files\LeetSolv\leetsolv.exe"
rmdir "C:\Program Files\LeetSolv"

# Remove configuration directory (optional)
rmdir /s "%USERPROFILE%\.leetsolv"
```

**Windows PowerShell:**
```powershell
# Remove the binary
Remove-Item "C:\Program Files\LeetSolv\leetsolv.exe"
Remove-Item "C:\Program Files\LeetSolv"

# Remove configuration directory (optional)
Remove-Item -Recurse "$env:USERPROFILE\.leetsolv"
```

## Support

- **Installation Help**: This guide
- **Development Guide**: [DEVELOPMENT_GUIDE.md](DEVELOPMENT_GUIDE.md)
- **Documentation**: [README.md](../README.md)
- **Issues**: [GitHub Issues](https://github.com/eannchen/leetsolv/issues)
- **Discussions**: [GitHub Discussions](https://github.com/eannchen/leetsolv/discussions)
