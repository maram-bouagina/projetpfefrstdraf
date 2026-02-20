package routes

import (
	handlers "projet/internal/handler"
	"projet/internal/repository"
	services "projet/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterProduitRoutes(app *fiber.App, db *gorm.DB) {
	repo := repository.NewRepo(db)
	service := services.NewProduitService(repo)
	handler := handlers.NewProduitHandler(service)

	produits := app.Group("/produits")
	produits.Post("/", handler.CreateProduit)
	produits.Get("/", handler.ListProduits)
	produits.Get("/search", handler.SearchProduits)
	produits.Get("/:id", handler.GetProduitByID)
	produits.Put("/:id", handler.UpdateProduit)
	produits.Delete("/:id", handler.DeleteProduit)
}
