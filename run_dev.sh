#!/bin/bash

# Test mode script for leetsolv
# This script sets environment variables to use test files and runs the application

export LEETSOLV_QUESTIONS_FILE="questions.test.json"
export LEETSOLV_DELTAS_FILE="deltas.test.json"
export LEETSOLV_INFO_LOG_FILE="info.test.log"
export LEETSOLV_ERROR_LOG_FILE="error.test.log"

echo "Running leetsolv in TEST mode with files:"
echo "  Questions: $LEETSOLV_QUESTIONS_FILE"
echo "  Deltas: $LEETSOLV_DELTAS_FILE"
echo "  Info Log: $LEETSOLV_INFO_LOG_FILE"
echo "  Error Log: $LEETSOLV_ERROR_LOG_FILE"
echo ""

# Run the application with any provided arguments
go run . "$@"