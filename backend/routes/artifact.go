package routes

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/zesilva15/report-api/database"
	"github.com/zesilva15/report-api/models"
)

type Artifact struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func CreateResponseArtifact(artifactModel models.Artifact) Artifact {
	return Artifact{ID: artifactModel.ID, Name: artifactModel.Name}
}

func findArtifact(id int, artifact *models.Artifact) error {
	database.Database.Db.Find(&artifact, "id = ?", id)
	if artifact.ID == 0 {
		return errors.New("Artifact not found")
	}
	return nil
}

func CreateArtifact(c *fiber.Ctx) error {
	var artifact models.Artifact

	if err := c.BodyParser(&artifact); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	database.Database.Db.Create(&artifact)
	response := CreateResponseArtifact(artifact)
	return c.Status(201).JSON(response)
}

func GetArtifacts(c *fiber.Ctx) error {
	artifacts := []models.Artifact{}

	database.Database.Db.Find(&artifacts)
	responseArtifacts := []Artifact{}
	for _, artifact := range artifacts {
		responseArtifact := CreateResponseArtifact(artifact)
		responseArtifacts = append(responseArtifacts, responseArtifact)
	}
	return c.Status(200).JSON(responseArtifacts)
}

func GetArtifact(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	var artifact models.Artifact

	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	if err := findArtifact(id, &artifact); err != nil {
		return c.Status(404).JSON(err.Error())
	}
	response := CreateResponseArtifact(artifact)
	return c.Status(200).JSON(response)
}

func UpdateArtifact(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	var artifact models.Artifact

	if err != nil {
		return c.Status(400).JSON(err.Error())
	}
	if err := findArtifact(id, &artifact); err != nil {
		return c.Status(404).JSON(err.Error())
	}
	type UpdateArtifact struct {
		Name string `json:"name"`
	}
	var updateData UpdateArtifact
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(err.Error())
	}
	artifact.Name = updateData.Name
	database.Database.Db.Save(&artifact)
	response := CreateResponseArtifact(artifact)
	return c.Status(200).JSON(response)
}
