package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zesilva15/report-api/database"
	"github.com/zesilva15/report-api/models"
)

type Artifact struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func createResponseArtifact(artifactModel models.Artifact) Artifact {
	return Artifact{ID: artifactModel.ID, Name: artifactModel.Name}
}

func createArtifact(c *fiber.Ctx) error {
	var artifact models.Artifact

	if err := c.BodyParser(&artifact); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	database.Database.Db.Create(&artifact)
	response := createResponseArtifact(artifact)
	return c.Status(200).JSON(response)
}
