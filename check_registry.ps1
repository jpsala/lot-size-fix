# This script checks for the registry keys created by install.ps1

Write-Host "--- Checking Registry for .mq5 File Context Menu ---" -ForegroundColor Cyan
reg query HKCU\Software\Classes\SystemFileAssociations\.mq5\shell\FixMQ5Scripts /s
Write-Host ""

Write-Host "--- Checking Registry for Folder Context Menu ---" -ForegroundColor Cyan
reg query HKCU\Software\Classes\Directory\shell\FixMQ5Scripts /s
Write-Host ""

Write-Host "--- Checking Registry for Folder Background Context Menu ---" -ForegroundColor Cyan
reg query HKCU\Software\Classes\Directory\Background\shell\FixMQ5Scripts /s
Write-Host ""

Write-Host "--- Check Complete ---" -ForegroundColor Green
Read-Host "Press Enter to exit"
