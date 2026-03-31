package models

import "time"

type RawDataPoint struct {
	ID           string     `db:"id"            json:"id"`
	RouteID      string     `db:"route_id"      json:"route_id"`
	OfficeFromID string     `db:"office_from_id" json:"office_from_id"`
	Timestamp    time.Time  `db:"timestamp"     json:"timestamp"`
	Status1      *float64   `db:"status_1"      json:"status_1"`
	Status2      *float64   `db:"status_2"      json:"status_2"`
	Status3      *float64   `db:"status_3"      json:"status_3"`
	Status4      *float64   `db:"status_4"      json:"status_4"`
	Status5      *float64   `db:"status_5"      json:"status_5"`
	Status6      *float64   `db:"status_6"      json:"status_6"`
	Status7      *float64   `db:"status_7"      json:"status_7"`
	Status8      *float64   `db:"status_8"      json:"status_8"`
	Target2H     *float64   `db:"target_2h"     json:"target_2h"`
	CreatedAt    time.Time  `db:"created_at"    json:"created_at"`
}

type RequestIngestData struct {
	DataPoints []RawDataPoint `json:"data_points" validate:"required,min=1"`
}

type ResponseIngestData struct {
	Inserted int `json:"inserted"`
}

type RequestRetrainModel struct {
	FromDate *string `json:"from_date"`
	ToDate   *string `json:"to_date"`
}

type ResponseRetrainModel struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}