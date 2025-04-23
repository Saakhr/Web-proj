package services

import (
	"crypto/rsa"
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"` // "admin" or "student"
	jwt.RegisteredClaims
}

func GenerateJWT(userID int, username, email, role string, privateKey *rsa.PrivateKey) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "school-management",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

func SetJWTCookie(c *fiber.Ctx, token string) {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(72 * time.Hour),
		HTTPOnly: true,
		Secure:   true, // Set to false in development, true in production
		SameSite: "Lax",
		Path:     "/",
	}
	c.Cookie(&cookie)
}

func ClearJWTCookie(c *fiber.Ctx) {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), // Expire immediately
		HTTPOnly: true,
		Path:     "/",
	}
	c.Cookie(&cookie)
}
func GetUserFromContext(c *fiber.Ctx) (*Claims, error) {
    // Safely get claims from context
    claims, ok := c.Locals("user").(*Claims)
    if !ok || claims == nil {
        return nil, errors.New("no valid user claims found in context")
    }

    // Create a copy to prevent modification
    user := &Claims{
        UserID:       claims.UserID,
        Username: claims.Username,
        Email:    claims.Email,
        Role:     claims.Role,
    }

    // Validate required fields
    if user.Role == "" {
        return nil, errors.New("invalid user claims")
    }

    return user, nil
}
func GetUserFromCookie(c *fiber.Ctx, privateKey *rsa.PrivateKey) (*Claims,error){
		token := c.Cookies("jwt")
		if token == "" {
			// Fallback to Authorization header
			authHeader := c.Get("Authorization")
			if len(authHeader) > 7 && strings.EqualFold(authHeader[:7], "Bearer ") {
				token = authHeader[7:]
			}
		}

		if token == "" {
			return nil,errors.New("Couldn't find a token")
		}

		// Parse the token with enhanced error handling
		parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate the algorithm
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.New("Couldn't find a token")
    }
			return privateKey.Public(), nil // Use the public key for verification
		})

		// Detailed error handling

		if !parsedToken.Valid || err != nil {
			return nil,errors.New("Error not valid token")
  }

		claims, ok := parsedToken.Claims.(*Claims)
		if !ok {
			return nil,errors.New("Bad Token")
		}

  return claims,nil
}
