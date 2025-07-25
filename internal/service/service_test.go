package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"simple-service/internal/dto"
	"simple-service/internal/repo"
	"simple-service/internal/repo/mocks"
)

// TestCreateTask - тестирование метода CreateTask
func TestCreateTask(t *testing.T) {
	// Создаем мок репозитория
	mockRepo := new(mocks.Repository)
	logger := zap.NewNop().Sugar() // Без вывода логов

	// Создаем экземпляр сервиса с мок-репозиторием
	s := NewService(mockRepo, logger)

	// Инициализируем Fiber-контекст
	app := fiber.New()
	app.Post("/tasks", s.CreateTask)

	t.Run("успешное создание задачи", func(t *testing.T) {
		task := TaskRequest{
			Title:       "Test Task",
			Description: "Test Description",
		}
		body, _ := json.Marshal(task)

		// Ожидаем, что вызов метода `CreateTask` в репозитории вернёт ID = 1
		mockRepo.On("CreateTask", mock.Anything, repo.Task{
			Title:       task.Title,
			Description: task.Description,
		}).Return(1, nil).Once()

		// Отправляем запрос
		req, err := http.NewRequest("POST", "/tasks", bytes.NewReader(body))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		// Выполняем запрос
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		// Проверяем ответ
		var response dto.Response
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "success", response.Status)

		// Проверяем вызов мок-методов
		mockRepo.AssertExpectations(t)
	})

	t.Run("ошибка валидации входных данных", func(t *testing.T) {
		body := []byte(`{}`) // Пустое тело, `title` обязателен

		req, err := http.NewRequest("POST", "/tasks", bytes.NewReader(body))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response dto.Response
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "error", response.Status)
	})

	t.Run("ошибка при создании задачи в БД", func(t *testing.T) {
		task := TaskRequest{
			Title:       "Test Task",
			Description: "Test Description",
		}
		body, _ := json.Marshal(task)

		// Ожидаем ошибку при вставке в БД
		mockRepo.On("CreateTask", mock.Anything, repo.Task{
			Title:       task.Title,
			Description: task.Description,
		}).Return(0, errors.New("DB error")).Once()

		req, err := http.NewRequest("POST", "/tasks", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		var response dto.Response
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "error", response.Status)

		mockRepo.AssertExpectations(t)
	})
}

// TestGetTaskById - тестирование метода GetTaskById
func TestGetTaskById(t *testing.T) {
	// Создаем мок репозитория
	mockRepo := new(mocks.Repository)
	logger := zap.NewNop().Sugar() // Без вывода логов

	// Создаем экземпляр сервиса с мок-репозиторием
	s := NewService(mockRepo, logger)

	// Инициализируем Fiber-контекст
	app := fiber.New()
	app.Get("/task/:id", s.GetTaskById)

	t.Run("успешное получение задачи", func(t *testing.T) {
		var id int = 1
		expectedTask := repo.TaskInfo{
			Id:          id,
			Title:       "Test Task",
			Description: "Test Description",
			Status:      "new",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Ожидаем, что вызов метода `GetTaskById` в репозитории вернёт expectedTask
		mockRepo.On("GetTaskById", mock.Anything, id).Return(expectedTask, nil).Once()

		// Отправляем запрос
		req, err := http.NewRequest("GET", "/task/1", nil)
		assert.NoError(t, err)

		// Выполняем запрос
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		// Проверяем ответ
		var response dto.TaskResponse
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, expectedTask.Id, response.Id)
		assert.Equal(t, expectedTask.Title, response.Title)
		assert.Equal(t, expectedTask.Description, response.Description)
		assert.Equal(t, expectedTask.Status, response.Status)

		// Проверяем вызов мок-методов
		mockRepo.AssertExpectations(t)
	})

	t.Run("пустой id", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/task/", nil)
		assert.NoError(t, err)

		resp, err := app.Test(req)
		assert.NoError(t, err)

		// Ожидаем ошибку 404
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("невалидный id", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/task/abc", nil)
		assert.NoError(t, err)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var response dto.Response
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "error", response.Status)
	})

	t.Run("ошибка получения задачи из БД", func(t *testing.T) {
		mockRepo.On("GetTaskById", mock.Anything, 2).Return(repo.TaskInfo{}, errors.New("db error")).Once()

		req, err := http.NewRequest("GET", "/task/2", nil)
		assert.NoError(t, err)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		var response dto.Response
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "error", response.Status)

		mockRepo.AssertExpectations(t)
	})
}
