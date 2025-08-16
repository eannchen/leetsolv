# LeetSolv Installation Guide

This guide explains how to install LeetSolv on your system using various methods.

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

1. The script will install the binary to `/usr/local/bin`.
2. Binary Selection:
   | Operating System   | Processor     | uname -s Output | uname -m Output | Final Binary Name     |
   | ------------------ | ------------- | --------------- | --------------- | --------------------- |
   | Ubuntu/Debian/etc. | Intel / AMD   | Linux           | x86_64          | leetsolv-linux-amd64  |
   | Ubuntu/Debian/etc. | ARM           | Linux           | aarch64         | leetsolv-linux-arm64  |
   | macOS              | Intel         | Darwin          | x86_64          | leetsolv-darwin-amd64 |
   | macOS              | Apple Silicon | Darwin          | arm64           | leetsolv-darwin-arm64 |

### Method 2: Manual Download (Linux/macOS/Windows)

1. Go to the [Releases page](https://github.com/eannchen/leetsolv/releases)
2. Download the appropriate binary for your platform
3. Make it executable (Linux/macOS): `chmod +x leetsolv-<platform>`
4. Move it to a directory in your PATH or run it directly

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

## Verification

After installation, verify that LeetSolv is working:

```bash
# Check if the command is available
leetsolv version

# Check the help
leetsolv help
```

## Configuration

LeetSolv will create its configuration directory at the location of the binary when you first run the application.

Configuration files:
- `questions.json` - Database of questions
- `deltas.json` - Database for action history
- `settings.json` - Settings for LeetSolv
- `info.log` - Informational logs
- `error.log` - Error logs


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

### Getting Help

If you encounter issues:

1. Check the [GitHub Issues](https://github.com/eannchen/leetsolv/issues)
2. Run `leetsolv help` to see available commands
3. Check the logs in the configuration directory

## Uninstallation

Remove the binary and configuration files from the installation directory.

Example:
```bash
# If you installed using the install.sh script
rm /usr/local/bin/leetsolv
```


## Support

- **Installation Help**: This guide
- **Development Guide**: [DEVELOPMENT_GUIDE.md](DEVELOPMENT_GUIDE.md)
- **Documentation**: [README.md](../README.md)
- **Issues**: [GitHub Issues](https://github.com/eannchen/leetsolv/issues)
- **Discussions**: [GitHub Discussions](https://github.com/eannchen/leetsolv/discussions)

---