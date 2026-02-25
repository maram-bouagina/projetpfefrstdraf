package handler

import (
	"projet/internal/dto"
	"projet/internal/service"

	"github.com/gofiber/fiber/v2"
)

type OptionProduitHandler struct {
	service *service.OptionProduitService
}

func NewOptionProduitHandler(service *service.OptionProduitService) *OptionProduitHandler {
	return &OptionProduitHandler{service: service}
}

func (h *OptionProduitHandler) getBoutiqueID(c *fiber.Ctx) (string, error) {
	boutiqueID := c.Get("X-Boutique-ID")
	if boutiqueID == "" {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Missing store context")
	}
	return boutiqueID, nil
}

// ============================================================
// OPTIONS
// ============================================================

// POST /api/produits/:produitId/options
func (h *OptionProduitHandler) CreateOption(c *fiber.Ctx) error {
	produitID := c.Params("produitId")
	if produitID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "ID produit requis"})
	}

	var req dto.RequeteCreationOption
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format JSON invalide"})
	}

	if _, err := h.getBoutiqueID(c); err != nil {
		return err
	}

	option, err := h.service.CreationOptionProduit(c.Context(), produitID, req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(option)
}

// GET /api/produits/:produitId/options
func (h *OptionProduitHandler) ListOptions(c *fiber.Ctx) error {
	produitID := c.Params("produitId")
	if produitID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "ID produit requis"})
	}

	if _, err := h.getBoutiqueID(c); err != nil {
		return err
	}

	options, err := h.service.ListOptionProduit(c.Context(), produitID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"options": options})
}

// GET /api/options/:optionId
func (h *OptionProduitHandler) GetOptionByID(c *fiber.Ctx) error {
	optionID := c.Params("optionId")
	if optionID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "ID option requis"})
	}

	if _, err := h.getBoutiqueID(c); err != nil {
		return err
	}

	option, err := h.service.GetByIDOptionProduit(c.Context(), optionID)
	if err != nil {
		if err.Error() == "option non trouvée" {
			return c.Status(404).JSON(fiber.Map{"error": "Option non trouvée"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(option)
}

// PUT /api/options/:optionId
func (h *OptionProduitHandler) UpdateOption(c *fiber.Ctx) error {
	optionID := c.Params("optionId")
	if optionID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "ID option requis"})
	}

	var req dto.RequeteUpdateOption
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format JSON invalide"})
	}

	if _, err := h.getBoutiqueID(c); err != nil {
		return err
	}

	option, err := h.service.Update(c.Context(), optionID, req)
	if err != nil {
		if err.Error() == "option non trouvée" {
			return c.Status(404).JSON(fiber.Map{"error": "Option non trouvée"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(option)
}

// DELETE /api/options/:optionId
func (h *OptionProduitHandler) DeleteOption(c *fiber.Ctx) error {
	optionID := c.Params("optionId")
	if optionID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "ID option requis"})
	}

	if _, err := h.getBoutiqueID(c); err != nil {
		return err
	}

	err := h.service.Delete(c.Context(), optionID)
	if err != nil {
		if err.Error() == "option non trouvée" {
			return c.Status(404).JSON(fiber.Map{"error": "Option non trouvée"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"ok": true})
}

// ============================================================
// VALEURS D'OPTIONS
// ============================================================

// POST /api/options/:optionId/valeurs
func (h *OptionProduitHandler) CreateValeur(c *fiber.Ctx) error {
	optionID := c.Params("optionId")
	if optionID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "ID option requis"})
	}

	var req dto.RequeteCreationValeurOption
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format JSON invalide"})
	}

	if _, err := h.getBoutiqueID(c); err != nil {
		return err
	}

	valeur, err := h.service.CreationValeurOption(c.Context(), optionID, req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(valeur)
}

// GET /api/options/:optionId/valeurs
func (h *OptionProduitHandler) ListValeurs(c *fiber.Ctx) error {
	optionID := c.Params("optionId")
	if optionID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "ID option requis"})
	}

	if _, err := h.getBoutiqueID(c); err != nil {
		return err
	}

	valeurs, err := h.service.ListValeursByOption(c.Context(), optionID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"valeurs": valeurs})
}

// GET /api/valeurs/:valeurId
func (h *OptionProduitHandler) GetValeurByID(c *fiber.Ctx) error {
	valeurID := c.Params("valeurId")
	if valeurID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "ID valeur requis"})
	}

	if _, err := h.getBoutiqueID(c); err != nil {
		return err
	}

	// Implémenter GetByIDValeurOption dans le service si nécessaire
	// Pour l'instant, on peut utiliser ListValeursByOption avec l'ID
	return c.Status(501).JSON(fiber.Map{"error": "Non implémenté"})
}

// PUT /api/valeurs/:valeurId
func (h *OptionProduitHandler) UpdateValeur(c *fiber.Ctx) error {
	valeurID := c.Params("valeurId")
	if valeurID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "ID valeur requis"})
	}

	var req dto.RequeteUpdateValeurOption
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format JSON invalide"})
	}

	if _, err := h.getBoutiqueID(c); err != nil {
		return err
	}

	valeur, err := h.service.UpdateValeur(c.Context(), valeurID, req)
	if err != nil {
		if err.Error() == "valeur non trouvée" {
			return c.Status(404).JSON(fiber.Map{"error": "Valeur non trouvée"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(valeur)
}

// DELETE /api/valeurs/:valeurId
func (h *OptionProduitHandler) DeleteValeur(c *fiber.Ctx) error {
	valeurID := c.Params("valeurId")
	if valeurID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "ID valeur requis"})
	}

	if _, err := h.getBoutiqueID(c); err != nil {
		return err
	}

	err := h.service.DeleteValeur(c.Context(), valeurID)
	if err != nil {
		if err.Error() == "valeur non trouvée" {
			return c.Status(404).JSON(fiber.Map{"error": "Valeur non trouvée"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"ok": true})
}
