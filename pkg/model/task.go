package model

import "time"

const (
	default_status = "todo"
)

type Task struct {
	Id        int        `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	Status    string     `db:"status" json:"status"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
}

func NewTask(n string) *Task {
	now := time.Now() // Сначала сохраняем в переменную
	return &Task{
		Name:      n,
		Status:    default_status,
		CreatedAt: &now,
	}
}
