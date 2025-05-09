@echo off

cd ..\..

REM Определяем переменные
set "project_dir=%cd%"          :: Текущая директория проекта
set "output_dir=%project_dir%\cmd\accrual" :: Директория для сборки и запуска
set "executable=accrual.exe" :: Имя исполняемого файла

REM Переходим в директорию cmd\accrual
cd /d "%output_dir%"

REM Команда для сборки программы
echo Building the executable...
go build -o %executable% "%project_dir%\cmd\accrual\main.go"
if %errorlevel% neq 0 (
    echo Failed to build the executable.
    cd /d "%project_dir%"
    pause
    exit /b %errorlevel%
)

REM Возвращаемся в project_dir для запуска
cd /d "%project_dir%"

REM Запуск программы с параметрами
echo Running the executable with parameters...
%output_dir%\%executable% -a=localhost:8081 -d=postgres://user:password@localhost:5432/postgres?sslmode=disable

REM Возвращаемся в исходную директорию
cd /d "%project_dir%"

pause