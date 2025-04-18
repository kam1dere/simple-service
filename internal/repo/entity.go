package repo

import "time"

// Task - структура, соответствующая таблице tasks
type Task struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// FullTask - структура, полностью описывающая task
type FullTask struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
