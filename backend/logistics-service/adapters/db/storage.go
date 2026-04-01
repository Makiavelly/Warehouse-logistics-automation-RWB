package db

import (
	"database/sql"
	"errors"
	coreErrors "logistics-service/logistics-service/core/errors"
	"logistics-service/logistics-service/core/ports"

	"github.com/jackc/pgconn"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	log  ports.Logger
	conn *sqlx.DB
}

func New(log ports.Logger, address string) (*DB, error) {
	const op = "adapters.db.New"

	log.Info("connecting to postgres", "address", address)
	db, err := sqlx.Connect("pgx", address)
	if err != nil {
		log.Error(op, "address", address, "error", err)
		return nil, coreErrors.ErrConnectDB
	}

	return &DB{log: log, conn: db}, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

func isNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}