package database

import (
	"log"
	"os"

	"github.com/zesilva15/report-api/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBInstance struct {
	Db *gorm.DB
}

var Database DBInstance

func Connect() {
	database, err := gorm.Open(mysql.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database \n", err.Error())
		os.Exit(2)
	}
	log.Println("Database connected")
	database.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Database running migrations")
	database.AutoMigrate(&models.Artifact{}, &models.Report{})
	Database = DBInstance{Db: database}
}
