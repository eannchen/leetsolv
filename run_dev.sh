#!/bin/bash

# Development mode script for leetsolv
# This script sets environment variables to use test files and runs the application

export LEETSOLV_QUESTIONS_FILE="questions.dev.json"
export LEETSOLV_DELTAS_FILE="deltas.dev.json"
export LEETSOLV_INFO_LOG_FILE="info.dev.log"
export LEETSOLV_ERROR_LOG_FILE="error.dev.log"
export LEETSOLV_SETTINGS_FILE="settings.dev.json"

echo "Running leetsolv in DEVELOPMENT mode with files:"
echo "  Questions: $LEETSOLV_QUESTIONS_FILE"
echo "  Deltas: $LEETSOLV_DELTAS_FILE"
echo "  Info Log: $LEETSOLV_INFO_LOG_FILE"
echo "  Error Log: $LEETSOLV_ERROR_LOG_FILE"
echo "  Settings: $LEETSOLV_SETTINGS_FILE"
echo ""

# Run the application with any provided arguments
go run . "$@"