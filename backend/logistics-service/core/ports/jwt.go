package ports

import "logistics-service/logistics-service/core/models"

type TokenGenerator interface {
	Generate(claims models.TokenClaims) (string, error)
	Validate(token string) (models.TokenClaims, error)
}