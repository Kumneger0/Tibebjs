#!/bin/bash

# Script to build the Go runtime
echo "Building JS Runtime..."
go build -o jsruntime ./cmd/main.go
