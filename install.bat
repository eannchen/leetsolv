@echo off
setlocal enabledelayedexpansion

REM LeetSolv Installation Script for Windows
REM This script installs LeetSolv CLI application

REM Configuration
set "REPO_OWNER=eannchen"
set "REPO_NAME=leetsolv"
set "BINARY_NAME=leetsolv.exe"
set "INSTALL_DIR=%USERPROFILE%\AppData\Local\Microsoft\WinGet\Packages"
set "BACKUP_DIR=%USERPROFILE%\.leetsolv\backup"

REM Colors (Windows 10+)
set "BLUE=[94m"
set "GREEN=[92m"
set "YELLOW=[93m"
set "RED=[91m"
set "NC=[0m"

REM Function to print colored output
:print_status
echo %BLUE%[INFO]%NC% %~1
goto :eof

:print_success
echo %GREEN%[SUCCESS]%NC% %~1
goto :eof

:print_warning
echo %YELLOW%[WARNING]%NC% %~1
goto :eof

:print_error
echo %RED%[ERROR]%NC% %~1
goto :eof

REM Function to detect architecture
:detect_architecture
if "%PROCESSOR_ARCHITECTURE%"=="AMD64" (
    set "ARCH=amd64"
) else if "%PROCESSOR_ARCHITECTURE%"=="ARM64" (
    set "ARCH=arm64"
) else (
    set "ARCH=amd64"
)
call :print_status "Detected architecture: %ARCH%"
goto :eof

REM Function to check if command exists
:command_exists
where %1 >nul 2>&1
if %errorlevel% equ 0 (
    set "EXISTS=1"
) else (
    set "EXISTS=0"
)
goto :eof

REM Function to backup existing installation
:backup_existing
call :command_exists "%BINARY_NAME%"
if "%EXISTS%"=="1" (
    call :print_status "Backing up existing installation..."
    if not exist "%BACKUP_DIR%" mkdir "%BACKUP_DIR%"
    set "BACKUP_PATH=%BACKUP_DIR%\%date:~-4,4%%date:~-10,2%%date:~-7,2%_%time:~0,2%%time:~3,2%%time:~6,2%_%BINARY_NAME%"
    set "BACKUP_PATH=%BACKUP_PATH: =0%"
    copy "%BINARY_NAME%" "%BACKUP_PATH%" >nul
    call :print_success "Backup created at: %BACKUP_PATH%"
)
goto :eof

REM Function to download latest release
:download_release
call :print_status "Checking for latest release..."

REM Try to get latest release from GitHub API
set "LATEST_RELEASE=latest"
if exist "powershell.exe" (
    for /f "tokens=*" %%i in ('powershell -Command "(Invoke-WebRequest -Uri 'https://api.github.com/repos/%REPO_OWNER%/%REPO_NAME%/releases/latest' -UseBasicParsing).Content | ConvertFrom-Json | Select-Object -ExpandProperty tag_name"') do set "LATEST_RELEASE=%%i"
)

call :print_status "Latest release: %LATEST_RELEASE%"

REM Download URL
set "DOWNLOAD_URL=https://github.com/%REPO_OWNER%/%REPO_NAME%/releases/download/%LATEST_RELEASE%/%BINARY_NAME%-windows-%ARCH%.exe"

call :print_status "Downloading from: %DOWNLOAD_URL%"

REM Create temporary directory
set "TEMP_DIR=%TEMP%\leetsolv_install"
if exist "%TEMP_DIR%" rmdir /s /q "%TEMP_DIR%"
mkdir "%TEMP_DIR%"
cd /d "%TEMP_DIR%"

REM Download binary using PowerShell
if exist "powershell.exe" (
    powershell -Command "Invoke-WebRequest -Uri '%DOWNLOAD_URL%' -OutFile '%BINARY_NAME%'"
) else (
    call :print_error "PowerShell not found. Please install PowerShell or download manually from: %DOWNLOAD_URL%"
    exit /b 1
)

if not exist "%BINARY_NAME%" (
    call :print_error "Download failed"
    exit /b 1
)

call :print_success "Download completed"
goto :eof

REM Function to install binary
:install_binary
call :print_status "Installing binary..."

REM Try to install to user directory first
set "USER_INSTALL_DIR=%USERPROFILE%\AppData\Local\Programs\leetsolv"
if not exist "%USER_INSTALL_DIR%" mkdir "%USER_INSTALL_DIR%"

copy "%BINARY_NAME%" "%USER_INSTALL_DIR%\" >nul

REM Add to PATH if not already there
set "PATH_ENTRY=%USER_INSTALL_DIR%"
echo %PATH% | find /i "%PATH_ENTRY%" >nul
if %errorlevel% neq 0 (
    call :print_status "Adding to PATH..."
    setx PATH "%PATH%;%PATH_ENTRY%" >nul
    call :print_warning "PATH updated. Please restart your terminal for changes to take effect."
)

call :print_success "Binary installed to %USER_INSTALL_DIR%"
goto :eof

REM Function to verify installation
:verify_installation
call :print_status "Verifying installation..."

set "FULL_PATH=%USER_INSTALL_DIR%\%BINARY_NAME%"
if exist "%FULL_PATH%" (
    call :print_success "LeetSolv installed successfully!"
    call :print_status "Location: %FULL_PATH%"
    call :print_status "You can now run: %BINARY_NAME%"
) else (
    call :print_error "Installation verification failed"
    exit /b 1
)
goto :eof

REM Function to create configuration directory
:setup_config
call :print_status "Setting up configuration directory..."
if not exist "%USERPROFILE%\.leetsolv" mkdir "%USERPROFILE%\.leetsolv"
call :print_success "Configuration directory created at %USERPROFILE%\.leetsolv"
goto :eof

REM Function to show post-install instructions
:show_post_install
echo.
call :print_success "Installation completed successfully!"
echo.
echo Next steps:
echo 1. Restart your terminal to update PATH
echo 2. Run '%BINARY_NAME%' to start the application
echo 3. Run '%BINARY_NAME% help' to see available commands
echo 4. Configuration files will be created in %USERPROFILE%\.leetsolv\
echo.
echo To uninstall, run: rmdir /s /q "%USER_INSTALL_DIR%"
echo.
goto :eof

REM Main installation function
:main
echo %BLUE%╭───────────────────────────────────────────────────╮%NC%
echo %BLUE%│                                                   │%NC%
echo %BLUE%│    ░▒▓   LeetSolv — CLI SRS for LeetCode   ▓▒░    │%NC%
echo %BLUE%│                                                   │%NC%
echo %BLUE%│                Installation Script               │%NC%
echo %BLUE%│                                                   │%NC%
echo %BLUE%╰───────────────────────────────────────────────────╯%NC%
echo.

REM Check prerequisites
call :detect_architecture

REM Installation steps
call :backup_existing
call :download_release
call :install_binary
call :setup_config
call :verify_installation
call :show_post_install
goto :eof

REM Handle command line arguments
if "%1"=="--help" goto :help
if "%1"=="-h" goto :help
if "%1"=="--version" goto :version
if "%1"=="-v" goto :version
if "%1"=="--uninstall" goto :uninstall
if "%1"=="" goto :main
goto :unknown_option

:help
echo Usage: %0 [OPTIONS]
echo Options:
echo   --help, -h     Show this help message
echo   --version, -v  Show version information
echo   --uninstall    Uninstall LeetSolv
exit /b 0

:version
echo LeetSolv Installer v1.0.0
exit /b 0

:uninstall
call :print_status "Uninstalling LeetSolv..."
set "USER_INSTALL_DIR=%USERPROFILE%\AppData\Local\Programs\leetsolv"
if exist "%USER_INSTALL_DIR%" (
    rmdir /s /q "%USER_INSTALL_DIR%"
    call :print_success "LeetSolv uninstalled successfully"
) else (
    call :print_warning "LeetSolv not found"
)
exit /b 0

:unknown_option
call :print_error "Unknown option: %1"
echo Use --help for usage information
exit /b 1
