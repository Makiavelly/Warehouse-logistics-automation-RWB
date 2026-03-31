package service

import (
	"context"
	coreErrors "logistics-service/logistics-service/core/errors"
	"logistics-service/logistics-service/core/models"
	"time"
)

func (s *Service) RequestForecast(ctx context.Context, req models.RequestPredict) (models.Forecast, error) {
	const op = "service.RequestForecast"

	warehouse, err := s.db.GetWarehouseByID(ctx, req.WarehouseID)
	if err != nil {
		return models.Forecast{}, coreErrors.ErrNotFoundWarehouse
	}

	route, err := s.db.GetRouteByID(ctx, req.RouteID)
	if err != nil {
		return models.Forecast{}, coreErrors.ErrNotFoundRoute
	}

	horizonHours := req.HorizonHours
	if horizonHours <= 0 {
		horizonHours = 2
	}

	mlResp, err := s.ml.Predict(ctx, models.MLPredictRequest{
		RouteID:      route.RouteID,
		OfficeFromID: warehouse.OfficeFromID,
		Timestamp:    req.ForecastTime,
		HorizonHours: horizonHours,
	})
	if err != nil {
		s.log.Error(op+" ml predict failed", "error", err)
		return models.Forecast{}, coreErrors.ErrMLServiceUnavailable
	}

	forecast := models.Forecast{
		WarehouseID:    req.WarehouseID,
		RouteID:        req.RouteID,
		ForecastTime:   req.ForecastTime,
		HorizonHours:   horizonHours,
		PredictedCount: mlResp.PredictedCount,
	}

	saved, err := s.db.SaveForecast(ctx, forecast)
	if err != nil {
		s.log.Error(op, "error", err)
		return models.Forecast{}, err
	}

	// auto-trigger truck call if forecast exceeds threshold
	go s.checkThresholdAndCallTruck(req.WarehouseID, req.RouteID, mlResp.PredictedCount)

	return saved, nil
}

func (s *Service) checkThresholdAndCallTruck(warehouseID, routeID string, predicted float64) {
	ctx := context.Background()

	threshold, err := s.db.GetThreshold(ctx, warehouseID, routeID)
	if err != nil {
		return
	}

	if predicted < threshold.Value {
		return
	}

	s.log.Info("threshold exceeded, creating truck call",
		"warehouse_id", warehouseID,
		"route_id", routeID,
		"predicted", predicted,
		"threshold", threshold.Value,
	)

	tc := models.TruckCall{
		WarehouseID:    warehouseID,
		RouteID:        routeID,
		ForecastValue:  predicted,
		ThresholdValue: threshold.Value,
		CalledAt:       time.Now(),
		Status:         models.TruckCallStatusPending,
	}

	if _, err = s.db.CreateTruckCall(ctx, tc); err != nil {
		s.log.Error("checkThresholdAndCallTruck create truck call failed", "error", err)
	}
}

func (s *Service) GetForecasts(ctx context.Context, q models.RequestForecastQuery) ([]models.Forecast, error) {
	const op = "service.GetForecasts"

	list, err := s.db.GetForecasts(ctx, q)
	if err != nil {
		s.log.Error(op, "error", err)
		return nil, err
	}
	return list, nil
}