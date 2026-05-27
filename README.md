# Notes App - Setup & Build Guide

## Prerequisites

You need **Go 1.21+** and **Wails** installed on your system.

### 1. Install Go
- Download from: https://golang.org/dl
- Add Go to your PATH (should be automatic)
- Verify: `go version`

### 2. Install Wails
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

On Windows, you may need to install additional dependencies:
- Download and install **WebView2 Runtime**: https://developer.microsoft.com/en-us/microsoft-edge/webview2/
- You may need Visual Studio Build Tools

### 3. Verify Wails Installation
```bash
wails version
```

## Building the App

### Development Mode (with hot reload)
```bash
wails dev
```
This opens the app in a window and auto-reloads as you make changes.

### Build to Executable
```bash
wails build
```
The standalone executable will be created in the `build/bin/` directory.

On Windows: `NotesApp.exe`
On Mac: `NotesApp.app`
On Linux: `NotesApp`

## Project Structure

```
.
├── main.go              # Application entry point
├── app.go               # Backend logic for note management
├── go.mod               # Go dependencies
├── wails.json           # Wails configuration
└── frontend/
    └── dist/
        ├── index.html   # HTML structure
        ├── style.css    # Styling
        └── app.js       # Frontend JavaScript logic
```

## Database

Notes are stored in an SQLite database at:
- Windows: `%USERPROFILE%\.notesapp\notes.db`
- Mac/Linux: `~/.notesapp/notes.db`

Each note has:
- **ID**: Unique identifier
- **Text**: The note content
- **Created At**: Timestamp of creation
- **Updated At**: Timestamp of last modification

## Features

✅ Add new notes with the "+ Add Note" button or press Enter
✅ Copy notes to clipboard with the "Copy" button
✅ Remove notes with the "Remove" button
✅ All changes automatically saved to JSON database
✅ Clean, modern UI with keyboard support
