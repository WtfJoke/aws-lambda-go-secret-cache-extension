#!/bin/bash

echo "Building for architecture ${ARCH:-amd64}..."
GOOS=linux GOARCH=${ARCH} go build -o bin/extensions/secret-cache-extension main.go
chmod +x bin/extensions/secret-cache-extension
echo "Building finished."
