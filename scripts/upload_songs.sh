#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCRIPT_PATH="$SCRIPT_DIR/upload_songs.sh"
DEFAULT_CONFIG="$SCRIPT_DIR/upload_songs.config.sh"

print_usage() {
  cat <<'EOF'
Usage:
  scripts/upload_songs.sh [config_file]

Examples:
  scripts/upload_songs.sh
  scripts/upload_songs.sh scripts/upload_songs.config.sh

Notes:
  - Set AUTH_TOKEN in the config file before running.
  - TOTAL_UPLOADS and CONCURRENCY control load volume.
EOF
}

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  print_usage
  exit 0
fi

run_single_upload() {
  local index="$1"

  local title="${TITLE_PREFIX}-${index}"
  local response_file
  response_file="$(mktemp)"

  local http_code
  http_code="$(curl -sS -o "$response_file" -w "%{http_code}" \
    -X POST "$API_URL" \
    -H "Authorization: Bearer ${AUTH_TOKEN}" \
    -F "file=@${FILE_PATH}" \
    -F "title=${title}" \
    -F "artist=${ARTIST}" \
    -F "album=${ALBUM}" \
    -F "genre=${GENRE}" \
    -F "duration=${DURATION}")"

  if [[ "$http_code" -ge 200 && "$http_code" -lt 300 ]]; then
    echo "OK index=${index} status=${http_code} title=${title}"
    rm -f "$response_file"
    return 0
  fi

  local body
  body="$(tr '\n' ' ' < "$response_file")"
  echo "FAIL index=${index} status=${http_code} body=${body}"
  rm -f "$response_file"
  return 1
}

if [[ "${1:-}" == "--single" ]]; then
  config_file="${2:-$DEFAULT_CONFIG}"
  index="${3:-}"

  if [[ -z "$index" ]]; then
    echo "missing index in --single mode"
    exit 2
  fi

  # shellcheck disable=SC1090
  source "$config_file"
  run_single_upload "$index"
  exit 0
fi

CONFIG_FILE="${1:-$DEFAULT_CONFIG}"

if [[ ! -f "$CONFIG_FILE" ]]; then
  echo "config file not found: $CONFIG_FILE"
  exit 1
fi

# shellcheck disable=SC1090
source "$CONFIG_FILE"

if ! command -v curl >/dev/null 2>&1; then
  echo "curl is required"
  exit 1
fi

if [[ -z "${AUTH_TOKEN:-}" ]]; then
  echo "AUTH_TOKEN is empty in config: $CONFIG_FILE"
  exit 1
fi

if [[ ! -f "$FILE_PATH" ]]; then
  echo "FILE_PATH not found: $FILE_PATH"
  exit 1
fi

if [[ "$TOTAL_UPLOADS" -lt 1 ]]; then
  echo "TOTAL_UPLOADS must be >= 1"
  exit 1
fi

if [[ "$CONCURRENCY" -lt 1 ]]; then
  echo "CONCURRENCY must be >= 1"
  exit 1
fi

echo "Starting uploads"
echo "  API_URL=$API_URL"
echo "  FILE_PATH=$FILE_PATH"
echo "  TOTAL_UPLOADS=$TOTAL_UPLOADS"
echo "  CONCURRENCY=$CONCURRENCY"

LOG_FILE="$(mktemp)"

set +e
seq 1 "$TOTAL_UPLOADS" | xargs -n 1 -P "$CONCURRENCY" "$SCRIPT_PATH" --single "$CONFIG_FILE" 2>&1 | tee "$LOG_FILE"
XARGS_STATUS=${PIPESTATUS[1]}
set -e

SUCCESS_COUNT="$(grep -c '^OK ' "$LOG_FILE" || true)"
FAIL_COUNT="$(grep -c '^FAIL ' "$LOG_FILE" || true)"

rm -f "$LOG_FILE"

echo "Upload summary: success=$SUCCESS_COUNT fail=$FAIL_COUNT"

if [[ "$XARGS_STATUS" -ne 0 || "$FAIL_COUNT" -gt 0 ]]; then
  exit 1
fi

exit 0
