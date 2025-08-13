# LeetSolv Release Workflow

This document explains how to create releases and trigger the CI/CD pipeline for LeetSolv.

## Release Process

### 1. Prepare for Release

Before creating a release, ensure:

- All tests pass: `make test`
- Code is formatted: `make lint`
- Version information is updated in `main.go`
- CHANGELOG.md is updated with new features/fixes

### 2. Create a Git Tag

For development releases, use pre-release tags:

```bash
# Alpha release (development version)
git tag -a v0.1.0-alpha -m "Alpha release - development version"
git push origin v0.1.0-alpha

# Beta release (testing version)
git tag -a v0.1.0-beta -m "Beta release - ready for testing"
git push origin v0.1.0-beta

# Stable release (production ready)
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

**Recommended for your current development stage:**
```bash
git tag -a v0.1.0-alpha -m "Alpha release - development version"
git push origin v0.1.0-alpha
```

### 3. Create GitHub Release

1. Go to [GitHub Releases](https://github.com/eannchen/leetsolv/releases)
2. Click "Create a new release"
3. Select the tag you just created
4. Add release notes describing changes
5. Click "Publish release"

### 4. CI/CD Pipeline

The GitHub Actions workflow will automatically:

1. **Test**: Run tests on multiple Go versions and platforms
2. **Build**: Compile binaries for all supported platforms
3. **Release**: Upload binaries to the GitHub release

## Supported Platforms

The CI/CD pipeline builds for:

- **Linux**: AMD64, ARM64
- **macOS**: AMD64, ARM64 (Intel & Apple Silicon)
- **Windows**: AMD64, ARM64

## Build Artifacts

Each release includes:

- Platform-specific binaries
- SHA256 checksums for verification
- Source code archive

## Manual Release

If you need to create a release manually:

```bash
# Build for all platforms
make build-all

# Create checksums
cd dist
sha256sum leetsolv-* > checksums.txt

# Upload to GitHub release manually
```

## Version Management

### Version Variables

The application includes these version variables:

```go
var (
    Version   = "dev"      // Set during build
    BuildTime = "unknown"  // Set during build
    GitCommit = "unknown"  // Set during build
)
```

### Build Flags

These are set during the CI build process:

```bash
go build -ldflags="-s -w -X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.GitCommit=$GIT_COMMIT"
```

## CI/CD Configuration

### GitHub Actions Workflow

The workflow file `.github/workflows/ci.yml` defines:

- **Triggers**: Push to main/develop, pull requests, releases
- **Jobs**: Test, Build, Release
- **Matrix**: Multiple Go versions and platforms
- **Artifacts**: Build outputs and checksums

### Workflow Steps

1. **Test Job**:
   - Matrix testing on multiple platforms
   - Code coverage reporting
   - Race condition detection

2. **Build Job**:
   - Cross-platform compilation
   - Version information injection
   - Checksum generation

3. **Release Job**:
   - Automatic release creation
   - Binary uploads
   - Release notes generation

## Troubleshooting

### Common Issues

#### Build Failures
- Check Go version compatibility
- Verify all dependencies are available
- Check for platform-specific code

#### Test Failures
- Run tests locally: `make test`
- Check for race conditions: `go test -race`
- Verify test data files

#### Release Failures
- Ensure GitHub token has release permissions
- Check tag format (should be semantic versioning)
- Verify workflow file syntax

### Debugging

```bash
# Check workflow runs
gh run list

# View workflow logs
gh run view <run-id>

# Re-run failed workflow
gh run rerun <run-id>
```

## Best Practices

1. **Semantic Versioning**: Use `v1.0.0` format
2. **Release Notes**: Include meaningful descriptions
3. **Testing**: Always test locally before releasing
4. **Documentation**: Update docs with new features
5. **Changelog**: Keep CHANGELOG.md updated

## Security

- Never commit sensitive information
- Use GitHub secrets for API keys
- Verify checksums before installing
- Keep dependencies updated

---

For more information, see the [GitHub Actions documentation](https://docs.github.com/en/actions).
