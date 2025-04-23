package handlers

import (
	"crypto/rsa"

	"github.com/Saakhr/Web-proj/pkg/models"
	v1 "github.com/Saakhr/Web-proj/pkg/v1/middlewares"
	"github.com/Saakhr/Web-proj/pkg/v1/services"
	"github.com/Saakhr/Web-proj/pkg/v1/utility"
	"github.com/Saakhr/Web-proj/templates"
	"github.com/gofiber/fiber/v2"
)

func AdminDashboard(c *fiber.Ctx, privateKey *rsa.PrivateKey) error {

	announcements, err := models.GetAllAnnouncements()
	if err != nil {
		return v1.NotFoundMiddleware(c)
	}
  wishlists,err:=models.GetAllWishlists()
	if err != nil {
		return v1.NotFoundMiddleware(c)
	}

  projects,err:=models.GetProjects()
	if err != nil {
		return v1.NotFoundMiddleware(c)
	}

  students,err:=models.GetStudents()
	if err != nil {
		return v1.NotFoundMiddleware(c)
	}
	user, err := services.GetUserFromCookie(c, privateKey)
	if err != nil {
		return v1.NotFoundMiddleware(c)
	}

	return utility.Render(c, templates.AdminDashBoard(&templates.DashBoardDataAdmin{
    Announcements:announcements,
    Projects: projects,
    Wishlists: wishlists,
    Students: students,
  } , user))
}
