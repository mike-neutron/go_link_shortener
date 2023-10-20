package models

import "time"

type Link struct {
	ID        uint      `gorm:"primaryKey" json:"-"`
	Original  string    `gorm:"varchar(256);not null" json:"original"`
	CreatedAt time.Time `json:"created_at" sql:"DEFAULT:CURRENT_TIMESTAMP"`
}
