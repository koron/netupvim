@ECHO OFF
SETLOCAL
SET BASEDIR=%~dp0
SET TARGET_DIR=%BASEDIR%
GOTO :MAIN

:DELETE
IF EXIST "%1" DEL /F /Q "%1"
EXIT /B 0

:MAIN

CALL :DELETE %TARGET_DIR%netupvim\var\vim74-win32-anchor.txt
CALL :DELETE %TARGET_DIR%netupvim\var\vim74-win32-recipe.txt
CALL :DELETE %TARGET_DIR%netupvim\var\vim74-win64-anchor.txt
CALL :DELETE %TARGET_DIR%netupvim\var\vim74-win64-recipe.txt

netupvim.exe %TARGET_DIR%
GOTO :END

:END
ECHO 約10秒後にこのウィンドウは自動的に閉じます。
PING localhost -n 10 > nul
