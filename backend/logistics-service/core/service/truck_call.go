package service

import (
	"context"
	coreErrors "logistics-service/logistics-service/core/errors"
	"logistics-service/logistics-service/core/models"
	"time"
)

func (s *Service) GetTruckCalls(ctx context.Context, warehouseID, routeID string, from, to time.Time) ([]models.TruckCall, error) {
	const op = "service.GetTruckCalls"

	list, err := s.db.GetTruckCalls(ctx, warehouseID, routeID, from, to)
	if err != nil {
		s.log.Error(op, "error", err)
		return nil, err
	}
	return list, nil
}

func (s *Service) GetTruckCallAccuracy(ctx context.Context, warehouseID, routeID string) (models.TruckCallAccuracy, error) {
	const op = "service.GetTruckCallAccuracy"

	acc, err := s.db.GetTruckCallAccuracy(ctx, warehouseID, routeID)
	if err != nil {
		s.log.Error(op, "error", err)
		return models.TruckCallAccuracy{}, err
	}
	return acc, nil
}

func (s *Service) ReportTimeliness(ctx context.Context, truckCallID string, req models.RequestReportTimeliness) (models.TruckCall, error) {
	const op = "service.ReportTimeliness"

	tc, err := s.db.GetTruckCallByID(ctx, truckCallID)
	if err != nil {
		return models.TruckCall{}, coreErrors.ErrNotFoundTruckCall
	}

	if err = s.db.UpdateTruckCallTimeliness(ctx, tc.ID, req.Timeliness, req.ActualContainers); err != nil {
		s.log.Error(op, "error", err)
		return models.TruckCall{}, err
	}

	tc.Timeliness = &req.Timeliness
	tc.ActualContainers = req.ActualContainers
	return tc, nil
}