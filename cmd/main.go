package main

import (
	"database/sql"
	"github.com/deplagene/snippetbox"
	"github.com/deplagene/snippetbox/cmd/api"
	"github.com/deplagene/snippetbox/cmd/migrate"
	"github.com/deplagene/snippetbox/configs"
	"github.com/deplagene/snippetbox/pkg/postgres"
	"flag"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// todo: add slog
	migrateCmd := flag.Bool("migrate", false, "run database migrations")
	flag.Parse()

	db, err := postgres.NewPostgresPool(configs.Envs.DbUrl, configs.Envs.MaxConns)
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}

	defer db.Close()

	if *migrateCmd {
		sqlDB, err := sql.Open("pgx", configs.Envs.DbUrl)
		if err != nil {
			log.Fatalf("Error connecting to database for migration: %s", err)
		}
		defer sqlDB.Close()

		if err := migrate.RunMigrations(sqlDB, snippetbox.SchemaFS, "sql/schema"); err != nil {
			log.Fatalf("Error running migrations: %s", err)
		}
		log.Println("Migrations applied successfully")
		return
	}

	api := api.NewAPI(configs.Envs.Port, db)
	if err := api.Run(); err != nil {
		log.Fatalf("Error running API: %s", err)
	}
}
