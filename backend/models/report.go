package models

import "time"

type Report struct {
	ID            uint `json:"id" gorm:"primaryKey"`
	CreatedAt     time.Time
	Type          string   `json:"type"`
	Status        string   `json:"status"`
	File          string   `json:"file"`
	ArtifactRefer uint     `json:"artifact_id"`
	Artifact      Artifact `gorm:"foreignKey:ArtifactRefer"`
}
