package routes

import (
	"encoding/base64"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zesilva15/report-api/database"
	"github.com/zesilva15/report-api/models"
)

type Report struct {
	ID        uint     `json:"id"`
	Artifact  Artifact `json:"artifact"`
	CreatedAt time.Time
	Type      string `json:"type"`
	Status    string `json:"status"`
	File      string `json:"file"`
}

func CreateResponseReport(artifact Artifact, report models.Report) Report {
	return Report{ID: report.ID, Artifact: artifact, Type: report.Type, Status: report.Status, File: report.File, CreatedAt: report.CreatedAt}
}

func validateBody(report models.Report) error {
	if report.Type == "" || report.Status == "" || report.File == "" {
		return fiber.ErrBadRequest
	}
	return nil
}

func CreateReport(c *fiber.Ctx) error {
	var report models.Report
	if err := c.BodyParser(&report); err != nil {
		return c.Status(400).JSON(fiber.ErrBadRequest)
	}
	if validateBody(report) != nil {
		return c.Status(400).JSON(fiber.ErrBadRequest)
	}
	var artifact models.Artifact
	if err := findArtifact(report.ArtifactRefer, &artifact); err != nil {
		return c.Status(404).JSON(fiber.ErrNotFound)
	}
	if err := database.Database.Db.Create(&report).Error; err != nil {
		return c.Status(500).JSON(fiber.ErrInternalServerError)
	}
	report.File = base64toFile(report.File)
	responseArtifact := CreateResponseArtifact(artifact)
	response := CreateResponseReport(responseArtifact, report)
	return c.Status(201).JSON(response)
}

func base64toFile(file string) string {
	decodedData, err := base64.StdEncoding.DecodeString(file)
	if err != nil {
		return ""
	}
	return string(decodedData)
}
