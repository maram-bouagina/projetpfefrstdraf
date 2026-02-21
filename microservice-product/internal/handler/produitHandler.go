package handlers

import (
	"projet/internal/dto"
	services "projet/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

type ProduitHandler struct {
	service *services.ProduitService
}

func NewProduitHandler(service *services.ProduitService) *ProduitHandler {
	return &ProduitHandler{service: service}
}

func (h *ProduitHandler) getBoutiqueID(c *fiber.Ctx) (string, error) {
	boutiqueID := c.Get("X-Boutique-ID")
	if boutiqueID == "" {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Missing store context")
	}
	return boutiqueID, nil
}

func (h *ProduitHandler) CreateProduit(c *fiber.Ctx) error {
	var req dto.RequeteCreationProduit
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	boutiqueID, err := h.getBoutiqueID(c)
	if err != nil {
		return err
	}

	produit, err := h.service.Create(c.Context(), boutiqueID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(produit)
}

func (h *ProduitHandler) ListProduits(c *fiber.Ctx) error {
	boutiqueID, err := h.getBoutiqueID(c)
	if err != nil {
		return err
	}

	produits, err := h.service.List(c.Context(), boutiqueID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch products"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"produits": produits})
}

func (h *ProduitHandler) GetProduitByID(c *fiber.Ctx) error {
	id := c.Params("id")
	boutiqueID, err := h.getBoutiqueID(c)
	if err != nil {
		return err
	}

	produit, err := h.service.GetByID(c.Context(), id, boutiqueID)
	if err != nil {
		if err.Error() == "product not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch product"})
	}
	return c.Status(fiber.StatusOK).JSON(produit)
}

func (h *ProduitHandler) UpdateProduit(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.RequeteUpdateProduit
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	boutiqueID, err := h.getBoutiqueID(c)
	if err != nil {
		return err
	}

	produit, err := h.service.Update(c.Context(), id, boutiqueID, req)
	if err != nil {
		if err.Error() == "product not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(produit)
}

func (h *ProduitHandler) DeleteProduit(c *fiber.Ctx) error {
	id := c.Params("id")
	boutiqueID, err := h.getBoutiqueID(c)
	if err != nil {
		return err
	}

	err = h.service.Delete(c.Context(), id, boutiqueID)
	if err != nil {
		if err.Error() == "product not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"ok": true})
}

func (h *ProduitHandler) SearchProduits(c *fiber.Ctx) error {
	var filter dto.FiltreProduit
	if err := c.QueryParser(&filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid filter parameters"})
	}

	boutiqueID, err := h.getBoutiqueID(c)
	if err != nil {
		return err
	}

	produits, page, limite, err := h.service.Search(c.Context(), boutiqueID, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to search products"})
	}
	return c.JSON(fiber.Map{
		"produits": produits,
		"page":     page,
		"limite":   limite,
	})
}
