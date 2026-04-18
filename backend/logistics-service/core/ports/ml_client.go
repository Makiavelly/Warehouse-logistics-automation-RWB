package ports

import (
	"context"
	"logistics-service/logistics-service/core/models"
)

type MLClient interface {
	Predict(ctx context.Context, req models.MLPredictRequest) (models.MLPredictResponse, error)
	Retrain(ctx context.Context, req models.RequestRetrainModel) (models.ResponseRetrainModel, error)
}