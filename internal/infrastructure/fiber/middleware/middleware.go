package middleware

import (
	"checker/internal/api/providers"
	"checker/internal/config"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		cfg, err := config.LoadConfig("config.json")
		if err != nil {
			fmt.Printf("failed to load config file: %v", err)
		}

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header is missing"})
		}

		encodedJWT := strings.TrimPrefix(authHeader, "Bearer ")
		if encodedJWT == authHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid authorization scheme"})
		}

		decodedJWT, err := jwt.Parse(encodedJWT, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWT.JwtSecretKey), nil
		})
		if err != nil || !decodedJWT.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		claims, ok := decodedJWT.Claims.(jwt.MapClaims)
		if !ok || claims["exp"] == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
		}

		expiration := int64(claims["exp"].(float64))
		if time.Now().Unix() > expiration {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token expired"})
		}

		c.Locals("jwtClaims", &providers.JWTClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Unix(expiration, 0)),
				Subject:   claims["sub"].(string),
				Issuer:    claims["iss"].(string),
			},
		})

		return c.Next()
	}
}
