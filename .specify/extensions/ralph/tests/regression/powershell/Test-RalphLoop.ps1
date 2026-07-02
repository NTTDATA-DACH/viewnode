<#
.SYNOPSIS
    Regression tests for ralph-loop.ps1 helper functions.

.DESCRIPTION
    Extracts and tests functions in isolation without running the full loop.
    Uses a lightweight assertion harness (no Pester dependency).

.EXAMPLE
    pwsh -File tests/regression/powershell/Test-RalphLoop.ps1
#>

$ErrorActionPreference = "Stop"

$ScriptDir = $PSScriptRoot
$RepoRoot = Resolve-Path (Join-Path $ScriptDir "..\..\..") | Select-Object -ExpandProperty Path
$FixtureDir = Join-Path $ScriptDir "..\fixtures"
$SourceScript = Join-Path $RepoRoot "scripts\powershell\ralph-loop.ps1"

# Test bookkeeping
$script:TestsRun = 0
$script:TestsPassed = 0
$script:TestsFailed = 0
$script:Failures = @()

#region Test Harness

function Assert-Equal {
    param([string]$TestName, $Expected, $Actual)
    $script:TestsRun++

    if ($Expected -eq $Actual) {
        Write-Host "  PASS " -NoNewline -ForegroundColor Green
        Write-Host $TestName
        $script:TestsPassed++
    } else {
        Write-Host "  FAIL " -NoNewline -ForegroundColor Red
        Write-Host $TestName
        Write-Host "         expected: [$Expected]"
        Write-Host "         actual:   [$Actual]"
        $script:TestsFailed++
        $script:Failures += $TestName
    }
}

function Assert-True {
    param([string]$TestName, [bool]$Condition)
    $script:TestsRun++

    if ($Condition) {
        Write-Host "  PASS " -NoNewline -ForegroundColor Green
        Write-Host $TestName
        $script:TestsPassed++
    } else {
        Write-Host "  FAIL " -NoNewline -ForegroundColor Red
        Write-Host $TestName
        $script:TestsFailed++
        $script:Failures += $TestName
    }
}

function Write-Section {
    param([string]$Name)
    Write-Host ""
    Write-Host "-- $Name --" -ForegroundColor Cyan
}

#endregion

#region Extract Functions

# Parse the source script to extract function definitions without executing the main body.
# We use AST parsing to safely extract only the function blocks.
$ast = [System.Management.Automation.Language.Parser]::ParseFile($SourceScript, [ref]$null, [ref]$null)
$functionDefs = $ast.FindAll({ param($node) $node -is [System.Management.Automation.Language.FunctionDefinitionAst] }, $false)

foreach ($funcDef in $functionDefs) {
    # Define each function in the current scope
    Invoke-Expression $funcDef.Extent.Text
}

#endregion

#region Tests: Get-IncompleteTaskCount

Write-Section "Get-IncompleteTaskCount"

# Missing file -> 0
$result = Get-IncompleteTaskCount -Path "C:\nonexistent_ralph_test_$PID.md"
Assert-Equal "missing file returns 0" 0 $result

# Empty file -> 0
$tmpFile = [System.IO.Path]::GetTempFileName()
Set-Content -Path $tmpFile -Value "" -Encoding UTF8
$result = Get-IncompleteTaskCount -Path $tmpFile
Assert-Equal "empty file returns 0" 0 $result
Remove-Item $tmpFile -Force

# No checkboxes -> 0
$result = Get-IncompleteTaskCount -Path (Join-Path $FixtureDir "tasks-empty.md")
Assert-Equal "no checkboxes returns 0" 0 $result

# All done -> 0
$result = Get-IncompleteTaskCount -Path (Join-Path $FixtureDir "tasks-all-done.md")
Assert-Equal "all done returns 0" 0 $result

# Mixed tasks -> correct count (3 incomplete with T\d+ pattern)
$result = Get-IncompleteTaskCount -Path (Join-Path $FixtureDir "tasks-mixed.md")
Assert-Equal "mixed tasks returns 3" 3 $result

#endregion

#region Tests: Test-IgnorableAgentOutputLine

Write-Section "Test-IgnorableAgentOutputLine"

Assert-True "filters defaultPrompt manifest warning" (
    Test-IgnorableAgentOutputLine -Line "2026-04-16T19:39:57.702566Z  WARN codex_core::plugins::manifest: ignoring interface.defaultPrompt: maximum of 3 prompts is supported path=/tmp/plugin.json"
)

Assert-True "filters defaultPrompt manifest warning from codex_core_plugins logger" (
    Test-IgnorableAgentOutputLine -Line "2026-04-22T20:16:53.964590Z  WARN codex_core_plugins::manifest: ignoring interface.defaultPrompt: prompt must be at most 128 characters path=/tmp/plugin.json"
)

Assert-True "filters indexed defaultPrompt manifest warning from codex_core_plugins logger" (
    Test-IgnorableAgentOutputLine -Line "2026-06-25T21:09:45.862247Z  WARN codex_core_plugins::manifest: ignoring interface.defaultPrompt[0]: prompt must be at most 128 characters path=/Users/adam.boczek/.codex/.tmp/plugins/plugins/ngs-analysis/.codex-plugin/plugin.json"
)

Assert-True "filters skill icon loader warning" (
    Test-IgnorableAgentOutputLine -Line "2026-05-23T14:19:28.001666Z  WARN codex_core_skills::loader: ignoring interface.icon_small: icon path must not contain '..'"
)

Assert-True "filters skill icon loader warning variants" (
    Test-IgnorableAgentOutputLine -Line "2026-05-23T14:19:28.001674Z  WARN codex_core_skills::loader: ignoring interface.icon_large: unsupported icon format"
)

Assert-True "filters state db migration warning" (
    Test-IgnorableAgentOutputLine -Line "2026-04-19T04:53:03.940369Z  WARN codex_state::runtime: failed to open state db at /Users/adam.boczek/.codex/state_5.sqlite: migration 23 was previously applied but is missing in the resolved migrations"
)

Assert-True "filters rollout fallback warning" (
    Test-IgnorableAgentOutputLine -Line "2026-04-19T04:53:03.963265Z  WARN codex_rollout::list: state db discrepancy during find_thread_path_by_id_str_in_subdir: falling_back"
)

Assert-True "filters shell snapshot deletion warning" (
    Test-IgnorableAgentOutputLine -Line '2026-04-19T04:53:05.067222Z  WARN codex_core::shell_snapshot: Failed to delete shell snapshot at "/Users/adam.boczek/.codex/shell_snapshots/019da415-c6fb-7d40-9860-ea76cf5d13e5.tmp-1776574383885855000": Os { code: 2, kind: NotFound, message: "No such file or directory" }'
)

Assert-True "keeps unrelated warning lines" (
    -not (Test-IgnorableAgentOutputLine -Line "2026-04-16T19:39:57.702566Z  WARN codex_core::plugins::manifest: failed to load plugin manifest")
)

Assert-True "keeps normal agent output" (
    -not (Test-IgnorableAgentOutputLine -Line "<promise>COMPLETE</promise>")
)

#endregion

#region Tests: Test-TerminalAgentErrorOutput

Write-Section "Test-TerminalAgentErrorOutput"

Assert-True "detects unavailable model error" (
    Test-TerminalAgentErrorOutput -Output 'Error: Model "gpt-5.5" from --model flag is not available.'
)

Assert-True "detects missing named agent error" (
    Test-TerminalAgentErrorOutput -Output "No such agent: speckit.ralph.iterate, available:"
)

Assert-True "keeps ordinary agent output" (
    -not (Test-TerminalAgentErrorOutput -Output "I updated one task and committed the change.")
)

Assert-True "keeps retryable-looking generic errors" (
    -not (Test-TerminalAgentErrorOutput -Output "Error: temporary network issue")
)

#endregion

#region Tests: Get-IncompleteTasks

Write-Section "Get-IncompleteTasks"

# Returns correct task IDs from mixed fixture
$tasks = Get-IncompleteTasks -Path (Join-Path $FixtureDir "tasks-mixed.md")
Assert-Equal "returns 3 incomplete tasks" 3 $tasks.Count

# Verify task ID content
Assert-True "first task contains T002" ($tasks[0] -like "T002*")
Assert-True "second task contains T003" ($tasks[1] -like "T003*")
Assert-True "third task contains T006" ($tasks[2] -like "T006*")

# All done -> empty array
$tasks = Get-IncompleteTasks -Path (Join-Path $FixtureDir "tasks-all-done.md")
Assert-Equal "all done returns empty" 0 $tasks.Count

# Missing file -> empty array
$tasks = Get-IncompleteTasks -Path "C:\nonexistent_ralph_test_$PID.md"
Assert-Equal "missing file returns empty" 0 $tasks.Count

#endregion

#region Tests: Test-CompletionSignal

Write-Section "Test-CompletionSignal"

Assert-True "detects COMPLETE signal" (Test-CompletionSignal -Output "Some output <promise>COMPLETE</promise> more text")

Assert-True "rejects output without signal" (-not (Test-CompletionSignal -Output "Some output without the signal"))

Assert-True "rejects empty string" (-not (Test-CompletionSignal -Output ""))

$multiLine = "line1`n<promise>COMPLETE</promise>`nline3"
Assert-True "detects signal on its own line" (Test-CompletionSignal -Output $multiLine)

#endregion

#region Tests: Get-AgentCommandArguments

Write-Section "Get-AgentCommandArguments"

$tmpIterateCommand = [System.IO.Path]::GetTempFileName()
Set-Content -Path $tmpIterateCommand -Value "ITERATE INSTRUCTIONS" -Encoding UTF8
$IterateCommandPath = $tmpIterateCommand
$AgentCli = "copilot"

$commandArgs = Get-AgentCommandArguments -Backend "copilot" -Model "test-model" -Prompt "Runtime prompt"
Assert-Equal "copilot command uses prompt mode" "-p" $commandArgs[0]
Assert-True "copilot command does not require named agent" (-not ($commandArgs -contains "--agent"))
Assert-True "copilot prompt embeds iterate instructions" ($commandArgs[1] -like "*ITERATE INSTRUCTIONS*")
Assert-True "copilot prompt embeds runtime input" ($commandArgs[1] -like "*Runtime prompt*")

Remove-Item $tmpIterateCommand -Force

#endregion

#region Tests: Read-RalphConfig

Write-Section "Read-RalphConfig"

# Create temp directory with config structure
$tmpRepo = Join-Path ([System.IO.Path]::GetTempPath()) "ralph-test-$PID"
$configDir = Join-Path $tmpRepo ".specify\extensions\ralph"
New-Item -ItemType Directory -Path $configDir -Force | Out-Null
Copy-Item (Join-Path $FixtureDir "ralph-config-valid.yml") (Join-Path $configDir "ralph-config.yml")

$config = Read-RalphConfig -RepoRoot $tmpRepo

Assert-Equal "loads model from config" "gpt-4o" $config["model"]
Assert-Equal "loads max_iterations from config" "5" $config["max_iterations"]
Assert-Equal "loads agent_cli from config" "my-custom-cli" $config["agent_cli"]

# Missing config -> empty hashtable
$config = Read-RalphConfig -RepoRoot "C:\nonexistent_ralph_test_$PID"
Assert-Equal "missing config returns empty" 0 $config.Count

# Local config overrides project config
@"
model: "local-model"
max_iterations: 20
"@ | Set-Content (Join-Path $configDir "ralph-config.local.yml") -Encoding UTF8

$config = Read-RalphConfig -RepoRoot $tmpRepo

Assert-Equal "local config overrides model" "local-model" $config["model"]
Assert-Equal "local config overrides max_iterations" "20" $config["max_iterations"]
Assert-Equal "local config inherits agent_cli" "my-custom-cli" $config["agent_cli"]

Remove-Item $tmpRepo -Recurse -Force

#endregion

#region Tests: Initialize-ProgressFile

Write-Section "Initialize-ProgressFile"

$tmpDir = Join-Path ([System.IO.Path]::GetTempPath()) "ralph-progress-$PID"
New-Item -ItemType Directory -Path $tmpDir -Force | Out-Null

# Creates file when missing
$progressFile = Join-Path $tmpDir "progress.md"
Initialize-ProgressFile -Path $progressFile -Feature "test-feature"
Assert-True "creates progress file" (Test-Path $progressFile)

$content = Get-Content $progressFile -Raw
Assert-True "contains feature name" ($content -match "Feature: test-feature")
Assert-True "contains codebase patterns section" ($content -match "## Codebase Patterns")

# Doesn't overwrite existing file
Set-Content -Path $progressFile -Value "custom content" -Encoding UTF8
Initialize-ProgressFile -Path $progressFile -Feature "other-feature"
$content = (Get-Content $progressFile -Raw).Trim()
Assert-Equal "does not overwrite existing file" "custom content" $content

Remove-Item $tmpDir -Recurse -Force

#endregion

#region Summary

Write-Host ""
Write-Host ("=" * 40) -ForegroundColor Cyan
Write-Host "  PowerShell Regression Test Summary" -ForegroundColor Cyan
Write-Host ("=" * 40) -ForegroundColor Cyan
Write-Host "  Total:  $script:TestsRun"
Write-Host "  Passed: " -NoNewline; Write-Host "$script:TestsPassed" -ForegroundColor Green
Write-Host "  Failed: " -NoNewline; Write-Host "$script:TestsFailed" -ForegroundColor Red

if ($script:TestsFailed -gt 0) {
    Write-Host ""
    Write-Host "Failed tests:" -ForegroundColor Red
    foreach ($f in $script:Failures) {
        Write-Host "  - $f"
    }
    exit 1
}

Write-Host ""
Write-Host "All tests passed." -ForegroundColor Green
exit 0

#endregion
