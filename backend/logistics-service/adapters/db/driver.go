package db

import (
	"context"
	coreErrors "logistics-service/logistics-service/core/errors"
	"logistics-service/logistics-service/core/models"
)

func (db *DB) GetDriverByUserID(ctx context.Context, userID string) (models.Driver, error) {
	const op = "db.GetDriverByUserID"

	var d models.Driver
	err := db.conn.QueryRowxContext(ctx,
		`SELECT d.id, d.user_id, u.username, d.warehouse_id, d.route_id, d.created_at
		 FROM drivers d JOIN users u ON u.id=d.user_id
		 WHERE d.user_id=$1`, userID,
	).StructScan(&d)
	if err != nil {
		if isNotFound(err) {
			return models.Driver{}, coreErrors.ErrNotFoundDriver
		}
		db.log.Error(op, "error", err)
		return models.Driver{}, coreErrors.ErrExecQuery
	}
	return d, nil
}

func (db *DB) GetDrivers(ctx context.Context) ([]models.Driver, error) {
	const op = "db.GetDrivers"

	var list []models.Driver
	if err := db.conn.SelectContext(ctx, &list,
		`SELECT d.id, d.user_id, u.username, d.warehouse_id, d.route_id, d.created_at
		 FROM drivers d JOIN users u ON u.id=d.user_id
		 ORDER BY d.created_at`); err != nil {
		db.log.Error(op, "error", err)
		return nil, coreErrors.ErrExecQuery
	}
	return list, nil
}

func (db *DB) AssignDriver(ctx context.Context, req models.RequestAssignDriver) (models.Driver, error) {
	const op = "db.AssignDriver"

	var d models.Driver
	err := db.conn.QueryRowxContext(ctx,
		`INSERT INTO drivers (user_id, warehouse_id, route_id)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (user_id) DO UPDATE
		   SET warehouse_id=$2, route_id=$3
		 RETURNING id, user_id, warehouse_id, route_id, created_at`,
		req.DriverID, req.WarehouseID, req.RouteID,
	).StructScan(&d)
	if err != nil {
		db.log.Error(op, "error", err)
		return models.Driver{}, coreErrors.ErrExecQuery
	}

	// fetch username separately (not returned by INSERT)
	u, err := db.GetUserByID(ctx, req.DriverID)
	if err == nil {
		d.Username = u.Username
	}

	return d, nil
}