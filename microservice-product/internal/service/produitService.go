// internal/services/produit_service.go
package services

import (
	"context"
	"errors"

	"projet/internal/dto"
	"projet/internal/models"
	"projet/internal/repository"
	"strings"
	"time"
)

type ProduitService struct {
	Repo *repository.ProduitRepo
}

func NewProduitService(Repo *repository.ProduitRepo) *ProduitService {
	return &ProduitService{Repo: Repo}
}

func (s *ProduitService) Create(ctx context.Context, req dto.RequeteCreationProduit) (*models.Produit, error) {
	if req.Slug == nil || *req.Slug == "" {
		slug := strings.ToLower(strings.ReplaceAll(req.Titre, " ", "-"))
		req.Slug = &slug
	}

	produit := &models.Produit{
		Titre:           req.Titre,
		Description:     req.Description,
		Slug:            *req.Slug,
		Statut:          req.Statut,
		PrixDefaut:      req.PrixDefaut,
		Devise:          req.Devise,
		SKU:             req.SKU,
		SuiviStock:      req.SuiviStock,
		QuantiteStock:   req.QuantiteStock,
		Poids:           req.Poids,
		Marque:          req.Marque,
		ClasseTaxe:      req.ClasseTaxe,
		Visibilite:      req.Visibilite,
		DatePublication: req.DatePublication,
	}

	return s.Repo.CreateProduit(ctx, produit)
}

func (s *ProduitService) GetByID(ctx context.Context, id string) (*models.Produit, error) {
	produit, err := s.Repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if produit == nil {
		return nil, errors.New("product not found")
	}
	return produit, nil
}

func (s *ProduitService) Update(ctx context.Context, id string, req dto.RequeteUpdateProduit) (*models.Produit, error) {
	_, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})

	if req.Titre != nil {
		updates["titre"] = *req.Titre
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Slug != nil {
		updates["slug"] = *req.Slug
	}
	if req.Statut != nil {
		updates["statut"] = *req.Statut
	}
	if req.PrixDefaut != nil {
		updates["prix_defaut"] = *req.PrixDefaut
	}
	if req.Devise != nil {
		updates["devise"] = *req.Devise
	}
	if req.SKU != nil {
		updates["sku"] = *req.SKU
	}
	if req.SuiviStock != nil {
		updates["suivi_stock"] = *req.SuiviStock
	}
	if req.QuantiteStock != nil {
		updates["quantite_stock"] = *req.QuantiteStock
	}
	if req.Poids != nil {
		updates["poids"] = *req.Poids
	}
	if req.Marque != nil {
		updates["marque"] = *req.Marque
	}
	if req.ClasseTaxe != nil {
		updates["classe_taxe"] = *req.ClasseTaxe
	}
	if req.Visibilite != nil {
		updates["visibilite"] = *req.Visibilite
	}
	if req.DatePublication != nil {
		updates["date_publication"] = *req.DatePublication
	}

	updates["mis_a_jour_le"] = time.Now()

	return s.Repo.Update(ctx, id, updates)
}

func (s *ProduitService) Delete(ctx context.Context, id string) error {
	deleted, err := s.Repo.DeleteById(ctx, id)
	if err != nil {
		return err
	}
	if !deleted {
		return errors.New("product not found")
	}
	return nil
}
