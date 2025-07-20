# Uninstall script for FixMQ5Scripts context menu items

Write-Host "Removing context menu for .mq5 files..."
Remove-Item -Path "HKCU:\Software\Classes\SystemFileAssociations\.mq5\shell\FixMQ5Scripts" -Recurse -Force -ErrorAction SilentlyContinue

Write-Host "Removing context menu for folders..."
Remove-Item -Path "HKCU:\Software\Classes\Directory\shell\FixMQ5Scripts" -Recurse -Force -ErrorAction SilentlyContinue

Write-Host "Removing context menu for folder backgrounds..."
Remove-Item -Path "HKCU:\Software\Classes\Directory\Background\shell\FixMQ5Scripts" -Recurse -Force -ErrorAction SilentlyContinue

Write-Host "Uninstallation complete."

# Self-delete the uninstaller
Remove-Item -Path "$PSCommandPath" -Force
