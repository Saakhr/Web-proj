package handlers

import (
	"crypto/rsa"

	"github.com/Saakhr/Web-proj/pkg/models"
	"github.com/Saakhr/Web-proj/pkg/v1/services"
	"github.com/Saakhr/Web-proj/pkg/v1/utility"
	"github.com/Saakhr/Web-proj/templates"
	"github.com/gofiber/fiber/v2"
)

func HandleLogin(c *fiber.Ctx, privateKey *rsa.PrivateKey) error {
  user, err := services.GetUserFromCookie(c, privateKey)
	if user != nil {
		return c.Redirect("/v1")
	}
	email := c.FormValue("email")
	password := c.FormValue("password")

	admin, err := models.GetAdminByEmail(email)
	if err == nil {
		if err := admin.CheckPassword(password); err != nil {
			return utility.Render(c, templates.Login("Invalid credentials", nil))
		}

		token, err := services.GenerateJWT(
			admin.ID,
			admin.FullName,
			admin.Email,
			"admin",
			privateKey,
		)
		if err != nil {
			return utility.Render(c, templates.Login("Internal server error", nil))
		}

		services.SetJWTCookie(c, token)
		return c.Redirect("/v1/admin/dashboard")
	}
	student, err := models.GetStudentByEmail(email)
	if err != nil {
		return utility.Render(c, templates.Login("Invalid credentials", nil))
	}

	if err := student.CheckPassword(password); err != nil {
		return utility.Render(c, templates.Login("Invalid credentials", nil))
	}

	token, err := services.GenerateJWT(
		student.ID,
		student.FirstName+" "+student.LastName,
		student.Email,
		"student",
		privateKey,
	)
	if err != nil {
		return utility.Render(c, templates.Login("Internal server error", nil))
	}

	services.SetJWTCookie(c, token)
	return c.Redirect("/v1/student/dashboard")

}
