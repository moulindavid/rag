package database

import (
	"context"
	"embed"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Migrate(ctx context.Context, pool *pgxpool.Pool, migrationsFS embed.FS) error {
	files, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return err
	}

	for _, file := range files {
		content, err := migrationsFS.ReadFile("migrations/" + file.Name())
		if err != nil {
			return err
		}

		_, err = pool.Exec(ctx, string(content))
		if err != nil {
			return err
		}
	}

	return nil
}
