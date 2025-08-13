# LeetSolv CI/CD & Installation System - Complete Guide

This document provides a comprehensive overview of the CI/CD pipeline and installation system for LeetSolv.

## 🚀 Quick Start

### 1. Create Your First Release
```bash
# Create and push an alpha tag (recommended for development)
git tag -a v0.1.0-alpha -m "Alpha release - development version"
git push origin v0.1.0-alpha

# Create a GitHub release from the tag
# Go to: https://github.com/eannchen/leetsolv/releases
# Click "Create a new release" and select the tag
```

### 2. Install LeetSolv
```bash
# Linux/macOS
curl -fsSL https://raw.githubusercontent.com/eannchen/leetsolv/main/install.sh | bash

# Windows (PowerShell)
Invoke-WebRequest -Uri 'https://raw.githubusercontent.com/eannchen/leetsolv/main/install.ps1' -OutFile 'install.ps1'
.\install.ps1
```

## 📁 File Structure

```
leetsolv/
├── .github/workflows/ci.yml          # GitHub Actions CI/CD pipeline
├── install.sh                        # Linux/macOS installation script
├── install.bat                       # Windows batch installation script
├── install.ps1                       # Windows PowerShell installation script
├── uninstall.sh                      # Linux/macOS uninstall script
├── uninstall.bat                     # Windows batch uninstall script
├── Makefile                          # Enhanced build system
├── INSTALL.md                        # Installation guide
├── RELEASE.md                        # Release workflow guide
└── CI_INSTALL_SUMMARY.md            # This file
```

## 🔄 CI/CD Pipeline

### GitHub Actions Workflow (`.github/workflows/ci.yml`)

**Triggers:**
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop`
- GitHub releases

**Jobs:**

#### 1. Test Job
- **Matrix**: Multiple Go versions (1.24.4, 1.25.0) × Platforms (Linux, macOS, Windows)
- **Features**: Race condition detection, code coverage, dependency caching
- **Output**: Coverage reports uploaded to Codecov

#### 2. Build Job
- **Condition**: Runs on push or release events
- **Output**: Cross-platform binaries for all supported architectures
- **Artifacts**: Uploaded to GitHub Actions artifacts

#### 3. Release Job
- **Condition**: Only on GitHub releases
- **Action**: Automatically creates GitHub release with binaries
- **Files**: All platform binaries + SHA256 checksums

### Supported Platforms
- **Linux**: AMD64, ARM64
- **macOS**: Intel (AMD64), Apple Silicon (ARM64)
- **Windows**: AMD64, ARM64

## 🛠️ Build System

### Enhanced Makefile
```bash
make build        # Build for current platform
make build-all    # Build for all platforms
make test         # Run tests with coverage
make lint         # Code formatting and validation
make install      # Install locally
make clean        # Clean build artifacts
```

### Version Management
The CI pipeline automatically injects version information:
- `Version`: From Git tag or commit
- `BuildTime`: Build timestamp
- `GitCommit`: Git commit hash

## 📱 Installation System

### Installation Scripts

#### Linux/macOS (`install.sh`)
- **Features**: Cross-platform detection, automatic download, PATH management
- **Installation**: `/usr/local/bin` (system-wide) or `~/.local/bin` (user)
- **Options**: `--help`, `--version`, `--uninstall`

#### Windows Command Prompt (`install.bat`)
- **Features**: Windows batch script, PowerShell fallback
- **Installation**: `%USERPROFILE%\AppData\Local\Programs\leetsolv`
- **PATH**: Automatically added to user PATH

#### Windows PowerShell (`install.ps1`)
- **Features**: Modern PowerShell, better error handling
- **Execution Policy**: Checks and guides user through restrictions
- **Options**: `-Help`, `-Version`, `-Uninstall`

### Uninstall Scripts

#### Linux/macOS (`uninstall.sh`)
- **Features**: Automatic detection, backup creation, PATH cleanup
- **Options**: `--help`, `--force`, `--config-only`
- **Safety**: Confirms before removal, creates backups

#### Windows (`uninstall.bat`)
- **Features**: Windows-specific cleanup, directory cleanup
- **Options**: `--help`, `--force`, `--config-only`
- **Safety**: Confirms before removal, creates backups

## 🔧 Configuration

### Installation Locations
- **Linux/macOS**: `/usr/local/bin` (system) or `~/.local/bin` (user)
- **Windows**: `%USERPROFILE%\AppData\Local\Programs\leetsolv`

### Configuration Directory
- **Linux/macOS**: `~/.leetsolv/`
- **Windows**: `%USERPROFILE%\.leetsolv\`

## 📋 Release Workflow

### 1. Development Phase
```bash
# Make changes and test
make test
make lint

# Commit and push
git add .
git commit -m "Feature: add new functionality"
git push origin main
```

### 2. Create Release
```bash
# Create tag
git tag -a v0.1.0-alpha -m "Alpha release - development version"
git push origin v0.1.0-alpha

# Create GitHub release (triggers CI/CD)
# Go to GitHub → Releases → Create new release
```

### 3. CI/CD Pipeline
1. **Tests run** on all platforms
2. **Builds create** cross-platform binaries
3. **Release job** uploads binaries to GitHub release

### 4. User Installation
Users can now install from the release:
```bash
curl -fsSL https://raw.githubusercontent.com/eannchen/leetsolv/main/install.sh | bash
```

## 🚨 Troubleshooting

### Common Issues

#### CI/CD Failures
- Check Go version compatibility
- Verify workflow file syntax
- Check GitHub Actions logs

#### Installation Issues
- Ensure scripts are executable (`chmod +x install.sh`)
- Check PATH environment variable
- Verify platform compatibility

#### Build Issues
- Run `make test` locally first
- Check Go version (`go version`)
- Verify dependencies (`go mod download`)

### Debugging Commands
```bash
# Check installation
which leetsolv
leetsolv version

# Check configuration
ls -la ~/.leetsolv/

# Test build locally
make build
./leetsolv help
```

## 🔒 Security Features

- **Checksums**: SHA256 verification for all binaries
- **Backups**: Automatic backup creation before uninstallation
- **Permissions**: Proper file permissions and PATH management
- **Verification**: Installation verification and integrity checks

## 📈 Future Enhancements

### Potential Improvements
1. **Code signing** for macOS and Windows
2. **Docker images** for containerized deployment
3. **Package managers** (Homebrew, Chocolatey, apt)
4. **Auto-updater** functionality
5. **CI/CD matrix** expansion (more Go versions)

### Monitoring
- **Code coverage** tracking
- **Build success rates** monitoring
- **Release download** statistics
- **User feedback** collection

## 📚 Documentation

- **INSTALL.md**: Complete installation guide
- **RELEASE.md**: Release workflow documentation
- **README.md**: Project overview and usage
- **This file**: CI/CD and installation system overview

## 🆘 Support

### Getting Help
1. Check the documentation files
2. Run `leetsolv help` for command information
3. Check GitHub Issues for known problems
4. Review CI/CD logs for build failures

### Contributing
1. Fork the repository
2. Create a feature branch
3. Make changes and test locally
4. Submit a pull request

---

**Note**: This system is designed to be robust and user-friendly. The CI/CD pipeline ensures quality, while the installation scripts provide a seamless user experience across all platforms.
