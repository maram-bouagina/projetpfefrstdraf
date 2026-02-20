package routes

import (
	"projet/internal/routes"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *fiber.App {
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"ok":     true,
			"status": "healthy",
		})
	})

	routes.RegisterProduitRoutes(app, db)

	return app
}
