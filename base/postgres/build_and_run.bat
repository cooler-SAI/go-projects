@echo off
SETLOCAL

REM Define the target directory for the build
SET "PROJECT_DIR=D:\Enterprise Development\Go-projects\go-projects\base\postgres"
SET "BUILD_NAME=transactions.exe"
SET "BUILD_PATH=%PROJECT_DIR%\%BUILD_NAME%"

echo.
echo --- Build Process ---
echo Current working directory for build: %PROJECT_DIR%
echo Target build path: %BUILD_PATH%

REM Clean up previous build first (optional, but good practice)
if exist "%BUILD_PATH%" (
    echo Deleting previous build: "%BUILD_PATH%"
    del /q "%BUILD_PATH%" >nul 2>&1
)

REM Build the Go project
pushd "%PROJECT_DIR%"
go build -o "%BUILD_NAME%" transactions.go
SET "BUILD_EXIT_CODE=%ERRORLEVEL%"
popd

if %BUILD_EXIT_CODE% neq 0 (
    echo Error: Go build failed with exit code %BUILD_EXIT_CODE%.
    exit /b %BUILD_EXIT_CODE%
)
echo Go build successful.

echo.
echo --- Launching Program ---
REM Small delay to ensure the file is written to disk before launching
ping -n 2 127.0.0.1 >nul

if not exist "%BUILD_PATH%" (
    echo Error: Compiled executable not found at "%BUILD_PATH%" after build.
    exit /b 1
)

echo Running: "%BUILD_PATH%"
start /wait "" "%BUILD_PATH%"
SET "PROGRAM_EXIT_CODE=%ERRORLEVEL%"

echo Program finished with exit code: %PROGRAM_EXIT_CODE%

echo.
echo --- Cleanup Process ---
if exist "%BUILD_PATH%" (
    echo Attempting to delete build: "%BUILD_PATH%"
    REM Retry deletion multiple times with a small delay
    for /L %%i in (1,1,5) do (
        del /q "%BUILD_PATH%"
        if %ERRORLEVEL% equ 0 (
            echo Build successfully deleted.
            goto :CleanupEnd
        ) else (
            echo Attempt %%i failed. Retrying in 0.5 seconds...
            ping -n 1 127.0.0.1 >nul
            ping -n 1 127.0.0.1 >nul
        )
    )
    echo Error: Could not delete build after multiple attempts.
) else (
    echo Build file not found or already deleted.
)

:CleanupEnd
ENDLOCAL
exit /b %PROGRAM_EXIT_CODE%