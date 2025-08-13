@echo off
setlocal enabledelayedexpansion

REM LeetSolv Uninstall Script for Windows
REM This script removes LeetSolv CLI application

REM Configuration
set "BINARY_NAME=leetsolv.exe"
set "USER_INSTALL_DIR=%USERPROFILE%\AppData\Local\Programs\leetsolv"
set "CONFIG_DIR=%USERPROFILE%\.leetsolv"
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

REM Function to check if LeetSolv is running
:check_running
tasklist /FI "IMAGENAME eq %BINARY_NAME%" 2>NUL | find /I /N "%BINARY_NAME%" >NUL
if %errorlevel% equ 0 (
    call :print_warning "LeetSolv is currently running. Please close it before uninstalling."
    set /p "CONTINUE=Continue anyway? (y/N): "
    if /i not "!CONTINUE!"=="y" (
        call :print_status "Uninstallation cancelled"
        exit /b 0
    )
)
goto :eof

REM Function to find LeetSolv installation
:find_installation
set "FOUND_PATH="
if exist "%USER_INSTALL_DIR%\%BINARY_NAME%" (
    set "FOUND_PATH=%USER_INSTALL_DIR%\%BINARY_NAME%"
) else (
    where %BINARY_NAME% >nul 2>&1
    if %errorlevel% equ 0 (
        for /f "tokens=*" %%i in ('where %BINARY_NAME%') do set "FOUND_PATH=%%i"
    )
)
goto :eof

REM Function to backup before uninstalling
:backup_binary
if exist "%~1" (
    call :print_status "Creating backup before uninstallation..."
    if not exist "%BACKUP_DIR%" mkdir "%BACKUP_DIR%"
    set "BACKUP_PATH=%BACKUP_DIR%\%date:~-4,4%%date:~-10,2%%date:~-7,2%_%time:~0,2%%time:~3,2%%time:~6,2%_%BINARY_NAME%"
    set "BACKUP_PATH=%BACKUP_PATH: =0%"
    copy "%~1" "%BACKUP_PATH%" >nul
    call :print_success "Backup created at: %BACKUP_PATH%"
)
goto :eof

REM Function to remove binary
:remove_binary
if exist "%~1" (
    call :print_status "Removing binary: %~1"
    del /f /q "%~1" >nul 2>&1
    if exist "%~1" (
        call :print_error "Failed to remove binary. You may need administrator privileges."
        call :print_status "Please remove manually: %~1"
    ) else (
        call :print_success "Binary removed successfully"
    )
)
goto :eof

REM Function to remove configuration
:remove_config
if exist "%CONFIG_DIR%" (
    call :print_status "Removing configuration directory: %CONFIG_DIR%"
    set /p "KEEP_CONFIG=Do you want to keep your configuration files? (y/N): "
    if /i not "!KEEP_CONFIG!"=="y" (
        rmdir /s /q "%CONFIG_DIR%"
        call :print_success "Configuration directory removed"
    ) else (
        call :print_status "Configuration directory kept at: %CONFIG_DIR%"
    )
) else (
    call :print_status "No configuration directory found"
)
goto :eof

REM Function to remove from PATH
:remove_from_path
if exist "%~1" (
    set "INSTALL_DIR=%~dp1"
    set "INSTALL_DIR=%INSTALL_DIR:~0,-1%"

    call :print_status "Checking PATH entries..."

    REM Check if the directory is in PATH
    echo %PATH% | find /i "%INSTALL_DIR%" >nul
    if %errorlevel% equ 0 (
        call :print_warning "Found PATH entry for: %INSTALL_DIR%"
        call :print_status "You may need to manually remove this from your PATH environment variable"
        call :print_status "Or restart your computer to clear the PATH"
    )
)
goto :eof

REM Function to show uninstall summary
:show_summary
echo.
call :print_success "Uninstallation completed!"
echo.
echo Summary of actions:
echo ✓ Binary removed
echo ✓ Configuration handled
echo ✓ PATH entries checked
echo.
echo If you want to reinstall later, run:
echo curl -fsSL https://raw.githubusercontent.com/eannchen/leetsolv/main/install.sh ^| bash
echo.
echo Backup files are stored in: %BACKUP_DIR%
goto :eof

REM Function to clean up empty directories
:cleanup_directories
if exist "%USER_INSTALL_DIR%" (
    dir "%USER_INSTALL_DIR%" /b >nul 2>&1
    if %errorlevel% equ 0 (
        call :print_status "Checking if install directory is empty..."
        dir "%USER_INSTALL_DIR%" /b | findstr /r "^" >nul
        if %errorlevel% neq 0 (
            call :print_status "Removing empty install directory..."
            rmdir "%USER_INSTALL_DIR%"
        )
    )
)
goto :eof

REM Main uninstall function
:main
echo %BLUE%╭───────────────────────────────────────────────────╮%NC%
echo %BLUE%│                                                   │%NC%
echo %BLUE%│    ░▒▓   LeetSolv — Uninstall Script        ▓▒░    │%NC%
echo %BLUE%│                                                   │%NC%
echo %BLUE%╰───────────────────────────────────────────────────╯%NC%
echo.

REM Check if LeetSolv is running
call :check_running

REM Find installation
call :print_status "Searching for LeetSolv installation..."
call :find_installation

if "%FOUND_PATH%"=="" (
    call :print_error "LeetSolv not found. It may already be uninstalled."
    exit /b 1
)

call :print_status "Found installation at: %FOUND_PATH%"

REM Confirm uninstallation
echo.
call :print_warning "About to uninstall LeetSolv from: %FOUND_PATH%"
set /p "CONFIRM=Are you sure you want to continue? (y/N): "
if /i not "%CONFIRM%"=="y" (
    call :print_status "Uninstallation cancelled"
    exit /b 0
)

REM Perform uninstallation
call :backup_binary "%FOUND_PATH%"
call :remove_binary "%FOUND_PATH%"
call :remove_from_path "%FOUND_PATH%"
call :remove_config
call :cleanup_directories

call :show_summary
goto :eof

REM Handle command line arguments
if "%1"=="--help" goto :help
if "%1"=="-h" goto :help
if "%1"=="--force" goto :force
if "%1"=="-f" goto :force
if "%1"=="--config-only" goto :config_only
if "%1"=="" goto :main
goto :unknown_option

:help
echo Usage: %0 [OPTIONS]
echo Options:
echo   --help, -h     Show this help message
echo   --force, -f    Skip confirmation prompts
echo   --config-only  Remove only configuration files
exit /b 0

:force
REM Skip confirmations (for automation)
set "FORCE=true"
goto :main

:config_only
REM Remove only configuration
call :remove_config
call :print_success "Configuration cleanup completed"
exit /b 0

:unknown_option
call :print_error "Unknown option: %1"
echo Use --help for usage information
exit /b 1
