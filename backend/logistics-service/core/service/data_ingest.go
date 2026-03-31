package service

import (
	"context"
	"logistics-service/logistics-service/core/models"
)

func (s *Service) IngestData(ctx context.Context, req models.RequestIngestData) (models.ResponseIngestData, error) {
	const op = "service.IngestData"

	n, err := s.db.InsertRawData(ctx, req.DataPoints)
	if err != nil {
		s.log.Error(op, "error", err)
		return models.ResponseIngestData{}, err
	}

	return models.ResponseIngestData{Inserted: n}, nil
}

func (s *Service) RetrainModel(ctx context.Context, req models.RequestRetrainModel) (models.ResponseRetrainModel, error) {
	const op = "service.RetrainModel"

	resp, err := s.ml.Retrain(ctx, req)
	if err != nil {
		s.log.Error(op, "error", err)
		return models.ResponseRetrainModel{}, err
	}

	return resp, nil
}