package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"golang.org/x/sys/windows"
)

func main() {
	hideConsoleWindow()

	if len(os.Args) < 1 {
		writeLog("This program requires file paths as input.")
		return
	}

	filePaths := os.Args[1:]

	if len(filePaths) == 0 {
		writeLog("No files are registered.")
		return
	}

	firstFilePath := filePaths[0]
	baseDir := filepath.Dir(firstFilePath)
	newFolderName := time.Now().Format("20060102_150405")
	destinationFolder := filepath.Join(baseDir, newFolderName)

	if err := os.Mkdir(destinationFolder, 0755); err != nil {
		if !os.IsExist(err) {
			writeLog(fmt.Sprintf("Failed to create new folder: %v", err))
			return
		}
	}

	writeLog(fmt.Sprintf("Destination folder created or already exists: %s", destinationFolder))

	for _, path := range filePaths {
		err := moveFileToFolder(path, destinationFolder)
		if err != nil {
			writeLog(fmt.Sprintf("Failed to move file %s: %v", path, err))
			continue
		}
		writeLog(fmt.Sprintf("File successfully moved: %s", filepath.Base(path)))
	}

	writeLog("All files have been moved.")
}

func moveFileToFolder(sourcePath, destinationFolder string) error {
	destinationPath := filepath.Join(destinationFolder, filepath.Base(sourcePath))
	if err := os.Rename(sourcePath, destinationPath); err != nil {
		return err
	}
	return nil
}

func hideConsoleWindow() {
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	getConsoleWindow := kernel32.NewProc("GetConsoleWindow")
	showWindow := kernel32.NewProc("ShowWindow")

	consoleWindow, _, _ := getConsoleWindow.Call()
	if consoleWindow != 0 {
		const SW_HIDE = 0
		ret, _, err := showWindow.Call(consoleWindow, uintptr(SW_HIDE))
		if err != nil && err != syscall.Errno(0) {
			writeLog(fmt.Sprintf("ShowWindow failed: %v", err))
		}
		_ = ret
	}
}

func writeLog(message string) {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Failed to get program location:", err)
		return
	}
	logDir := filepath.Dir(exePath)
	logFilePath := filepath.Join(logDir, "app.log")

	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		return
	}
	defer logFile.Close()

	logMessage := fmt.Sprintf("[%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), message)
	if _, err := logFile.WriteString(logMessage); err != nil {
		fmt.Println("Failed to write to log file:", err)
	}
}
