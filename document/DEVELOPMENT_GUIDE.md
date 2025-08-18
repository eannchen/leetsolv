# LeetSolv Development Guide

This comprehensive guide covers the CI/CD pipeline, installation system, and GitHub Flow workflow for LeetSolv development.

## Quick Start for Developers

### 1. Install Go

The project requires Go 1.25.0 or later.

```bash
# Install Go (macOS)
brew install go
```

Or download from [official website](https://golang.org/dl/).

### 2. Run in Development Mode
```bash
# Clone the repository
git clone https://github.com/eannchen/leetsolv.git
cd leetsolv

# Run in development mode
make dev

# Run tests
make test

# Build locally
make build
```

## File Structure

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

## GitHub Flow Workflow

### Branch Structure

#### `main` Branch
- **Content**: Only tested, verified, and production-ready code
- **Protection**: Protected, require PR reviews
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

## Complete Development Workflow

### 1. Initial Setup

```bash
# Clone the repository
git clone https://github.com/eannchen/leetsolv.git
cd leetsolv

# Set up remote tracking
git checkout main
git pull origin main
```

### 2. Development

```bash
# Start from main branch
git checkout main
git pull origin main

# Create feature branch
git checkout -b feature/new-command
# Or create bug fix branch `git checkout -b hotfix/test-failures`

# ... make your changes ...

# Run local tests to ensure everything works
make test
make lint
make build

# Commit with conventional commit format
git add .
git commit -m "feat: add new command functionality"

# Push feature branch
git push origin feature/new-command

# Create Pull Request on GitHub
# Source: feature/new-command → Target: main
```

### 3. After PR Merged

```bash
# Switch back to main branch
git checkout main
git pull origin main

# Delete local feature branch
git branch -d feature/new-command

# Delete remote feature branch
git push origin --delete feature/new-command
```

### 4. Release Preparation

```bash
# When ready for release, create a tag from main
git checkout main
git pull origin main

# Create release tag
git tag -a v0.1.0-beta -m "Beta release - ready for testing"
git push origin v0.1.0-beta

# Create GitHub Release
# 1. Go to https://github.com/eannchen/leetsolv/releases
# 2. Click "Create a new release"
# 3. Select the tag you just created
# 4. Add release notes describing changes
# 5. Click "Publish release"
```

The tag style follows [Semantic Versioning](https://semver.org/).


## CI/CD Pipeline

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

## Build System

### Makefile
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

## Tagging Strategy

### Semantic Versioning

The versioning follows [Semantic Versioning](https://semver.org/).

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

## Branch Protection Rules

### GitHub Repository Settings

#### `main` Branch Protection
- ✅ Require pull request reviews before merging
- ✅ Require status checks to pass before merging
- ✅ Require branches to be up to date before merging
- ✅ Restrict pushes that create files larger than 100 MB
- ✅ Include administrators in these restrictions

## Best Practices

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
- Keep release history clean and organized
- Always test locally before releasing
- Update documentation with new features
- Keep CHANGELOG.md updated

## Future Enhancements

### Potential Improvements
1. **Package managers** (Homebrew, Chocolatey, apt)
2. **Auto-updater** functionality
3. **CI/CD matrix** expansion (more Go versions)

## Support

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
