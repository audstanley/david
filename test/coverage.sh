#!/bin/bash
# Run all tests with coverage reporting

set -e

echo "Running tests with coverage..."

# Run all tests
go test ./... -race -coverprofile=coverage.out -covermode=atomic

# Generate coverage report
go tool cover -html=coverage.out -o coverage.html

echo ""
echo "Coverage report generated: coverage.html"
echo ""
echo "Overall coverage:"
go tool cover -func=coverage.out | grep total:

# Clean up
rm -f coverage.out
