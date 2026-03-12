#!/bin/bash
# Ralph doctor - sanity checks for running Ralph in a project

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CODEX_BIN="${CODEX_BIN:-codex}"
PRD_FILE="$SCRIPT_DIR/prd.json"
PROMPT_FILE="$SCRIPT_DIR/prompt.md"

fail() {
  echo "ERROR: $1" >&2
  exit 1
}

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || fail "Missing required command: $1"
}

echo "Ralph doctor"
echo "ralph dir: $SCRIPT_DIR"

require_cmd git
require_cmd jq
require_cmd "$CODEX_BIN"

if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  fail "Not inside a git repository. Run this from within your project repo."
fi

if [ ! -f "$PROMPT_FILE" ]; then
  fail "Missing $PROMPT_FILE"
fi

if [ ! -f "$PRD_FILE" ]; then
  echo "WARN: Missing $PRD_FILE"
  echo "      Create it (or copy from prd.json.example) before running ralph.sh."
fi

if ! "$CODEX_BIN" exec --help >/dev/null 2>&1; then
  fail "Codex exec help failed. Check your Codex installation."
fi

if "$CODEX_BIN" --yolo exec --help 2>&1 | grep -qi "unexpected argument '--yolo'"; then
  echo "WARN: Your Codex does not support --yolo; ralph.sh will use a safe fallback."
else
  echo "OK: codex --yolo available"
fi

echo "OK: prerequisites present"
