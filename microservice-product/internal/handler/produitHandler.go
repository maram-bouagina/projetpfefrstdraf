package handler

import (
	"projet/internal/dto"
	services "projet/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// bch tchuf les donéées respctiw dto walla same fields and conditions etc
var validate = validator.New()

// hedhi injection de dépendance
type ProduitHandler struct {
	service *services.ProduitService
}

// Hedha constructeur
// we want a pointer to modify the original instance
func NewProduitHandler(service *services.ProduitService) *ProduitHandler {
	return &ProduitHandler{service: service}
}

func (h *ProduitHandler) getBoutiqueID(c *fiber.Ctx) (string, error) {
	//bech t9ra l headers mt3 http
	boutiqueID := c.Get("X-Boutique-ID")
	if boutiqueID == "" {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Missing store context")
	}
	return boutiqueID, nil
}

func (h *ProduitHandler) CreateProduit(c *fiber.Ctx) error {
	var req dto.RequeteCreationProduit
	//decodi min json li struct
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
	//for path parameters kima produits/1<-
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
	//for path parameters
	id := c.Params("id")

	var req dto.RequeteUpdateProduit
	//decodi min json li struct
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

/*ylwj par id tore ou id produit ou faama tests pour les erreurs simple pas besoin de more explanations*/
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

/*concentre toi sur ce que retourne la méthode --->error*/
func (h *ProduitHandler) SearchProduits(c *fiber.Ctx) error {
	//declarit FiltreProduit mi dto
	var filter dto.FiltreProduit

	//fl postman famma akl les parametres fl url ismhm query parameters
	//tht fl query parameter l'addresse mtaa akl object mtaa dto bech t3abbiha bl query parameter
	if err := c.QueryParser(&filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid filter parameters"})
	}

	//tcherchi par store id
	boutiqueID, err := h.getBoutiqueID(c)
	if err != nil {
		return err
	}

	//t3yt ll func illi fi service
	produits, page, limite, err := h.service.Search(c.Context(), boutiqueID, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to search products"})
	}

	//houni trj3 nil khtr fmch erreur, ou howwa yriturni erorr ou zeda yiktb response
	return c.JSON(fiber.Map{
		"produits": produits,
		"page":     page,
		"limite":   limite,
	})
}
