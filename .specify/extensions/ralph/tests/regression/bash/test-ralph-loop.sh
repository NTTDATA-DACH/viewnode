#!/usr/bin/env bash
#
# Regression tests for ralph-loop.sh helper functions.
# Extracts and tests functions in isolation without running the full loop.
#
# Usage:
#   bash tests/regression/bash/test-ralph-loop.sh
#
# Exit codes:
#   0 - All tests passed
#   1 - One or more tests failed

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
FIXTURE_DIR="$SCRIPT_DIR/../fixtures"
SOURCE_SCRIPT="$REPO_ROOT/scripts/bash/ralph-loop.sh"

# Test bookkeeping
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0
FAILURES=()

#region Test Harness

assert_eq() {
    local test_name="$1"
    local expected="$2"
    local actual="$3"
    ((TESTS_RUN++))

    if [[ "$expected" == "$actual" ]]; then
        echo -e "  \033[32mPASS\033[0m $test_name"
        ((TESTS_PASSED++))
    else
        echo -e "  \033[31mFAIL\033[0m $test_name"
        echo "         expected: [$expected]"
        echo "         actual:   [$actual]"
        ((TESTS_FAILED++))
        FAILURES+=("$test_name")
    fi
}

assert_true() {
    local test_name="$1"
    shift
    ((TESTS_RUN++))

    if "$@" >/dev/null 2>&1; then
        echo -e "  \033[32mPASS\033[0m $test_name"
        ((TESTS_PASSED++))
    else
        echo -e "  \033[31mFAIL\033[0m $test_name"
        echo "         command returned non-zero: $*"
        ((TESTS_FAILED++))
        FAILURES+=("$test_name")
    fi
}

assert_false() {
    local test_name="$1"
    shift
    ((TESTS_RUN++))

    if "$@" >/dev/null 2>&1; then
        echo -e "  \033[31mFAIL\033[0m $test_name"
        echo "         command should have failed but returned 0: $*"
        ((TESTS_FAILED++))
        FAILURES+=("$test_name")
    else
        echo -e "  \033[32mPASS\033[0m $test_name"
        ((TESTS_PASSED++))
    fi
}

section() {
    echo ""
    echo -e "\033[36m── $1 ──\033[0m"
}

#endregion

#region Extract Functions

# Source only the helper functions from ralph-loop.sh (not the main script body).
# We extract from the Helper Functions region to avoid triggering argument parsing.
extract_functions() {
    # Extract get_incomplete_task_count
    sed -n '/^get_incomplete_task_count()/,/^}/p' "$SOURCE_SCRIPT"
    # Extract is_ignorable_agent_output_line
    sed -n '/^is_ignorable_agent_output_line()/,/^}/p' "$SOURCE_SCRIPT"
    # Extract is_terminal_agent_error_output
    sed -n '/^is_terminal_agent_error_output()/,/^}/p' "$SOURCE_SCRIPT"
    # Extract initialize_progress_file
    sed -n '/^initialize_progress_file()/,/^}/p' "$SOURCE_SCRIPT"
    # Extract agent prompt/command builders
    sed -n '/^build_iteration_prompt()/,/^}/p' "$SOURCE_SCRIPT"
    sed -n '/^build_agent_command()/,/^}/p' "$SOURCE_SCRIPT"
    # Extract test_completion_signal
    sed -n '/^test_completion_signal()/,/^}/p' "$SOURCE_SCRIPT"
    # Extract load_ralph_config
    sed -n '/^load_ralph_config()/,/^}/p' "$SOURCE_SCRIPT"
}

eval "$(extract_functions)"

#endregion

#region Tests: get_incomplete_task_count

section "get_incomplete_task_count"

# Missing file → 0
result=$(get_incomplete_task_count "/tmp/_ralph_test_nonexistent_$$")
assert_eq "missing file returns 0" "0" "$result"

# Empty file → 0
TMPFILE=$(mktemp)
result=$(get_incomplete_task_count "$TMPFILE")
assert_eq "empty file returns 0" "0" "$result"
rm -f "$TMPFILE"

# File with no task checkboxes → 0
result=$(get_incomplete_task_count "$FIXTURE_DIR/tasks-empty.md")
assert_eq "no checkboxes returns 0" "0" "$result"

# All tasks done → 0
result=$(get_incomplete_task_count "$FIXTURE_DIR/tasks-all-done.md")
assert_eq "all done returns 0" "0" "$result"

# Mixed tasks → correct incomplete count (3)
result=$(get_incomplete_task_count "$FIXTURE_DIR/tasks-mixed.md")
assert_eq "mixed tasks returns 3" "3" "$result"

# Result is a valid integer for arithmetic comparison
result=$(get_incomplete_task_count "$FIXTURE_DIR/tasks-empty.md")
assert_true "result is arithmetic-safe (0)" test "$result" -eq 0

result=$(get_incomplete_task_count "$FIXTURE_DIR/tasks-mixed.md")
assert_true "result is arithmetic-safe (3)" test "$result" -eq 3

# Single-line result (no double output regression)
result=$(get_incomplete_task_count "$FIXTURE_DIR/tasks-empty.md")
line_count=$(echo "$result" | wc -l | tr -d ' ')
assert_eq "single-line output (regression #1)" "1" "$line_count"

#endregion

#region Tests: is_ignorable_agent_output_line

section "is_ignorable_agent_output_line"

assert_true "filters defaultPrompt manifest warning" \
    is_ignorable_agent_output_line \
    "2026-04-16T19:39:57.702566Z  WARN codex_core::plugins::manifest: ignoring interface.defaultPrompt: maximum of 3 prompts is supported path=/tmp/plugin.json"

assert_true "filters defaultPrompt manifest warning from codex_core_plugins logger" \
    is_ignorable_agent_output_line \
    "2026-04-22T20:16:53.964590Z  WARN codex_core_plugins::manifest: ignoring interface.defaultPrompt: prompt must be at most 128 characters path=/tmp/plugin.json"

assert_true "filters indexed defaultPrompt manifest warning from codex_core_plugins logger" \
    is_ignorable_agent_output_line \
    "2026-06-25T21:09:45.862247Z  WARN codex_core_plugins::manifest: ignoring interface.defaultPrompt[0]: prompt must be at most 128 characters path=/Users/adam.boczek/.codex/.tmp/plugins/plugins/ngs-analysis/.codex-plugin/plugin.json"

assert_true "filters skill icon loader warning" \
    is_ignorable_agent_output_line \
    "2026-05-23T14:19:28.001666Z  WARN codex_core_skills::loader: ignoring interface.icon_small: icon path must not contain '..'"

assert_true "filters skill icon loader warning variants" \
    is_ignorable_agent_output_line \
    "2026-05-23T14:19:28.001674Z  WARN codex_core_skills::loader: ignoring interface.icon_large: unsupported icon format"

assert_true "filters state db migration warning" \
    is_ignorable_agent_output_line \
    "2026-04-19T04:53:03.940369Z  WARN codex_state::runtime: failed to open state db at /Users/adam.boczek/.codex/state_5.sqlite: migration 23 was previously applied but is missing in the resolved migrations"

assert_true "filters rollout fallback warning" \
    is_ignorable_agent_output_line \
    "2026-04-19T04:53:03.963265Z  WARN codex_rollout::list: state db discrepancy during find_thread_path_by_id_str_in_subdir: falling_back"

assert_true "filters shell snapshot deletion warning" \
    is_ignorable_agent_output_line \
    '2026-04-19T04:53:05.067222Z  WARN codex_core::shell_snapshot: Failed to delete shell snapshot at "/Users/adam.boczek/.codex/shell_snapshots/019da415-c6fb-7d40-9860-ea76cf5d13e5.tmp-1776574383885855000": Os { code: 2, kind: NotFound, message: "No such file or directory" }'

assert_false "keeps unrelated warning lines" \
    is_ignorable_agent_output_line \
    "2026-04-16T19:39:57.702566Z  WARN codex_core::plugins::manifest: failed to load plugin manifest"

assert_false "keeps normal agent output" \
    is_ignorable_agent_output_line \
    "<promise>COMPLETE</promise>"

#endregion

#region Tests: is_terminal_agent_error_output

section "is_terminal_agent_error_output"

assert_true "detects unavailable model error" \
    is_terminal_agent_error_output \
    'Error: Model "gpt-5.5" from --model flag is not available.'

assert_true "detects missing named agent error" \
    is_terminal_agent_error_output \
    "No such agent: speckit.ralph.iterate, available:"

assert_false "keeps ordinary agent output" \
    is_terminal_agent_error_output \
    "I updated one task and committed the change."

assert_false "keeps retryable-looking generic errors" \
    is_terminal_agent_error_output \
    "Error: temporary network issue"

#endregion

#region Tests: test_completion_signal

section "test_completion_signal"

assert_false "rejects COMPLETE signal embedded in prose" test_completion_signal "Some output <promise>COMPLETE</promise> more text"

assert_false "rejects output without signal" test_completion_signal "Some output without the signal"

assert_false "rejects empty string" test_completion_signal ""

assert_true "detects signal on its own line" test_completion_signal "line1
<promise>COMPLETE</promise>
line3"

#endregion

#region Tests: build_agent_command

section "build_agent_command"

TMP_ITERATE_COMMAND=$(mktemp)
printf '%s\n' "ITERATE INSTRUCTIONS" > "$TMP_ITERATE_COMMAND"
ITERATE_COMMAND_PATH="$TMP_ITERATE_COMMAND"
AGENT_CLI="copilot"

build_agent_command "copilot" "test-model" "Runtime prompt"
assert_eq "copilot command uses configured CLI" "copilot" "${AGENT_COMMAND[0]}"
assert_eq "copilot command uses prompt mode" "-p" "${AGENT_COMMAND[1]}"
assert_false "copilot command does not require named agent" bash -c 'for arg in "$@"; do [[ "$arg" == "--agent" ]] && exit 0; done; exit 1' _ "${AGENT_COMMAND[@]}"
assert_true "copilot prompt embeds iterate instructions" bash -c '[[ "$1" == *"ITERATE INSTRUCTIONS"* ]]' _ "${AGENT_COMMAND[2]}"
assert_true "copilot prompt embeds runtime input" bash -c '[[ "$1" == *"Runtime prompt"* ]]' _ "${AGENT_COMMAND[2]}"

rm -f "$TMP_ITERATE_COMMAND"

#endregion

#region Tests: load_ralph_config

section "load_ralph_config"

# Create a temp directory mimicking the expected config structure
TMP_REPO=$(mktemp -d)
CONFIG_DIR="$TMP_REPO/.specify/extensions/ralph"
mkdir -p "$CONFIG_DIR"
cp "$FIXTURE_DIR/ralph-config-valid.yml" "$CONFIG_DIR/ralph-config.yml"

# Reset config variables
CONFIG_MODEL=""
CONFIG_MAX_ITERATIONS=""
CONFIG_AGENT_CLI=""

load_ralph_config "$TMP_REPO"

assert_eq "loads model from config" "gpt-4o" "$CONFIG_MODEL"
assert_eq "loads max_iterations from config" "5" "$CONFIG_MAX_ITERATIONS"
assert_eq "loads agent_cli from config" "my-custom-cli" "$CONFIG_AGENT_CLI"

# Reset and test with missing config
CONFIG_MODEL=""
CONFIG_MAX_ITERATIONS=""
CONFIG_AGENT_CLI=""

load_ralph_config "/tmp/_ralph_test_no_config_$$"

assert_eq "missing config leaves model empty" "" "$CONFIG_MODEL"
assert_eq "missing config leaves max_iterations empty" "" "$CONFIG_MAX_ITERATIONS"
assert_eq "missing config leaves agent_cli empty" "" "$CONFIG_AGENT_CLI"

# Test local config overrides project config
cat > "$CONFIG_DIR/ralph-config.local.yml" << 'LOCALCFG'
model: "local-model"
max_iterations: 20
LOCALCFG

CONFIG_MODEL=""
CONFIG_MAX_ITERATIONS=""
CONFIG_AGENT_CLI=""

load_ralph_config "$TMP_REPO"

assert_eq "local config overrides model" "local-model" "$CONFIG_MODEL"
assert_eq "local config overrides max_iterations" "20" "$CONFIG_MAX_ITERATIONS"
assert_eq "local config inherits agent_cli from project" "my-custom-cli" "$CONFIG_AGENT_CLI"

rm -rf "$TMP_REPO"

#endregion

#region Tests: initialize_progress_file

section "initialize_progress_file"

TMP_PROGRESS=$(mktemp -d)

# Creates file when missing
PROGRESS_FILE="$TMP_PROGRESS/progress.md"
initialize_progress_file "$PROGRESS_FILE" "test-feature" >/dev/null 2>&1
assert_true "creates progress file" test -f "$PROGRESS_FILE"

# File contains expected header content
assert_true "contains feature name" grep -q "Feature: test-feature" "$PROGRESS_FILE"
assert_true "contains codebase patterns section" grep -q "## Codebase Patterns" "$PROGRESS_FILE"

# Doesn't overwrite existing file
echo "custom content" > "$PROGRESS_FILE"
initialize_progress_file "$PROGRESS_FILE" "other-feature" >/dev/null 2>&1
content=$(cat "$PROGRESS_FILE")
assert_eq "does not overwrite existing file" "custom content" "$content"

rm -rf "$TMP_PROGRESS"

#endregion

#region Summary

echo ""
echo -e "\033[36m══════════════════════════════════\033[0m"
echo -e "\033[36m  Bash Regression Test Summary\033[0m"
echo -e "\033[36m══════════════════════════════════\033[0m"
echo -e "  Total:  $TESTS_RUN"
echo -e "  Passed: \033[32m$TESTS_PASSED\033[0m"
echo -e "  Failed: \033[31m$TESTS_FAILED\033[0m"

if [[ $TESTS_FAILED -gt 0 ]]; then
    echo ""
    echo -e "\033[31mFailed tests:\033[0m"
    for f in "${FAILURES[@]}"; do
        echo "  - $f"
    done
    exit 1
fi

echo ""
echo -e "\033[32mAll tests passed.\033[0m"
exit 0

#endregion
