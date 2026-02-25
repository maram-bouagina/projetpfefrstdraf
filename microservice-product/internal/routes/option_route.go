package routes

import (
	handlers "projet/internal/handler"
	"projet/internal/repository"
	services "projet/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterOptionRoutes(app *fiber.App, db *gorm.DB) {
	// Repositories
	optionRepo := repository.NewOptionProduitValeurRepo(db)

	// Services
	optionService := services.NewOptionProduitService(optionRepo)

	// Handlers
	optionHandler := handlers.NewOptionProduitHandler(optionService)

	options := app.Group("/produits/:produitId/options")
	options.Post("/", optionHandler.CreateOption)
	options.Get("/", optionHandler.ListOptions)

	option := app.Group("/options/:optionId")
	option.Get("/", optionHandler.GetOptionByID)
	option.Put("/", optionHandler.UpdateOption)
	option.Delete("/", optionHandler.DeleteOption)

	valeurs := app.Group("/options/:optionId/valeurs")
	valeurs.Post("/", optionHandler.CreateValeur)
	valeurs.Get("/", optionHandler.ListValeurs)

	valeur := app.Group("/valeurs/:valeurId")
	valeur.Put("/", optionHandler.UpdateValeur)
	valeur.Delete("/", optionHandler.DeleteValeur)
}
