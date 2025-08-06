#!/bin/bash

# Main run script for leetsolv
# Allows choosing between test and production modes

echo "LeetSolv - Choose your mode:"
echo "1) Test mode (safe for development)"
echo "2) Production mode (uses real data)"
echo "3) Exit"
echo ""

read -p "Enter your choice (1-3): " choice

case $choice in
    1)
        echo "Starting in TEST mode..."
        ./run_test.sh "$@"
        ;;
    2)
        echo "Starting in PRODUCTION mode..."
        ./run_prod.sh "$@"
        ;;
    3)
        echo "Exiting..."
        exit 0
        ;;
    *)
        echo "Invalid choice. Exiting..."
        exit 1
        ;;
esac