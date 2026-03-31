package ml_client_grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	grpcStatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	coreErrors "logistics-service/logistics-service/core/errors"
	"logistics-service/logistics-service/core/models"
	"logistics-service/logistics-service/core/ports"
	pb "logistics-service/proto/logistics"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.MLServiceClient
}

func New(address string) (ports.MLClient, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:   conn,
		client: pb.NewMLServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Predict(ctx context.Context, req models.MLPredictRequest) (models.MLPredictResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	reply, err := c.client.Predict(ctx, &pb.PredictRequest{
		RouteId:      req.RouteID,
		OfficeFromId: req.OfficeFromID,
		TimestampUnix: req.Timestamp.Unix(),
		HorizonHours: int32(req.HorizonHours),
	})
	if err != nil {
		if grpcStatus.Code(err) == codes.Unavailable {
			return models.MLPredictResponse{}, coreErrors.ErrMLServiceUnavailable
		}
		return models.MLPredictResponse{}, coreErrors.ErrMLServiceError
	}

	return models.MLPredictResponse{PredictedCount: reply.PredictedCount}, nil
}

func (c *Client) Retrain(ctx context.Context, req models.RequestRetrainModel) (models.ResponseRetrainModel, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	pbReq := &pb.RetrainRequest{}
	if req.FromDate != nil {
		if t, err := time.Parse(time.RFC3339, *req.FromDate); err == nil {
			pbReq.FromDateUnix = t.Unix()
		}
	}
	if req.ToDate != nil {
		if t, err := time.Parse(time.RFC3339, *req.ToDate); err == nil {
			pbReq.ToDateUnix = t.Unix()
		}
	}

	reply, err := c.client.Retrain(ctx, pbReq)
	if err != nil {
		if grpcStatus.Code(err) == codes.Unavailable {
			return models.ResponseRetrainModel{}, coreErrors.ErrMLServiceUnavailable
		}
		return models.ResponseRetrainModel{}, coreErrors.ErrMLServiceError
	}

	return models.ResponseRetrainModel{
		Status:  reply.Status,
		Message: reply.Message,
	}, nil
}

func (c *Client) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := c.client.Ping(ctx, &emptypb.Empty{})
	if err != nil {
		return coreErrors.ErrMLServiceUnavailable
	}
	return nil
}