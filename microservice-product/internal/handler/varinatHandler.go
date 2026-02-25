package handler

import (
	"projet/internal/dto"
	"projet/internal/service"

	"github.com/gofiber/fiber/v2"
)

type VarianteHandler struct {
	service        *service.VarianteService
	produitService *service.ProduitService
}

func NewVarianteHandler(
	service *service.VarianteService,
	produitService *service.ProduitService,
) *VarianteHandler {
	return &VarianteHandler{
		service:        service,
		produitService: produitService,
	}
}

func (h *VarianteHandler) getBoutiqueID(c *fiber.Ctx) (string, error) {
	boutiqueID := c.Get("X-Boutique-ID")
	if boutiqueID == "" {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Missing store context")
	}
	return boutiqueID, nil
}

// Récupérer le prix d'un produit
func (h *VarianteHandler) getPrixProduit(c *fiber.Ctx, produitID string) (float64, error) {
	boutiqueID, err := h.getBoutiqueID(c)
	if err != nil {
		return 0, err
	}

	produit, err := h.produitService.GetByID(c.Context(), produitID, boutiqueID)
	if err != nil {
		return 0, err
	}
	return produit.PrixDefaut, nil
}

// POST /api/produits/:produitId/variantes
func (h *VarianteHandler) CreateVariante(c *fiber.Ctx) error {
	produitID := c.Params("produitId")
	if produitID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "ID produit requis"})
	}

	// Vérifier la boutique
	if _, err := h.getBoutiqueID(c); err != nil {
		return err
	}

	var req dto.RequeteCreationVariante
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "JSON invalide"})
	}

	prixProduit, err := h.getPrixProduit(c, produitID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erreur récupération produit: " + err.Error()})
	}

	variante, err := h.service.Create(c.Context(), produitID, req, prixProduit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(variante)
}

// GET /api/produits/:produitId/variantes
func (h *VarianteHandler) ListVariantes(c *fiber.Ctx) error {
	produitID := c.Params("produitId")
	if produitID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "ID produit requis"})
	}

	if _, err := h.getBoutiqueID(c); err != nil {
		return err
	}

	prixProduit, err := h.getPrixProduit(c, produitID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erreur récupération produit: " + err.Error()})
	}

	variantes, err := h.service.ListByProduit(c.Context(), produitID, prixProduit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"variantes": variantes})
}

// GET /api/variantes/:varianteId
func (h *VarianteHandler) GetVarianteByID(c *fiber.Ctx) error {
	varianteID := c.Params("varianteId")
	if varianteID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "ID variante requis"})
	}

	boutiqueID, err := h.getBoutiqueID(c)
	if err != nil {
		return err
	}

	// D'abord récupérer la variante sans prix
	temp, err := h.service.GetByID(c.Context(), varianteID, 0)
	if err != nil {
		if err.Error() == "variante non trouvée" {
			return c.Status(404).JSON(fiber.Map{"error": "Variante non trouvée"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Récupérer le produit pour avoir son prix par défaut
	produit, err := h.produitService.GetByID(c.Context(), temp.ProduitID, boutiqueID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erreur récupération produit: " + err.Error()})
	}

	// Récupérer la variante avec le bon prix
	variante, err := h.service.GetByID(c.Context(), varianteID, produit.PrixDefaut)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(variante)
}

// PUT /api/variantes/:varianteId
func (h *VarianteHandler) UpdateVariante(c *fiber.Ctx) error {
	varianteID := c.Params("varianteId")
	if varianteID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "ID variante requis"})
	}

	boutiqueID, err := h.getBoutiqueID(c)
	if err != nil {
		return err
	}

	var req dto.RequeteUpdateVariante
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "JSON invalide"})
	}

	// Récupérer la variante sans prix
	temp, err := h.service.GetByID(c.Context(), varianteID, 0)
	if err != nil {
		if err.Error() == "variante non trouvée" {
			return c.Status(404).JSON(fiber.Map{"error": "Variante non trouvée"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Récupérer le produit
	produit, err := h.produitService.GetByID(c.Context(), temp.ProduitID, boutiqueID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Erreur récupération produit: " + err.Error()})
	}

	// Mettre à jour
	variante, err := h.service.Update(c.Context(), varianteID, req, produit.PrixDefaut)
	if err != nil {
		if err.Error() == "variante non trouvée" {
			return c.Status(404).JSON(fiber.Map{"error": "Variante non trouvée"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(variante)
}

// DELETE /api/variantes/:varianteId
func (h *VarianteHandler) DeleteVariante(c *fiber.Ctx) error {
	varianteID := c.Params("varianteId")
	if varianteID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "ID variante requis"})
	}

	if _, err := h.getBoutiqueID(c); err != nil {
		return err
	}

	err := h.service.Delete(c.Context(), varianteID)
	if err != nil {
		if err.Error() == "variante non trouvée" {
			return c.Status(404).JSON(fiber.Map{"error": "Variante non trouvée"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"ok": true})
}
