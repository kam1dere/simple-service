package repo

import "time"

// Task - структура, соответствующая таблице tasks
type Task struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type GetTask struct {
	ID          int
	Title       string
	Description string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
