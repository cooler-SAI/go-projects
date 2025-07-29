@echo off
SETLOCAL

REM Define the target directory for the build
REM This should be the directory containing your transactions.go file
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
REM Use the -o flag to specify the output path and name
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
REM --- Removed 'timeout /t 1 /nobreak >nul' ---
REM This line was causing "Input redirection is not supported" error.
REM If you encounter "file not found" issues again after removing this,
REM you might need to find an alternative way to ensure the file is ready,
REM or re-evaluate the need for a delay.

if not exist "%BUILD_PATH%" (
    echo Error: Compiled executable not found at "%BUILD_PATH%" after build.
    exit /b 1
)

echo Launching: "%BUILD_PATH%"
start /wait "" "%BUILD_PATH%"
SET "PROGRAM_EXIT_CODE=%ERRORLEVEL%"

echo Program finished with exit code: %PROGRAM_EXIT_CODE%

echo.
echo --- Cleanup Process ---
if exist "%BUILD_PATH%" (
    echo Deleting build: "%BUILD_PATH%"
    del /q "%BUILD_PATH%"
    if %ERRORLEVEL% equ 0 (
        echo Build successfully deleted.
    ) else (
        echo Error deleting build.
    )
) else (
    echo Build file not found or already deleted.
)

ENDLOCAL
exit /b %PROGRAM_EXIT_CODE%