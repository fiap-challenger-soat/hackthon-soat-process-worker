package domain

import (
	"os"
	"time"
)

type VideoStatus string

const (
	VideoStatusProcessing VideoStatus = "processing"
	VideoStatusCompleted  VideoStatus = "completed"
	VideoStatusFailed     VideoStatus = "failed"
)

type VideoJobDTO struct {
	ID         string      `gorm:"primaryKey;type:uuid;" json:"job_id"`
	Status     VideoStatus `gorm:"type:varchar(20);not null;" json:"status"`
	CreatedAt  string      `gorm:"type:timestamp;not null;" json:"created_at"`
	OutputPath *string     `gorm:"type:varchar(255);" json:"output_path"`
	UserID     string      `gorm:"not null;" json:"user_id"`
	Email      string      `gorm:"type:varchar(255);not null;" json:"email"`
	VideoPath  string      `gorm:"type:varchar(255);not null;" json:"video_path"`
}

type VideoJob struct {
	ID         string      `gorm:"primaryKey;type:uuid;" json:"job_id"`
	Status     VideoStatus `gorm:"type:varchar(20);not null;" json:"status"`
	CreatedAt  string      `gorm:"type:timestamp;not null;" json:"created_at"`
	OutputPath *string     `gorm:"type:varchar(255);" json:"output_path"`
	UserID     string      `gorm:"not null;" json:"user_id"`
	VideoPath  string      `gorm:"type:varchar(255);not null;" json:"video_path"`
}

type DownloadedFile struct {
	Path string
	File *os.File
}

type JobStatusHistory struct {
	ID        string      `gorm:"primaryKey;type:uuid;" json:"id"`
	JobID     string      `gorm:"type:uuid;not null;" json:"job_id"`
	Status    VideoStatus `gorm:"type:varchar(20);not null;" json:"status"`
	CreatedAt time.Time   `json:"created_at" gorm:"type:timestamptz;default:now()"`
}

func (VideoJob) TableName() string {
	return "tb_video_jobs"
}
