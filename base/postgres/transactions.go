package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func scheduleSelfDelete() {
	if runtime.GOOS != "windows" {
		return
	}

	exePath, err := os.Executable()
	if err != nil {
		return
	}

	batContent := `@echo off
:loop
tasklist | find "` + filepath.Base(exePath) + `" >nul
if not errorlevel 1 (
	timeout /t 1 /nobreak >nul
	goto loop
)
del "` + exePath + `"
del "%~f0"`

	batPath := exePath + "_delete.bat"
	_ = os.WriteFile(batPath, []byte(batContent), 0666)

	cmd := exec.Command("cmd.exe", "/C", batPath)
	_ = cmd.Start()
}

func main() {
	scheduleSelfDelete()

	println("Program Working...2")
}
