package db

import (
	"context"
	coreErrors "logistics-service/logistics-service/core/errors"
	"logistics-service/logistics-service/core/models"
)

func (db *DB) CreateWarehouse(ctx context.Context, req models.RequestCreateWarehouse) (models.Warehouse, error) {
	const op = "db.CreateWarehouse"

	var w models.Warehouse
	err := db.conn.QueryRowxContext(ctx,
		`INSERT INTO warehouses (name, office_from_id, address)
		 VALUES ($1, $2, $3)
		 RETURNING id, name, office_from_id, address, created_at`,
		req.Name, req.OfficeFromID, req.Address,
	).StructScan(&w)
	if err != nil {
		if isUniqueViolation(err) {
			return models.Warehouse{}, coreErrors.ErrDuplicateWarehouse
		}
		db.log.Error(op, "error", err)
		return models.Warehouse{}, coreErrors.ErrExecQuery
	}
	return w, nil
}

func (db *DB) GetWarehouses(ctx context.Context) ([]models.Warehouse, error) {
	const op = "db.GetWarehouses"

	var list []models.Warehouse
	if err := db.conn.SelectContext(ctx, &list,
		`SELECT id, name, office_from_id, address, created_at FROM warehouses ORDER BY created_at`); err != nil {
		db.log.Error(op, "error", err)
		return nil, coreErrors.ErrExecQuery
	}
	return list, nil
}

func (db *DB) GetWarehouseByID(ctx context.Context, id string) (models.Warehouse, error) {
	const op = "db.GetWarehouseByID"

	var w models.Warehouse
	err := db.conn.QueryRowxContext(ctx,
		`SELECT id, name, office_from_id, address, created_at FROM warehouses WHERE id=$1`, id,
	).StructScan(&w)
	if err != nil {
		if isNotFound(err) {
			return models.Warehouse{}, coreErrors.ErrNotFoundWarehouse
		}
		db.log.Error(op, "error", err)
		return models.Warehouse{}, coreErrors.ErrExecQuery
	}
	return w, nil
}

func (db *DB) DeleteWarehouse(ctx context.Context, id string) error {
	const op = "db.DeleteWarehouse"

	res, err := db.conn.ExecContext(ctx, `DELETE FROM warehouses WHERE id=$1`, id)
	if err != nil {
		db.log.Error(op, "error", err)
		return coreErrors.ErrExecQuery
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return coreErrors.ErrNotFoundWarehouse
	}
	return nil
}