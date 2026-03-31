package models

import "time"

type Threshold struct {
	ID          string    `db:"id"           json:"id"`
	WarehouseID string    `db:"warehouse_id" json:"warehouse_id"`
	RouteID     string    `db:"route_id"     json:"route_id"`
	Value       float64   `db:"value"        json:"value"`
	UpdatedAt   time.Time `db:"updated_at"   json:"updated_at"`
}

type RequestSetThreshold struct {
	WarehouseID string  `json:"warehouse_id" validate:"required"`
	RouteID     string  `json:"route_id"     validate:"required"`
	Value       float64 `json:"value"        validate:"required,gt=0"`
}