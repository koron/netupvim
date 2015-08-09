@ECHO OFF
SETLOCAL
SET BASEDIR=%~dp0

netupvim.exe %BASEDIR%

:END
ECHO 約10秒後にこのウィンドウは自動的に閉じます。
PING localhost -n 10 > nul
