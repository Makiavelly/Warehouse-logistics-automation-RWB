package db

import (
	"context"
	coreErrors "logistics-service/logistics-service/core/errors"
	"logistics-service/logistics-service/core/models"
)

func (db *DB) CreateUser(ctx context.Context, username, passwordHash, role string) (models.User, error) {
	const op = "db.CreateUser"

	var u models.User
	err := db.conn.QueryRowxContext(ctx,
		`INSERT INTO users (username, password_hash, role)
		 VALUES ($1, $2, $3)
		 RETURNING id, username, password_hash, role, created_at`,
		username, passwordHash, role,
	).StructScan(&u)
	if err != nil {
		if isUniqueViolation(err) {
			return models.User{}, coreErrors.ErrDuplicateUser
		}
		db.log.Error(op, "error", err)
		return models.User{}, coreErrors.ErrExecQuery
	}
	return u, nil
}

func (db *DB) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	const op = "db.GetUserByUsername"

	var u models.User
	err := db.conn.QueryRowxContext(ctx,
		`SELECT id, username, password_hash, role, created_at FROM users WHERE username=$1`, username,
	).StructScan(&u)
	if err != nil {
		if isNotFound(err) {
			return models.User{}, coreErrors.ErrNotFoundUser
		}
		db.log.Error(op, "error", err)
		return models.User{}, coreErrors.ErrExecQuery
	}
	return u, nil
}

func (db *DB) GetUserByID(ctx context.Context, id string) (models.User, error) {
	const op = "db.GetUserByID"

	var u models.User
	err := db.conn.QueryRowxContext(ctx,
		`SELECT id, username, password_hash, role, created_at FROM users WHERE id=$1`, id,
	).StructScan(&u)
	if err != nil {
		if isNotFound(err) {
			return models.User{}, coreErrors.ErrNotFoundUser
		}
		db.log.Error(op, "error", err)
		return models.User{}, coreErrors.ErrExecQuery
	}
	return u, nil
}