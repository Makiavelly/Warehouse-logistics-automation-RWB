package db

import (
	"context"
	"strconv"
	"time"

	coreErrors "logistics-service/logistics-service/core/errors"
	"logistics-service/logistics-service/core/models"
)

func (db *DB) CreateTruckCall(ctx context.Context, tc models.TruckCall) (models.TruckCall, error) {
	const op = "db.CreateTruckCall"

	var saved models.TruckCall
	err := db.conn.QueryRowxContext(ctx,
		`INSERT INTO truck_calls
		   (warehouse_id, route_id, driver_id, forecast_value, threshold_value, called_at, status)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, warehouse_id, route_id, driver_id, forecast_value, threshold_value,
		           called_at, status, timeliness, timeliness_at, actual_containers`,
		tc.WarehouseID, tc.RouteID, tc.DriverID, tc.ForecastValue, tc.ThresholdValue,
		tc.CalledAt, tc.Status,
	).StructScan(&saved)
	if err != nil {
		db.log.Error(op, "error", err)
		return models.TruckCall{}, coreErrors.ErrExecQuery
	}
	return saved, nil
}

func (db *DB) GetTruckCalls(ctx context.Context, warehouseID, routeID string, from, to time.Time) ([]models.TruckCall, error) {
	const op = "db.GetTruckCalls"

	query := `SELECT id, warehouse_id, route_id, driver_id, forecast_value, threshold_value,
	                 called_at, status, timeliness, timeliness_at, actual_containers
	          FROM truck_calls WHERE 1=1`
	args := []any{}

	if warehouseID != "" {
		args = append(args, warehouseID)
		query += " AND warehouse_id=$" + strconv.Itoa(len(args))
	}
	if routeID != "" {
		args = append(args, routeID)
		query += " AND route_id=$" + strconv.Itoa(len(args))
	}
	if !from.IsZero() {
		args = append(args, from)
		query += " AND called_at>=$" + strconv.Itoa(len(args))
	}
	if !to.IsZero() {
		args = append(args, to)
		query += " AND called_at<=$" + strconv.Itoa(len(args))
	}
	query += " ORDER BY called_at DESC"

	var list []models.TruckCall
	if err := db.conn.SelectContext(ctx, &list, query, args...); err != nil {
		db.log.Error(op, "error", err)
		return nil, coreErrors.ErrExecQuery
	}
	return list, nil
}

func (db *DB) GetTruckCallByID(ctx context.Context, id string) (models.TruckCall, error) {
	const op = "db.GetTruckCallByID"

	var tc models.TruckCall
	err := db.conn.QueryRowxContext(ctx,
		`SELECT id, warehouse_id, route_id, driver_id, forecast_value, threshold_value,
		        called_at, status, timeliness, timeliness_at, actual_containers
		 FROM truck_calls WHERE id=$1`, id,
	).StructScan(&tc)
	if err != nil {
		if isNotFound(err) {
			return models.TruckCall{}, coreErrors.ErrNotFoundTruckCall
		}
		db.log.Error(op, "error", err)
		return models.TruckCall{}, coreErrors.ErrExecQuery
	}
	return tc, nil
}

func (db *DB) GetPendingTruckCallForDriver(ctx context.Context, driverID string) (*models.TruckCall, error) {
	const op = "db.GetPendingTruckCallForDriver"

	driver, err := db.GetDriverByUserID(ctx, driverID)
	if err != nil {
		return nil, err
	}
	if driver.WarehouseID == nil || driver.RouteID == nil {
		return nil, nil
	}

	var tc models.TruckCall
	err = db.conn.QueryRowxContext(ctx,
		`SELECT id, warehouse_id, route_id, driver_id, forecast_value, threshold_value,
		        called_at, status, timeliness, timeliness_at, actual_containers
		 FROM truck_calls
		 WHERE warehouse_id=$1 AND route_id=$2 AND status='pending'
		 ORDER BY called_at DESC LIMIT 1`,
		*driver.WarehouseID, *driver.RouteID,
	).StructScan(&tc)
	if err != nil {
		if isNotFound(err) {
			return nil, nil
		}
		db.log.Error(op, "error", err)
		return nil, coreErrors.ErrExecQuery
	}
	return &tc, nil
}

func (db *DB) UpdateTruckCallTimeliness(ctx context.Context, id, timeliness string, actual *int) error {
	const op = "db.UpdateTruckCallTimeliness"

	_, err := db.conn.ExecContext(ctx,
		`UPDATE truck_calls
		 SET timeliness=$1, timeliness_at=NOW(), actual_containers=$2, status='completed'
		 WHERE id=$3`,
		timeliness, actual, id)
	if err != nil {
		db.log.Error(op, "error", err)
		return coreErrors.ErrExecQuery
	}
	return nil
}

func (db *DB) GetTruckCallAccuracy(ctx context.Context, warehouseID, routeID string) (models.TruckCallAccuracy, error) {
	const op = "db.GetTruckCallAccuracy"

	query := `SELECT
		COUNT(*) AS total_calls,
		COUNT(*) FILTER (WHERE timeliness='on_time') AS on_time_calls,
		COUNT(*) FILTER (WHERE timeliness='late')    AS late_calls,
		COUNT(*) FILTER (WHERE timeliness='early')   AS early_calls,
		COALESCE(AVG(forecast_value), 0)             AS avg_forecast,
		COALESCE(AVG(actual_containers), 0)          AS avg_actual
	FROM truck_calls WHERE 1=1`
	args := []any{}

	if warehouseID != "" {
		args = append(args, warehouseID)
		query += " AND warehouse_id=$" + strconv.Itoa(len(args))
	}
	if routeID != "" {
		args = append(args, routeID)
		query += " AND route_id=$" + strconv.Itoa(len(args))
	}

	row := db.conn.QueryRowxContext(ctx, query, args...)

	var total, onTime, late, early int
	var avgForecast, avgActual float64
	if err := row.Scan(&total, &onTime, &late, &early, &avgForecast, &avgActual); err != nil {
		db.log.Error(op, "error", err)
		return models.TruckCallAccuracy{}, coreErrors.ErrExecQuery
	}

	var accuracy float64
	if total > 0 {
		accuracy = float64(onTime) / float64(total) * 100
	}

	return models.TruckCallAccuracy{
		WarehouseID:  warehouseID,
		RouteID:      routeID,
		TotalCalls:   total,
		OnTimeCalls:  onTime,
		LateCalls:    late,
		EarlyCalls:   early,
		AccuracyRate: accuracy,
		AvgForecast:  avgForecast,
		AvgActual:    avgActual,
	}, nil
}