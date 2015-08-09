@ECHO OFF
SETLOCAL
SET BASEDIR=%~dp0
SET TARGET_DIR=%BASEDIR%
GOTO :MAIN

:MAIN

netupvim.exe %TARGET_DIR%
GOTO :END

:END
ECHO 約10秒後にこのウィンドウは自動的に閉じます。
PING localhost -n 10 > nul
