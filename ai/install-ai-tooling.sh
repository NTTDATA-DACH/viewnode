#!/usr/bin/env bash
#
# Installs a private ai-tooling version into this repository.
#
# Usage examples:
#   ./ai/install-ai-tooling.sh
#   ./ai/install-ai-tooling.sh v0.2.0
#   AI_TOOLING_REPO=/path/to/local/ai-tooling ./ai/install-ai-tooling.sh
#   AI_TOOLING_REPO=git@github.com:your-org/private-ai-tooling.git ./ai/install-ai-tooling.sh
#
# This script only runs on a clean git worktree. After a successful install,
# it records the installed ai-tooling tag in ai/installed-ai-tooling-version
# and commits that metadata automatically.

set -euo pipefail

if [ $# -gt 1 ]; then
    echo "Usage: $0 [tag]" >&2
    exit 1
fi

SCRIPT_DIR="$(CDPATH="" cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(CDPATH="" cd "$SCRIPT_DIR/.." && pwd)"
DEFAULT_LOCAL_AI_TOOLING_REPO="$(CDPATH="" cd "$REPO_ROOT/../.." && pwd)/ai-tooling"
VERSION_MARKER_FILE="$REPO_ROOT/ai/installed-ai-tooling-version"

if [ -n "$(git -C "$REPO_ROOT" status --porcelain)" ]; then
    echo "Refusing to install ai-tooling into a dirty repository." >&2
    echo "Commit or stash existing changes first, then rerun this script." >&2
    exit 1
fi

if [ -n "${AI_TOOLING_REPO:-}" ]; then
    RESOLVED_AI_TOOLING_REPO="$AI_TOOLING_REPO"
elif [ -d "$DEFAULT_LOCAL_AI_TOOLING_REPO/.git" ]; then
    RESOLVED_AI_TOOLING_REPO="$DEFAULT_LOCAL_AI_TOOLING_REPO"
else
    echo "Could not resolve ai-tooling source." >&2
    echo "Tried default local checkout: $DEFAULT_LOCAL_AI_TOOLING_REPO" >&2
    echo "Set AI_TOOLING_REPO to a local ai-tooling checkout or private git URL." >&2
    exit 1
fi

latest_tag() {
    local repo="$1"

    if [ -d "$repo/.git" ]; then
        git -C "$repo" tag --sort=-version:refname | head -n 1
        return
    fi

    git ls-remote --tags --refs "$repo" \
        | sed 's#.*refs/tags/##' \
        | sort -V -r \
        | head -n 1
}

TAG="${1:-$(latest_tag "$RESOLVED_AI_TOOLING_REPO")}"
if [ -z "$TAG" ]; then
    echo "Could not resolve an ai-tooling tag from $RESOLVED_AI_TOOLING_REPO" >&2
    exit 1
fi

TMP_DIR="$(mktemp -d)"
cleanup() {
    rm -rf "$TMP_DIR"
}
trap cleanup EXIT

CLONE_SOURCE="$RESOLVED_AI_TOOLING_REPO"
if [ -d "$RESOLVED_AI_TOOLING_REPO/.git" ]; then
    CLONE_SOURCE="file://$RESOLVED_AI_TOOLING_REPO"
fi

git -c advice.detachedHead=false clone --quiet --depth 1 --branch "$TAG" "$CLONE_SOURCE" "$TMP_DIR"
"$REPO_ROOT/ai/cleanup-ai-tooling-sync.sh" "$TMP_DIR"
"$TMP_DIR/install/sync.sh" "$REPO_ROOT"

printf '%s\n' "$TAG" > "$VERSION_MARKER_FILE"

if [ -n "$(git -C "$REPO_ROOT" status --porcelain -- "$VERSION_MARKER_FILE")" ]; then
    git -C "$REPO_ROOT" add "$VERSION_MARKER_FILE"
    git -C "$REPO_ROOT" commit -m "chore(ai): install ai-tooling $TAG"
fi

echo "Installed ai-tooling $TAG from $RESOLVED_AI_TOOLING_REPO"
