package ml_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	coreErrors "logistics-service/logistics-service/core/errors"
	"logistics-service/logistics-service/core/models"
	"logistics-service/logistics-service/core/ports"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func New(baseURL string) ports.MLClient {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) Predict(ctx context.Context, req models.MLPredictRequest) (models.MLPredictResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return models.MLPredictResponse{}, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/predict", bytes.NewReader(body))
	if err != nil {
		return models.MLPredictResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return models.MLPredictResponse{}, coreErrors.ErrMLServiceUnavailable
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.MLPredictResponse{}, fmt.Errorf("%w: status %d", coreErrors.ErrMLServiceError, resp.StatusCode)
	}

	var result models.MLPredictResponse
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return models.MLPredictResponse{}, coreErrors.ErrMLServiceError
	}

	return result, nil
}

func (c *Client) Retrain(ctx context.Context, req models.RequestRetrainModel) (models.ResponseRetrainModel, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return models.ResponseRetrainModel{}, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/retrain", bytes.NewReader(body))
	if err != nil {
		return models.ResponseRetrainModel{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return models.ResponseRetrainModel{}, coreErrors.ErrMLServiceUnavailable
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.ResponseRetrainModel{}, fmt.Errorf("%w: status %d", coreErrors.ErrMLServiceError, resp.StatusCode)
	}

	var result models.ResponseRetrainModel
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return models.ResponseRetrainModel{}, coreErrors.ErrMLServiceError
	}

	return result, nil
}