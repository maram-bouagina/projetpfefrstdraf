package routes

import (
	handlers "projet/internal/handler"
	"projet/internal/repository"
	services "projet/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterVarianteRoutes(app *fiber.App, db *gorm.DB) {
	// Repositories
	varianteRepo := repository.NewVarianteRepo(db)
	produitRepo := repository.NewRepo(db)

	// Services
	produitService := services.NewProduitService(produitRepo)
	varianteService := services.NewVarianteService(varianteRepo)

	// Handlers
	varianteHandler := handlers.NewVarianteHandler(varianteService, produitService)

	variantes := app.Group("/produits/:produitId/variantes")
	variantes.Post("/", varianteHandler.CreateVariante)
	variantes.Get("/", varianteHandler.ListVariantes)

	variante := app.Group("/variantes/:varianteId")
	variante.Get("/", varianteHandler.GetVarianteByID)
	variante.Put("/", varianteHandler.UpdateVariante)
	variante.Delete("/", varianteHandler.DeleteVariante)
}
