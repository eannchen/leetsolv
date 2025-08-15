# LeetSolv Development Guide

This comprehensive guide covers the CI/CD pipeline, installation system, and GitHub Flow workflow for LeetSolv development.

## 🚀 Quick Start for Developers

### 1. Setup Development Environment
```bash
# Clone the repository
git clone https://github.com/eannchen/leetsolv.git
cd leetsolv

# Install dependencies
go mod download

# Run tests
make test

# Build locally
make build
```

### 2. Create Your First Release
```bash
# Create and push a tag (recommended for development)
git tag -a v0.1.0-alpha -m "Alpha release - development version"
git push origin v0.1.0-alpha

# Create a GitHub release from the tag
# Go to: https://github.com/eannchen/leetsolv/releases
# Click "Create a new release" and select the tag
```

> **Note**: For installation instructions, see [INSTALL.md](INSTALL.md)

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
├── document/
│   ├── INSTALL.md                    # Installation guide
│   └── DEVELOPMENT_GUIDE.md          # This file
└── README.md                         # Project overview
```

## 🌿 GitHub Flow Workflow

### Branch Structure

#### `main` Branch
- **Purpose**: Production-ready, stable code
- **Content**: Only tested, verified, and released code
- **Protection**: Should be protected, require PR reviews
- **Tags**: Release tags are created from this branch
- **Merges**: Feature branches merge directly to main

### Supporting Branches

#### Feature Branches
- **Naming**: `feature/descriptive-name`
- **Source**: Branch from `main`
- **Target**: Merge back to `main` via Pull Request
- **Lifetime**: Short-lived, delete after merge

#### Hotfix Branches
- **Naming**: `hotfix/critical-issue`
- **Source**: Branch from `main`
- **Target**: Merge back to `main` via Pull Request
- **Lifetime**: Very short-lived, delete after merge

## 🔄 Complete Development Workflow

### 1. Initial Setup

```bash
# Clone the repository
git clone https://github.com/eannchen/leetsolv.git
cd leetsolv

# Set up remote tracking
git checkout main
git pull origin main
```

### 2. Feature Development

```bash
# Start from main branch
git checkout main
git pull origin main

# Create feature branch
git checkout -b feature/new-command
# ... make your changes ...

# Commit with conventional commit format
git add .
git commit -m "feat: add new command functionality

- Implements new 'search' command
- Adds fuzzy search capabilities
- Updates help documentation"

# Push feature branch
git push origin feature/new-command

# Create Pull Request on GitHub
# Source: feature/new-command → Target: main
```

### 3. Code Review and Integration

```bash
# After PR is approved and merged
git checkout main
git pull origin main

# Delete local feature branch
git branch -d feature/new-command

# Delete remote feature branch
git push origin --delete feature/new-command
```

### 4. Testing and Validation

```bash
# CI/CD automatically runs on main branch
# Run local tests to ensure everything works
make test
make lint
make build

# If issues found, create a new feature branch to fix them
git checkout -b fix/test-failures
# ... fix issues ...
git add .
git commit -m "fix: resolve test failures in new command"
git push origin fix/test-failures
# Create PR to merge fix back to main
```

### 5. Release Preparation

Before creating a release, ensure:
- All tests pass: `make test`
- Code is formatted: `make lint`
- Version information is updated in `main.go`
- CHANGELOG.md is updated with new features/fixes

```bash
# When ready for release, create a tag from main
git checkout main
git pull origin main

# Create release tag
git tag -a v0.1.0-alpha -m "Alpha release - development version

- New search command functionality
- Improved error handling
- Updated documentation

This is an alpha release for testing purposes."
git push origin v0.1.0-alpha

# Create GitHub Release
# 1. Go to https://github.com/eannchen/leetsolv/releases
# 2. Click "Create a new release"
# 3. Select the tag you just created
# 4. Add release notes describing changes
# 5. Click "Publish release"
```

#### Manual Release Process

If you need to create a release manually:

```bash
# Build for all platforms
make build-all

# Create checksums
cd dist
sha256sum leetsolv-* > checksums.txt

# Upload to GitHub release manually
```

### 6. Hotfix Process (if needed)

```bash
# For critical bugs that need immediate fix
git checkout main
git pull origin main
git checkout -b hotfix/critical-bug
# ... fix the critical bug ...

# Test the fix
make test
make build

# Commit the fix
git add .
git commit -m "fix: resolve critical bug in search functionality

- Fixes crash when searching with empty query
- Adds input validation
- Updates error messages"

# Push and create PR
git push origin hotfix/critical-bug
# Create PR: hotfix/critical-bug → main

# After PR is merged, create hotfix tag
git checkout main
git pull origin main
git tag -a v0.1.1 -m "Hotfix release - critical bug fix"
git push origin v0.1.1

# Delete hotfix branch
git branch -d hotfix/critical-bug
git push origin --delete hotfix/critical-bug
```

## 🔄 CI/CD Pipeline

### GitHub Actions Workflow (`.github/workflows/ci.yml`)

**Triggers:**
- Push to `main` branch
- Pull requests to `main`
- GitHub releases
- Manual workflow dispatch

**Jobs:**

#### 1. Test Job
- **Matrix**: Multiple Go versions (1.25.0) × Platforms (Linux, macOS, Windows)
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

### Branch-Specific Behavior

#### `main` Branch
- **Tests**: Full test suite runs
- **Builds**: Creates production builds
- **Artifacts**: Uploaded to GitHub Actions artifacts
- **Deployment**: No automatic deployment

#### Release Tags
- **Tests**: Full test suite runs
- **Builds**: Creates release builds with tag version
- **Release**: Automatically creates GitHub release
- **Artifacts**: Binaries uploaded to release

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

#### Version Variables
The application includes these version variables:

```go
var (
    Version   = "dev"      // Set during build
    BuildTime = "unknown"  // Set during build
    GitCommit = "unknown"  // Set during build
)
```

#### Build Flags
These are set during the CI build process:

```bash
go build -ldflags="-s -w -X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.GitCommit=$GIT_COMMIT"
```

The CI pipeline automatically injects version information:
- `Version`: From Git tag or commit
- `BuildTime`: Build timestamp
- `GitCommit`: Git commit hash

## 📱 Installation System Overview

The project includes automated installation scripts for all platforms. For detailed installation instructions, see [INSTALL.md](INSTALL.md).

### Installation Script Features
- **Cross-platform**: Linux, macOS, Windows support
- **Automatic detection**: Platform and architecture detection
- **PATH management**: Automatic PATH configuration
- **Uninstall support**: Dedicated uninstall scripts with backup creation

### Installation Locations
- **Linux/macOS**: `/usr/local/bin` (system) or `~/.local/bin` (user)
- **Windows**: `%USERPROFILE%\AppData\Local\Programs\leetsolv`

### Configuration Directory
- **Linux/macOS**: `~/.leetsolv/`
- **Windows**: `%USERPROFILE%\.leetsolv\`

## 🏷️ Tagging Strategy

### Version Naming Convention

#### Pre-release Tags
```bash
# Alpha release (development/testing)
git tag -a v0.1.0-alpha -m "Alpha release - development version"

# Beta release (testing)
git tag -a v0.1.0-beta -m "Beta release - ready for testing"

# Release candidate
git tag -a v0.1.0-rc.1 -m "Release candidate 1"
```

#### Stable Release Tags
```bash
# Minor version (new features)
git tag -a v0.1.0 -m "Release version 0.1.0"

# Patch version (bug fixes)
git tag -a v0.1.1 -m "Patch release - bug fixes"

# Major version (breaking changes)
git tag -a v1.0.0 -m "Major release - stable version"
```

### Tag Message Format

```bash
git tag -a v0.1.0-alpha -m "Alpha release - development version

## What's New
- New search command functionality
- Improved error handling
- Updated documentation

## Breaking Changes
- None

## Known Issues
- Search may be slow with large datasets

## Testing
- Run 'make test' to verify functionality
- Test search with various query types"
```

## 📋 Branch Protection Rules

### GitHub Repository Settings

#### `main` Branch Protection
- ✅ Require pull request reviews before merging
- ✅ Require status checks to pass before merging
- ✅ Require branches to be up to date before merging
- ✅ Restrict pushes that create files larger than 100 MB
- ✅ Include administrators in these restrictions

## 🔍 Useful Git Commands

### Branch Management
```bash
# List all branches
git branch -a

# See branch relationships
git log --graph --oneline --all

# Check which branch you're on
git branch

# See remote tracking
git branch -vv
```

### Tag Management
```bash
# List all tags
git tag -l

# See tag details
git show v0.1.0-alpha

# Delete local tag
git tag -d v0.1.0-alpha

# Delete remote tag
git push origin --delete v0.1.0-alpha
```

### History and Comparison
```bash
# See commit history
git log --oneline -10

# Compare feature branch with main
git diff main..feature/my-feature

# See what's in feature branch but not in main
git log main..feature/my-feature

# See what's in main but not in feature branch
git log feature/my-feature..main
```

## 📚 Best Practices

### 1. Commit Messages
- Use conventional commit format: `type(scope): description`
- Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`
- Keep descriptions clear and concise
- Reference issues when applicable

### 2. Branch Naming
- Use descriptive names: `feature/user-authentication`
- Use kebab-case for multi-word names
- Prefix with type: `feature/`, `fix/`, `hotfix/`

### 3. Pull Requests
- Create PRs for all changes
- Use descriptive titles and descriptions
- Include testing instructions
- Request reviews from team members

### 4. Testing
- Always run tests before pushing
- Ensure CI/CD passes before merging
- Test locally before creating releases
- Document known issues in release notes

### 5. Release Management
- Use semantic versioning (v1.0.0 format)
- Create detailed release notes with meaningful descriptions
- Test release builds before publishing
- Keep release history clean and organized
- Always test locally before releasing
- Update documentation with new features
- Keep CHANGELOG.md updated

## 🚨 Common Pitfalls

### 1. Direct Commits to Main
- ❌ Never commit directly to main
- ✅ Always use pull requests
- ✅ Require code review

### 2. Large Feature Branches
- ❌ Long-lived feature branches that diverge too much
- ✅ Keep feature branches small and focused
- ✅ Merge frequently to main

### 3. Incomplete Testing
- ❌ Releasing without proper testing
- ✅ Run full test suite
- ✅ Test on multiple platforms

### 4. Poor Tag Messages
- ❌ Generic tag messages
- ✅ Detailed release notes
- ✅ Include what's new and breaking changes

## 🔧 Troubleshooting

### Merge Conflicts
```bash
# Abort merge if needed
git merge --abort

# See conflict files
git status

# Resolve conflicts manually
# Edit conflicted files, then:
git add .
git commit -m "resolve merge conflicts"
```

### Lost Commits
```bash
# Find lost commits
git reflog

# Recover lost commit
git checkout -b recovery <commit-hash>
```

### Branch Cleanup
```bash
# Delete merged branches
git branch --merged | grep -v "\*" | xargs -n 1 git branch -d

# Delete remote branches that no longer exist
git remote prune origin
```

## 🚨 Troubleshooting

### Common Issues

#### CI/CD Failures
- Check Go version compatibility
- Verify workflow file syntax
- Check GitHub Actions logs

#### Build Failures
- Check Go version compatibility
- Verify all dependencies are available
- Check for platform-specific code

#### Test Failures
- Run tests locally: `make test` (with race detection)
- For platforms without race support: `make test-no-race`
- Check for race conditions: `go test -race` (requires CGO_ENABLED=1)
- Verify test data files

#### Release Failures
- Ensure GitHub token has release permissions
- Check tag format (should be semantic versioning)
- Verify workflow file syntax

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

# Check workflow runs
gh run list

# View workflow logs
gh run view <run-id>

# Re-run failed workflow
gh run rerun <run-id>
```

## 🔒 Security Features

- **Checksums**: SHA256 verification for all binaries
- **Backups**: Automatic backup creation before uninstallation
- **Permissions**: Proper file permissions and PATH management
- **Verification**: Installation verification and integrity checks

### Security Best Practices

- Never commit sensitive information
- Use GitHub secrets for API keys
- Verify checksums before installing
- Keep dependencies updated

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

## 🆘 Support

### Getting Help
1. **Installation**: See [INSTALL.md](INSTALL.md)
2. **Development**: This guide
3. Run `leetsolv help` for command information
4. Check GitHub Issues for known problems
5. Review CI/CD logs for build failures

### Contributing
1. Fork the repository
2. Create a feature branch from `main`
3. Make changes and test locally
4. Submit a pull request to `main`

---

**Note**: This system uses GitHub Flow for simplified development workflow. The CI/CD pipeline ensures quality, while the installation scripts provide a seamless user experience across all platforms.
