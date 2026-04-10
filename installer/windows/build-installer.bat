@echo off
REM =============================================================================
REM  Flang Programming Language - Script de Construcao do Instalador
REM  Versao 0.5.0
REM
REM  O que este script faz:
REM    1. Compila o flang.exe para Windows amd64
REM    2. (Opcional) Cria um ZIP auto-extraivel com tudo
REM    3. (Opcional) Gera o instalador .exe via Inno Setup (setup.iss)
REM
REM  Uso:
REM    build-installer.bat            -> compila e instala direto
REM    build-installer.bat zip        -> compila e gera ZIP
REM    build-installer.bat inno       -> compila e gera .exe via Inno Setup
REM    build-installer.bat tudo       -> faz tudo acima
REM =============================================================================

setlocal EnableDelayedExpansion

REM --- Cores para o terminal (requer Windows 10+) ---
set "COR_OK=[92m"
set "COR_ERRO=[91m"
set "COR_INFO=[96m"
set "COR_RESET=[0m"

REM --- Variaveis de configuracao ---
set "VERSAO=0.5.0"
set "RAIZ=%~dp0..\.."
set "EXE_SAIDA=%RAIZ%\flang.exe"
set "PASTA_INSTALLER=%~dp0"
set "ZIP_SAIDA=%PASTA_INSTALLER%flang-windows-amd64-v%VERSAO%.zip"
set "INNO_EXE=C:\Program Files (x86)\Inno Setup 6\ISCC.exe"
set "MODO=%~1"

echo.
echo  ███████╗██╗      █████╗ ███╗   ██╗ ██████╗
echo  ██╔════╝██║     ██╔══██╗████╗  ██║██╔════╝
echo  █████╗  ██║     ███████║██╔██╗ ██║██║  ███╗
echo  ██╔══╝  ██║     ██╔══██║██║╚██╗██║██║   ██║
echo  ██║     ███████╗██║  ██║██║ ╚████║╚██████╔╝
echo  ╚═╝     ╚══════╝╚═╝  ╚═╝╚═╝  ╚═══╝ ╚═════╝
echo.
echo  Build do Instalador - Flang v%VERSAO%
echo  =========================================
echo.

REM =============================================================================
REM  PASSO 1: Verificar dependencias
REM =============================================================================
echo %COR_INFO%  [1/4] Verificando dependencias...%COR_RESET%

REM Verifica Go
where go >nul 2>&1
if errorlevel 1 (
    echo %COR_ERRO%  [ERRO] Go nao encontrado. Instale em: https://go.dev/dl/%COR_RESET%
    exit /b 1
)
for /f "tokens=3" %%v in ('go version') do set "GO_VER=%%v"
echo         Go encontrado: !GO_VER!

REM Verifica PowerShell
where powershell >nul 2>&1
if errorlevel 1 (
    echo %COR_ERRO%  [ERRO] PowerShell nao encontrado.%COR_RESET%
    exit /b 1
)
echo         PowerShell encontrado.

REM =============================================================================
REM  PASSO 2: Compilar o flang.exe para Windows amd64
REM =============================================================================
echo.
echo %COR_INFO%  [2/4] Compilando flang.exe para Windows amd64...%COR_RESET%

cd /d "%RAIZ%"

REM Define variaveis de build para um binario limpo e menor
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=0

REM -ldflags "-s -w" remove simbolos de debug (binario menor)
REM -ldflags "-X main.Version=..." injeta a versao no binario
go build -ldflags "-s -w -X main.Version=%VERSAO%" -o "%EXE_SAIDA%" .

if errorlevel 1 (
    echo %COR_ERRO%  [ERRO] Falha na compilacao. Veja os erros acima.%COR_RESET%
    exit /b 1
)

REM Mostra tamanho do binario gerado
for %%F in ("%EXE_SAIDA%") do set "TAMANHO=%%~zF"
set /a TAMANHO_MB=!TAMANHO! / 1048576
echo %COR_OK%  [OK]  flang.exe compilado com sucesso (!TAMANHO_MB! MB)%COR_RESET%

REM =============================================================================
REM  PASSO 3: Criar ZIP auto-extraivel (se solicitado)
REM =============================================================================
if /i "%MODO%"=="zip" goto :fazer_zip
if /i "%MODO%"=="tudo" goto :fazer_zip
goto :pular_zip

:fazer_zip
echo.
echo %COR_INFO%  [3/4] Criando pacote ZIP para distribuicao...%COR_RESET%

REM Cria pasta temporaria para o ZIP
set "ZIP_TMP=%TEMP%\flang-installer-tmp"
if exist "%ZIP_TMP%" rmdir /s /q "%ZIP_TMP%"
mkdir "%ZIP_TMP%\Flang"
mkdir "%ZIP_TMP%\Flang\bin"
mkdir "%ZIP_TMP%\Flang\exemplos"
mkdir "%ZIP_TMP%\Flang\installer"

REM Copia arquivos
copy "%EXE_SAIDA%" "%ZIP_TMP%\Flang\bin\flang.exe" >nul
if exist "%RAIZ%\exemplos" xcopy "%RAIZ%\exemplos" "%ZIP_TMP%\Flang\exemplos" /E /I /Q >nul
copy "%PASTA_INSTALLER%install.ps1" "%ZIP_TMP%\Flang\installer\" >nul
copy "%PASTA_INSTALLER%uninstall.ps1" "%ZIP_TMP%\Flang\installer\" >nul
copy "%RAIZ%\LICENSE" "%ZIP_TMP%\Flang\" >nul 2>&1

REM Cria README de instalacao
(
echo Flang Programming Language v%VERSAO%
echo ======================================
echo.
echo INSTALACAO RAPIDA:
echo   1. Extraia esta pasta para onde quiser
echo   2. Abra o PowerShell como usuario normal
echo   3. Execute: powershell -ExecutionPolicy Bypass -File installer\install.ps1
echo.
echo REQUISITOS:
echo   - Windows 10 ou superior
echo   - PowerShell 5.0+
echo.
echo Nao e necessario ser administrador.
) > "%ZIP_TMP%\Flang\INSTALAR.txt"

REM Compacta usando PowerShell (disponivel no Windows 10+)
powershell -NoProfile -Command ^
    "Compress-Archive -Path '%ZIP_TMP%\Flang' -DestinationPath '%ZIP_SAIDA%' -Force"

if errorlevel 1 (
    echo %COR_ERRO%  [ERRO] Falha ao criar ZIP.%COR_RESET%
) else (
    echo %COR_OK%  [OK]  ZIP criado: %ZIP_SAIDA%%COR_RESET%
)

REM Limpa temporarios
rmdir /s /q "%ZIP_TMP%" >nul 2>&1

:pular_zip

REM =============================================================================
REM  PASSO 4: Gerar instalador .exe via Inno Setup (se solicitado)
REM =============================================================================
if /i "%MODO%"=="inno" goto :fazer_inno
if /i "%MODO%"=="tudo" goto :fazer_inno
goto :pular_inno

:fazer_inno
echo.
echo %COR_INFO%  [4/4] Gerando instalador .exe com Inno Setup...%COR_RESET%

if not exist "%INNO_EXE%" (
    echo %COR_ERRO%  [ERRO] Inno Setup nao encontrado em:%COR_RESET%
    echo          %INNO_EXE%
    echo.
    echo          Baixe em: https://jrsoftware.org/isdl.php
    goto :pular_inno
)

"%INNO_EXE%" "%PASTA_INSTALLER%setup.iss"

if errorlevel 1 (
    echo %COR_ERRO%  [ERRO] Falha ao compilar setup.iss%COR_RESET%
) else (
    echo %COR_OK%  [OK]  Instalador .exe gerado com sucesso!%COR_RESET%
    echo          Procure o arquivo em: %PASTA_INSTALLER%Output\
)

:pular_inno

REM =============================================================================
REM  PASSO FINAL: Instalar direto (modo padrao sem argumentos)
REM =============================================================================
if "%MODO%"=="" (
    echo.
    echo %COR_INFO%  [4/4] Executando instalador PowerShell...%COR_RESET%
    echo.
    powershell -ExecutionPolicy Bypass -File "%PASTA_INSTALLER%install.ps1"
)

REM =============================================================================
REM  Resumo
REM =============================================================================
echo.
echo  =========================================
echo %COR_OK%   Build concluido!%COR_RESET%
echo  =========================================
echo.
echo   flang.exe : %EXE_SAIDA%
if /i "%MODO%"=="zip" echo   ZIP       : %ZIP_SAIDA%
if /i "%MODO%"=="tudo" echo   ZIP       : %ZIP_SAIDA%
echo.
echo   Para instalar manualmente:
echo     powershell -ExecutionPolicy Bypass -File installer\windows\install.ps1
echo.

endlocal
