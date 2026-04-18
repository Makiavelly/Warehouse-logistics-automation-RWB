package service

import (
	"context"
	"logistics-service/logistics-service/core/models"
)

func (s *Service) SetThreshold(ctx context.Context, req models.RequestSetThreshold) (models.Threshold, error) {
	const op = "service.SetThreshold"

	t, err := s.db.SetThreshold(ctx, req)
	if err != nil {
		s.log.Error(op, "error", err)
		return models.Threshold{}, err
	}
	return t, nil
}

func (s *Service) GetThresholds(ctx context.Context, warehouseID, routeID string) ([]models.Threshold, error) {
	const op = "service.GetThresholds"

	list, err := s.db.GetThresholds(ctx, warehouseID, routeID)
	if err != nil {
		s.log.Error(op, "error", err)
		return nil, err
	}
	return list, nil
}