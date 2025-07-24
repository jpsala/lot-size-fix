go clean
go build -ldflags "-H=windowsgui" -o fix-SQ-scripts.exe main.go

# Create the destination directory if it doesn't exist
$destinationDir = "G:\My Drive\SQX\sync\lot-size-fix"
if (-not (Test-Path $destinationDir)) {
    New-Item -ItemType Directory -Path $destinationDir -Force
}

# Copy the executable and PowerShell scripts to the destination
Copy-Item -Path "fix-SQ-scripts.exe" -Destination $destinationDir -Force
Copy-Item -Path "install.ps1" -Destination $destinationDir -Force
Copy-Item -Path "uninstall.ps1" -Destination $destinationDir -Force

Write-Host "Build completed and files copied to: $destinationDir"