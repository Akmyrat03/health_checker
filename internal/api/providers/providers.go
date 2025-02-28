package providers

import (
	"checker/internal/infrastructure/pgx"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

func GetJwtClaims(c *fiber.Ctx) *JWTClaims {
	claims := c.Locals("jwtClaims").(*JWTClaims)
	return claims
}

func GetDbPool() (*pgxpool.Pool, error) {
	pool, err := pgx.PostgresPool()
	if err != nil {
		return nil, err
	}

	return pool, err
}
