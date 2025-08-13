# LeetSolv Git Workflow Guide

This guide explains the recommended Git branching strategy and workflow for LeetSolv development.

## 🌿 Branch Structure

### Main Branches

#### `main` Branch
- **Purpose**: Production-ready, stable code
- **Content**: Only tested, verified, and released code
- **Protection**: Should be protected, require PR reviews
- **Tags**: Release tags are created from this branch
- **Merges**: Only from `develop` or hotfix branches

#### `develop` Branch
- **Purpose**: Integration branch for features and fixes
- **Content**: Code that has passed basic testing
- **CI/CD**: Automated testing runs on this branch
- **Stability**: Should be relatively stable but may contain bugs
- **Merges**: From feature branches and hotfix branches

### Supporting Branches

#### Feature Branches
- **Naming**: `feature/descriptive-name`
- **Source**: Branch from `develop`
- **Target**: Merge back to `develop`
- **Lifetime**: Short-lived, delete after merge

#### Hotfix Branches
- **Naming**: `hotfix/critical-issue`
- **Source**: Branch from `main`
- **Target**: Merge to both `main` and `develop`
- **Lifetime**: Very short-lived, delete after merge

#### Release Branches (Optional)
- **Naming**: `release/v1.0.0`
- **Source**: Branch from `develop`
- **Target**: Merge to `main` and back to `develop`
- **Purpose**: Final testing and preparation for release

## 🔄 Complete Development Workflow

### 1. Initial Setup

```bash
# Clone the repository
git clone https://github.com/eannchen/leetsolv.git
cd leetsolv

# Set up remote tracking
git checkout -b develop origin/develop
git checkout -b main origin/main
```

### 2. Feature Development

```bash
# Start from develop branch
git checkout develop
git pull origin develop

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
# Source: feature/new-command → Target: develop
```

### 3. Code Review and Integration

```bash
# After PR is approved and merged
git checkout develop
git pull origin develop

# Delete local feature branch
git branch -d feature/new-command

# Delete remote feature branch
git push origin --delete feature/new-command
```

### 4. Testing and Validation

```bash
# CI/CD automatically runs on develop branch
# Run local tests to ensure everything works
make test
make lint
make build

# If issues found, fix them and push to develop
git add .
git commit -m "fix: resolve test failures in new command"
git push origin develop
```

### 5. Release Preparation

```bash
# When ready for release, merge develop to main
git checkout main
git pull origin main
git merge develop
git push origin main

# Create release tag
git tag -a v0.1.0-alpha -m "Alpha release - development version

- New search command functionality
- Improved error handling
- Updated documentation

This is an alpha release for testing purposes."
git push origin v0.1.0-alpha
```

### 6. Hotfix Process (if needed)

```bash
# For critical bugs that need immediate fix
git checkout main
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

# Merge to main
git checkout main
git merge hotfix/critical-bug
git push origin main

# Create hotfix tag
git tag -a v0.1.1 -m "Hotfix release - critical bug fix"
git push origin v0.1.1

# Merge to develop to keep it in sync
git checkout develop
git merge hotfix/critical-bug
git push origin develop

# Delete hotfix branch
git branch -d hotfix/critical-bug
git push origin --delete hotfix/critical-bug
```

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

#### `develop` Branch Protection
- ✅ Require pull request reviews before merging
- ✅ Require status checks to pass before merging
- ✅ Allow force pushes (for development flexibility)

## 🚀 CI/CD Integration

### Branch-Specific Behavior

#### `develop` Branch
- **Tests**: Full test suite runs
- **Builds**: Creates development builds with `-dev` suffix
- **Artifacts**: Uploaded to GitHub Actions artifacts
- **Deployment**: No automatic deployment

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

# Compare branches
git diff develop..main

# See what's in develop but not in main
git log main..develop

# See what's in main but not in develop
git log develop..main
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
- Create PRs for all changes (except hotfixes)
- Use descriptive titles and descriptions
- Include testing instructions
- Request reviews from team members

### 4. Testing
- Always run tests before pushing
- Ensure CI/CD passes before merging
- Test locally before creating releases
- Document known issues in release notes

### 5. Release Management
- Use semantic versioning
- Create detailed release notes
- Test release builds before publishing
- Keep release history clean and organized

## 🚨 Common Pitfalls

### 1. Direct Commits to Main
- ❌ Never commit directly to main
- ✅ Always use pull requests
- ✅ Require code review

### 2. Forgetting to Sync Branches
- ❌ Letting develop and main diverge too much
- ✅ Regular merges from develop to main
- ✅ Keep branches in sync

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

---

**Remember**: This workflow is designed to maintain code quality while enabling rapid development. Adapt it to your team's needs and preferences.
