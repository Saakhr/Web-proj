package handlers

import (
	"crypto/rsa"
	"time"

	"github.com/Saakhr/Web-proj/pkg/models"
	"github.com/Saakhr/Web-proj/pkg/v1/services"
	"github.com/Saakhr/Web-proj/pkg/v1/utility"
	"github.com/Saakhr/Web-proj/templates"
	"github.com/gofiber/fiber/v2"
)

func AnnouncementList(c *fiber.Ctx,privateKey *rsa.PrivateKey,dept string) error {
	// Get department from query or user context
	
	announcements, err := models.GetAnnouncementsByDepartment(dept)
	if err != nil {
    return c.Redirect("/")
  }

  user, err := services.GetUserFromCookie(c, privateKey)
  if err!=nil{
    return utility.Render(c,templates.AnnouncementList(announcements,nil))
  }

	return utility.Render(c,templates.AnnouncementList(announcements,user))
}
func ListOfAnnouncements(c *fiber.Ctx,dept string)error {
	announcements, err := models.GetAnnouncementsByDepartment(dept)
	if err != nil {
    return c.Redirect("/")
	}
  return utility.Render(c, templates.ListOfAnnouncements(announcements))
}

func CreateAnnouncement(c *fiber.Ctx) error {
	var announcement models.Announcement
	if err := c.BodyParser(&announcement); err != nil {
    return c.Redirect("/")
	}

	announcement.DateTime = time.Now()
	if err := models.CreateAnnouncement(&announcement); err != nil {
    return c.Redirect("/")
	}

	return c.Redirect("/v1/admin/dashboard")
}

func DeleteAnnouncement(c *fiber.Ctx) error {
  id:=c.QueryInt("id")
  if err:= models.DeleteAnnouncement(id); err!=nil{
    return err
  }
  c.Set("HX-Refresh", "true")
  return c.SendStatus(302)
}
