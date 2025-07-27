package service

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"simple-service/internal/dto"
	"simple-service/internal/repo"
	"simple-service/pkg/validator"
	"strconv"
)

// Слой бизнес-логики. Тут должна быть основная логика сервиса

// Service - интерфейс для бизнес-логики
type Service interface {
	CreateTask(ctx *fiber.Ctx) error
	GetTaskByID(ctx *fiber.Ctx) error
}

type service struct {
	repo repo.Repository
	log  *zap.SugaredLogger
}

// NewService - конструктор сервиса
func NewService(repo repo.Repository, logger *zap.SugaredLogger) Service {
	return &service{
		repo: repo,
		log:  logger,
	}
}

// CreateTask - обработчик запроса на создание задачи
func (s *service) CreateTask(ctx *fiber.Ctx) error {
	var req TaskRequest

	// Десериализация JSON-запроса
	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		s.log.Error("Invalid request body", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid request body")
	}

	// Валидация входных данных
	if vErr := validator.Validate(ctx.Context(), req); vErr != nil {
		return dto.BadResponseError(ctx, dto.FieldIncorrect, vErr.Error())
	}

	// Вставка задачи в БД через репозиторий
	task := repo.Task{
		Title:       req.Title,
		Description: req.Description,
	}
	taskID, err := s.repo.CreateTask(ctx.Context(), task)
	if err != nil {
		s.log.Error("Failed to insert task", zap.Error(err))
		return dto.InternalServerError(ctx)
	}

	// Формирование ответа
	response := dto.Response{
		Status: "success",
		Data:   map[string]int{"task_id": taskID},
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (s *service) GetTaskByID(ctx *fiber.Ctx) error {
	// Достаем ID из URL и преобразуем
	idParam := ctx.Params("id")
	taskID, err := strconv.Atoi(idParam)
	if err != nil {
		s.log.Error("Invalid id parameter", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid task id")
	}

	// Вызываем метод в репозитории
	task, err := s.repo.GetTaskById(ctx.Context(), taskID)
	if err != nil {
		s.log.Error("Failed to get task", zap.Error(err))
		return dto.InternalServerError(ctx)
	}
	if task == (repo.Task{}) {
		return dto.NotFoundError(ctx, "task not found")
	}

	taskResponse := dto.TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
	response := dto.Response{
		Status: "success",
		Data:   taskResponse,
	}

	return ctx.Status(fiber.StatusOK).JSON(response)

}
