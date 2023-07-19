package routes

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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
	report.File = base64toFile(report.File, artifact.Name, report.Type)
	if report.File == "" {
		return c.Status(500).JSON(fiber.ErrInternalServerError)
	}
	if err := database.Database.Db.Create(&report).Error; err != nil {
		return c.Status(500).JSON(fiber.ErrInternalServerError)
	}
	responseArtifact := CreateResponseArtifact(artifact)
	response := CreateResponseReport(responseArtifact, report)
	return c.Status(201).JSON(response)
}

func base64toFile(fileContent string, artifactName string, reportType string) string {
	decodedData, err := base64.StdEncoding.DecodeString(fileContent)
	if err != nil {
		return ""
	}
	filename := time.Now().Format("2006-01-02T15:04:05Z") + ".json"
	file, err := os.Create(filename)
	if err != nil {
		return ""
	}
	fmt.Fprintf(file, string(decodedData))
	reportFile := uploadToMinio(filename, artifactName, reportType)
	return string(reportFile)
}

func uploadToMinio(filename string, artifactName string, reportType string) string {
	ctx := context.Background()
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	useSSL := false
	bucketName := os.Getenv("MINIO_BUCKET")

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
		return ""
	}
	destinationFilename := artifactName + "-" + reportType + "-" + filename
	info, err := minioClient.FPutObject(ctx, bucketName, destinationFilename, filename, minio.PutObjectOptions{ContentType: "application/json"})
	if err != nil {
		log.Fatalln(err)
		return ""
	}
	log.Printf("Successfully uploaded %s of size %d\n", destinationFilename, info.Size)
	err = os.Remove(filename)
	if err != nil {
		log.Fatalln(err)
		return ""
	}
	return destinationFilename
}
