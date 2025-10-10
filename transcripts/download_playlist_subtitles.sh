#!/bin/bash

# Script to download subtitles from YouTube playlist
# Usage: ./download_playlist_subtitles.sh <playlist_url> [output_dir]

# Set variables
PLAYLIST_URL="$1"
OUTPUT_DIR="${2:-playlist-subtitles}"
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}YouTube Playlist Subtitle Downloader${NC}"
echo "======================================"

# Check if playlist URL is provided
if [ -z "$PLAYLIST_URL" ]; then
    echo -e "${RED}Error: Please provide a YouTube playlist URL${NC}"
    echo "Usage: $0 <playlist_url> [output_dir]"
    echo "Example: $0 'https://youtube.com/playlist?list=PLxxx' my_subtitles"
    exit 1
fi

# Check if yt-dlp is installed
if ! command -v yt-dlp &> /dev/null; then
    echo -e "${RED}Error: yt-dlp is not installed or not in PATH${NC}"
    echo "Please install yt-dlp first: pip install yt-dlp"
    exit 1
fi

# Create output directory if it doesn't exist
cd "$SCRIPT_DIR"
if [ ! -d "$OUTPUT_DIR" ]; then
    echo -e "${YELLOW}Creating directory: $OUTPUT_DIR${NC}"
    mkdir -p "$OUTPUT_DIR"
fi

cd "$OUTPUT_DIR"

echo -e "${GREEN}Starting subtitle download from playlist...${NC}"
echo "Playlist URL: $PLAYLIST_URL"
echo "Output directory: $(pwd)"
echo ""

# Download subtitles with robust options
yt-dlp \
    --skip-download \
    --write-sub \
    --write-auto-sub \
    --sub-lang "en" \
    --convert-subs "srt" \
    --no-warnings \
    --ignore-errors \
    --no-abort-on-error \
    --output "%(playlist_index)03d-%(title)s.%(ext)s" \
    --restrict-filenames \
    "$PLAYLIST_URL"

# Check if any files were downloaded
SRT_COUNT=$(find . -name "*.srt" -type f 2>/dev/null | wc -l)
VTT_COUNT=$(find . -name "*.vtt" -type f 2>/dev/null | wc -l)

echo ""
echo "======================================"
if [ $SRT_COUNT -gt 0 ] || [ $VTT_COUNT -gt 0 ]; then
    echo -e "${GREEN}Download completed successfully!${NC}"
    echo "SRT files downloaded: $SRT_COUNT"
    echo "VTT files downloaded: $VTT_COUNT"
    echo ""
    echo "Files are located in: $(pwd)"

    # List downloaded files
    echo ""
    echo "Downloaded subtitle files:"
    echo "--------------------------"
    ls -la *.srt 2>/dev/null || ls -la *.vtt 2>/dev/null || echo "No subtitle files found"
else
    echo -e "${YELLOW}Warning: No subtitle files were downloaded${NC}"
    echo ""
    echo "Possible reasons:"
    echo "1. Videos might not have subtitles available"
    echo "2. YouTube might be blocking the requests"
    echo "3. yt-dlp might need updating"
    echo ""
    echo "Try running with alternative command:"
    echo -e "${YELLOW}yt-dlp --list-subs \"$PLAYLIST_URL\"${NC}"
fi

echo ""
echo "Script completed."