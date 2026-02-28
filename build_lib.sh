#!/bin/bash
# Build the Go code as a shared library

# Ensure we're in the go-html-to-md directory
cd "$(dirname "$0")"

# Build the shared library
echo "Building shared library..."
go build -buildmode=c-shared -o libhtmltomd.so html-to-markdown-lib.go

# Check if the build was successful
if [ $? -eq 0 ]; then
    echo "Shared library built successfully."
    echo "Output: libhtmltomd.so"
else
    echo "Failed to build shared library."
    exit 1
fi