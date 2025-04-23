package handlers

import (
	"github.com/Saakhr/Web-proj/pkg/models"
	"github.com/gofiber/fiber/v2"
)

func CreateProject(c *fiber.Ctx) error {
	var project models.Projects
	if err := c.BodyParser(&project); err != nil {
    return c.SendString("ss")
	}

	if err := models.CreateProject(&project); err != nil {
    return c.SendString("ss")
	}

	return c.Redirect("/v1/admin/dashboard")
}

func DeleteProject(c *fiber.Ctx) error {
  id:=c.QueryInt("id")
  if err:= models.DeleteProject(id); err!=nil{
    return err
  }
  c.Set("HX-Refresh", "true")
  return c.SendStatus(302)
}
