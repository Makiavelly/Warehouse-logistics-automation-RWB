package db

import (
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	coreErrors "logistics-service/logistics-service/core/errors"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func (db *DB) Migrate() error {
	db.log.Info("running migrations")

	files, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		db.log.Error("failed to load migration files", "err", err)
		return coreErrors.ErrMigrationFailed
	}

	driver, err := pgx.WithInstance(db.conn.DB, &pgx.Config{})
	if err != nil {
		db.log.Error("failed to create pgx driver", "err", err)
		return coreErrors.ErrMigrationFailed
	}

	m, err := migrate.NewWithInstance("iofs", files, "pgx", driver)
	if err != nil {
		db.log.Error("failed to create migrate instance", "err", err)
		return coreErrors.ErrMigrationFailed
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		db.log.Error("migration failed", "err", err)
		return coreErrors.ErrMigrationFailed
	}

	db.log.Info("migrations done")
	return nil
}