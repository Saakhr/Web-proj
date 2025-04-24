package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"time"

	// "github.com/joho/godotenv"

	"github.com/Saakhr/Web-proj/pkg/database"
	"github.com/Saakhr/Web-proj/pkg/models"
	v1middlewares "github.com/Saakhr/Web-proj/pkg/v1/middlewares"
	v1routes "github.com/Saakhr/Web-proj/pkg/v1/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var (
	privateKey *rsa.PrivateKey
)

func main() {
	var err error
	// err = godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	err = database.InitDB()
	if err != nil {
		log.Fatal("Database initialization failed: ", err)
	}
	defer database.DB.Close()
  createInitialAdmin()

	app := fiber.New()

  app.Static("/static", "./static")

  // JWT initialization
  err = getKey()
  if err != nil {
    log.Fatal("Couldn't Load JWT RSA key" + err.Error())
  }

  // Middleware
  app.Use(logger.New())

  // Routes
  v1 := v1routes.GetRoutes(privateKey)
  app.Mount("/v1", v1)

  // Root redirect
  app.Get("/", func(c *fiber.Ctx) error {
    return c.Redirect("/v1")
  })

  // 404 Handler (make sure this comes after all other routes)
  app.Use(v1middlewares.NotFoundMiddleware)


	log.Fatal(app.Listen(":8080"))
}
func getKey() error {

	xs := os.Getenv("JWT_RS_SECRET")

	// Decode the PEM block
	block, _ := pem.Decode([]byte(xs))
	if block == nil {
		return errors.New("Error decoding PEM block")
	}

	// Parse the private key
	privateKeys, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return errors.New("Error parsing private key")
	}
	var ok bool
	privateKey, ok = privateKeys.(*rsa.PrivateKey)
	if !ok {
		return errors.New("Error: parsed key is not an RSA private key")
	}
	return nil
}
func createInitialAdmin() {
	// Check if admin exists
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM admins").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
  err=GenerateDummyData()
	if err != nil {
		log.Fatal(err)
	}

	if count == 0 {
		admin := &models.Admin{
			Email:    "admin@admin.com",
			Password: "123456",
			FullName: "System Administrator",
		}

		if err := admin.HashPassword(); err != nil {
			log.Fatal("Failed to hash admin password: ", err)
		}

		_, err := database.DB.Exec(
			"INSERT INTO admins (email, password, full_name) VALUES (?, ?, ?)",
			admin.Email, admin.Password, admin.FullName)
		if err != nil {
			log.Fatal("Failed to create initial admin: ", err)
		}
		log.Println("Initial admin user created")
	}
}

func GenerateDummyData() error {
	// Generate dummy students
	students := []models.Student{
    { FirstName:"John",
      LastName:"Doe",
      Email: "john.doe@school.edu",
      Password: "123"},
    {FirstName: "Jane",LastName:"Smith", Email: "jane.smith@school.edu",Password: "123"},
    {FirstName:"Alex", LastName: "Johnson", Email: "alex.johnson@school.edu", Password: "123"},
	}

	for _, s := range students {
    err:=s.HashPassword()
    if err!=nil{
			return fmt.Errorf("error hashing student password: %v", err)
    }
		_, err = database.DB.Exec(
			"INSERT INTO students (first_name, last_name, email, password) VALUES (?, ?, ?, ?)",
			s.FirstName, s.LastName, s.Email, s.Password)
		if err != nil {
			return fmt.Errorf("error inserting student: %v", err)
		}
	}

	// Generate dummy projects
	projects := []struct {
		title       string
		description string
	}{
		{
			"School Website Redesign",
			"Redesign the school website with modern UI/UX principles using HTML, CSS, and JavaScript",
		},
		{
			"Library Management System",
			"Develop a digital system to track books, borrowers, and due dates",
		},
		{
			"Science Fair Robot",
			"Build an autonomous robot for the annual science fair competition",
		},
		{
			"Math Tutoring App",
			"Create a mobile app to help students with math problems",
		},
	}

	for _, p := range projects {
		_, err := database.DB.Exec(
			"INSERT INTO projects (title, description) VALUES (?, ?)",
			p.title, p.description)
		if err != nil {
			return fmt.Errorf("error inserting project: %v", err)
		}
	}

	// Generate dummy announcements (2 per display type)
	displayTypes := []string{"general", "computer_science", "physics", "chemistry", "math"}
	now := time.Now()

	for _, display := range displayTypes {
		for i := 1; i <= 2; i++ {
			_, err := database.DB.Exec(
				"INSERT INTO announcements (title, content, display, datetime) VALUES (?, ?, ?, ?)",
				fmt.Sprintf("%s Announcement %d", display, i),
				fmt.Sprintf("This is a sample announcement for %s department. Details about upcoming events, deadlines, or important information would go here.", display),
				display,
				now.Add(-time.Duration(i)*24*time.Hour),
			)
			if err != nil {
				return fmt.Errorf("error inserting announcement: %v", err)
			}
		}
	}

	// Generate dummy wishlist entries
	var studentIDs []int
	rows, err := database.DB.Query("SELECT id FROM students")
	if err != nil {
		return fmt.Errorf("error fetching student IDs: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("error scanning student ID: %v", err)
		}
		studentIDs = append(studentIDs, id)
	}

	var projectIDs []int
	rows, err = database.DB.Query("SELECT id FROM projects")
	if err != nil {
		return fmt.Errorf("error fetching project IDs: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("error scanning project ID: %v", err)
		}
		projectIDs = append(projectIDs, id)
	}

	// Each student gets 2 random projects in their wishlist
	for _, studentID := range studentIDs {
		// Shuffle project IDs
		rand.Shuffle(len(projectIDs), func(i, j int) {
			projectIDs[i], projectIDs[j] = projectIDs[j], projectIDs[i]
		})

		for i := 0; i < 2 && i < len(projectIDs); i++ {
			_, err := database.DB.Exec(
				"INSERT INTO student_project_wishlist (student_id, project_id) VALUES (?, ?)",
				studentID, projectIDs[i])
			if err != nil {
				return fmt.Errorf("error inserting wishlist: %v", err)
			}
		}
	}

	return nil
}
