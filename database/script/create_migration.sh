#!/usr/bin/env bash

# -------------------------------
# CONFIGURATION
# -------------------------------
# Hardcoded base migrations folder
BASE_MIGRATIONS_DIR="./database/migrations"

# -------------------------------
# SCRIPT START
# -------------------------------

set -e

# Usage check
if [[ -z "$1" || -z "$2" ]]; then
  echo "Usage: $0 <feature_name> <migration_name>"
  echo "Example: $0 rbac add_user_role_mapping"
  exit 1
fi

FEATURE_NAME="$1"
MIGRATION_NAME="$2"

# Sanitize feature name: lowercase, replace non-alphanum with underscore
SANITIZED_FEATURE=$(echo "$FEATURE_NAME" \
  | tr '[:upper:]' '[:lower:]' \
  | sed 's/[^a-z0-9]/_/g')

# Full folder path for this feature
FEATURE_DIR="${BASE_MIGRATIONS_DIR}/${SANITIZED_FEATURE}"

# Create base + feature folder if missing
mkdir -p "$FEATURE_DIR"

# Generate UTC timestamp (YYYYMMDDHHMMSS)
UTC_TIMESTAMP=$(date -u +"%Y%m%d%H%M%S")

# Sanitize migration name: lowercase, replace non-alphanum with underscore
SANITIZED_MIGRATION=$(echo "$MIGRATION_NAME" \
  | tr '[:upper:]' '[:lower:]' \
  | sed 's/[^a-z0-9]/_/g')

# Construct filename
FILEPATH="${FEATURE_DIR}/${UTC_TIMESTAMP}_${SANITIZED_MIGRATION}.sql"

# Prevent overwrite
if [[ -f "$FILEPATH" ]]; then
  echo "Migration already exists: $FILEPATH"
  exit 1
fi

# Create empty migration file
touch "$FILEPATH"

echo "Created migration:"
echo "  $FILEPATH"
