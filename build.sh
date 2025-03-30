#!/bin/bash

mkdir -p bin

echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o bin/obs-cli-linux-amd64

echo "Building for Windows..."
GOOS=windows GOARCH=amd64 go build -o bin/obs-cli-windows-amd64.exe

echo "Building for macOS..."
GOOS=darwin GOARCH=amd64 go build -o bin/obs-cli-darwin-amd64

echo "Build completed! Binaries are in the bin directory."
