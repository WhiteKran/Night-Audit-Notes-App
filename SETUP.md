# Detailed Setup Instructions for Windows

## Step 1: Install Go

1. Download Go from: https://golang.org/dl (choose Windows installer)
2. Run the installer and follow the prompts
3. Open Command Prompt and verify:
   ```
   go version
   ```
   You should see something like `go version go1.21.x windows/amd64`

## Step 2: Install WebView2 Runtime

1. Download from: https://developer.microsoft.com/en-us/microsoft-edge/webview2/
2. Click "Download Evergreen Runtime" 
3. Run the installer
4. Restart your computer after installation

## Step 3: Install Visual Studio Build Tools

1. Go to: https://visualstudio.microsoft.com/downloads/
2. Scroll down and find "Visual Studio Build Tools"
3. Click "Free Download"
4. Run the installer
5. Select "Desktop development with C++" workload
6. Click Install

## Step 4: Install Wails

Open Command Prompt and run:
```
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

Verify installation:
```
wails version
```

## Step 5: Test the Project

Navigate to the project directory and run:
```
wails dev
```

This should open a window with your Notes App running.

## Troubleshooting

### "wails: command not found"
- Make sure Go's bin directory is in your PATH
- Usually at: `C:\Users\YourUsername\go\bin`
- Restart Command Prompt or PowerShell after adding to PATH

### "WebView2 error"
- Download and install WebView2 Runtime from the link above
- Restart your computer

### Build fails with C++ errors
- Install Visual Studio Build Tools (see Step 3)
- Ensure you selected "Desktop development with C++"

### Still having issues?
- Check Wails docs: https://wails.io/docs/gettingstarted/installation
