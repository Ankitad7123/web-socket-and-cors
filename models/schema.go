package models

import "time"

type Users struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Username  string    `json:"username" gorm:"unique;not null"`
    Email     string    `json:"email" gorm:"unique;not null"`
    Password  string    `json:"password" gorm:"not null"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

