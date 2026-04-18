package models

import "time"

type Route struct {
	ID          string    `db:"id"           json:"id"`
	WarehouseID string    `db:"warehouse_id" json:"warehouse_id"`
	RouteID     string    `db:"route_id"     json:"route_id"`
	Name        string    `db:"name"         json:"name"`
	CreatedAt   time.Time `db:"created_at"   json:"created_at"`
}

type RequestCreateRoute struct {
	RouteID string `json:"route_id" validate:"required"`
	Name    string `json:"name"`
}