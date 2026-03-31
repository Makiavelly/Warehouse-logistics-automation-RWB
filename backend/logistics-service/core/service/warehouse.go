package service

import (
	"context"
	"logistics-service/logistics-service/core/models"
)

func (s *Service) CreateWarehouse(ctx context.Context, req models.RequestCreateWarehouse) (models.Warehouse, error) {
	const op = "service.CreateWarehouse"

	w, err := s.db.CreateWarehouse(ctx, req)
	if err != nil {
		s.log.Error(op, "error", err)
		return models.Warehouse{}, err
	}
	return w, nil
}

func (s *Service) GetWarehouses(ctx context.Context) ([]models.Warehouse, error) {
	const op = "service.GetWarehouses"

	list, err := s.db.GetWarehouses(ctx)
	if err != nil {
		s.log.Error(op, "error", err)
		return nil, err
	}
	return list, nil
}

func (s *Service) DeleteWarehouse(ctx context.Context, id string) error {
	const op = "service.DeleteWarehouse"

	if err := s.db.DeleteWarehouse(ctx, id); err != nil {
		s.log.Error(op, "error", err)
		return err
	}
	return nil
}