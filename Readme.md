# Fix SQ Scripts

This is a simple tool to patch `.mq5` files, replacing a specific line of code with a corrected version. It can be run from the command line or via a graphical user interface (GUI) by installing context menu shortcuts.

## The Problem

The original `.mq5` scripts contain an incorrect calculation for `PointValue`. This tool replaces the faulty line with a corrected version that uses `SymbolInfoDouble(correctedSymbol, SYMBOL_TRADE_CONTRACT_SIZE)`.

## Features

- **Command-Line Interface (CLI):** Process files directly from the terminal.
- **Graphical User Interface (GUI):** A simple window that shows the status of each file being processed.
- **Context Menu Integration:** Right-click on `.mq5` files, folders, or folder backgrounds in Windows Explorer to patch files.

## Installation

1.  Make sure you have PowerShell and Go installed on your system.
2.  Run the `build.ps1` script to compile the application. This will create a `fix-SQ-scripts.exe` executable in the project directory.
3.  Right-click on the `install.ps1` script and select "Run with PowerShell". This will install the context menu shortcuts.

## How to Use

### GUI Mode (via Context Menu)

-   **For a single `.mq5` file:** Right-click the file and select "Fix MQ5 Scripts".
-   **For a folder:** Right-click the folder and select "Fix MQ5 Scripts" to process all `.mq5` files within that folder and its subdirectories.
-   **For the current folder:** Right-click on the background of a folder in Explorer and select "Fix MQ5 Scripts" to process all `.mq5` files in the current directory and its subdirectories.

### Command-Line Mode

Open a terminal and run the executable with the path to the file or a glob pattern for multiple files.

```sh
./fix-SQ-scripts.exe "path/to/your/file.mq5"
./fix-SQ-scripts.exe "*.mq5"
```

## Uninstallation

Run the `uninstall.ps1` script to remove the context menu shortcuts. This script is generated automatically when you run the installer.

## Technical Details

The application is written in Go and uses the Fyne library for the GUI. The `install.ps1` script creates registry entries to add the context menu shortcuts.

### Context Menu Fix

The original `install.ps1` script used `"%1"` to pass the path to the application for all context menu entries. This works for files and folders, but not for the folder background context menu. The corrected script now uses `"%V"` for the folder background, which correctly passes the current directory path to the application.