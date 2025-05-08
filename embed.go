package migrations

import "embed"

var (
	//go:embed migrations/*.sql
	EmbedMigrations embed.FS

	DirMigrations string = "migrations"
)
