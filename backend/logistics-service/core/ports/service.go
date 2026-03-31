package ports

import (
	"context"
	"logistics-service/logistics-service/core/models"
	"time"
)

type Service interface {
	// Auth
	Login(ctx context.Context, req models.RequestLogin) (models.ResponseLogin, error)
	Register(ctx context.Context, req models.RequestRegister) (models.User, error)

	// Admin — Warehouses
	CreateWarehouse(ctx context.Context, req models.RequestCreateWarehouse) (models.Warehouse, error)
	GetWarehouses(ctx context.Context) ([]models.Warehouse, error)
	DeleteWarehouse(ctx context.Context, id string) error

	// Admin — Routes
	CreateRoute(ctx context.Context, warehouseID string, req models.RequestCreateRoute) (models.Route, error)
	GetRoutesByWarehouse(ctx context.Context, warehouseID string) ([]models.Route, error)
	DeleteRoute(ctx context.Context, warehouseID, routeID string) error

	// Admin — Thresholds
	SetThreshold(ctx context.Context, req models.RequestSetThreshold) (models.Threshold, error)
	GetThresholds(ctx context.Context, warehouseID, routeID string) ([]models.Threshold, error)

	// Admin — Forecasts
	RequestForecast(ctx context.Context, req models.RequestPredict) (models.Forecast, error)
	GetForecasts(ctx context.Context, q models.RequestForecastQuery) ([]models.Forecast, error)

	// Admin — Analytics
	GetTruckCalls(ctx context.Context, warehouseID, routeID string, from, to time.Time) ([]models.TruckCall, error)
	GetTruckCallAccuracy(ctx context.Context, warehouseID, routeID string) (models.TruckCallAccuracy, error)

	// Admin — Drivers
	GetDrivers(ctx context.Context) ([]models.Driver, error)
	AssignDriver(ctx context.Context, req models.RequestAssignDriver) (models.Driver, error)

	// Driver
	GetDriverSignal(ctx context.Context, userID string) (*models.TruckCall, error)
	ReportTimeliness(ctx context.Context, truckCallID string, req models.RequestReportTimeliness) (models.TruckCall, error)
	GetDriverStats(ctx context.Context, userID string) (models.TruckCallAccuracy, error)

	// External integration
	IngestData(ctx context.Context, req models.RequestIngestData) (models.ResponseIngestData, error)
	RetrainModel(ctx context.Context, req models.RequestRetrainModel) (models.ResponseRetrainModel, error)
}