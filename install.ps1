# Check for Administrator privileges
if (-not ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    Write-Host "This script needs to be run as Administrator."
    Write-Host "Attempting to re-launch as Administrator..."
    Start-Process pwsh -ArgumentList "-NoProfile -File `"$PSCommandPath`"" -Verb RunAs
    exit
}

Write-Host "Running as Administrator..."

# Discover the absolute path of the script's directory
$scriptPath = $PSScriptRoot
$executablePath = Join-Path -Path $scriptPath -ChildPath "fix-SQ-scripts.exe"

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