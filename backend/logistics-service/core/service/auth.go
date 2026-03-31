package service

import (
	"context"
	coreErrors "logistics-service/logistics-service/core/errors"
	"logistics-service/logistics-service/core/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Login(ctx context.Context, req models.RequestLogin) (models.ResponseLogin, error) {
	const op = "service.Login"

	user, err := s.db.GetUserByUsername(ctx, req.Username)
	if err != nil {
		s.log.Warn(op, "error", err)
		return models.ResponseLogin{}, coreErrors.ErrInvalidCredentials
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return models.ResponseLogin{}, coreErrors.ErrInvalidCredentials
	}

	token, err := s.tokenGen.Generate(models.TokenClaims{
		UserID: user.ID,
		Role:   user.Role,
		Exp:    time.Now().Add(72 * time.Hour),
	})
	if err != nil {
		s.log.Error(op, "error", err)
		return models.ResponseLogin{}, err
	}

	return models.ResponseLogin{Token: token, Role: user.Role}, nil
}

func (s *Service) Register(ctx context.Context, req models.RequestRegister) (models.User, error) {
	const op = "service.Register"

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error(op, "error", err)
		return models.User{}, err
	}

	user, err := s.db.CreateUser(ctx, req.Username, string(hash), req.Role)
	if err != nil {
		s.log.Error(op, "error", err)
		return models.User{}, err
	}

	return user, nil
}