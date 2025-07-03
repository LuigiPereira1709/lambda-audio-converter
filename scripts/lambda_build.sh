#!/bin/bash
set -euo pipefail

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m'

ARCH="arm64"
GOOS="linux"
echo -e "${YELLOW}Building the lambda function for architecture: $GOOS-$ARCH...${NC}"
env GOOS=$GOOS GOARCH=$ARCH go build -tags lambda.norpc -o bootstrap main.go
echo -e "${GREEN}Build successful!${NC}"

ZIP_NAME="lambda.zip"

echo -e "${YELLOW}Creating zip file: ${ZIP_NAME}...${NC}"
zip -j "$ZIP_NAME" bootstrap
echo -e "${GREEN}Zip file created successfully: ${ZIP_NAME}${NC}"

echo -e "${YELLOW}Cleaning up bootstrap file...${NC}"
rm bootstrap
echo -e "${GREEN}Cleanup done!${NC}"

echo -e "${GREEN}Lambda function build and packaging completed successfully!${NC}"
