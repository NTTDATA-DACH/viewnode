<#
.SYNOPSIS
    Ralph loop orchestrator for autonomous implementation.

.DESCRIPTION
    Executes a supported agent CLI in a controlled loop, processing tasks from tasks.md.
    Each iteration spawns a fresh agent context with the Ralph iteration instructions.
    
    The loop terminates when:
    - Agent outputs <promise>COMPLETE</promise>
    - Max iterations reached
    - All tasks in tasks.md are complete
    - User interrupts with Ctrl+C

    Configuration precedence (highest to lowest):
    1. Script parameters (always win when explicitly provided)
    2. Environment variables (SPECKIT_RALPH_MODEL, SPECKIT_RALPH_MAX_ITERATIONS, SPECKIT_RALPH_AGENT_CLI)
    3. Local config (.specify/extensions/ralph/ralph-config.local.yml)
    4. Project config (.specify/extensions/ralph/ralph-config.yml)
    5. Extension defaults (hardcoded parameter defaults)

.PARAMETER FeatureName
    Name of the feature being implemented (e.g., "001-ralph-loop-implement")

.PARAMETER TasksPath
    Path to tasks.md file

.PARAMETER SpecDir
    Path to the spec directory containing plan.md, spec.md, etc.

.PARAMETER MaxIterations
    Maximum number of iterations before stopping (default: 10)

.PARAMETER Model
    AI model to use (default: claude-sonnet-4.6)

.PARAMETER AgentCli
    Path or name of the agent CLI binary (default: copilot)

.PARAMETER AgentBackend
    Agent backend profile to use (auto, copilot, codex). Defaults to auto-detect from AgentCli.

.PARAMETER DetailedOutput
    Show detailed iteration output

.EXAMPLE
    .\ralph-loop.ps1 -FeatureName "001-feature" -TasksPath "specs/001-feature/tasks.md" -SpecDir "specs/001-feature" -MaxIterations 10 -Model "claude-sonnet-4.6"
#>

[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [string]$FeatureName,
    
    [Parameter(Mandatory = $true)]
    [string]$TasksPath,
    
    [Parameter(Mandatory = $true)]
    [string]$SpecDir,
    
    [Parameter(Mandatory = $false)]
    [int]$MaxIterations = 10,
    
    [Parameter(Mandatory = $false)]
    [string]$Model = "claude-sonnet-4.6",
    
    [Parameter(Mandatory = $false)]
    [string]$AgentCli = "copilot",

    [Parameter(Mandatory = $false)]
    [string]$AgentBackend = "auto",
    
    [Parameter(Mandatory = $false)]
    [string]$WorkingDirectory = "",
    
    [switch]$DetailedOutput
)

# Resolve working directory - if not provided, use the directory containing tasks.md
if (-not $WorkingDirectory) {
    # Infer from TasksPath - go up to find the repo root (directory with .git or .specify)
    $taskDir = Split-Path -Parent ([System.IO.Path]::GetFullPath($TasksPath))
    $searchDir = $taskDir
    while ($searchDir -and -not (Test-Path (Join-Path $searchDir ".git")) -and -not (Test-Path (Join-Path $searchDir ".specify"))) {
        $parent = Split-Path -Parent $searchDir
        if ($parent -eq $searchDir) { break }
        $searchDir = $parent
    }
    if ($searchDir -and ((Test-Path (Join-Path $searchDir ".git")) -or (Test-Path (Join-Path $searchDir ".specify")))) {
        $WorkingDirectory = $searchDir
    } else {
        $WorkingDirectory = (Get-Location).Path
    }
}

# Resolve paths
$RepoRoot = $WorkingDirectory
$TasksPath = [System.IO.Path]::GetFullPath($TasksPath)
$SpecDir = [System.IO.Path]::GetFullPath($SpecDir)
$ProgressPath = Join-Path $SpecDir "progress.md"
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ExtensionRoot = [System.IO.Path]::GetFullPath((Join-Path $ScriptDir "../.."))
$IterateCommandPath = Join-Path $ExtensionRoot "commands/iterate.md"

# Load config from extension config file
function Read-RalphConfig {
    param([string]$RepoRoot)
    
    $config = @{}
    $configPaths = @(
        (Join-Path $RepoRoot ".specify/extensions/ralph/ralph-config.yml"),
        (Join-Path $RepoRoot ".specify/extensions/ralph/ralph-config.local.yml")
    )
    
    foreach ($configPath in $configPaths) {
        if (Test-Path $configPath) {
            Get-Content $configPath | ForEach-Object {
                $line = $_.Trim()
                if ($line -and -not $line.StartsWith('#') -and $line -match '^(\w+)\s*:\s*"?(.+?)"?\s*$') {
                    $config[$Matches[1]] = $Matches[2]
                }
            }
        }
    }
    
    return $config
}

# Apply config defaults (only when script parameter was not explicitly provided)
$config = Read-RalphConfig -RepoRoot $RepoRoot

# Check if parameters were explicitly provided via PSBoundParameters
if (-not $PSBoundParameters.ContainsKey('Model') -and $config.ContainsKey('model')) {
    $Model = $config['model']
}
if (-not $PSBoundParameters.ContainsKey('MaxIterations') -and $config.ContainsKey('max_iterations')) {
    $MaxIterations = [int]$config['max_iterations']
}
if (-not $PSBoundParameters.ContainsKey('AgentCli') -and $config.ContainsKey('agent_cli')) {
    $AgentCli = $config['agent_cli']
}
if (-not $PSBoundParameters.ContainsKey('AgentBackend') -and $config.ContainsKey('agent_backend')) {
    $AgentBackend = $config['agent_backend']
}

# Environment variable overrides (higher priority than config, lower than explicit params)
if (-not $PSBoundParameters.ContainsKey('Model') -and $env:SPECKIT_RALPH_MODEL) {
    $Model = $env:SPECKIT_RALPH_MODEL
}
if (-not $PSBoundParameters.ContainsKey('MaxIterations') -and $env:SPECKIT_RALPH_MAX_ITERATIONS) {
    $MaxIterations = [int]$env:SPECKIT_RALPH_MAX_ITERATIONS
}
if (-not $PSBoundParameters.ContainsKey('AgentCli') -and $env:SPECKIT_RALPH_AGENT_CLI) {
    $AgentCli = $env:SPECKIT_RALPH_AGENT_CLI
}
if (-not $PSBoundParameters.ContainsKey('AgentBackend') -and $env:SPECKIT_RALPH_AGENT_BACKEND) {
    $AgentBackend = $env:SPECKIT_RALPH_AGENT_BACKEND
}

#region Helper Functions

function Write-RalphHeader {
    param([int]$Iteration, [int]$Max)
    
    $border = "=" * 60
    Write-Host ""
    Write-Host $border -ForegroundColor Cyan
    Write-Host "  Ralph Loop - $FeatureName" -ForegroundColor Cyan
    Write-Host "  Iteration $Iteration of $Max" -ForegroundColor White
    Write-Host $border -ForegroundColor Cyan
    Write-Host ""
}

function Write-IterationStatus {
    param(
        [int]$Iteration,
        [string]$Status,
        [string]$Message
    )
    
    $timestamp = Get-Date -Format "HH:mm:ss"
    $statusIcon = switch ($Status) {
        "running"   { "o" }
        "success"   { "*" }
        "failure"   { "x" }
        "skipped"   { "-" }
        default     { "o" }
    }
    $statusColor = switch ($Status) {
        "running"   { "Cyan" }
        "success"   { "Green" }
        "failure"   { "Red" }
        "skipped"   { "Yellow" }
        default     { "White" }
    }
    
    Write-Host "[$timestamp] " -NoNewline -ForegroundColor DarkGray
    Write-Host $statusIcon -NoNewline -ForegroundColor $statusColor
    Write-Host " Iteration $Iteration" -NoNewline -ForegroundColor White
    if ($Message) {
        Write-Host " - $Message" -ForegroundColor DarkGray
    } else {
        Write-Host ""
    }
}

function Get-IncompleteTasks {
    param([string]$Path)
    
    if (-not (Test-Path $Path)) {
        return @()
    }
    
    $content = Get-Content $Path -Raw
    $taskMatches = [regex]::Matches($content, '- \[ \] (T\d+.*?)(?=\r?\n|$)')
    
    return $taskMatches | ForEach-Object { $_.Groups[1].Value }
}

function Get-IncompleteTaskCount {
    param([string]$Path)
    
    return (Get-IncompleteTasks -Path $Path).Count
}

function Test-IgnorableAgentOutputLine {
    param([string]$Line)

    # Codex may emit non-fatal startup/runtime warnings on stderr before the
    # agent starts producing useful output. Keep this filter narrowly scoped
    # to known noisy warnings so real agent failures still surface.
    return (
        $Line -like "*WARN codex_core::plugins::manifest: ignoring interface.defaultPrompt*" -or
        $Line -like "*WARN codex_core_plugins::manifest: ignoring interface.defaultPrompt*" -or
        $Line -like "*WARN codex_core_skills::loader: ignoring interface.icon_*" -or
        $Line -like "*WARN codex_state::runtime: failed to open state db at *: migration 23 was previously applied but is missing in the resolved migrations*" -or
        $Line -like "*WARN codex_rollout::list: state db discrepancy during find_thread_path_by_id_str_in_subdir: falling_back*" -or
        $Line -like '*WARN codex_core::shell_snapshot: Failed to delete shell snapshot at *: Os { code: 2, kind: NotFound, message: "No such file or directory" }*'
    )
}

function Test-TerminalAgentErrorOutput {
    param([string]$Output)

    # Some agent CLIs report configuration errors in output while still
    # returning zero. These cannot be fixed by retrying another iteration.
    return (
        $Output -like '*Error: Model "*"*is not available.*' -or
        $Output -like '*No such agent:*'
    )
}

function Initialize-ProgressFile {
    param([string]$Path, [string]$Feature)
    
    if (-not (Test-Path $Path)) {
        $timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
        $header = @"
# Ralph Progress Log

Feature: $Feature
Started: $timestamp

## Codebase Patterns

[Patterns discovered during implementation - updated by agent]

---

"@
        Set-Content -Path $Path -Value $header -Encoding UTF8
        Write-Host "Created progress log: $Path" -ForegroundColor DarkGray
    }
}

function Resolve-AgentBackend {
    param(
        [string]$ConfiguredBackend,
        [string]$CliPath
    )

    # Allow an explicit backend override for wrapper scripts or renamed binaries.
    if ($ConfiguredBackend -and $ConfiguredBackend -ne "auto") {
        return $ConfiguredBackend.ToLowerInvariant()
    }

    $cliName = [System.IO.Path]::GetFileName($CliPath).ToLowerInvariant()
    switch ($cliName) {
        "copilot" { return "copilot" }
        "codex" { return "codex" }
        default { return "unknown" }
    }
}

function Get-IterationPrompt {
    param([string]$Prompt)

    # Backends receive the iterate instructions inline so the loop does not
    # depend on a named agent being registered in the target repository.
    if (-not (Test-Path $IterateCommandPath)) {
        return $Prompt
    }

    $instructions = Get-Content $IterateCommandPath -Raw
    return @"
You are executing the speckit.ralph.iterate workflow for this repository.
Follow the instructions below exactly.

$instructions

Runtime input for `$ARGUMENTS:
$Prompt
"@
}

function Get-AgentCommandArguments {
    param(
        [string]$Backend,
        [string]$Model,
        [string]$Prompt
    )

    # Map Ralph's generic loop contract onto each supported CLI's invocation shape.
    switch ($Backend) {
        "copilot" {
            $fullPrompt = Get-IterationPrompt -Prompt $Prompt
            return @("-p", $fullPrompt, "--model", $Model, "--yolo", "-s")
        }
        "codex" {
            $fullPrompt = Get-IterationPrompt -Prompt $Prompt
            return @("exec", "--full-auto", "--model", $Model, $fullPrompt)
        }
        default {
            throw "Unsupported Ralph agent backend '$Backend' for CLI '$AgentCli'. Supported backends: copilot, codex."
        }
    }
}

function Invoke-AgentIteration {
    param(
        [string]$Model,
        [int]$Iteration,
        [string]$WorkDir,
        [string]$Backend,
        [switch]$Verbose
    )
    
    # Simple prompt - backend-specific wrappers provide the full Ralph instructions when needed
    $prompt = "Iteration $Iteration - Complete one work unit from tasks.md"
    $commandArgs = Get-AgentCommandArguments -Backend $Backend -Model $Model -Prompt $prompt
    
    # Only show debug info when verbose
    if ($Verbose) {
        Write-Host "DEBUG: Prompt = $prompt" -ForegroundColor Magenta
        Write-Host "DEBUG: WorkDir = $WorkDir" -ForegroundColor Magenta
        Write-Host "DEBUG: AgentCli = $AgentCli" -ForegroundColor Magenta
        Write-Host "DEBUG: AgentBackend = $Backend" -ForegroundColor Magenta
    }
    
    try {
        # Change to working directory before invoking the agent CLI
        $originalDir = Get-Location
        if ($WorkDir -and (Test-Path $WorkDir)) {
            Push-Location $WorkDir
            if ($Verbose) {
                Write-Host "DEBUG: Changed to $WorkDir" -ForegroundColor Magenta
            }
        }
        
        # Refresh PATH to ensure shell-installed CLIs remain discoverable
        $env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path", "User")
        
        # Ensure UTF-8 so agent output renders correctly
        $prevOutputEncoding = $OutputEncoding
        $prevConsoleEncoding = [Console]::OutputEncoding
        $OutputEncoding = [System.Text.Encoding]::UTF8
        [Console]::OutputEncoding = [System.Text.Encoding]::UTF8
        
        try {
            # Always stream agent output in real-time so user can see what the agent is doing
            Write-Host ""
            $backendLabel = (Get-Culture).TextInfo.ToTitleCase($Backend)
            Write-Host "--- $backendLabel Agent Output ---" -ForegroundColor DarkCyan
            $outputLines = @()
            & $AgentCli @commandArgs 2>&1 | ForEach-Object {
                # Stderr lines arrive as ErrorRecord objects; extract the message string
                $line = if ($_ -is [System.Management.Automation.ErrorRecord]) { $_.Exception.Message } else { $_ }
                if (Test-IgnorableAgentOutputLine -Line $line) {
                    return
                }
                Write-Host $line
                $outputLines += $line
            }
            $output = $outputLines -join "`n"
            $exitCode = $LASTEXITCODE
            Write-Host "--- End Agent Output ---" -ForegroundColor DarkCyan
            Write-Host ""
        }
        finally {
            $OutputEncoding = $prevOutputEncoding
            [Console]::OutputEncoding = $prevConsoleEncoding
            if ($WorkDir -and (Test-Path $WorkDir)) {
                Pop-Location
            }
        }
        
        if ($Verbose) {
            Write-Host "DEBUG: $AgentCli exit code = $exitCode" -ForegroundColor Magenta
        }
    }
    catch {
        $output = "Error invoking agent CLI: $_"
        $exitCode = 1
    }
    
    return @{
        Output = $output
        ExitCode = $exitCode
    }
}

function Test-CompletionSignal {
    param([string]$Output)
    
    return $Output -match '<promise>COMPLETE</promise>'
}

#endregion

#region Main Loop

# Initialize progress file
Initialize-ProgressFile -Path $ProgressPath -Feature $FeatureName

# Check initial task count
$initialTasks = Get-IncompleteTaskCount -Path $TasksPath
if ($initialTasks -eq 0) {
    Write-Host "All tasks are already complete!" -ForegroundColor Green
    Write-Host "<promise>COMPLETE</promise>"
    exit 0
}

Write-Host "Found $initialTasks incomplete task(s)" -ForegroundColor White

$resolvedAgentBackend = Resolve-AgentBackend -ConfiguredBackend $AgentBackend -CliPath $AgentCli
# Fail fast when the configured CLI cannot be matched to a supported invocation profile.
if ($resolvedAgentBackend -eq "unknown") {
    Write-Host "Unable to determine a supported agent backend for '$AgentCli'." -ForegroundColor Red
    Write-Host "Set agent_backend to 'copilot' or 'codex' in config, or via SPECKIT_RALPH_AGENT_BACKEND." -ForegroundColor Red
    exit 1
}

Write-Host "Using agent backend: $resolvedAgentBackend ($AgentCli)" -ForegroundColor White

# Iteration tracking
$iteration = 1
$consecutiveFailures = 0
$maxConsecutiveFailures = 3
$completed = $false
$circuitBreaker = $false

# Register Ctrl+C handler
$interrupted = $false
[Console]::TreatControlCAsInput = $false

try {
    while ($iteration -le $MaxIterations -and -not $completed -and -not $interrupted -and -not $circuitBreaker) {
        Write-RalphHeader -Iteration $iteration -Max $MaxIterations
        Write-IterationStatus -Iteration $iteration -Status "running" -Message "Starting iteration"
        
        # Invoke the configured agent backend for one Ralph iteration
        $verboseSwitch = @{}
        if ($DetailedOutput) { $verboseSwitch['Verbose'] = $true }
        $result = Invoke-AgentIteration -Model $Model -Iteration $iteration -WorkDir $WorkingDirectory -Backend $resolvedAgentBackend @verboseSwitch
        
        if ($DetailedOutput -or $result.ExitCode -ne 0) {
            Write-Host $result.Output
        }

        if (Test-TerminalAgentErrorOutput -Output $result.Output) {
            Write-IterationStatus -Iteration $iteration -Status "failure" -Message "Terminal agent configuration error"
            Write-Host "Stopping loop because the agent reported a non-retryable configuration error." -ForegroundColor Red
            $circuitBreaker = $true
            break
        }
        
        # Check for completion signal
        if (Test-CompletionSignal -Output $result.Output) {
            Write-IterationStatus -Iteration $iteration -Status "success" -Message "COMPLETE signal received"
            $completed = $true
            break
        }
        
        # Check exit code
        if ($result.ExitCode -ne 0) {
            $consecutiveFailures++
            Write-IterationStatus -Iteration $iteration -Status "failure" -Message "Exit code $($result.ExitCode) (failure $consecutiveFailures/$maxConsecutiveFailures)"
            
            if ($consecutiveFailures -ge $maxConsecutiveFailures) {
                Write-Host "Too many consecutive failures. Stopping loop." -ForegroundColor Red
                $circuitBreaker = $true
                break
            }
        } else {
            $consecutiveFailures = 0
            Write-IterationStatus -Iteration $iteration -Status "success" -Message "Iteration completed"
        }
        
        # Check remaining tasks
        $remainingTasks = Get-IncompleteTaskCount -Path $TasksPath
        if ($remainingTasks -eq 0) {
            Write-Host "All tasks complete!" -ForegroundColor Green
            $completed = $true
            break
        }
        
        Write-Host "$remainingTasks task(s) remaining" -ForegroundColor DarkGray
        
        $iteration++
    }
}
catch {
    if ($_.Exception.GetType().Name -eq "PipelineStoppedException") {
        $interrupted = $true
        Write-Host "`nInterrupted by user" -ForegroundColor Yellow
    } else {
        throw
    }
}
finally {
    [Console]::TreatControlCAsInput = $false
}

#endregion

#region Summary

Write-Host ""
Write-Host ("=" * 60) -ForegroundColor Cyan
Write-Host "  Ralph Loop Summary" -ForegroundColor Cyan
Write-Host ("=" * 60) -ForegroundColor Cyan

$finalTasks = Get-IncompleteTaskCount -Path $TasksPath
$tasksCompleted = $initialTasks - $finalTasks
$iterationsRun = if ($completed) {
    $iteration
} elseif ($circuitBreaker) {
    $iteration
} else {
    [Math]::Max(0, $iteration - 1)
}

Write-Host "  Iterations run: $iterationsRun" -ForegroundColor White
Write-Host "  Tasks completed: $tasksCompleted" -ForegroundColor White
Write-Host "  Tasks remaining: $finalTasks" -ForegroundColor White

if ($completed) {
    Write-Host "  Status: " -NoNewline -ForegroundColor White
    Write-Host "COMPLETED" -ForegroundColor Green
    exit 0
} elseif ($interrupted) {
    Write-Host "  Status: " -NoNewline -ForegroundColor White
    Write-Host "INTERRUPTED" -ForegroundColor Yellow
    exit 130
} elseif ($circuitBreaker) {
    Write-Host "  Status: " -NoNewline -ForegroundColor White
    Write-Host "FAILED" -ForegroundColor Red
    exit 1
} else {
    Write-Host "  Status: " -NoNewline -ForegroundColor White
    Write-Host "ITERATION LIMIT REACHED" -ForegroundColor Yellow
    exit 1
}

#endregion
