#!/bin/bash

# Generate Swagger documentation
echo "Generating Swagger documentation..."

swag init -g cmd/api/main.go -o docs/swagger --parseDependency --parseInternal

if [ $? -eq 0 ]; then
    echo "✓ Swagger documentation generated successfully"
    echo "  - docs/swagger/docs.go"
    echo "  - docs/swagger/swagger.json"
    echo "  - docs/swagger/swagger.yaml"
else
    echo "✗ Failed to generate Swagger documentation"
    exit 1
fi
