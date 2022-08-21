#!/bin/bash

echo "Building for architecture ${ARCH:-amd64}..."
GOOS=linux GOARCH=${ARCH} go build -o bin/extensions/go-example-extension main.go
chmod +x bin/extensions/go-example-extension
echo "Building finished."
