#!/bin/bash

# Production mode script for leetsolv
# This script explicitly sets environment variables to use production files

export LEETSOLV_QUESTIONS_FILE="questions.json"
export LEETSOLV_DELTAS_FILE="deltas.json"
export LEETSOLV_INFO_LOG_FILE="info.log"
export LEETSOLV_ERROR_LOG_FILE="error.log"

echo "Running leetsolv in PRODUCTION mode with files:"
echo "  Questions: $LEETSOLV_QUESTIONS_FILE"
echo "  Deltas: $LEETSOLV_DELTAS_FILE"
echo "  Info Log: $LEETSOLV_INFO_LOG_FILE"
echo "  Error Log: $LEETSOLV_ERROR_LOG_FILE"
echo ""

# Run the application with any provided arguments
go run . "$@"