package migrate

import (
	"database/sql"
	"fmt"
	"io/fs"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func RunMigrations(db *sql.DB, migrationsFS fs.FS, dir string) error {
	const op = "db.RunMigrations"
	goose.SetBaseFS(migrationsFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := goose.Up(db, dir); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func MigrateDown(db *sql.DB, migrationsFS fs.FS, dir string) error {
	const op = "db.MigrateDown"
	goose.SetBaseFS(migrationsFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := goose.Down(db, dir); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func MigrateStatus(db *sql.DB, migrationsFS fs.FS, dir string) error {
	const op = "db.MigrateStatus"
	goose.SetBaseFS(migrationsFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := goose.Status(db, dir); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
