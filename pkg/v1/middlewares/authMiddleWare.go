package v1

import (
	"crypto/rsa"
	"fmt"
	"strings"

	"github.com/Saakhr/Web-proj/pkg/v1/services"
	"github.com/Saakhr/Web-proj/pkg/v1/utility"
	templates "github.com/Saakhr/Web-proj/templates"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func NewAuthMiddleware(privateKey *rsa.PrivateKey, role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// First try to get token from cookie
		token := c.Cookies("jwt")
		if token == "" {
			// Fallback to Authorization header
			authHeader := c.Get("Authorization")
			if len(authHeader) > 7 && strings.EqualFold(authHeader[:7], "Bearer ") {
				token = authHeader[7:]
			}
		}

		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization token required",
			})
		}

		// Parse the token with enhanced error handling
		parsedToken, err := jwt.ParseWithClaims(token, &services.Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate the algorithm
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return privateKey.Public(), nil // Use the public key for verification
		})

		// Detailed error handling

		if !parsedToken.Valid || err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		claims, ok := parsedToken.Claims.(*services.Claims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		if claims.Role != role {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient permissions",
			})
		}

    c.Locals("user", claims) // Store the actual claims, not the token
		return c.Next()
	}
}
func NotFoundMiddleware(c *fiber.Ctx) error {
	if !strings.HasPrefix(c.Path(), "/static/") &&
		!strings.HasPrefix(c.Path(), "/css/") {
		c.Status(fiber.StatusNotFound)
		return utility.Render(c, templates.NotFound())
	}
	return c.Next()
}
