#!/usr/bin/env bash
#
# Purpose:
#   Install a private ai-tooling version into a target repository.
#
# Usage:
#   ./ai/install-ai-tooling.sh [tag]
#
# Parameters:
#   tag      Optional ai-tooling tag or branch to install. If omitted, the
#            latest tag from the recorded ai-tooling source is used.
#
# Notes:
#   The target repository must have a clean git worktree, except for feature
#   specs under specs/. After a successful install, this script records the
#   installed tag, writes an install manifest, validates copied assets, and
#   commits the resulting installer-managed changes automatically.

set -euo pipefail

TAG_ARG=""

while [ $# -gt 0 ]; do
    case "$1" in
        -h|--help)
            echo "Usage: $0 [tag]" >&2
            exit 0
            ;;
        -*)
            echo "Unknown option: $1" >&2
            echo "Usage: $0 [tag]" >&2
            exit 1
            ;;
        *)
            if [ -n "$TAG_ARG" ]; then
                echo "Usage: $0 [tag]" >&2
                exit 1
            fi
            TAG_ARG="$1"
            shift
            ;;
    esac
done

SCRIPT_DIR="$(CDPATH="" cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(CDPATH="" cd "$SCRIPT_DIR/.." && pwd)"
VERSION_MARKER_FILE="$REPO_ROOT/ai/installed-ai-tooling-version"
MANIFEST_FILE="$REPO_ROOT/ai/installed-ai-tooling-manifest.txt"
SOURCE_FILE="$REPO_ROOT/ai/ai-tooling-source"

is_local_git_repo() {
    [ -n "$1" ] && git -C "$1" rev-parse --git-dir >/dev/null 2>&1
}

recorded_source() {
    if [ -f "$SOURCE_FILE" ]; then
        sed -n '1p' "$SOURCE_FILE"
    fi
}

RECORDED_SOURCE_REPO="$(recorded_source)"
if is_local_git_repo "$RECORDED_SOURCE_REPO"; then
    RESOLVED_SOURCE_REPO="$RECORDED_SOURCE_REPO"
else
    echo "Could not resolve ai-tooling source." >&2
    if [ -n "$RECORDED_SOURCE_REPO" ]; then
        echo "Tried recorded source: $RECORDED_SOURCE_REPO" >&2
    else
        echo "No recorded source found at ai/ai-tooling-source." >&2
    fi
    echo "Run ./ai/bootstrap-ai-tooling.sh /path/to/target-repo from the ai-tooling checkout first." >&2
    exit 1
fi

latest_tag() {
    local repo="$1"

    git -C "$repo" tag --sort=-version:refname | head -n 1
}

hash_file() {
    local path="$1"

    if command -v sha256sum >/dev/null 2>&1; then
        sha256sum "$path" | awk '{print $1}'
        return
    fi

    if command -v shasum >/dev/null 2>&1; then
        shasum -a 256 "$path" | awk '{print $1}'
        return
    fi

    if command -v openssl >/dev/null 2>&1; then
        openssl dgst -sha256 -r "$path" | awk '{print $1}'
        return
    fi

    echo "No SHA-256 hashing tool available (expected sha256sum, shasum, or openssl)." >&2
    exit 1
}

append_file_if_present() {
    local path="$1"
    local list_file="$2"

    if [ -f "$path" ]; then
        printf '%s\n' "${path#"$REPO_ROOT"/}" >> "$list_file"
    fi
}

append_tree_files_if_present() {
    local path="$1"
    local list_file="$2"

    if [ -d "$path" ]; then
        find "$path" -type f | sed "s#^$REPO_ROOT/##" >> "$list_file"
    fi
}

is_manifest_excluded_path() {
    case "$1" in
        ai/ai-tooling-source|ai/installed-ai-tooling-manifest.txt|ai/installed-ai-tooling-version|.specify/extensions/*/*-config.yml|.specify/extensions/*/config.yml|.specify/extensions/*/local-config.yml|.specify/extensions/*/*.local.yml)
            return 0
            ;;
        *)
            return 1
            ;;
    esac
}

is_installer_dirty_path() {
    case "$1" in
        ai/bootstrap-ai-tooling.sh|ai/cleanup-ai-tooling-sync.sh|ai/install-ai-tooling.sh|ai/ai-tooling-source|ai/installed-ai-tooling-manifest.txt|ai/installed-ai-tooling-version)
            return 0
            ;;
    esac

    if [ -f "$MANIFEST_FILE" ]; then
        awk -v path="$1" '$1 == path { found = 1 } END { exit found ? 0 : 1 }' "$MANIFEST_FILE"
        return $?
    fi

    return 1
}

is_target_local_allowed_dirty_path() {
    case "$1" in
        specs|specs/*)
            return 0
            ;;
        *)
            return 1
            ;;
    esac
}

validate_project_artifact_tracking() {
    local ignored_paths=()
    local path=""
    local required_paths=(
        specs/__ai_tooling_tracking_probe__/spec.md
        specs/__ai_tooling_tracking_probe__/checklists/requirements.md
        .specify/extensions.yml
        .specify/memory/constitution.md
        .specify/memory/architecture.md
        .specify/memory/architecture-scenario-view.md
        .specify/memory/architecture-logical-view.md
        .specify/memory/architecture-process-view.md
        .specify/memory/architecture-development-view.md
        .specify/memory/architecture-physical-view.md
        .specify/memory/architecture-repo-facts.md
    )

    for path in "${required_paths[@]}"; do
        if git -C "$REPO_ROOT" check-ignore -q -- "$path"; then
            ignored_paths+=("$path")
        fi
    done

    if [ "${#ignored_paths[@]}" -eq 0 ]; then
        return 0
    fi

    echo "Refusing to install ai-tooling while project-owned Spec Kit artifacts are ignored." >&2
    echo "These paths must be trackable because they carry specs, constitution, and architecture memory:" >&2
    for path in "${ignored_paths[@]}"; do
        echo "  - $path" >&2
    done
    cat >&2 <<'EOF'
Update the target repository .gitignore so it does not blanket-ignore Spec Kit project state. For example:
  .specify/*
  !.specify/
  !.specify/memory/
  !.specify/memory/**
  !.specify/extensions.yml
  !specs/
  !specs/**
EOF
    exit 1
}

write_git_status_snapshot() {
    local output_file="$1"

    git -C "$REPO_ROOT" status --porcelain --untracked-files=all > "$output_file"
}

check_target_dirty_state() {
    local blockers_file="$1"
    local status_file="${blockers_file}.status"
    local line=""
    local path=""

    : > "$blockers_file"
    write_git_status_snapshot "$status_file"

    while IFS= read -r line; do
        path="${line:3}"
        case "$path" in
            *" -> "*)
                path="${path##* -> }"
                ;;
        esac

        if ! is_installer_dirty_path "$path" && ! is_target_local_allowed_dirty_path "$path"; then
            printf '%s\n' "$line" >> "$blockers_file"
        fi
    done < "$status_file"

    if [ -s "$blockers_file" ]; then
        echo "Refusing to install ai-tooling into a dirty repository." >&2
        sed 's/^/  - /' "$blockers_file" >&2
        echo "Commit or stash target-local changes outside specs/ first, then rerun this script." >&2
        echo "Installer-managed files may be dirty; they will be refreshed by this install." >&2
        exit 1
    fi
}

stage_install_changes() {
    local status_file="$TMP_DIR/install-status.txt"
    local line=""
    local path=""

    write_git_status_snapshot "$status_file"

    while IFS= read -r line; do
        path="${line:3}"
        case "$path" in
            *" -> "*)
                path="${path##* -> }"
                ;;
        esac

        if is_target_local_allowed_dirty_path "$path"; then
            continue
        fi

        if git -C "$REPO_ROOT" ls-files --error-unmatch -- "$path" >/dev/null 2>&1; then
            git -C "$REPO_ROOT" add -u -- "$path"
        elif ! git -C "$REPO_ROOT" check-ignore -q -- "$path"; then
            git -C "$REPO_ROOT" add -A -- "$path"
        fi
    done < "$status_file"
}

stage_manifest_files() {
    local line=""
    local path=""

    if [ ! -f "$MANIFEST_FILE" ]; then
        return 0
    fi

    while IFS= read -r line; do
        [ -n "$line" ] || continue

        path="${line%% *}"
        case "$path" in
            \#*|tag=*|source=*|commit=*|synced_at_utc=*)
                continue
                ;;
        esac

        if [ -e "$REPO_ROOT/$path" ]; then
            git -C "$REPO_ROOT" add -f -- "$path"
        fi
    done < "$MANIFEST_FILE"
}

write_install_manifest() {
    local source_repo="$1"
    local source_commit="$2"
    local tag="$3"
    local synced_at="$4"
    local list_file="$5"
    local sorted_file="$6"

    : > "$list_file"

    append_file_if_present "$REPO_ROOT/AGENTS.md" "$list_file"
    append_tree_files_if_present "$REPO_ROOT/.agents/skills" "$list_file"
    append_file_if_present "$REPO_ROOT/.specify/init-options.json" "$list_file"
    append_file_if_present "$REPO_ROOT/.specify/integration.json" "$list_file"
    append_file_if_present "$REPO_ROOT/.specify/extensions.yml" "$list_file"
    append_file_if_present "$REPO_ROOT/.specify/extension-catalogs.yml" "$list_file"
    append_tree_files_if_present "$REPO_ROOT/.specify/integrations" "$list_file"
    append_tree_files_if_present "$REPO_ROOT/.specify/workflows" "$list_file"
    append_tree_files_if_present "$REPO_ROOT/.specify/scripts" "$list_file"
    append_tree_files_if_present "$REPO_ROOT/.specify/templates" "$list_file"
    append_tree_files_if_present "$REPO_ROOT/.specify/presets" "$list_file"
    append_tree_files_if_present "$REPO_ROOT/.specify/extensions" "$list_file"
    append_file_if_present "$REPO_ROOT/.specify/memory/constitution.md" "$list_file"
    append_tree_files_if_present "$REPO_ROOT/extensions" "$list_file"
    append_tree_files_if_present "$REPO_ROOT/integrations" "$list_file"
    append_tree_files_if_present "$REPO_ROOT/presets" "$list_file"
    append_tree_files_if_present "$REPO_ROOT/workflows" "$list_file"
    append_tree_files_if_present "$REPO_ROOT/ai" "$list_file"

    sort -u "$list_file" > "$sorted_file"

    {
        printf 'tag=%s\n' "$tag"
        printf 'source=%s\n' "$source_repo"
        printf 'commit=%s\n' "$source_commit"
        printf 'synced_at_utc=%s\n' "$synced_at"
        printf '\n'
        printf '# path sha256\n'

        while IFS= read -r relative_path; do
            [ -n "$relative_path" ] || continue
            if is_manifest_excluded_path "$relative_path"; then
                continue
            fi
            printf '%s %s\n' "$relative_path" "$(hash_file "$REPO_ROOT/$relative_path")"
        done < "$sorted_file"
    } > "$MANIFEST_FILE"
}

validate_installed_artifacts() {
    python3 - "$REPO_ROOT" <<'PY'
import re
import sys
from pathlib import Path

repo_root = Path(sys.argv[1])
skills_root = repo_root / ".agents" / "skills"
extensions_root = repo_root / ".specify" / "extensions"
pattern = re.compile(r"\.specify/extensions/([A-Za-z0-9._-]+)/([A-Za-z0-9._/\-]+)")
missing: set[str] = set()

if not skills_root.is_dir():
    raise SystemExit(0)

for skill_file in skills_root.glob("*/SKILL.md"):
    try:
        content = skill_file.read_text(encoding="utf-8")
    except OSError:
        continue

    for match in pattern.finditer(content):
        extension_id = match.group(1)
        relative_path = match.group(2)
        target_path = extensions_root / extension_id / relative_path
        template_path = None

        if target_path.suffix == ".yml":
            template_path = target_path.with_name(target_path.name[:-4] + ".template.yml")

        if target_path.exists():
            continue
        if template_path is not None and template_path.exists():
            continue

        missing.add(f".specify/extensions/{extension_id}/{relative_path}")

if missing:
    sys.stderr.write("Installed ai-tooling is missing extension artifacts referenced by skills:\n")
    for item in sorted(missing):
        sys.stderr.write(f"  - {item}\n")
    raise SystemExit(1)
PY
}

TAG="${TAG_ARG:-$(latest_tag "$RESOLVED_SOURCE_REPO")}"
if [ -z "$TAG" ]; then
    echo "Could not resolve an ai-tooling tag from $RESOLVED_SOURCE_REPO" >&2
    exit 1
fi

TMP_DIR="$(mktemp -d)"
CLONE_DIR="$TMP_DIR/source"
cleanup() {
    rm -rf "$TMP_DIR"
}
trap cleanup EXIT

check_target_dirty_state "$TMP_DIR/dirty-blockers.txt"
validate_project_artifact_tracking

CLONE_SOURCE="$RESOLVED_SOURCE_REPO"
if is_local_git_repo "$RESOLVED_SOURCE_REPO"; then
    CLONE_SOURCE="file://$RESOLVED_SOURCE_REPO"
fi

git -c advice.detachedHead=false clone --quiet --depth 1 --branch "$TAG" "$CLONE_SOURCE" "$CLONE_DIR"
"$CLONE_DIR/ai/cleanup-ai-tooling-sync.sh" "$REPO_ROOT" "$CLONE_DIR"
"$CLONE_DIR/install/sync.sh" "$REPO_ROOT"
validate_installed_artifacts

SOURCE_COMMIT="$(git -C "$CLONE_DIR" rev-parse HEAD)"
SYNCED_AT="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
MANIFEST_LIST_FILE="$TMP_DIR/install-manifest-paths.txt"
MANIFEST_SORTED_FILE="$TMP_DIR/install-manifest-paths.sorted.txt"
write_install_manifest "$RESOLVED_SOURCE_REPO" "$SOURCE_COMMIT" "$TAG" "$SYNCED_AT" "$MANIFEST_LIST_FILE" "$MANIFEST_SORTED_FILE"

printf '%s\n' "$TAG" > "$VERSION_MARKER_FILE"

stage_manifest_files
stage_install_changes

modified_files=()
while IFS= read -r modified_file; do
    if ! is_target_local_allowed_dirty_path "$modified_file"; then
        modified_files+=("$modified_file")
    fi
done < <(git -C "$REPO_ROOT" diff --cached --name-only)

if [ "${#modified_files[@]}" -gt 0 ]; then
    commit_message="chore(ai): install ai-tooling $TAG"

    commit_message+=$'\n\nModified files:\n'

    for modified_file in "${modified_files[@]}"; do
        commit_message+="- ${modified_file}"$'\n'
    done

    git -C "$REPO_ROOT" commit -m "$commit_message" -- "${modified_files[@]}"
fi

echo "Installed ai-tooling $TAG from $RESOLVED_SOURCE_REPO"
