package models

import "gorm.io/gorm"

// Todo represents a task with a title and a completion status
type Todo struct {
	gorm.Model
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}
