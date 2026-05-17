#!/usr/bin/env bash

# API target
API_URL="http://localhost:8080/api/v1/songs/upload"

# Required auth token for /api/v1 routes
AUTH_TOKEN=""

# Mock upload file (requested default)
FILE_PATH="/Users/shawonkanji/Documents/Ahonnis_Mukherjee__-_-1.pdf"

# Load settings
TOTAL_UPLOADS=300
CONCURRENCY=20

# Form fields
TITLE_PREFIX="bulk-song"
ARTIST="eminem"
ALBUM="mock-album"
GENRE="rap"
DURATION=180
