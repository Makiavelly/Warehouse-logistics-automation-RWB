package models

import "time"

type Forecast struct {
	ID             string     `db:"id"              json:"id"`
	WarehouseID    string     `db:"warehouse_id"    json:"warehouse_id"`
	RouteID        string     `db:"route_id"        json:"route_id"`
	ForecastTime   time.Time  `db:"forecast_time"   json:"forecast_time"`
	HorizonHours   int        `db:"horizon_hours"   json:"horizon_hours"`
	PredictedCount float64    `db:"predicted_count" json:"predicted_count"`
	ActualCount    *float64   `db:"actual_count"    json:"actual_count"`
	CreatedAt      time.Time  `db:"created_at"      json:"created_at"`
}

type RequestPredict struct {
	WarehouseID  string    `json:"warehouse_id"  validate:"required"`
	RouteID      string    `json:"route_id"      validate:"required"`
	ForecastTime time.Time `json:"forecast_time" validate:"required"`
	HorizonHours int       `json:"horizon_hours"`
}

type MLPredictRequest struct {
	RouteID      string    `json:"route_id"`
	OfficeFromID string    `json:"office_from_id"`
	Timestamp    time.Time `json:"timestamp"`
	HorizonHours int       `json:"horizon_hours"`
}

type MLPredictResponse struct {
	PredictedCount float64 `json:"predicted_count"`
}

type RequestForecastQuery struct {
	WarehouseID string
	RouteID     string
	From        time.Time
	To          time.Time
}