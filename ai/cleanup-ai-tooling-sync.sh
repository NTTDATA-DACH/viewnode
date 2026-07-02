#!/usr/bin/env bash
#
# Purpose:
#   Remove stale ai-tooling-managed files from a target repository before sync.
#
# Usage:
#   ./ai/cleanup-ai-tooling-sync.sh [--dry-run] <target-repo> <ai-tooling-repo>
#
# Parameters:
#   --dry-run          Print removals without deleting files.
#   target-repo       Target repository to clean.
#   ai-tooling-repo   Source checkout to compare against.

set -euo pipefail

usage() {
    cat >&2 <<'EOF'
Usage: cleanup-ai-tooling-sync.sh [--dry-run] <target-repo> <ai-tooling-repo>
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

if [ $# -ne 2 ]; then
    usage
fi

REPO_ROOT="$(CDPATH="" cd "$1" && pwd)"
SOURCE_ROOT="$(CDPATH="" cd "$2" && pwd)"

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

cleanup_child_dirs_only() {
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

        if [ ! -d "$target_entry" ]; then
            delete_path "$target_entry"
            continue
        fi

        if [ ! -d "$source_dir/$name" ]; then
            delete_path "$target_entry"
            continue
        fi

        cleanup_tree "$source_dir/$name" "$target_entry"
    done

    if [ "$DRY_RUN" -eq 0 ]; then
        find "$target_dir" -depth -type d -empty -delete
    fi
}

cleanup_tree "$SOURCE_ROOT/.agents/skills" "$REPO_ROOT/.agents/skills"
cleanup_child_dirs_only "$SOURCE_ROOT/extensions" "$REPO_ROOT/extensions"
cleanup_child_dirs_only "$SOURCE_ROOT/integrations" "$REPO_ROOT/integrations"
cleanup_child_dirs_only "$SOURCE_ROOT/presets" "$REPO_ROOT/presets"
cleanup_child_dirs_only "$SOURCE_ROOT/workflows" "$REPO_ROOT/workflows"
cleanup_tree "$SOURCE_ROOT/.specify/scripts" "$REPO_ROOT/.specify/scripts"
cleanup_tree "$SOURCE_ROOT/.specify/templates" "$REPO_ROOT/.specify/templates"
cleanup_tree "$SOURCE_ROOT/.specify/integrations" "$REPO_ROOT/.specify/integrations"
cleanup_tree "$SOURCE_ROOT/.specify/presets" "$REPO_ROOT/.specify/presets"
cleanup_tree "$SOURCE_ROOT/.specify/workflows" "$REPO_ROOT/.specify/workflows"

for managed_metadata_file in \
    .specify/extension-catalogs.yml \
    .specify/extensions.yml \
    .specify/init-options.json \
    .specify/integration.json \
    .specify/workflow-catalogs.yml
do
    if [ ! -e "$SOURCE_ROOT/$managed_metadata_file" ]; then
        delete_path "$REPO_ROOT/$managed_metadata_file"
    fi
done

for legacy_managed_file in \
    scripts/ralph/ralph.sh \
    scripts/ralph/doctor.sh \
    scripts/ralph/prompt.md \
    scripts/ralph/prd.json.example
do
    delete_path "$REPO_ROOT/$legacy_managed_file"
done

if [ "$DRY_RUN" -eq 0 ] && [ -d "$REPO_ROOT/scripts/ralph" ]; then
    find "$REPO_ROOT/scripts/ralph" -depth -type d -empty -delete
fi

if [ "$DRY_RUN" -eq 1 ]; then
    echo "Dry run complete."
else
    echo "Cleanup complete."
fi
