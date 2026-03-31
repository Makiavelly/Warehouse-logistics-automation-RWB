package db

import (
	"context"
	"strconv"

	coreErrors "logistics-service/logistics-service/core/errors"
	"logistics-service/logistics-service/core/models"
)

func (db *DB) SaveForecast(ctx context.Context, f models.Forecast) (models.Forecast, error) {
	const op = "db.SaveForecast"

	var saved models.Forecast
	err := db.conn.QueryRowxContext(ctx,
		`INSERT INTO forecasts (warehouse_id, route_id, forecast_time, horizon_hours, predicted_count, actual_count)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, warehouse_id, route_id, forecast_time, horizon_hours, predicted_count, actual_count, created_at`,
		f.WarehouseID, f.RouteID, f.ForecastTime, f.HorizonHours, f.PredictedCount, f.ActualCount,
	).StructScan(&saved)
	if err != nil {
		db.log.Error(op, "error", err)
		return models.Forecast{}, coreErrors.ErrExecQuery
	}
	return saved, nil
}

func (db *DB) GetForecasts(ctx context.Context, q models.RequestForecastQuery) ([]models.Forecast, error) {
	const op = "db.GetForecasts"

	query := `SELECT id, warehouse_id, route_id, forecast_time, horizon_hours, predicted_count, actual_count, created_at
	          FROM forecasts WHERE 1=1`
	args := []any{}

	if q.WarehouseID != "" {
		args = append(args, q.WarehouseID)
		query += " AND warehouse_id=$" + strconv.Itoa(len(args))
	}
	if q.RouteID != "" {
		args = append(args, q.RouteID)
		query += " AND route_id=$" + strconv.Itoa(len(args))
	}
	if !q.From.IsZero() {
		args = append(args, q.From)
		query += " AND forecast_time>=$" + strconv.Itoa(len(args))
	}
	if !q.To.IsZero() {
		args = append(args, q.To)
		query += " AND forecast_time<=$" + strconv.Itoa(len(args))
	}
	query += " ORDER BY forecast_time"

	var list []models.Forecast
	if err := db.conn.SelectContext(ctx, &list, query, args...); err != nil {
		db.log.Error(op, "error", err)
		return nil, coreErrors.ErrExecQuery
	}
	return list, nil
}

func (db *DB) UpdateForecastActual(ctx context.Context, id string, actual float64) error {
	const op = "db.UpdateForecastActual"

	_, err := db.conn.ExecContext(ctx,
		`UPDATE forecasts SET actual_count=$1 WHERE id=$2`, actual, id)
	if err != nil {
		db.log.Error(op, "error", err)
		return coreErrors.ErrExecQuery
	}
	return nil
}