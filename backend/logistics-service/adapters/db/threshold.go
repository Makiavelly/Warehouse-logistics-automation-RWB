package db

import (
	"context"
	"strconv"

	coreErrors "logistics-service/logistics-service/core/errors"
	"logistics-service/logistics-service/core/models"
)

func (db *DB) SetThreshold(ctx context.Context, req models.RequestSetThreshold) (models.Threshold, error) {
	const op = "db.SetThreshold"

	var t models.Threshold
	err := db.conn.QueryRowxContext(ctx,
		`INSERT INTO thresholds (warehouse_id, route_id, value)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (warehouse_id, route_id) DO UPDATE
		   SET value=$3, updated_at=NOW()
		 RETURNING id, warehouse_id, route_id, value, updated_at`,
		req.WarehouseID, req.RouteID, req.Value,
	).StructScan(&t)
	if err != nil {
		db.log.Error(op, "error", err)
		return models.Threshold{}, coreErrors.ErrExecQuery
	}
	return t, nil
}

func (db *DB) GetThresholds(ctx context.Context, warehouseID, routeID string) ([]models.Threshold, error) {
	const op = "db.GetThresholds"

	query := `SELECT id, warehouse_id, route_id, value, updated_at FROM thresholds WHERE 1=1`
	args := []any{}

	if warehouseID != "" {
		args = append(args, warehouseID)
		query += " AND warehouse_id=$" + itoa(len(args))
	}
	if routeID != "" {
		args = append(args, routeID)
		query += " AND route_id=$" + itoa(len(args))
	}
	query += " ORDER BY updated_at"

	var list []models.Threshold
	if err := db.conn.SelectContext(ctx, &list, query, args...); err != nil {
		db.log.Error(op, "error", err)
		return nil, coreErrors.ErrExecQuery
	}
	return list, nil
}

func (db *DB) GetThreshold(ctx context.Context, warehouseID, routeID string) (models.Threshold, error) {
	const op = "db.GetThreshold"

	var t models.Threshold
	err := db.conn.QueryRowxContext(ctx,
		`SELECT id, warehouse_id, route_id, value, updated_at FROM thresholds
		 WHERE warehouse_id=$1 AND route_id=$2`,
		warehouseID, routeID,
	).StructScan(&t)
	if err != nil {
		if isNotFound(err) {
			return models.Threshold{}, coreErrors.ErrNotFoundThreshold
		}
		db.log.Error(op, "error", err)
		return models.Threshold{}, coreErrors.ErrExecQuery
	}
	return t, nil
}

func itoa(n int) string { return strconv.Itoa(n) }