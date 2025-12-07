#!/bin/bash
set -e

OUTPUT_FILE=${1:-"lambda-functions.zip"}

echo "Creating Lambda deployment package..."
echo "Output: $OUTPUT_FILE"

# Build first
./build-lambda.sh

# Create zip of entire lambda-functions folder
echo "Zipping lambda-functions directory..."
cd lambda-functions
zip -r ../"$OUTPUT_FILE" . -x "*.git*" -x "*.DS_Store" -x "*__pycache__*"
cd ..

echo "Package created: $OUTPUT_FILE"
echo "Size: $(du -h "$OUTPUT_FILE" | cut -f1)"
