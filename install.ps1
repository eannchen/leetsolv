# LeetSolv Installation Script for Windows (PowerShell)
# This script installs LeetSolv CLI application

param(
    [switch]$Help,
    [switch]$Version,
    [switch]$Uninstall
)

# Configuration
$RepoOwner = "eannchen"  # Change this to your GitHub username
$RepoName = "leetsolv"
$BinaryName = "leetsolv.exe"
$InstallDir = "$env:USERPROFILE\AppData\Local\Programs\leetsolv"
$BackupDir = "$env:USERPROFILE\.leetsolv\backup"

# Colors
$Blue = "`e[94m"
$Green = "`e[92m"
$Yellow = "`e[93m"
$Red = "`e[91m"
$NC = "`e[0m"

# Function to print colored output
function Write-Status {
    param([string]$Message)
    Write-Host "$Blue[INFO]$NC $Message"
}

function Write-Success {
    param([string]$Message)
    Write-Host "$Green[SUCCESS]$NC $Message"
}

function Write-Warning {
    param([string]$Message)
    Write-Host "$Yellow[WARNING]$NC $Message"
}

function Write-Error {
    param([string]$Message)
    Write-Host "$Red[ERROR]$NC $Message"
}

# Function to detect architecture
function Get-Architecture {
    if ($env:PROCESSOR_ARCHITECTURE -eq "AMD64") {
        return "amd64"
    } elseif ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") {
        return "arm64"
    } else {
        return "amd64"
    }
}

# Function to backup existing installation
function Backup-Existing {
    if (Get-Command $BinaryName -ErrorAction SilentlyContinue) {
        Write-Status "Backing up existing installation..."
        if (!(Test-Path $BackupDir)) {
            New-Item -ItemType Directory -Path $BackupDir -Force | Out-Null
        }
        $BackupPath = Join-Path $BackupDir "$(Get-Date -Format 'yyyyMMdd_HHmmss')_$BinaryName"
        Copy-Item (Get-Command $BinaryName).Source $BackupPath
        Write-Success "Backup created at: $BackupPath"
    }
}

# Function to download latest release
function Download-Release {
    Write-Status "Checking for latest release..."

    try {
        $Response = Invoke-RestMethod -Uri "https://api.github.com/repos/$RepoOwner/$RepoName/releases/latest"
        $LatestRelease = $Response.tag_name
        Write-Status "Latest release: $LatestRelease"
    } catch {
        Write-Warning "Could not determine latest release. Using 'latest' tag."
        $LatestRelease = "latest"
    }

    $Arch = Get-Architecture
    $DownloadUrl = "https://github.com/$RepoOwner/$RepoName/releases/download/$LatestRelease/$BinaryName-windows-$Arch.exe"

    Write-Status "Downloading from: $DownloadUrl"

    # Create temporary directory
    $TempDir = Join-Path $env:TEMP "leetsolv_install"
    if (Test-Path $TempDir) {
        Remove-Item $TempDir -Recurse -Force
    }
    New-Item -ItemType Directory -Path $TempDir -Force | Out-Null

    # Download binary
    $OutputPath = Join-Path $TempDir $BinaryName
    try {
        Invoke-WebRequest -Uri $DownloadUrl -OutFile $OutputPath
        Write-Success "Download completed"
    } catch {
        Write-Error "Download failed: $($_.Exception.Message)"
        exit 1
    }
}

# Function to install binary
function Install-Binary {
    Write-Status "Installing binary..."

    # Create install directory
    if (!(Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }

    # Copy binary
    $SourcePath = Join-Path $env:TEMP "leetsolv_install\$BinaryName"
    Copy-Item $SourcePath $InstallDir -Force

    # Add to PATH if not already there
    $CurrentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($CurrentPath -notlike "*$InstallDir*") {
        Write-Status "Adding to PATH..."
        $NewPath = "$CurrentPath;$InstallDir"
        [Environment]::SetEnvironmentVariable("PATH", $NewPath, "User")
        Write-Warning "PATH updated. Please restart your terminal for changes to take effect."
    }

    Write-Success "Binary installed to $InstallDir"
}

# Function to verify installation
function Test-Installation {
    Write-Status "Verifying installation..."

    $FullPath = Join-Path $InstallDir $BinaryName
    if (Test-Path $FullPath) {
        Write-Success "LeetSolv installed successfully!"
        Write-Status "Location: $FullPath"
        Write-Status "You can now run: $BinaryName"
    } else {
        Write-Error "Installation verification failed"
        exit 1
    }
}

# Function to create configuration directory
function Setup-Config {
    Write-Status "Setting up configuration directory..."
    $ConfigDir = "$env:USERPROFILE\.leetsolv"
    if (!(Test-Path $ConfigDir)) {
        New-Item -ItemType Directory -Path $ConfigDir -Force | Out-Null
    }
    Write-Success "Configuration directory created at $ConfigDir"
}

# Function to show post-install instructions
function Show-PostInstall {
    Write-Host ""
    Write-Success "Installation completed successfully!"
    Write-Host ""
    Write-Host "Next steps:"
    Write-Host "1. Restart your terminal to update PATH"
    Write-Host "2. Run '$BinaryName' to start the application"
    Write-Host "3. Run '$BinaryName help' to see available commands"
    Write-Host "4. Configuration files will be created in $env:USERPROFILE\.leetsolv\"
    Write-Host ""
    Write-Host "To uninstall, run: Remove-Item '$InstallDir' -Recurse -Force"
    Write-Host ""
}

# Function to uninstall
function Uninstall-LeetSolv {
    Write-Status "Uninstalling LeetSolv..."

    if (Test-Path $InstallDir) {
        Remove-Item $InstallDir -Recurse -Force
        Write-Success "LeetSolv uninstalled successfully"
    } else {
        Write-Warning "LeetSolv not found"
    }
}

# Main installation function
function Install-LeetSolv {
    Write-Host "$Blue╭───────────────────────────────────────────────────╮$NC"
    Write-Host "$Blue│                                                   │$NC"
    Write-Host "$Blue│    ░▒▓   LeetSolv — CLI SRS for LeetCode   ▓▒░    │$NC"
    Write-Host "$Blue│                                                   │$NC"
    Write-Host "$Blue│                Installation Script               │$NC"
    Write-Host "$Blue│                                                   │$NC"
    Write-Host "$Blue╰───────────────────────────────────────────────────╯$NC"
    Write-Host ""

    # Check prerequisites
    $Arch = Get-Architecture

    # Installation steps
    Backup-Existing
    Download-Release
    Install-Binary
    Setup-Config
    Test-Installation
    Show-PostInstall
}

# Handle command line arguments
if ($Help) {
    Write-Host "Usage: .\install.ps1 [OPTIONS]"
    Write-Host "Options:"
    Write-Host "  -Help         Show this help message"
    Write-Host "  -Version      Show version information"
    Write-Host "  -Uninstall    Uninstall LeetSolv"
    exit 0
}

if ($Version) {
    Write-Host "LeetSolv Installer v1.0.0"
    exit 0
}

if ($Uninstall) {
    Uninstall-LeetSolv
    exit 0
}

# Check execution policy
$ExecutionPolicy = Get-ExecutionPolicy
if ($ExecutionPolicy -eq "Restricted") {
    Write-Warning "PowerShell execution policy is restricted."
    Write-Host "To run this script, you may need to change the execution policy:"
    Write-Host "Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser"
    Write-Host "Or run: Set-ExecutionPolicy -ExecutionPolicy Bypass -Scope Process"
    exit 1
}

# Run installation
Install-LeetSolv
