package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/zesilva15/report-api/database"
	"github.com/zesilva15/report-api/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	database.Connect()
	app := fiber.New()
	setupRoutes(app)
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}

func welcome(c *fiber.Ctx) error {
	return c.SendString("Welcome to Fiber!")
}

func setupRoutes(app *fiber.App) {
	app.Get("/", welcome)
	// Artifact Endpoints
	app.Post("/artifacts", routes.CreateArtifact)
	app.Get("/artifacts", routes.GetArtifacts)
	app.Get("/artifacts/:id", routes.GetArtifact)
	app.Patch("/artifacts/:id", routes.UpdateArtifact)
}
