package service

import (
	"context"
	coreErrors "logistics-service/logistics-service/core/errors"
	"logistics-service/logistics-service/core/models"
)

func (s *Service) GetDrivers(ctx context.Context) ([]models.Driver, error) {
	const op = "service.GetDrivers"

	list, err := s.db.GetDrivers(ctx)
	if err != nil {
		s.log.Error(op, "error", err)
		return nil, err
	}
	return list, nil
}

func (s *Service) AssignDriver(ctx context.Context, req models.RequestAssignDriver) (models.Driver, error) {
	const op = "service.AssignDriver"

	d, err := s.db.AssignDriver(ctx, req)
	if err != nil {
		s.log.Error(op, "error", err)
		return models.Driver{}, err
	}
	return d, nil
}

func (s *Service) GetDriverSignal(ctx context.Context, userID string) (*models.TruckCall, error) {
	const op = "service.GetDriverSignal"

	driver, err := s.db.GetDriverByUserID(ctx, userID)
	if err != nil {
		return nil, coreErrors.ErrNotFoundDriver
	}

	tc, err := s.db.GetPendingTruckCallForDriver(ctx, driver.ID)
	if err != nil {
		s.log.Error(op, "error", err)
		return nil, err
	}
	return tc, nil
}

func (s *Service) GetDriverStats(ctx context.Context, userID string) (models.TruckCallAccuracy, error) {
	const op = "service.GetDriverStats"

	driver, err := s.db.GetDriverByUserID(ctx, userID)
	if err != nil {
		return models.TruckCallAccuracy{}, coreErrors.ErrNotFoundDriver
	}

	warehouseID := ""
	routeID := ""
	if driver.WarehouseID != nil {
		warehouseID = *driver.WarehouseID
	}
	if driver.RouteID != nil {
		routeID = *driver.RouteID
	}

	acc, err := s.db.GetTruckCallAccuracy(ctx, warehouseID, routeID)
	if err != nil {
		s.log.Error(op, "error", err)
		return models.TruckCallAccuracy{}, err
	}
	return acc, nil
}