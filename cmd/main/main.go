package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/joho/godotenv"
	"log"
	"os"

	v1middlewares "github.com/Saakhr/Web-proj/pkg/v1/middlewares"
	v1routes "github.com/Saakhr/Web-proj/pkg/v1/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var (
	privateKey *rsa.PrivateKey
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app := fiber.New()


  app.Static("/static", "./static")
  // app.Static("/css", "./css")

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
