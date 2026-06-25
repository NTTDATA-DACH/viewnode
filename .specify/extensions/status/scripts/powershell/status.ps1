#Requires -Version 5.1
Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

# ============================================================================
# Spec Kit Status — Show project status and SDD workflow progress
# ============================================================================

function Find-RepoRoot {
    param([string]$StartDir)
    $dir = $StartDir
    while ($dir -and $dir -ne [System.IO.Path]::GetPathRoot($dir)) {
        if ((Test-Path "$dir/.git") -or (Test-Path "$dir/.specify")) {
            return $dir
        }
        $dir = Split-Path $dir -Parent
    }
    return $null
}

function ConvertTo-JsonSafe {
    param([string]$Value)
    $Value = $Value -replace '\\', '\\' -replace '"', '\"' -replace "`n", '\n' -replace "`r", '\r' -replace "`t", '\t'
    return $Value
}

# --- Locate project root ----------------------------------------------------

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$HasGit = $false
$RepoRoot = $null

try {
    $gitRoot = git rev-parse --show-toplevel 2>$null
    if ($LASTEXITCODE -eq 0 -and $gitRoot) {
        $RepoRoot = $gitRoot.Trim()
        $HasGit = $true
    }
} catch {}

if (-not $RepoRoot) {
    $RepoRoot = Find-RepoRoot $ScriptDir
    if (-not $RepoRoot) {
        Write-Error '{"error":"Could not determine repository root"}'
        exit 1
    }
}

$SpecifyDir = Join-Path $RepoRoot '.specify'
if (-not (Test-Path $SpecifyDir)) {
    Write-Error '{"error":"Not a spec-kit project (no .specify/ directory)"}'
    exit 1
}

$ProjectName = Split-Path $RepoRoot -Leaf

# --- Detect AI agent(s) -----------------------------------------------------

$AgentDefs = @(
    @{ Key='copilot';      Name='GitHub Copilot'; Folder='.github';        CmdDir='agents' },
    @{ Key='claude';       Name='Claude Code';    Folder='.claude';        CmdDir='commands' },
    @{ Key='gemini';       Name='Gemini CLI';     Folder='.gemini';        CmdDir='commands' },
    @{ Key='cursor-agent'; Name='Cursor';         Folder='.cursor';        CmdDir='commands' },
    @{ Key='qwen';         Name='Qwen Code';      Folder='.qwen';          CmdDir='commands' },
    @{ Key='opencode';     Name='opencode';        Folder='.opencode';      CmdDir='command' },
    @{ Key='codex';        Name='Codex CLI';      Folder='.codex';         CmdDir='prompts' },
    @{ Key='windsurf';     Name='Windsurf';        Folder='.windsurf';      CmdDir='workflows' },
    @{ Key='kilocode';     Name='Kilo Code';      Folder='.kilocode';      CmdDir='workflows' },
    @{ Key='auggie';       Name='Auggie CLI';     Folder='.augment';       CmdDir='commands' },
    @{ Key='codebuddy';    Name='CodeBuddy';       Folder='.codebuddy';     CmdDir='commands' },
    @{ Key='qodercli';     Name='Qoder CLI';      Folder='.qoder';         CmdDir='commands' },
    @{ Key='roo';          Name='Roo Code';        Folder='.roo';           CmdDir='commands' },
    @{ Key='kiro-cli';     Name='Kiro CLI';       Folder='.kiro';          CmdDir='prompts' },
    @{ Key='amp';          Name='Amp';              Folder='.agents';        CmdDir='commands' },
    @{ Key='shai';         Name='SHAI';             Folder='.shai';          CmdDir='commands' },
    @{ Key='tabnine';      Name='Tabnine CLI';    Folder='.tabnine/agent'; CmdDir='commands' },
    @{ Key='agy';          Name='Antigravity';     Folder='.agent';         CmdDir='workflows' },
    @{ Key='bob';          Name='IBM Bob';         Folder='.bob';           CmdDir='commands' },
    @{ Key='vibe';         Name='Mistral Vibe';   Folder='.vibe';          CmdDir='prompts' },
    @{ Key='kimi';         Name='Kimi Code';      Folder='.kimi';          CmdDir='skills' }
)

$DetectedAgents = @()
foreach ($agent in $AgentDefs) {
    $agentDir = Join-Path $RepoRoot $agent.Folder
    if (Test-Path $agentDir) {
        $cmdDir = Join-Path $agentDir $agent.CmdDir
        $hasCmds = (Test-Path $cmdDir) -and (@(Get-ChildItem $cmdDir -ErrorAction SilentlyContinue).Count -gt 0)
        $DetectedAgents += @{
            key = $agent.Key
            name = $agent.Name
            folder = "$($agent.Folder)/"
            has_commands = $hasCmds
        }
    }
}

# --- Detect script type ------------------------------------------------------

$ScriptType = 'none'
$hasBash = Test-Path (Join-Path $SpecifyDir 'scripts/bash')
$hasPs = Test-Path (Join-Path $SpecifyDir 'scripts/powershell')
if ($hasBash -and $hasPs) { $ScriptType = 'sh + ps' }
elseif ($hasBash) { $ScriptType = 'sh' }
elseif ($hasPs) { $ScriptType = 'ps' }

# --- Detect current feature --------------------------------------------------

$CurrentBranch = ''
$SpecsDir = Join-Path $RepoRoot 'specs'

# 1. SPECIFY_FEATURE env var
$envFeature = $env:SPECIFY_FEATURE
if ($envFeature) { $CurrentBranch = $envFeature.Trim() }

# 2. Try git
if (-not $CurrentBranch -and $HasGit) {
    try {
        $branch = git rev-parse --abbrev-ref HEAD 2>$null
        if ($LASTEXITCODE -eq 0 -and $branch) { $CurrentBranch = $branch.Trim() }
    } catch {}
}

# 3. Fallback: scan specs/ for highest-numbered dir
if (-not $CurrentBranch -or $CurrentBranch -eq 'main' -or $CurrentBranch -eq 'master') {
    if (Test-Path $SpecsDir) {
        $highest = 0; $latest = ''
        foreach ($d in Get-ChildItem $SpecsDir -Directory -ErrorAction SilentlyContinue) {
            if ($d.Name -match '^(\d{3})-') {
                $num = [int]$Matches[1]
                if ($num -gt $highest) { $highest = $num; $latest = $d.Name }
            }
        }
        if ($latest -and (-not $CurrentBranch -or $CurrentBranch -eq 'main' -or $CurrentBranch -eq 'master')) {
            $CurrentBranch = $latest
        }
    }
}

# --- Resolve feature directory -----------------------------------------------

$FeatureDir = ''
$FeatureName = ''
if ($CurrentBranch -match '^(\d{3})-' -and (Test-Path $SpecsDir)) {
    $prefix = $Matches[1]
    $matches2 = @(Get-ChildItem $SpecsDir -Directory -ErrorAction SilentlyContinue | Where-Object { $_.Name -match "^$prefix-" })
    if ($matches2.Count -eq 1) {
        $FeatureDir = $matches2[0].FullName
    } elseif ($matches2.Count -gt 1) {
        $exact = Join-Path $SpecsDir $CurrentBranch
        if (Test-Path $exact) { $FeatureDir = $exact }
        else { $FeatureDir = $matches2[0].FullName }
    }
    if ($FeatureDir) { $FeatureName = Split-Path $FeatureDir -Leaf }
}

# --- Check SDD artifacts -----------------------------------------------------

$hasSpec = $false; $hasPlan = $false; $hasTasks = $false
$hasResearch = $false; $hasDataModel = $false; $hasQuickstart = $false
$hasContracts = $false; $contractCount = 0
$hasChecklists = $false; $checklistCount = 0; $checklistsAllPass = $true
$tasksCompleted = 0; $tasksTotal = 0

if ($FeatureDir -and (Test-Path $FeatureDir)) {
    $hasSpec = Test-Path (Join-Path $FeatureDir 'spec.md')
    $hasPlan = Test-Path (Join-Path $FeatureDir 'plan.md')
    $hasTasks = Test-Path (Join-Path $FeatureDir 'tasks.md')
    $hasResearch = Test-Path (Join-Path $FeatureDir 'research.md')
    $hasDataModel = Test-Path (Join-Path $FeatureDir 'data-model.md')
    $hasQuickstart = Test-Path (Join-Path $FeatureDir 'quickstart.md')

    $contractsDir = Join-Path $FeatureDir 'contracts'
    if ((Test-Path $contractsDir) -and (@(Get-ChildItem $contractsDir -File -ErrorAction SilentlyContinue).Count -gt 0)) {
        $hasContracts = $true
        $contractCount = @(Get-ChildItem $contractsDir -File -ErrorAction SilentlyContinue).Count
    }

    $checklistsDir = Join-Path $FeatureDir 'checklists'
    if (Test-Path $checklistsDir) {
        foreach ($cl in Get-ChildItem $checklistsDir -Filter '*.md' -File -ErrorAction SilentlyContinue) {
            $checklistCount++
            $hasChecklists = $true
            $content = Get-Content $cl.FullName -Raw -ErrorAction SilentlyContinue
            if ($content -match '^\s*-\s*\[ \]') { $checklistsAllPass = $false }
        }
    }

    # Parse task progress
    if ($hasTasks) {
        $tasksContent = Get-Content (Join-Path $FeatureDir 'tasks.md') -ErrorAction SilentlyContinue
        foreach ($line in $tasksContent) {
            $stripped = $line.TrimStart()
            if ($stripped -match '^-\s+\[([ xX])\]') {
                $tasksTotal++
                if ($stripped -match '^-\s+\[[xX]\]') { $tasksCompleted++ }
            }
        }
    }
}

# --- Detect workflow phase ---------------------------------------------------

$Phase = 'unknown'; $PhaseHint = ''
if (-not $FeatureDir -or -not (Test-Path $FeatureDir -ErrorAction SilentlyContinue)) {
    $Phase = 'no_feature'; $PhaseHint = 'Run /speckit.specify to create a feature'
} elseif ($hasTasks -and $tasksTotal -gt 0 -and $tasksCompleted -eq $tasksTotal) {
    $Phase = 'complete'; $PhaseHint = 'All tasks done. Review your implementation.'
} elseif ($hasTasks) {
    $Phase = 'implement'; $PhaseHint = 'Ready for /speckit.implement'
} elseif ($hasPlan) {
    $Phase = 'tasks'; $PhaseHint = 'Ready for /speckit.tasks'
} elseif ($hasSpec) {
    $Phase = 'plan'; $PhaseHint = 'Ready for /speckit.clarify or /speckit.plan'
} else {
    $Phase = 'not_started'; $PhaseHint = 'Run /speckit.specify to create a spec'
}

# --- Extensions summary ------------------------------------------------------

$installedCount = 0
$registryFile = Join-Path $SpecifyDir 'extensions/.registry'
if (Test-Path $registryFile) {
    try {
        $regData = Get-Content $registryFile -Raw | ConvertFrom-Json
        $installedCount = @($regData.extensions.PSObject.Properties).Count
    } catch { $installedCount = 0 }
}

# --- Feature count -----------------------------------------------------------

$featureCount = 0
if (Test-Path $SpecsDir) {
    $featureCount = @(Get-ChildItem $SpecsDir -Directory -ErrorAction SilentlyContinue).Count
}

# --- Output JSON -------------------------------------------------------------

$agentsJson = '['
$first = $true
foreach ($a in $DetectedAgents) {
    if (-not $first) { $agentsJson += ',' }
    $hc = if ($a.has_commands) { 'true' } else { 'false' }
    $agentsJson += "{`"key`":`"$(ConvertTo-JsonSafe $a.key)`",`"name`":`"$(ConvertTo-JsonSafe $a.name)`",`"folder`":`"$(ConvertTo-JsonSafe $a.folder)`",`"has_commands`":$hc}"
    $first = $false
}
$agentsJson += ']'

$bools = @{
    has_git = if ($HasGit) { 'true' } else { 'false' }
    spec = if ($hasSpec) { 'true' } else { 'false' }
    plan = if ($hasPlan) { 'true' } else { 'false' }
    tasks = if ($hasTasks) { 'true' } else { 'false' }
    research = if ($hasResearch) { 'true' } else { 'false' }
    data_model = if ($hasDataModel) { 'true' } else { 'false' }
    quickstart = if ($hasQuickstart) { 'true' } else { 'false' }
    contracts = if ($hasContracts) { 'true' } else { 'false' }
    checklists = if ($hasChecklists) { 'true' } else { 'false' }
    checklists_all_pass = if ($checklistsAllPass) { 'true' } else { 'false' }
}

Write-Output "{`"project_name`":`"$(ConvertTo-JsonSafe $ProjectName)`",`"has_git`":$($bools.has_git),`"current_branch`":`"$(ConvertTo-JsonSafe $CurrentBranch)`",`"feature_name`":`"$(ConvertTo-JsonSafe $FeatureName)`",`"script_type`":`"$(ConvertTo-JsonSafe $ScriptType)`",`"agents`":$agentsJson,`"artifacts`":{`"spec`":$($bools.spec),`"plan`":$($bools.plan),`"tasks`":$($bools.tasks),`"research`":$($bools.research),`"data_model`":$($bools.data_model),`"quickstart`":$($bools.quickstart),`"contracts`":$($bools.contracts),`"contract_count`":$contractCount,`"checklists`":$($bools.checklists),`"checklist_count`":$checklistCount,`"checklists_all_pass`":$($bools.checklists_all_pass)},`"tasks`":{`"completed`":$tasksCompleted,`"total`":$tasksTotal},`"phase`":`"$(ConvertTo-JsonSafe $Phase)`",`"phase_hint`":`"$(ConvertTo-JsonSafe $PhaseHint)`",`"installed_extensions`":$installedCount,`"feature_count`":$featureCount}"
