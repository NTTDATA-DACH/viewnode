#!/usr/bin/env bash
#
# Purpose:
#   Compose repository-local AGENTS guidance after ai-tooling syncs.
#
# Usage:
#   ./ai/compose-agents-md.sh [target-repo]
#
# Parameters:
#   target-repo  Optional repository root containing AGENTS.md and
#                AGENTS.local.md. Defaults to the current directory.
#
# Notes:
#   If AGENTS.local.md is absent, the script exits without changing anything.
#   Existing generated local guidance between the REPO-LOCAL-AGENTS markers is
#   replaced so repeated runs stay idempotent.

set -euo pipefail

START_MARKER="<!-- REPO-LOCAL-AGENTS START -->"
END_MARKER="<!-- REPO-LOCAL-AGENTS END -->"

REPO_ROOT="${1:-$(pwd)}"
AGENTS_PATH="$REPO_ROOT/AGENTS.md"
LOCAL_PATH="$REPO_ROOT/AGENTS.local.md"

if [ ! -f "$AGENTS_PATH" ] || [ ! -f "$LOCAL_PATH" ]; then
    exit 0
fi

TMP_DIR="$(mktemp -d)"
cleanup() {
    rm -rf "$TMP_DIR"
}
trap cleanup EXIT

trim_trailing_blank_lines() {
    local input_path="$1"
    local output_path="$2"

    awk '
        { lines[NR] = $0 }
        END {
            last = NR
            while (last > 0 && lines[last] ~ /^[[:space:]]*$/) {
                last--
            }
            for (i = 1; i <= last; i++) {
                print lines[i]
            }
        }
    ' "$input_path" > "$output_path"
}

awk -v start="$START_MARKER" -v end="$END_MARKER" '
    $0 == start {
        skip = 1
        next
    }
    skip && $0 == end {
        skip = 0
        next
    }
    !skip {
        print
    }
' "$AGENTS_PATH" > "$TMP_DIR/base.raw"

trim_trailing_blank_lines "$TMP_DIR/base.raw" "$TMP_DIR/base"
trim_trailing_blank_lines "$LOCAL_PATH" "$TMP_DIR/local"

{
    cat "$TMP_DIR/base"
    printf '\n\n%s\n' "$START_MARKER"
    cat "$TMP_DIR/local"
    printf '\n%s\n' "$END_MARKER"
} > "$AGENTS_PATH"
