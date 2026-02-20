package handlers

import (
	"projet/internal/dto"
	services "projet/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ProduitHandler struct {
	service *services.ProduitService
}

func NewProduitHandler(service *services.ProduitService) *ProduitHandler {
	return &ProduitHandler{service: service}
}

func (h *ProduitHandler) CreateProduit(c *fiber.Ctx) error {
	var req dto.RequeteCreationProduit

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	produit, err := h.service.Create(c.Context(), req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(produit)
}

func (h *ProduitHandler) ListProduits(c *fiber.Ctx) error {
	produits, err := h.service.Repo.ListProduits(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch products"})
	}

	return c.Status(200).JSON(fiber.Map{"produits": produits})
}

func (h *ProduitHandler) GetProduitByID(c *fiber.Ctx) error {
	id := c.Params("id")

	produit, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		if err.Error() == "product not found" {
			return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch product"})
	}

	return c.Status(200).JSON(produit)
}

func (h *ProduitHandler) UpdateProduit(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.RequeteUpdateProduit

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	produit, err := h.service.Update(c.Context(), id, req)
	if err != nil {
		if err.Error() == "product not found" {
			return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(produit)
}

func (h *ProduitHandler) DeleteProduit(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.service.Delete(c.Context(), id)
	if err != nil {
		if err.Error() == "product not found" {
			return c.Status(404).JSON(fiber.Map{"error": "Product not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"ok": true})
}
func (h *ProduitHandler) SearchProduits(c *fiber.Ctx) error {
	var filter dto.FiltreProduit

	if err := c.QueryParser(&filter); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid filter parameters"})
	}

	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.Limite == 0 {
		filter.Limite = 20
	}

	filterMap := make(map[string]interface{})
	if filter.Statut != nil {
		filterMap["statut"] = *filter.Statut
	}
	if filter.Visibilite != nil {
		filterMap["visibilite"] = *filter.Visibilite
	}
	if filter.Marque != nil {
		filterMap["marque"] = *filter.Marque
	}
	if filter.Recherche != nil && *filter.Recherche != "" {
		filterMap["titre ILIKE ?"] = "%" + *filter.Recherche + "%"
	}
	if !filter.InclureSupprime {
		filterMap["supprime_le"] = nil
	}

	// Add pagination (adjust key names to match your repository)
	filterMap["limit"] = filter.Limite
	filterMap["offset"] = (filter.Page - 1) * filter.Limite

	produits, err := h.service.Repo.GetWithFilter(c.Context(), filterMap)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to search products"})
	}

	return c.Status(200).JSON(fiber.Map{
		"produits": produits,
		"page":     filter.Page,
		"limite":   filter.Limite,
	})
}
