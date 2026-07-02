#!/usr/bin/env bash
set -e

# ============================================================================
# Spec Kit Status — Show project status and SDD workflow progress
# ============================================================================

# --- Helpers ----------------------------------------------------------------

json_escape() {
    local s="$1"
    s="${s//\\/\\\\}"
    s="${s//\"/\\\"}"
    s="${s//$'\n'/\\n}"
    s="${s//$'\t'/\\t}"
    s="${s//$'\r'/\\r}"
    printf '%s' "$s"
}

has_jq() { command -v jq >/dev/null 2>&1; }

# --- Locate project root ----------------------------------------------------

find_repo_root() {
    local dir="$1"
    while [ "$dir" != "/" ]; do
        if [ -d "$dir/.git" ] || [ -d "$dir/.specify" ]; then
            echo "$dir"
            return 0
        fi
        dir="$(dirname "$dir")"
    done
    return 1
}

SCRIPT_DIR="$(CDPATH="" cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if git rev-parse --show-toplevel >/dev/null 2>&1; then
    REPO_ROOT=$(git rev-parse --show-toplevel)
    HAS_GIT=true
else
    REPO_ROOT="$(find_repo_root "$SCRIPT_DIR")"
    if [ -z "$REPO_ROOT" ]; then
        echo '{"error":"Could not determine repository root"}' >&2
        exit 1
    fi
    HAS_GIT=false
fi

SPECIFY_DIR="$REPO_ROOT/.specify"
if [ ! -d "$SPECIFY_DIR" ]; then
    echo '{"error":"Not a spec-kit project (no .specify/ directory)"}' >&2
    exit 1
fi

PROJECT_NAME=$(basename "$REPO_ROOT")

# --- Detect AI agent(s) -----------------------------------------------------

AGENTS_JSON="[]"
declare -a AGENT_KEYS=("copilot" "claude" "gemini" "cursor-agent" "qwen" "opencode" "codex" "windsurf" "kilocode" "auggie" "codebuddy" "qodercli" "roo" "kiro-cli" "amp" "shai" "tabnine" "agy" "bob" "vibe" "kimi")
declare -a AGENT_NAMES=("GitHub Copilot" "Claude Code" "Gemini CLI" "Cursor" "Qwen Code" "opencode" "Codex CLI" "Windsurf" "Kilo Code" "Auggie CLI" "CodeBuddy" "Qoder CLI" "Roo Code" "Kiro CLI" "Amp" "SHAI" "Tabnine CLI" "Antigravity" "IBM Bob" "Mistral Vibe" "Kimi Code")
declare -a AGENT_FOLDERS=(".github" ".claude" ".gemini" ".cursor" ".qwen" ".opencode" ".codex" ".windsurf" ".kilocode" ".augment" ".codebuddy" ".qoder" ".roo" ".kiro" ".agents" ".shai" ".tabnine/agent" ".agent" ".bob" ".vibe" ".kimi")
declare -a AGENT_CMD_DIRS=("agents" "commands" "commands" "commands" "commands" "command" "prompts" "workflows" "workflows" "commands" "commands" "commands" "commands" "prompts" "commands" "commands" "commands" "workflows" "commands" "prompts" "skills")

detected_agents=""
for i in "${!AGENT_KEYS[@]}"; do
    folder="${AGENT_FOLDERS[$i]}"
    agent_dir="$REPO_ROOT/$folder"
    if [ -d "$agent_dir" ]; then
        cmd_dir="$agent_dir/${AGENT_CMD_DIRS[$i]}"
        has_commands=false
        if [ -d "$cmd_dir" ] && [ -n "$(ls -A "$cmd_dir" 2>/dev/null)" ]; then
            has_commands=true
        fi
        entry="$(printf '{"key":"%s","name":"%s","folder":"%s","has_commands":%s}' \
            "$(json_escape "${AGENT_KEYS[$i]}")" \
            "$(json_escape "${AGENT_NAMES[$i]}")" \
            "$(json_escape "$folder/")" \
            "$has_commands")"
        if [ -z "$detected_agents" ]; then
            detected_agents="$entry"
        else
            detected_agents="$detected_agents,$entry"
        fi
    fi
done
AGENTS_JSON="[$detected_agents]"

# --- Detect script type ------------------------------------------------------

SCRIPT_TYPE="none"
has_bash=false
has_ps=false
[ -d "$SPECIFY_DIR/scripts/bash" ] && has_bash=true
[ -d "$SPECIFY_DIR/scripts/powershell" ] && has_ps=true

if $has_bash && $has_ps; then
    SCRIPT_TYPE="sh + ps"
elif $has_bash; then
    SCRIPT_TYPE="sh"
elif $has_ps; then
    SCRIPT_TYPE="ps"
fi

# --- Detect current feature --------------------------------------------------

CURRENT_BRANCH=""
FEATURE_DIR=""
SPECS_DIR="$REPO_ROOT/specs"

# 1. Check SPECIFY_FEATURE env var
if [ -n "${SPECIFY_FEATURE:-}" ]; then
    CURRENT_BRANCH="$SPECIFY_FEATURE"
fi

# 2. Try git
if [ -z "$CURRENT_BRANCH" ] && [ "$HAS_GIT" = true ]; then
    CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "")
fi

# 3. Fallback: scan specs/ for highest-numbered dir
if [ -z "$CURRENT_BRANCH" ] || [ "$CURRENT_BRANCH" = "main" ] || [ "$CURRENT_BRANCH" = "master" ]; then
    if [ -d "$SPECS_DIR" ]; then
        highest=0
        latest=""
        for dir in "$SPECS_DIR"/*/; do
            [ -d "$dir" ] || continue
            dirname=$(basename "$dir")
            if [[ "$dirname" =~ ^([0-9]{3})- ]]; then
                num=$((10#${BASH_REMATCH[1]}))
                if [ "$num" -gt "$highest" ]; then
                    highest=$num
                    latest=$dirname
                fi
            fi
        done
        if [ -n "$latest" ] && { [ -z "$CURRENT_BRANCH" ] || [ "$CURRENT_BRANCH" = "main" ] || [ "$CURRENT_BRANCH" = "master" ]; }; then
            CURRENT_BRANCH="$latest"
        fi
    fi
fi

# --- Resolve feature directory (prefix-based) --------------------------------

if [[ "$CURRENT_BRANCH" =~ ^([0-9]{3})- ]] && [ -d "$SPECS_DIR" ]; then
    prefix="${BASH_REMATCH[1]}"
    matches=()
    for dir in "$SPECS_DIR"/"$prefix"-*/; do
        [ -d "$dir" ] && matches+=("$(basename "$dir")")
    done
    if [ ${#matches[@]} -eq 1 ]; then
        FEATURE_DIR="$SPECS_DIR/${matches[0]}"
    elif [ ${#matches[@]} -gt 1 ]; then
        if [ -d "$SPECS_DIR/$CURRENT_BRANCH" ]; then
            FEATURE_DIR="$SPECS_DIR/$CURRENT_BRANCH"
        else
            FEATURE_DIR="$SPECS_DIR/${matches[0]}"
        fi
    fi
fi

# --- Check SDD artifacts -----------------------------------------------------

has_spec=false; has_plan=false; has_tasks=false
has_research=false; has_datamodel=false; has_quickstart=false
has_contracts=false; contract_count=0
has_checklists=false; checklist_count=0; checklists_all_pass=true

tasks_completed=0; tasks_total=0

if [ -n "$FEATURE_DIR" ] && [ -d "$FEATURE_DIR" ]; then
    [ -f "$FEATURE_DIR/spec.md" ] && has_spec=true
    [ -f "$FEATURE_DIR/plan.md" ] && has_plan=true
    [ -f "$FEATURE_DIR/tasks.md" ] && has_tasks=true
    [ -f "$FEATURE_DIR/research.md" ] && has_research=true
    [ -f "$FEATURE_DIR/data-model.md" ] && has_datamodel=true
    [ -f "$FEATURE_DIR/quickstart.md" ] && has_quickstart=true

    if [ -d "$FEATURE_DIR/contracts" ] && [ -n "$(ls -A "$FEATURE_DIR/contracts" 2>/dev/null)" ]; then
        has_contracts=true
        contract_count=$(find "$FEATURE_DIR/contracts" -maxdepth 1 -type f | wc -l | tr -d ' ')
    fi

    if [ -d "$FEATURE_DIR/checklists" ]; then
        for cl in "$FEATURE_DIR/checklists"/*.md; do
            [ -f "$cl" ] || continue
            checklist_count=$((checklist_count + 1))
            has_checklists=true
            if grep -q '^\s*-\s*\[ \]' "$cl" 2>/dev/null; then
                checklists_all_pass=false
            fi
        done
    fi

    # Parse task progress
    if [ "$has_tasks" = true ]; then
        while IFS= read -r line; do
            stripped=$(echo "$line" | sed 's/^[[:space:]]*//')
            if echo "$stripped" | grep -qE '^\-\s+\[([ xX])\]'; then
                tasks_total=$((tasks_total + 1))
                if echo "$stripped" | grep -qE '^\-\s+\[[xX]\]'; then
                    tasks_completed=$((tasks_completed + 1))
                fi
            fi
        done < "$FEATURE_DIR/tasks.md"
    fi
fi

# --- Detect workflow phase ---------------------------------------------------

PHASE="unknown"
PHASE_HINT=""

if [ -z "$FEATURE_DIR" ] || [ ! -d "$FEATURE_DIR" ]; then
    PHASE="no_feature"
    PHASE_HINT="Run /speckit.specify to create a feature"
elif [ "$has_tasks" = true ] && [ "$tasks_total" -gt 0 ] && [ "$tasks_completed" -eq "$tasks_total" ]; then
    PHASE="complete"
    PHASE_HINT="All tasks done. Review your implementation."
elif [ "$has_tasks" = true ]; then
    PHASE="implement"
    PHASE_HINT="Ready for /speckit.implement"
elif [ "$has_plan" = true ]; then
    PHASE="tasks"
    PHASE_HINT="Ready for /speckit.tasks"
elif [ "$has_spec" = true ]; then
    PHASE="plan"
    PHASE_HINT="Ready for /speckit.clarify or /speckit.plan"
else
    PHASE="not_started"
    PHASE_HINT="Run /speckit.specify to create a spec"
fi

# --- Extensions summary ------------------------------------------------------

installed_count=0
EXTENSIONS_DIR="$SPECIFY_DIR/extensions"
REGISTRY_FILE="$EXTENSIONS_DIR/.registry"
if [ -f "$REGISTRY_FILE" ]; then
    if has_jq; then
        installed_count=$(jq '.extensions | length' "$REGISTRY_FILE" 2>/dev/null || echo 0)
    else
        installed_count=$(grep -c '"installed_at"' "$REGISTRY_FILE" 2>/dev/null || echo 0)
    fi
fi

# --- Feature count -----------------------------------------------------------

feature_count=0
if [ -d "$SPECS_DIR" ]; then
    for dir in "$SPECS_DIR"/*/; do
        [ -d "$dir" ] && feature_count=$((feature_count + 1))
    done
fi

# --- Output JSON -------------------------------------------------------------

FEATURE_NAME=""
if [ -n "$FEATURE_DIR" ] && [ -d "$FEATURE_DIR" ]; then
    FEATURE_NAME=$(basename "$FEATURE_DIR")
fi

if has_jq; then
    jq -cn \
        --arg project_name "$PROJECT_NAME" \
        --arg has_git "$HAS_GIT" \
        --arg current_branch "$CURRENT_BRANCH" \
        --arg feature_name "$FEATURE_NAME" \
        --arg script_type "$SCRIPT_TYPE" \
        --argjson agents "$AGENTS_JSON" \
        --arg has_spec "$has_spec" \
        --arg has_plan "$has_plan" \
        --arg has_tasks "$has_tasks" \
        --arg has_research "$has_research" \
        --arg has_datamodel "$has_datamodel" \
        --arg has_quickstart "$has_quickstart" \
        --arg has_contracts "$has_contracts" \
        --argjson contract_count "$contract_count" \
        --arg has_checklists "$has_checklists" \
        --argjson checklist_count "$checklist_count" \
        --arg checklists_all_pass "$checklists_all_pass" \
        --argjson tasks_completed "$tasks_completed" \
        --argjson tasks_total "$tasks_total" \
        --arg phase "$PHASE" \
        --arg phase_hint "$PHASE_HINT" \
        --argjson installed_extensions "$installed_count" \
        --argjson feature_count "$feature_count" \
        '{
            project_name: $project_name,
            has_git: ($has_git == "true"),
            current_branch: $current_branch,
            feature_name: $feature_name,
            script_type: $script_type,
            agents: $agents,
            artifacts: {
                spec: ($has_spec == "true"),
                plan: ($has_plan == "true"),
                tasks: ($has_tasks == "true"),
                research: ($has_research == "true"),
                data_model: ($has_datamodel == "true"),
                quickstart: ($has_quickstart == "true"),
                contracts: ($has_contracts == "true"),
                contract_count: $contract_count,
                checklists: ($has_checklists == "true"),
                checklist_count: $checklist_count,
                checklists_all_pass: ($checklists_all_pass == "true")
            },
            tasks: {
                completed: $tasks_completed,
                total: $tasks_total
            },
            phase: $phase,
            phase_hint: $phase_hint,
            installed_extensions: $installed_extensions,
            feature_count: $feature_count
        }'
else
    cat <<ENDJSON
{"project_name":"$(json_escape "$PROJECT_NAME")","has_git":$HAS_GIT,"current_branch":"$(json_escape "$CURRENT_BRANCH")","feature_name":"$(json_escape "$FEATURE_NAME")","script_type":"$(json_escape "$SCRIPT_TYPE")","agents":$AGENTS_JSON,"artifacts":{"spec":$has_spec,"plan":$has_plan,"tasks":$has_tasks,"research":$has_research,"data_model":$has_datamodel,"quickstart":$has_quickstart,"contracts":$has_contracts,"contract_count":$contract_count,"checklists":$has_checklists,"checklist_count":$checklist_count,"checklists_all_pass":$checklists_all_pass},"tasks":{"completed":$tasks_completed,"total":$tasks_total},"phase":"$(json_escape "$PHASE")","phase_hint":"$(json_escape "$PHASE_HINT")","installed_extensions":$installed_count,"feature_count":$feature_count}
ENDJSON
fi
