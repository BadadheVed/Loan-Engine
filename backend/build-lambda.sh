#!/bin/bash
set -e

echo "Building Lambda functions with shared dependencies..."

cd lambda-functions
echo "Downloading dependencies..."
go mod download

echo "Building health function..."
cd health
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap main.go
cd ..
echo "Building uploadcsv function..."
cd uploadcsv
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap main.go
cd ..

cd ..

echo "Build complete!"
