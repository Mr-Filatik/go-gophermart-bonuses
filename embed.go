package migrations_root

import "embed"

var (
	//go:embed migrations/*.sql
	EmbedMigrations embed.FS

	DirMigrations string = "migrations"
)
