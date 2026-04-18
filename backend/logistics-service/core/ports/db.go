package ports

import (
	"context"
	"logistics-service/logistics-service/core/models"
	"time"
)

type DB interface {
	// Warehouses
	CreateWarehouse(ctx context.Context, req models.RequestCreateWarehouse) (models.Warehouse, error)
	GetWarehouses(ctx context.Context) ([]models.Warehouse, error)
	GetWarehouseByID(ctx context.Context, id string) (models.Warehouse, error)
	DeleteWarehouse(ctx context.Context, id string) error

	// Routes
	CreateRoute(ctx context.Context, warehouseID string, req models.RequestCreateRoute) (models.Route, error)
	GetRoutesByWarehouse(ctx context.Context, warehouseID string) ([]models.Route, error)
	GetRouteByID(ctx context.Context, id string) (models.Route, error)
	DeleteRoute(ctx context.Context, id string) error

	// Users
	CreateUser(ctx context.Context, username, passwordHash, role string) (models.User, error)
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
	GetUserByID(ctx context.Context, id string) (models.User, error)

	// Drivers
	GetDriverByUserID(ctx context.Context, userID string) (models.Driver, error)
	GetDrivers(ctx context.Context) ([]models.Driver, error)
	AssignDriver(ctx context.Context, req models.RequestAssignDriver) (models.Driver, error)

	// Thresholds
	SetThreshold(ctx context.Context, req models.RequestSetThreshold) (models.Threshold, error)
	GetThresholds(ctx context.Context, warehouseID, routeID string) ([]models.Threshold, error)
	GetThreshold(ctx context.Context, warehouseID, routeID string) (models.Threshold, error)

	// Forecasts
	SaveForecast(ctx context.Context, f models.Forecast) (models.Forecast, error)
	GetForecasts(ctx context.Context, q models.RequestForecastQuery) ([]models.Forecast, error)
	UpdateForecastActual(ctx context.Context, id string, actual float64) error

	// Truck calls
	CreateTruckCall(ctx context.Context, tc models.TruckCall) (models.TruckCall, error)
	GetTruckCalls(ctx context.Context, warehouseID, routeID string, from, to time.Time) ([]models.TruckCall, error)
	GetTruckCallByID(ctx context.Context, id string) (models.TruckCall, error)
	GetPendingTruckCallForDriver(ctx context.Context, driverID string) (*models.TruckCall, error)
	UpdateTruckCallTimeliness(ctx context.Context, id, timeliness string, actual *int) error
	GetTruckCallAccuracy(ctx context.Context, warehouseID, routeID string) (models.TruckCallAccuracy, error)

	// Raw data
	InsertRawData(ctx context.Context, points []models.RawDataPoint) (int, error)
	GetRawData(ctx context.Context, from, to *time.Time) ([]models.RawDataPoint, error)
}