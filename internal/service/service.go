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
	GetTasks(ctx *fiber.Ctx) error
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

// GetTask - обработчик запросов на получение задач
func (s *service) GetTasks(c *fiber.Ctx) error {
	tasks, err := s.repo.GetTasks(c.Context())
	if err != nil {
		s.log.Error("failed to get tasks: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get tasks",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   tasks,
	})
}

func (s *service) GetTaskByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	task, err := s.repo.GetTaskByID(c.Context(), id)
	if err != nil {
		s.log.Error("failed to get task by id:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get task",
		})
	}
	if task == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "task not found",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   task,
	})
}
