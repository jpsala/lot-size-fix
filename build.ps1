go clean
go build -ldflags "-H=windowsgui" -o fix-SQ-scripts.exe main.go

$destinationDirs = @(
    "G:\My Drive\SQX\sync\lot-size-fix",
    "G:\My Drive\shared\fix-lot-size"
)

# List of files to copy. We copy settings.json from the project root now.
$filesToCopy = @("fix-SQ-scripts.exe", "install.ps1", "settings.json")

foreach ($dir in $destinationDirs) {
    if (-not (Test-Path $dir)) {
        New-Item -ItemType Directory -Path $dir -Force
    }
    foreach ($file in $filesToCopy) {
        # Ensure the source file exists before copying
        if (Test-Path $file) {
            Copy-Item -Path $file -Destination $dir -Force
        }
    }
    Write-Host "Files copied to: $dir"
}