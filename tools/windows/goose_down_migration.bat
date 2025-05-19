@echo off

cd ..\..

REM Проверка наличия goose
where goose >nul 2>&1
if %errorlevel% neq 0 (
    echo Goose is not installed or not in PATH. Please install it first.
    pause
    exit /b 1
)

REM Название папки для миграций db/migrations
set "migration_dir=migrations"

REM Выполнение команды goose
echo Up migrations...
goose -dir %migration_dir% postgres "postgres://user:password@localhost:5432/postgres?sslmode=disable" down

REM Проверка успешности выполнения команды
if %errorlevel% neq 0 (
    echo Failed to up migration.
    pause
    exit /b %errorlevel%
)

echo Migration upped successfully!
pause