go clean
go build -ldflags "-H=windowsgui" -o fix-SQ-scripts.exe main.go

$destinationDirs = @(
    "G:\My Drive\SQX\sync\lot-size-fix",
    "G:\My Drive\shared\fix-lot-size"
)

# List of files to copy
$filesToCopy = @("fix-SQ-scripts.exe", "install.ps1")

foreach ($dir in $destinationDirs) {
    if (-not (Test-Path $dir)) {
        New-Item -ItemType Directory -Path $dir -Force
    }
    foreach ($file in $filesToCopy) {
        Copy-Item -Path $file -Destination $dir -Force
    }
    Write-Host "Files copied to: $dir"
}