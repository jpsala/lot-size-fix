# --- Pre-flight Checks ---

# 1. Check for Administrator privileges
if (-not ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    Write-Host "----------------------------------------------------------------"
    Write-Host "ERROR: Administrator privileges are required." -ForegroundColor Red
    Write-Host "Please right-click this script and select 'Run as Administrator'."
    Write-Host "Attempting to re-launch automatically..."
    Write-Host "----------------------------------------------------------------"
    # Attempt to re-launch with elevated privileges
    Start-Process pwsh -ArgumentList "-NoProfile -File `"$PSCommandPath`"" -Verb RunAs
    exit
}
Write-Host "SUCCESS: Running with Administrator privileges." -ForegroundColor Green

# 2. Check PowerShell Execution Policy
$policy = Get-ExecutionPolicy
Write-Host "Current PowerShell Execution Policy: $policy"
if ($policy -ne 'Unrestricted' -and $policy -ne 'RemoteSigned' -and $policy -ne 'Bypass') {
    Write-Host "----------------------------------------------------------------"
    Write-Host "WARNING: Your PowerShell Execution Policy might prevent this script from running correctly." -ForegroundColor Yellow
    Write-Host "If the script fails, please run the following command in an Administrator PowerShell window:"
    Write-Host "Set-ExecutionPolicy RemoteSigned -Scope CurrentUser"
    Write-Host "----------------------------------------------------------------"
}

# --- Script Main Logic ---

Write-Host "Starting context menu installation..."

# Discover the absolute path of the script's directory
$scriptPath = $PSScriptRoot
$executablePath = Join-Path -Path $scriptPath -ChildPath "fix-SQ-scripts.exe"

# 3. Verify executable exists
if (-not (Test-Path $executablePath)) {
    Write-Host "----------------------------------------------------------------"
    Write-Host "ERROR: 'fix-SQ-scripts.exe' not found!" -ForegroundColor Red
    Write-Host "Please ensure 'fix-SQ-scripts.exe' is in the same directory as this install script."
    Write-Host "Script directory: $scriptPath"
    Write-Host "----------------------------------------------------------------"
    Read-Host "Press Enter to exit"
    exit
}
Write-Host "SUCCESS: Found executable at: $executablePath" -ForegroundColor Green

# Define Registry Paths
$mq5KeyPath = "HKCU:\Software\Classes\SystemFileAssociations\.mq5\shell\FixMQ5Scripts"
$dirKeyPath = "HKCU:\Software\Classes\Directory\shell\FixMQ5Scripts"
$dirBgKeyPath = "HKCU:\Software\Classes\Directory\Background\shell\FixMQ5Scripts"
$commandKeyPath = "command"

# Command to be executed by the context menu
$commandForFilesAndDirs = "`"$executablePath`" --gui `"%1`""
$commandForDirBg = "`"$executablePath`" --gui `"%V`""

# Create context menu for .mq5 files
Write-Host "Creating context menu for .mq5 files..."
if (-not (Test-Path $mq5KeyPath)) {
    New-Item -Path $mq5KeyPath -Force | Out-Null
}
New-ItemProperty -Path $mq5KeyPath -Name "(default)" -Value "Fix MQ5 Scripts" -Force | Out-Null
$mq5CommandPath = Join-Path -Path $mq5KeyPath -ChildPath $commandKeyPath
if (-not (Test-Path $mq5CommandPath)) {
    New-Item -Path $mq5CommandPath -Force | Out-Null
}
New-ItemProperty -Path $mq5CommandPath -Name "(default)" -Value $commandForFilesAndDirs -Force | Out-Null

# Create context menu for folders
Write-Host "Creating context menu for folders..."
if (-not (Test-Path $dirKeyPath)) {
    New-Item -Path $dirKeyPath -Force | Out-Null
}
New-ItemProperty -Path $dirKeyPath -Name "(default)" -Value "Fix MQ5 Scripts" -Force | Out-Null
$dirCommandPath = Join-Path -Path $dirKeyPath -ChildPath $commandKeyPath
if (-not (Test-Path $dirCommandPath)) {
    New-Item -Path $dirCommandPath -Force | Out-Null
}
New-ItemProperty -Path $dirCommandPath -Name "(default)" -Value $commandForFilesAndDirs -Force | Out-Null

# Create context menu for folder backgrounds
Write-Host "Creating context menu for folder backgrounds..."
if (-not (Test-Path $dirBgKeyPath)) {
    New-Item -Path $dirBgKeyPath -Force | Out-Null
}
New-ItemProperty -Path $dirBgKeyPath -Name "(default)" -Value "Fix MQ5 Scripts" -Force | Out-Null
$dirBgCommandPath = Join-Path -Path $dirBgKeyPath -ChildPath $commandKeyPath
if (-not (Test-Path $dirBgCommandPath)) {
    New-Item -Path $dirBgCommandPath -Force | Out-Null
}
New-ItemProperty -Path $dirBgCommandPath -Name "(default)" -Value $commandForDirBg -Force | Out-Null

# Generate uninstall script
Write-Host "Generating uninstall.ps1..."
$uninstallScriptContent = @"
# Uninstall script for FixMQ5Scripts context menu items

Write-Host "Removing context menu for .mq5 files..."
Remove-Item -Path "$mq5KeyPath" -Recurse -Force -ErrorAction SilentlyContinue

Write-Host "Removing context menu for folders..."
Remove-Item -Path "$dirKeyPath" -Recurse -Force -ErrorAction SilentlyContinue

Write-Host "Removing context menu for folder backgrounds..."
Remove-Item -Path "$dirBgKeyPath" -Recurse -Force -ErrorAction SilentlyContinue

Write-Host "Uninstallation complete."

# Self-delete the uninstaller
Remove-Item -Path "`$PSCommandPath" -Force
"@

$uninstallScriptPath = Join-Path -Path $scriptPath -ChildPath "uninstall.ps1"
Set-Content -Path $uninstallScriptPath -Value $uninstallScriptContent

Write-Host "Installation complete. Context menu items have been added."
Write-Host "An uninstall.ps1 script has been created in the same directory."

Read-Host "Press Enter to exit"