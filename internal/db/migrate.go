package db

import (
	"embed"
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func RunMigrations(pgPool *pgxpool.Pool) error {
	// pgxpool.Pool implements only the pgx interface, not database/sql.
	// OpenDBFromPool wraps it in a *sql.DB adapter so golang-migrate can use
	// it. Per pgx docs, closing the wrapping *sql.DB does NOT close the
	// underlying pgxpool — the pool stays alive for the rest of the app.
	// Closing it here releases the connection golang-migrate held for the
	// migration run, which is required so that `pool.Close()` later does not
	// block waiting for it (matters in tests with short-lived pools).
	db := stdlib.OpenDBFromPool(pgPool)
	defer func() { _ = db.Close() }()

	src, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to load migration files: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to initialise migrator: %w", err)
	}
	// Close releases the pinned connection + the source. Errors from Close
	// don't affect migration correctness; surface them only in logs.
	defer func() {
		if srcErr, dbErr := m.Close(); srcErr != nil || dbErr != nil {
			slog.Error("migrator close returned errors",
				"source_error", srcErr,
				"db_error", dbErr,
			)
		}
	}()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	slog.Info("database migrations applied successfully")
	return nil
}
