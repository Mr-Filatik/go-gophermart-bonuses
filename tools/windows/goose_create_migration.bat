@echo off

cd ..\..

REM Проверка наличия goose
where goose >nul 2>&1
if %errorlevel% neq 0 (
    echo Goose is not installed or not in PATH. Please install it first.
    pause
    exit /b 1
)

REM Запрос имени миграции у пользователя
set /p migrations_name="Enter the name of the migration: "

REM Проверка, что имя миграции не пустое
if "%migrations_name%"=="" (
    echo Migration name cannot be empty. Please try again.
    pause
    exit /b 1
)

REM Название папки для миграций db/migrations
set "migration_dir=migrations"

REM Проверка существования папки дирректории
if not exist "%migration_dir%" (
    echo Folder "%migration_dir%" does not exist. Creating it...
    mkdir "%migration_dir%"
    if %errorlevel% neq 0 (
        echo Failed to create folder "%migration_dir%".
        pause
        exit /b 1
    )
    echo Folder "%migration_dir%" created successfully.
)

REM Выполнение команды goose
echo Creating migration: %migrations_name%
goose -dir %migration_dir% create %migrations_name% sql

REM Проверка успешности выполнения команды
if %errorlevel% neq 0 (
    echo Failed to create migration.
    pause
    exit /b %errorlevel%
)

echo Migration created successfully!
pause