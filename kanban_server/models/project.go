package models

import (
	"time"
)

type Project struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	Status      string    `json:"status" gorm:"not null;default:'active'"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Users       []User    `json:"users,omitempty" gorm:"many2many:user_projects;"`
	Tasks       []Task    `json:"tasks,omitempty"`
}

type Task struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	Status      string    `json:"status" gorm:"not null;default:'todo'"`     // todo, in_progress, done
	Priority    string    `json:"priority" gorm:"not null;default:'medium'"` // low, medium, high
	DueDate     time.Time `json:"due_date"`
	ProjectID   uint      `json:"project_id"`
	Project     Project   `json:"-"`
	AssignedTo  uint      `json:"assigned_to"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
