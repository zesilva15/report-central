package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app := fiber.New()
	app.Get("/", welcome)
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}

func welcome(c *fiber.Ctx) error {
	return c.SendString("Welcome to Fiber!")
}
