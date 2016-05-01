@ECHO OFF
SETLOCAL
SET TARGET_DIR=%~dp0

netupvim.exe -t %TARGET_DIR%

ECHO This window will be closed after 10 seconds.
PING localhost -n 10 > nul
