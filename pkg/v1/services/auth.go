package services

// import (
// 	"crypto/rsa"
// 	"time"
//
// 	"github.com/Saakhr/Web-proj/pkg/models"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/golang-jwt/jwt/v5"
// )
//
// func GenerateJWTAdmin(user *models.Admin, privateKey *rsa.PrivateKey) (string, error) {
// 	claims := jwt.MapClaims{
// 		"email": user.Email,
//     "admin":true,
// 		"exp":   time.Now().Add(time.Hour * 72).Unix(),
// 	}
//
// 	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
// 	return token.SignedString(privateKey)
// }
//
// func GetAdminFromToken(c *fiber.Ctx) (*models.Admin, error) {
// 	userToken := c.Locals("user").(*jwt.Token)
// 	claims := userToken.Claims.(jwt.MapClaims)
//
// 	user := &models.Admin{
// 		ID:       int(claims["id"].(float64)),
// 		FullName: claims["name"].(string),
// 		Email:    claims["email"].(string),
// 	}
//
// 	return user, nil
// }
