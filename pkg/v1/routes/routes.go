package v1

import (
	"crypto/rsa"

	"github.com/Saakhr/Web-proj/pkg/v1/handlers"
	"github.com/Saakhr/Web-proj/pkg/v1/services"
	"github.com/Saakhr/Web-proj/pkg/v1/utility"
	"github.com/Saakhr/Web-proj/templates"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var (
	privateKey *rsa.PrivateKey
)

func GetRoutes(Key *rsa.PrivateKey) *fiber.App {
	v1 := fiber.New()
	privateKey = Key

	v1.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/v1/announcements")
	})

	v1.Get("/announcements", func(c *fiber.Ctx)error {
    dept:=c.Query("dept")
    if dept == ""{
      return handlers.AnnouncementList(c,privateKey,"general")
    }else{
      return handlers.ListOfAnnouncements(c,dept)
    }
	})
	//login
	v1.Get("/logout", func(c *fiber.Ctx) error {
    services.ClearJWTCookie(c)
    return c.Redirect("/v1")
  })
	v1.Get("/login", func(c *fiber.Ctx) error {

    _, err := services.GetUserFromCookie(c, privateKey)
		
		if err != nil {
			return utility.Render(c, templates.Login("", nil))
		}
    return c.Redirect("/v1")
	})
	v1.Post("/login", func(c *fiber.Ctx) error {
		return handlers.HandleLogin(c, privateKey)
	})
	SetupAdminRoutes(v1, Key)
  SetupStudentRoutes(v1,Key)


	// Restricted Routes
	v1.Get("/2", accessible)

	return v1
}

func accessible(c *fiber.Ctx) error {
	return c.SendString("Accessible")
}

func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome " + name)
}
