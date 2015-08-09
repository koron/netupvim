@ECHO OFF
SETLOCAL
SET BASEDIR=%~dp0
SET TARGET_DIR=%BASEDIR%

REM データファイル削除
IF EXIST "%TARGET_DIR%"netupvim\var\vim74-anchor.txt DEL /F /Q "%TARGET_DIR%"netupvim\var\vim74-anchor.txt
IF EXIST "%TARGET_DIR%"netupvim\var\vim74-receipe.txt DEL /F /Q "%TARGET_DIR%"netupvim\var\vim74-receipe.txt

netupvim.exe %TARGET_DIR%

:END
ECHO 約10秒後にこのウィンドウは自動的に閉じます。
PING localhost -n 10 > nul
