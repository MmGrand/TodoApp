package tasks_transport_http

import (
	"time"

	"github.com/MmGrand/TodoApp/internal/core/domain"
)

type TaskDTOResponse struct {
	ID           int        `json:"id" example:"15"`
	Version      int        `json:"version" example:"2"`
	Title        string     `json:"title" example:"Купить продукты"`
	Description  *string    `json:"description" example:"Молоко, хлеб, яйца"`
	Completed    bool       `json:"completed" example:"true"`
	CreatedAt    time.Time  `json:"created_at" example:"2026-09-01T10:00:00Z"`
	CompletedAt  *time.Time `json:"completed_at" example:"2026-09-02T18:30:00Z"`
	AuthorUserID int        `json:"author_user_id" example:"10"`
}

func taskDTOFromDomain(task domain.Task) TaskDTOResponse {
	return TaskDTOResponse{
		ID:           task.ID,
		Version:      task.Version,
		Title:        task.Title,
		Description:  task.Description,
		Completed:    task.Completed,
		CreatedAt:    task.CreatedAt,
		CompletedAt:  task.CompletedAt,
		AuthorUserID: task.AuthorUserID,
	}
}

func taskDTOsFromDomains(tasks []domain.Task) []TaskDTOResponse {
	dtos := make([]TaskDTOResponse, len(tasks))

	for i, task := range tasks {
		dtos[i] = taskDTOFromDomain(task)
	}

	return dtos
}
