package models

import "time"

type Warehouse struct {
	ID           string    `db:"id"             json:"id"`
	Name         string    `db:"name"           json:"name"`
	OfficeFromID string    `db:"office_from_id" json:"office_from_id"`
	Address      string    `db:"address"        json:"address"`
	CreatedAt    time.Time `db:"created_at"     json:"created_at"`
}

type RequestCreateWarehouse struct {
	Name         string `json:"name"           validate:"required"`
	OfficeFromID string `json:"office_from_id" validate:"required"`
	Address      string `json:"address"`
}