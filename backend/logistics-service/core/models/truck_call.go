package models

import "time"

const (
	TruckCallStatusPending   = "pending"
	TruckCallStatusAccepted  = "accepted"
	TruckCallStatusCompleted = "completed"
	TruckCallStatusMissed    = "missed"

	TimelinessOnTime = "on_time"
	TimelinessLate   = "late"
	TimelinessEarly  = "early"
)

type TruckCall struct {
	ID               string     `db:"id"                 json:"id"`
	WarehouseID      string     `db:"warehouse_id"       json:"warehouse_id"`
	RouteID          string     `db:"route_id"           json:"route_id"`
	DriverID         *string    `db:"driver_id"          json:"driver_id"`
	ForecastValue    float64    `db:"forecast_value"     json:"forecast_value"`
	ThresholdValue   float64    `db:"threshold_value"    json:"threshold_value"`
	CalledAt         time.Time  `db:"called_at"          json:"called_at"`
	Status           string     `db:"status"             json:"status"`
	Timeliness       *string    `db:"timeliness"         json:"timeliness"`
	TimelinessAt     *time.Time `db:"timeliness_at"      json:"timeliness_at"`
	ActualContainers *int       `db:"actual_containers"  json:"actual_containers"`
}

type RequestReportTimeliness struct {
	Timeliness       string `json:"timeliness"        validate:"required,oneof=on_time late early"`
	ActualContainers *int   `json:"actual_containers"`
}

type TruckCallAccuracy struct {
	WarehouseID    string  `json:"warehouse_id"`
	RouteID        string  `json:"route_id"`
	TotalCalls     int     `json:"total_calls"`
	OnTimeCalls    int     `json:"on_time_calls"`
	LateCalls      int     `json:"late_calls"`
	EarlyCalls     int     `json:"early_calls"`
	AccuracyRate   float64 `json:"accuracy_rate"`
	AvgForecast    float64 `json:"avg_forecast"`
	AvgActual      float64 `json:"avg_actual"`
}