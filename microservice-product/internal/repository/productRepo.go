package repository

import (
	"context"
	"fmt"
	"projet/internal/dto"
	"projet/internal/models"
	"time"

	"gorm.io/gorm"
)

type ProduitRepo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *ProduitRepo {
	return &ProduitRepo{db: db}
}

func (r *ProduitRepo) CreateProduit(ctx context.Context, produit *models.Produit) (*models.Produit, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := r.db.WithContext(opCtx).Create(produit).Error; err != nil {
		return nil, fmt.Errorf("failed to insert product: %w", err)
	}
	return produit, nil
}

func (r *ProduitRepo) ListProduits(ctx context.Context, boutiqueID string) ([]models.Produit, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var produits []models.Produit
	if err := r.db.WithContext(opCtx).Where("boutique_id = ?", boutiqueID).Find(&produits).Error; err != nil {
		return nil, fmt.Errorf("find products failed: %w", err)
	}
	return produits, nil
}

func (r *ProduitRepo) GetByID(ctx context.Context, id, boutiqueID string) (*models.Produit, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var produit models.Produit
	err := r.db.WithContext(opCtx).Where("id = ? AND boutique_id = ?", id, boutiqueID).First(&produit).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching product: %w", err)
	}
	return &produit, nil
}

func (r *ProduitRepo) Update(ctx context.Context, id, boutiqueID string, updates map[string]interface{}) (*models.Produit, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result := r.db.WithContext(opCtx).Model(&models.Produit{}).
		Where("id = ? AND boutique_id = ?", id, boutiqueID).
		Updates(updates)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update product: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	var produit models.Produit
	if err := r.db.WithContext(opCtx).Where("id = ?", id).First(&produit).Error; err != nil {
		return nil, fmt.Errorf("product updated but failed to fetch: %w", err)
	}
	return &produit, nil
}

func (r *ProduitRepo) DeleteById(ctx context.Context, id, boutiqueID string) (bool, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result := r.db.WithContext(opCtx).Where("id = ? AND boutique_id = ?", id, boutiqueID).Delete(&models.Produit{})
	if result.Error != nil {
		return false, fmt.Errorf("failed to delete product: %w", result.Error)
	}
	return result.RowsAffected > 0, nil
}

func (r *ProduitRepo) GetWithFilter(ctx context.Context, boutiqueID string, filter dto.FiltreProduit) ([]models.Produit, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var produits []models.Produit
	query := r.db.WithContext(opCtx).Where("boutique_id = ?", boutiqueID)

	if filter.Statut != nil {
		query = query.Where("statut = ?", *filter.Statut)
	}
	if filter.Visibilite != nil {
		query = query.Where("visibilite = ?", *filter.Visibilite)
	}
	if filter.Marque != nil {
		query = query.Where("marque = ?", *filter.Marque)
	}
	if filter.Recherche != nil && *filter.Recherche != "" {
		query = query.Where("titre ILIKE ?", "%"+*filter.Recherche+"%")
	}
	if !filter.InclureSupprime {
		query = query.Where("supprime_le IS NULL")
	}

	offset := (filter.Page - 1) * filter.Limite
	query = query.Limit(filter.Limite).Offset(offset)

	if err := query.Find(&produits).Error; err != nil {
		return nil, fmt.Errorf("filter query failed: %w", err)
	}
	return produits, nil
}
