#!/bin/bash
set -euo pipefail
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m'

echo -e "${YELLOW}Building the Lambda layer for ARM64 architecture...${NC}"

# Create a directory structure for the Lambda layer
mkdir -p build/layer/bin
rm -rf build/ffmpeg*
cd build

# Download and extract the ffmpeg binary for ARM64 architecture
echo -e "${YELLOW}Downloading and extracting ffmpeg for ARM64 architecture...${NC}"
curl -sL https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-arm64-static.tar.xz | tar -xJ
echo -e "${GREEN}Download and extraction complete!${NC}"

# Check if the extracted directory exists
echo -e "${YELLOW}Moving ffmpeg and ffprobe to the layer directory...${NC}"
EXTRACTED_DIR=$(find . -maxdepth 1 -type d -name "ffmpeg-*static*" | head -n 1)
if [[ -d "$EXTRACTED_DIR" ]]; then
  mv "$EXTRACTED_DIR/ffmpeg" "$EXTRACTED_DIR/ffprobe" layer/bin/
else
  echo -e "${RED}Error: Extracted ffmpeg directory not found!${NC}"
  exit 1
fi
echo -e "${GREEN}ffmpeg and ffprobe moved successfully!${NC}"

rm -rf ${EXTRACTED_DIR}

# Zip the layer
echo -e "${YELLOW}Creating zip file for the Lambda layer...${NC}"
cd layer
zip -r ../ffmpeg.zip bin

# Move the zip file to the root directory
cd ..
mv ffmpeg.zip ../ffmpeg.zip
echo -e "${GREEN}Lambda layer zip file created successfully: ffmpeg.zip${NC}"

# Clean up the build directory
echo -e "${YELLOW}Cleaning up build directory...${NC}"
cd ..
rm -rf build
echo -e "${GREEN}Build and packaging of the Lambda layer completed successfully!${NC}"
