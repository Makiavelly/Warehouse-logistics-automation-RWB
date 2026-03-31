package service

import (
	"context"
	coreErrors "logistics-service/logistics-service/core/errors"
	"logistics-service/logistics-service/core/models"
)

func (s *Service) CreateRoute(ctx context.Context, warehouseID string, req models.RequestCreateRoute) (models.Route, error) {
	const op = "service.CreateRoute"

	if _, err := s.db.GetWarehouseByID(ctx, warehouseID); err != nil {
		return models.Route{}, coreErrors.ErrNotFoundWarehouse
	}

	r, err := s.db.CreateRoute(ctx, warehouseID, req)
	if err != nil {
		s.log.Error(op, "error", err)
		return models.Route{}, err
	}
	return r, nil
}

func (s *Service) GetRoutesByWarehouse(ctx context.Context, warehouseID string) ([]models.Route, error) {
	const op = "service.GetRoutesByWarehouse"

	list, err := s.db.GetRoutesByWarehouse(ctx, warehouseID)
	if err != nil {
		s.log.Error(op, "error", err)
		return nil, err
	}
	return list, nil
}

func (s *Service) DeleteRoute(ctx context.Context, warehouseID, routeID string) error {
	const op = "service.DeleteRoute"

	route, err := s.db.GetRouteByID(ctx, routeID)
	if err != nil {
		return coreErrors.ErrNotFoundRoute
	}
	if route.WarehouseID != warehouseID {
		return coreErrors.ErrNotFoundRoute
	}

	if err = s.db.DeleteRoute(ctx, routeID); err != nil {
		s.log.Error(op, "error", err)
		return err
	}
	return nil
}