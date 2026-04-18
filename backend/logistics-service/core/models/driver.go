package models

import "time"

type Driver struct {
	ID          string     `db:"id"           json:"id"`
	UserID      string     `db:"user_id"      json:"user_id"`
	Username    string     `db:"username"     json:"username"`
	WarehouseID *string    `db:"warehouse_id" json:"warehouse_id"`
	RouteID     *string    `db:"route_id"     json:"route_id"`
	CreatedAt   time.Time  `db:"created_at"   json:"created_at"`
}

type RequestAssignDriver struct {
	DriverID    string `json:"driver_id"    validate:"required"`
	WarehouseID string `json:"warehouse_id" validate:"required"`
	RouteID     string `json:"route_id"     validate:"required"`
}