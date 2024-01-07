package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

type Task struct {
	gorm.Model
	Name        string    `json:"name"`
	Deadline    time.Time `json:"deadline"`
	Description string    `json:"description"`
}
