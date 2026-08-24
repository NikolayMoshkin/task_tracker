package model

import "time"

const (
	default_status = "todo"
)

type Task struct {
	Id        int
	Name      string
	Status    string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func NewTask(n string) *Task {
	now := time.Now() // Сначала сохраняем в переменную
	return &Task{
		Name:      n,
		Status:    default_status,
		CreatedAt: &now,
	}
}
