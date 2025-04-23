package v1

import (
	"crypto/rsa"

	"github.com/Saakhr/Web-proj/pkg/v1/handlers"
	v1middlewares "github.com/Saakhr/Web-proj/pkg/v1/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SetupStudentRoutes(app *fiber.App, privateKey *rsa.PrivateKey) {
	student := app.Group("/student")
	
	// Public routes
	// admin.Get("/login", func(c *fiber.Ctx) error {
	// 	return utility.Render(c, templates.Login("","admin"))
	// })
	// admin.Post("/login",func(c *fiber.Ctx) error { 
	//    return handlers.AdminLogin(c,privateKey)
	//  })

	// Protected routes
	student.Use(v1middlewares.NewAuthMiddleware(privateKey,"student"))
	{ 
    student.Get("/dashboard", func(c *fiber.Ctx)error{
      return handlers.StudentDashboard(c,privateKey)
    })
    student.Delete("/wishlist", func(c *fiber.Ctx)error{
      return handlers.DeleteStudentWishlistItem(c,privateKey)
    })

    student.Post("/wishlist", func(c *fiber.Ctx)error{
      return handlers.CreateStudentWishlistItem(c,privateKey)
    })

		// admin.Get("/announcements", handlers.AnnouncementList)
		// admin.Post("/announcements", handlers.CreateAnnouncement)
	}
}
