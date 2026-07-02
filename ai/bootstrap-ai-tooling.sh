#!/usr/bin/env bash
#
# Purpose:
#   Bootstrap ai-tooling in two small steps:
#   - from ai-tooling: copy this bootstrap script into a target repository
#   - from a target: copy ai/install-ai-tooling.sh from the recorded source
#
# Usage:
#   ./ai/bootstrap-ai-tooling.sh [--allow-dirty] /path/to/target-repo
#   ./ai/bootstrap-ai-tooling.sh
#
# Parameters:
#   --allow-dirty  When copying into a target, allow copying from a dirty
#                  ai-tooling working tree.
#   target-repo    Optional target repository. If present, bootstrap is copied
#                  there. If omitted, install-ai-tooling.sh is copied into the
#                  current target repository from the recorded source.

set -euo pipefail

usage() {
    cat >&2 <<'EOF'
Usage:
  bootstrap-ai-tooling.sh [--allow-dirty] /path/to/target-repo
  bootstrap-ai-tooling.sh
EOF
    exit 1
}

ALLOW_DIRTY=0
TARGET_ARG=""

while [ $# -gt 0 ]; do
    case "$1" in
        --allow-dirty)
            ALLOW_DIRTY=1
            shift
            ;;
        -h|--help)
            usage
            ;;
        -*)
            echo "Unknown option: $1" >&2
            usage
            ;;
        *)
            if [ -n "$TARGET_ARG" ]; then
                usage
            fi
            TARGET_ARG="$1"
            shift
            ;;
    esac
done

SCRIPT_DIR="$(CDPATH="" cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(CDPATH="" cd "$SCRIPT_DIR/.." && pwd)"
SOURCE_FILE="$REPO_ROOT/ai/ai-tooling-source"

is_git_worktree() {
    git -C "$1" rev-parse --is-inside-work-tree >/dev/null 2>&1
}

install_bootstrap_into_target() {
    local target_root="$1"

    target_root="$(CDPATH="" cd "$target_root" && pwd)"
    if ! is_git_worktree "$target_root"; then
        echo "Target is not a git repository: $target_root" >&2
        exit 1
    fi

    if [ -n "$(git -C "$REPO_ROOT" status --porcelain)" ] && [ "$ALLOW_DIRTY" -ne 1 ]; then
        echo "Refusing to sync bootstrap from a dirty ai-tooling working tree." >&2
        echo "Commit or stash ai-tooling changes first, or rerun with --allow-dirty to copy the worktree bootstrap intentionally:" >&2
        echo "  ./ai/bootstrap-ai-tooling.sh --allow-dirty $target_root" >&2
        exit 1
    fi

    mkdir -p "$target_root/ai"
    cp "$REPO_ROOT/ai/bootstrap-ai-tooling.sh" "$target_root/ai/bootstrap-ai-tooling.sh"
    chmod +x "$target_root/ai/bootstrap-ai-tooling.sh"
    printf '%s\n' "$REPO_ROOT" > "$target_root/ai/ai-tooling-source"

    echo "Installed ai/bootstrap-ai-tooling.sh into $target_root"
}

install_installer_from_source() {
    local source_root=""

    if [ -f "$SOURCE_FILE" ]; then
        source_root="$(sed -n '1p' "$SOURCE_FILE")"
    fi

    if [ -z "$source_root" ] || [ ! -f "$source_root/ai/install-ai-tooling.sh" ]; then
        echo "Could not resolve ai-tooling source for install-ai-tooling.sh." >&2
        echo "Run this from the ai-tooling checkout first:" >&2
        echo "  ./ai/bootstrap-ai-tooling.sh $REPO_ROOT" >&2
        exit 1
    fi

    mkdir -p "$REPO_ROOT/ai"
    cp "$source_root/ai/install-ai-tooling.sh" "$REPO_ROOT/ai/install-ai-tooling.sh"
    chmod +x "$REPO_ROOT/ai/install-ai-tooling.sh"
    printf '%s\n' "$source_root" > "$SOURCE_FILE"

    echo "Installed ai/install-ai-tooling.sh from $source_root"
}

if [ -n "$TARGET_ARG" ]; then
    install_bootstrap_into_target "$TARGET_ARG"
else
    install_installer_from_source
fi
