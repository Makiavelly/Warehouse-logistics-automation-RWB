package JWT

import (
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"logistics-service/logistics-service/config"
	coreErrors "logistics-service/logistics-service/core/errors"
	"logistics-service/logistics-service/core/models"
	"logistics-service/logistics-service/core/ports"
)

type JWT struct {
	secret []byte
	ttl    time.Duration
}

func NewJWT(cfg config.AuthConfig) ports.TokenGenerator {
	return &JWT{
		secret: []byte(cfg.SecretKey),
		ttl:    cfg.TTLHours * time.Hour,
	}
}

func (a *JWT) Generate(claims models.TokenClaims) (string, error) {
	jwtClaims := jwt.MapClaims{
		"user_id": claims.UserID,
		"role":    claims.Role,
		"exp":     claims.Exp.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	return token.SignedString(a.secret)
}

func (a *JWT) Validate(tokenStr string) (models.TokenClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			slog.Warn("invalid signing method")
			return nil, coreErrors.NewErrInvalidAuthToken("invalid signature")
		}
		return a.secret, nil
	})
	if err != nil || !token.Valid {
		return models.TokenClaims{}, coreErrors.NewErrInvalidAuthToken("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return models.TokenClaims{}, coreErrors.NewErrInvalidAuthToken("invalid claims")
	}

	userID, _ := claims["user_id"].(string)
	role, _ := claims["role"].(string)
	expUnix, _ := claims["exp"].(float64)
	exp := time.Unix(int64(expUnix), 0)

	return models.TokenClaims{
		UserID: userID,
		Role:   role,
		Exp:    exp,
	}, nil
}