package handlers

import (
	"crypto/rsa"

	"github.com/Saakhr/Web-proj/pkg/models"
	"github.com/Saakhr/Web-proj/pkg/v1/services"
	"github.com/Saakhr/Web-proj/pkg/v1/utility"
	"github.com/Saakhr/Web-proj/templates"
	"github.com/gofiber/fiber/v2"
)

func AdminDashboard(c *fiber.Ctx, privateKey *rsa.PrivateKey) error {

	announcements, err := models.GetAllAnnouncements()
	if err != nil {
    return c.Redirect("/v1")
	}
  wishlists,err:=models.GetAllWishlists()
	if err != nil {
    return c.Redirect("/v1")
	}

  projects,err:=models.GetProjects()
	if err != nil {
    return c.Redirect("/v1")
	}

  students,err:=models.GetStudents()
	if err != nil {
    return c.Redirect("/v1")
	}
	user, err := services.GetUserFromCookie(c, privateKey)
	if err != nil {
    return c.Redirect("/v1")
	}

	return utility.Render(c, templates.AdminDashBoard(&templates.DashBoardDataAdmin{
    Announcements:announcements,
    Projects: projects,
    Wishlists: wishlists,
    Students: students,
  } , user))
}
