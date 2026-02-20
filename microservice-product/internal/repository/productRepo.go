// internal/repository/produit_repo.go
package repository

import (
	"context"
	"fmt"
	"projet/internal/models"
	"time"

	"gorm.io/gorm"
)

type ProduitRepo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *ProduitRepo {
	return &ProduitRepo{
		db: db,
	}
}

func (r *ProduitRepo) CreateProduit(ctx context.Context, produit *models.Produit) (*models.Produit, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := r.db.WithContext(opCtx).Create(produit).Error; err != nil {
		return nil, fmt.Errorf("failed to insert product: %v", err)
	}
	return produit, nil
}

func (r *ProduitRepo) ListProduits(ctx context.Context) ([]models.Produit, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var produits []models.Produit
	if err := r.db.WithContext(opCtx).Find(&produits).Error; err != nil {
		return nil, fmt.Errorf("find products failed: %w", err)
	}
	return produits, nil
}

func (r *ProduitRepo) GetByID(ctx context.Context, id string) (*models.Produit, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var produit models.Produit
	if err := r.db.WithContext(opCtx).Where("id = ?", id).First(&produit).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error in fetching: %w", err)
	}
	return &produit, nil
}

func (r *ProduitRepo) Update(ctx context.Context, id string, updates map[string]interface{}) (*models.Produit, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result := r.db.WithContext(opCtx).Model(&models.Produit{}).Where("id = ?", id).Updates(updates)
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

func (r *ProduitRepo) DeleteById(ctx context.Context, id string) (bool, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result := r.db.WithContext(opCtx).Where("id = ?", id).Delete(&models.Produit{})
	if result.Error != nil {
		return false, fmt.Errorf("failed to delete product: %w", result.Error)
	}
	return result.RowsAffected > 0, nil
}

func (r *ProduitRepo) GetWithFilter(ctx context.Context, filter map[string]interface{}) ([]models.Produit, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var produits []models.Produit
	query := r.db.WithContext(opCtx)

	for field, value := range filter {
		query = query.Where(field, value)
	}

	if err := query.Find(&produits).Error; err != nil {
		return nil, fmt.Errorf("filter query failed: %w", err)
	}
	return produits, nil
}
