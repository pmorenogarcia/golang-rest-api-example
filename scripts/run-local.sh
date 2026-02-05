#!/bin/bash

# Load environment variables from .env if it exists
if [ -f .env ]; then
    echo "Loading environment variables from .env..."
    export $(cat .env | grep -v '^#' | xargs)
fi

# Run the application
echo "Starting Pokemon REST API..."
go run cmd/api/main.go
