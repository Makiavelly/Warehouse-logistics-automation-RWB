package service

import "logistics-service/logistics-service/core/ports"

type Service struct {
	log      ports.Logger
	tokenGen ports.TokenGenerator
	db       ports.DB
	ml       ports.MLClient
}

func NewService(log ports.Logger, tokenGen ports.TokenGenerator, db ports.DB, ml ports.MLClient) *Service {
	return &Service{
		log:      log,
		tokenGen: tokenGen,
		db:       db,
		ml:       ml,
	}
}