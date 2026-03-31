package db

import (
	"context"
	coreErrors "logistics-service/logistics-service/core/errors"
	"logistics-service/logistics-service/core/models"
)

func (db *DB) CreateRoute(ctx context.Context, warehouseID string, req models.RequestCreateRoute) (models.Route, error) {
	const op = "db.CreateRoute"

	var r models.Route
	err := db.conn.QueryRowxContext(ctx,
		`INSERT INTO routes (warehouse_id, route_id, name)
		 VALUES ($1, $2, $3)
		 RETURNING id, warehouse_id, route_id, name, created_at`,
		warehouseID, req.RouteID, req.Name,
	).StructScan(&r)
	if err != nil {
		if isUniqueViolation(err) {
			return models.Route{}, coreErrors.ErrDuplicateRoute
		}
		db.log.Error(op, "error", err)
		return models.Route{}, coreErrors.ErrExecQuery
	}
	return r, nil
}

func (db *DB) GetRoutesByWarehouse(ctx context.Context, warehouseID string) ([]models.Route, error) {
	const op = "db.GetRoutesByWarehouse"

	var list []models.Route
	if err := db.conn.SelectContext(ctx, &list,
		`SELECT id, warehouse_id, route_id, name, created_at FROM routes WHERE warehouse_id=$1 ORDER BY created_at`,
		warehouseID); err != nil {
		db.log.Error(op, "error", err)
		return nil, coreErrors.ErrExecQuery
	}
	return list, nil
}

func (db *DB) GetRouteByID(ctx context.Context, id string) (models.Route, error) {
	const op = "db.GetRouteByID"

	var r models.Route
	err := db.conn.QueryRowxContext(ctx,
		`SELECT id, warehouse_id, route_id, name, created_at FROM routes WHERE id=$1`, id,
	).StructScan(&r)
	if err != nil {
		if isNotFound(err) {
			return models.Route{}, coreErrors.ErrNotFoundRoute
		}
		db.log.Error(op, "error", err)
		return models.Route{}, coreErrors.ErrExecQuery
	}
	return r, nil
}

func (db *DB) DeleteRoute(ctx context.Context, id string) error {
	const op = "db.DeleteRoute"

	res, err := db.conn.ExecContext(ctx, `DELETE FROM routes WHERE id=$1`, id)
	if err != nil {
		db.log.Error(op, "error", err)
		return coreErrors.ErrExecQuery
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return coreErrors.ErrNotFoundRoute
	}
	return nil
}