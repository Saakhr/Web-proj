package v1

import (
	"crypto/rsa"

	"github.com/Saakhr/Web-proj/pkg/v1/handlers"
	v1middlewares "github.com/Saakhr/Web-proj/pkg/v1/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SetupAdminRoutes(app *fiber.App, privateKey *rsa.PrivateKey) {
	admin := app.Group("/admin")
	
	// Public routes
	// admin.Get("/login", func(c *fiber.Ctx) error {
	// 	return utility.Render(c, templates.Login("","admin"))
	// })
	// admin.Post("/login",func(c *fiber.Ctx) error { 
	//    return handlers.AdminLogin(c,privateKey)
	//  })

	// Protected routes
	admin.Use(v1middlewares.NewAuthMiddleware(privateKey,"admin"))
	{
		 admin.Get("/dashboard", func(c *fiber.Ctx)error{
      return handlers.AdminDashboard(c,privateKey)
    })
    admin.Post("/announcement",handlers.CreateAnnouncement)
    admin.Post("/project",handlers.CreateProject)
    admin.Post("/student",handlers.CreateStudent)

    admin.Delete("/student", handlers.DeleteStudent)
    admin.Delete("/wishlist", handlers.DeleteWishlistItem)
    admin.Delete("/project", handlers.DeleteProject)
    admin.Delete("/announcement", handlers.DeleteAnnouncement)
		// admin.Get("/announcements", handlers.AnnouncementList)
		// admin.Post("/announcements", handlers.CreateAnnouncement)
	}
}
