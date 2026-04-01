#!/usr/bin/env bash
#
# Removes stale ai-tooling-managed files from this repository before a sync.
#
# Usage examples:
#   ./ai/cleanup-ai-tooling-sync.sh
#   ./ai/cleanup-ai-tooling-sync.sh --dry-run
#   AI_TOOLING_REPO=/path/to/local/ai-tooling ./ai/cleanup-ai-tooling-sync.sh
#   ./ai/cleanup-ai-tooling-sync.sh /path/to/local/ai-tooling

set -euo pipefail

usage() {
    cat >&2 <<'EOF'
Usage: cleanup-ai-tooling-sync.sh [--dry-run] [ai-tooling-repo]
EOF
    exit 1
}

DRY_RUN=0

while [ $# -gt 0 ]; do
    case "$1" in
        --dry-run)
            DRY_RUN=1
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
            break
            ;;
    esac
done

if [ $# -gt 1 ]; then
    usage
fi

SCRIPT_DIR="$(CDPATH="" cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(CDPATH="" cd "$SCRIPT_DIR/.." && pwd)"
DEFAULT_LOCAL_AI_TOOLING_REPO="$(CDPATH="" cd "$REPO_ROOT/../.." && pwd)/ai-tooling"

if [ -n "${AI_TOOLING_REPO:-}" ]; then
    SOURCE_ROOT="$AI_TOOLING_REPO"
elif [ $# -eq 1 ]; then
    SOURCE_ROOT="$1"
else
    SOURCE_ROOT="$DEFAULT_LOCAL_AI_TOOLING_REPO"
fi

if [ ! -d "$SOURCE_ROOT" ]; then
    echo "ai-tooling source repository does not exist: $SOURCE_ROOT" >&2
    exit 1
fi

delete_path() {
    local path="$1"

    if [ ! -e "$path" ]; then
        return
    fi

    if [ "$DRY_RUN" -eq 1 ]; then
        echo "Would remove $path"
        return
    fi

    rm -rf "$path"
    echo "Removed $path"
}

cleanup_tree() {
    local source_dir="$1"
    local target_dir="$2"

    if [ ! -d "$target_dir" ]; then
        return
    fi

    if [ ! -d "$source_dir" ]; then
        delete_path "$target_dir"
        return
    fi

    local target_entry
    for target_entry in "$target_dir"/*; do
        if [ "$target_entry" = "$target_dir/*" ]; then
            break
        fi

        local name
        name="$(basename "$target_entry")"

        if [ ! -e "$source_dir/$name" ]; then
            delete_path "$target_entry"
            continue
        fi

        if [ -d "$target_entry" ] && [ -d "$source_dir/$name" ]; then
            cleanup_tree "$source_dir/$name" "$target_entry"
        fi
    done

    if [ "$DRY_RUN" -eq 0 ]; then
        find "$target_dir" -depth -type d -empty -delete
    fi
}

cleanup_tree "$SOURCE_ROOT/.agents/skills" "$REPO_ROOT/.agents/skills"
cleanup_tree "$SOURCE_ROOT/.specify/scripts" "$REPO_ROOT/.specify/scripts"
cleanup_tree "$SOURCE_ROOT/.specify/templates" "$REPO_ROOT/.specify/templates"

for managed_file in scripts/ralph/ralph.sh scripts/ralph/doctor.sh scripts/ralph/prompt.md; do
    if [ ! -e "$SOURCE_ROOT/$managed_file" ]; then
        delete_path "$REPO_ROOT/$managed_file"
    fi
done

if [ "$DRY_RUN" -eq 1 ]; then
    echo "Dry run complete."
else
    echo "Cleanup complete."
fi
