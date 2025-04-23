package handlers

import (
	"crypto/rsa"
	"strconv"

	"github.com/Saakhr/Web-proj/pkg/models"
	"github.com/Saakhr/Web-proj/pkg/v1/services"
	"github.com/Saakhr/Web-proj/pkg/v1/utility"
	"github.com/Saakhr/Web-proj/templates"
	"github.com/gofiber/fiber/v2"
)

func CreateStudent(c *fiber.Ctx) error {
	var project models.Student
	if err := c.BodyParser(&project); err != nil {
		return c.SendString("ss")
	}
	project.FirstName = c.FormValue("first_name")
	project.LastName = c.FormValue("last_name")
	if err := project.CreateStudent(); err != nil {
		return c.SendString("ss")
	}
	return c.Redirect("/v1/admin/dashboard")
}

func DeleteWishlistItem(c *fiber.Ctx) error {
	id := c.QueryInt("id")
	if err := models.DeleteWishlistItem(id); err != nil {
		return err
	}
	c.Set("HX-Refresh", "true")
	return c.SendStatus(302)
}

func CreateStudentWishlistItem(c *fiber.Ctx, privateKey *rsa.PrivateKey) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid form data")
	}

	// Get selected project IDs (checkboxes)
	selectedProjects := form.Value["selectedProjects"] // []string of project IDs

	// Convert string IDs to integers
	var projectIDs []int
	for _, idStr := range selectedProjects {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			continue // Skip invalid IDs
		}

		projectIDs = append(projectIDs, id)
	}

	user, err := services.GetUserFromCookie(c, privateKey)
	if err != nil {
    return c.Redirect("/v1")
	}
	for _, id := range projectIDs {
		models.CreateWish(user.UserID, id)
	}

	c.Set("HX-Refresh", "true")
	return c.SendStatus(302)
}
func DeleteStudentWishlistItem(c *fiber.Ctx, privateKey *rsa.PrivateKey) error {
	id := c.QueryInt("id")
	user, err := services.GetUserFromCookie(c, privateKey)
	if err != nil {
    return c.Redirect("/v1")
	}
	if err := models.DeleteStudentWishlistItem(id, user.UserID); err != nil {
		return err
	}
	c.Set("HX-Refresh", "true")
	return c.SendStatus(302)
}

func DeleteStudent(c *fiber.Ctx) error {
	id := c.QueryInt("id")
	if err := models.DeleteStudent(id); err != nil {
		return err
	}
	c.Set("HX-Refresh", "true")
	return c.SendStatus(302)
}

func StudentDashboard(c *fiber.Ctx, privateKey *rsa.PrivateKey) error {


	user, err := services.GetUserFromCookie(c, privateKey)
	if err != nil {

    return c.Redirect("/v1")
	}

	wishlists, err := models.GetStudentWishlists(user.UserID)
	if err != nil {
    return c.Redirect("/v1")
	}
	projects, err := models.GetStudentProjectsRec(user.UserID)

	return utility.Render(c, templates.StudentDashBoard(&templates.DashBoardDataStudent{
		MyWishlists: wishlists,
		Projects:    projects,
	}, user))
}
