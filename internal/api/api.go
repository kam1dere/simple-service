package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"simple-service/internal/api/middleware"
	"simple-service/internal/service"
)

// Routers - структура для хранения зависимостей роутов
type Routers struct {
	Service service.Service
}

// NewRouters - конструктор для настройки API
func NewRouters(r *Routers, token string) *fiber.App {
	app := fiber.New()

	// Настройка CORS (разрешенные методы, заголовки, авторизация)
	app.Use(cors.New(cors.Config{
		AllowMethods:  "GET, POST, PUT, DELETE",
		AllowHeaders:  "Accept, Authorization, Content-Type, X-CSRF-Token, X-REQUEST-ID",
		ExposeHeaders: "Link",
		MaxAge:        300,
	}))

	// Группа маршрутов с авторизацией
	apiGroup := app.Group("/v1", middleware.Authorization(token))

	// Роут для создания задачи
	apiGroup.Post("/create_task", r.Service.CreateTask)
	// Роут для получения задачи по ID
	apiGroup.Get("/tasks/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid ID format",
			})
		}

		task, err := r.Service.GetTaskByID(c.Context(), id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if task == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Task not found",
			})
		}

		return c.JSON(task)
	})
	return app
}
