package v1

import (
	"crypto/rsa"
	"strings"

	"github.com/Saakhr/Web-proj/pkg/v1/utility"
	templates "github.com/Saakhr/Web-proj/templates"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func NewAuthMiddleware(privateKey *rsa.PrivateKey) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: privateKey.Public(),
			JWTAlg: jwtware.RS256},
	})
}
func NotFoundMiddleware(c *fiber.Ctx) error {
  if !strings.HasPrefix(c.Path(), "/static/") && 
  !strings.HasPrefix(c.Path(), "/css/") {
    c.Status(fiber.StatusNotFound)
    return utility.Render(c, templates.NotFound())
  }
  return c.Next()
}
