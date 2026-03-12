#!/bin/bash
# Ralph Wiggum - Long-running AI agent loop
# Usage: ./ralph.sh [max_iterations]

set -euo pipefail

MAX_ITERATIONS=${1:-10}
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PRD_FILE="$SCRIPT_DIR/prd.json"
PROGRESS_FILE="$SCRIPT_DIR/progress.txt"
ARCHIVE_DIR="$SCRIPT_DIR/archive"
LAST_BRANCH_FILE="$SCRIPT_DIR/.last-branch"

WORKSPACE_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
CODEX_BIN="${CODEX_BIN:-codex}"
CODEX_LAST_MESSAGE_LATEST_FILE="$SCRIPT_DIR/.codex-last-message.txt"

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Missing required command: $1" >&2
    exit 1
  fi
}

supports_codex_yolo() {
  local out
  out="$("$CODEX_BIN" --yolo exec --help 2>&1 || true)"
  if echo "$out" | grep -qi "unexpected argument '--yolo'"; then
    return 1
  fi
  if echo "$out" | grep -qi "Run Codex non-interactively"; then
    return 0
  fi
  return 1
}

escape_sed_replacement() {
  printf '%s' "$1" | sed -e 's/[\\/&]/\\&/g'
}

render_prompt() {
  local ralph_dir prd_file progress_file
  ralph_dir="$(escape_sed_replacement "$SCRIPT_DIR")"
  prd_file="$(escape_sed_replacement "$PRD_FILE")"
  progress_file="$(escape_sed_replacement "$PROGRESS_FILE")"

  sed \
    -e "s|{{RALPH_DIR}}|$ralph_dir|g" \
    -e "s|{{PRD_FILE}}|$prd_file|g" \
    -e "s|{{PROGRESS_FILE}}|$progress_file|g" \
    "$SCRIPT_DIR/prompt.md"
}

has_complete_token() {
  local path="$1"
  [ -f "$path" ] || return 1
  grep -q "<promise>COMPLETE</promise>" "$path"
}

prd_all_passes() {
  [ -f "$PRD_FILE" ] || return 1
  jq -e '(.userStories | length) > 0 and all(.userStories[]; .passes == true)' "$PRD_FILE" >/dev/null 2>&1
}

require_cmd jq
require_cmd git
require_cmd sed
require_cmd tee
require_cmd "$CODEX_BIN"

if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  echo "Ralph must be run inside a git repository." >&2
  echo "Tip: run from your project root (not from $SCRIPT_DIR)." >&2
  exit 1
fi

if ! [[ "$MAX_ITERATIONS" =~ ^[0-9]+$ ]] || [ "$MAX_ITERATIONS" -lt 1 ]; then
  echo "Invalid max_iterations: $MAX_ITERATIONS" >&2
  exit 1
fi

if [ ! -f "$SCRIPT_DIR/prompt.md" ]; then
  echo "Missing prompt file: $SCRIPT_DIR/prompt.md" >&2
  exit 1
fi

if [ ! -f "$PRD_FILE" ]; then
  echo "Missing PRD file: $PRD_FILE" >&2
  exit 1
fi

# Archive previous run if branch changed
if [ -f "$PRD_FILE" ] && [ -f "$LAST_BRANCH_FILE" ]; then
  CURRENT_BRANCH=$(jq -r '.branchName // empty' "$PRD_FILE" 2>/dev/null || echo "")
  LAST_BRANCH=$(cat "$LAST_BRANCH_FILE" 2>/dev/null || echo "")
  
  if [ -n "$CURRENT_BRANCH" ] && [ -n "$LAST_BRANCH" ] && [ "$CURRENT_BRANCH" != "$LAST_BRANCH" ]; then
    # Archive the previous run
    DATE=$(date +%Y-%m-%d)
    # Strip "ralph/" prefix from branch name for folder
    FOLDER_NAME=$(echo "$LAST_BRANCH" | sed 's|^ralph/||')
    ARCHIVE_FOLDER="$ARCHIVE_DIR/$DATE-$FOLDER_NAME"
    
    echo "Archiving previous run: $LAST_BRANCH"
    mkdir -p "$ARCHIVE_FOLDER"
    [ -f "$PRD_FILE" ] && cp "$PRD_FILE" "$ARCHIVE_FOLDER/"
    [ -f "$PROGRESS_FILE" ] && cp "$PROGRESS_FILE" "$ARCHIVE_FOLDER/"
    echo "   Archived to: $ARCHIVE_FOLDER"
    
    # Reset progress file for new run
    echo "# Ralph Progress Log" > "$PROGRESS_FILE"
    echo "Started: $(date)" >> "$PROGRESS_FILE"
    echo "---" >> "$PROGRESS_FILE"
  fi
fi

# Track current branch
if [ -f "$PRD_FILE" ]; then
  CURRENT_BRANCH=$(jq -r '.branchName // empty' "$PRD_FILE" 2>/dev/null || echo "")
  if [ -n "$CURRENT_BRANCH" ]; then
    echo "$CURRENT_BRANCH" > "$LAST_BRANCH_FILE"
  fi
fi

# Initialize progress file if it doesn't exist
if [ ! -f "$PROGRESS_FILE" ]; then
  echo "# Ralph Progress Log" > "$PROGRESS_FILE"
  echo "Started: $(date)" >> "$PROGRESS_FILE"
  echo "---" >> "$PROGRESS_FILE"
fi

echo "Starting Ralph - Max iterations: $MAX_ITERATIONS"
echo "Workspace root: $WORKSPACE_ROOT"

for ((i=1; i<=MAX_ITERATIONS; i++)); do
  echo ""
  echo "═══════════════════════════════════════════════════════"
  echo "  Ralph Iteration $i of $MAX_ITERATIONS"
  echo "═══════════════════════════════════════════════════════"

  CODEX_ITER_LAST_MESSAGE_FILE="$SCRIPT_DIR/.codex-last-message-iter-$i.txt"

  CODEX_ARGS=(exec -C "$WORKSPACE_ROOT" - --output-last-message "$CODEX_ITER_LAST_MESSAGE_FILE")
  if supports_codex_yolo; then
    CODEX_ARGS=(--yolo "${CODEX_ARGS[@]}")
  else
    CODEX_ARGS=(exec --dangerously-bypass-approvals-and-sandbox -C "$WORKSPACE_ROOT" - --output-last-message "$CODEX_ITER_LAST_MESSAGE_FILE")
  fi

  # Run Codex with the Ralph prompt (fresh context every iteration)
  OUTPUT=$(render_prompt | "$CODEX_BIN" "${CODEX_ARGS[@]}" 2>&1 | tee /dev/stderr) || true

  cp -f "$CODEX_ITER_LAST_MESSAGE_FILE" "$CODEX_LAST_MESSAGE_LATEST_FILE" >/dev/null 2>&1 || true
  
  # Check for completion signal
  if echo "$OUTPUT" | grep -q "<promise>COMPLETE</promise>" || has_complete_token "$CODEX_ITER_LAST_MESSAGE_FILE" || prd_all_passes; then
    echo ""
    echo "Ralph completed all tasks!"
    echo "Completed at iteration $i of $MAX_ITERATIONS"
    exit 0
  fi
  
  echo "Iteration $i complete. Continuing..."
  sleep 2
done

echo ""
echo "Ralph reached max iterations ($MAX_ITERATIONS) without completing all tasks."
echo "Check $PROGRESS_FILE for status."
exit 1
