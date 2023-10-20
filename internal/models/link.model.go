package models

import "time"

type Link struct {
	Short     string    `gorm:"type:varchar(256);not null" json:"short"`
	Original  string    `gorm:"type:varchar(1000);not null" json:"original"`
	CreatedAt time.Time `json:"created_at" sql:"DEFAULT:CURRENT_TIMESTAMP"`
}
