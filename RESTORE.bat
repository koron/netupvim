@ECHO OFF
SETLOCAL
SET TARGET_DIR=%~dp0

cd "%TARGET_DIR%"
"%TARGET_DIR%netupvim.exe" -restore

ECHO This window will be closed after 10 seconds.
PING localhost -n 10 > nul
