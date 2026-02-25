package repository

import (
	"context"
	"fmt"
	"projet/internal/models"
	"time"

	"gorm.io/gorm"
)

type VarianteRepo struct {
	db *gorm.DB
}

func NewVarianteRepo(db *gorm.DB) *VarianteRepo {
	return &VarianteRepo{db: db}
}

func (r *VarianteRepo) CreationVariant(ctx context.Context, variante *models.Variante) (*models.Variante, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := r.db.WithContext(opCtx).Create(variante).Error; err != nil {
		return nil, fmt.Errorf("failed to insert Variante: %w", err)
	}
	return variante, nil
}

func (r *VarianteRepo) ListeProduitID(ctx context.Context, produitID string) ([]models.Variante, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var variantes []models.Variante
	if err := r.db.WithContext(opCtx).Where("produit_id = ?", produitID).Find(&variantes).Error; err != nil {
		return nil, fmt.Errorf("find variantes failed: %w", err)
	}
	return variantes, nil
}

func (r *VarianteRepo) GetByID(ctx context.Context, id string) (*models.Variante, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var variante models.Variante
	err := r.db.WithContext(opCtx).Where("id = ?", id).First(&variante).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching Variante: %w", err)
	}
	return &variante, nil
}

func (r *VarianteRepo) Update(ctx context.Context, id string, updates map[string]interface{}) (*models.Variante, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result := r.db.WithContext(opCtx).Model(&models.Variante{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to update Variante: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	var variante models.Variante
	if err := r.db.WithContext(opCtx).Where("id = ?", id).First(&variante).Error; err != nil {
		return nil, fmt.Errorf("Variante updated but failed to fetch: %w", err)
	}
	return &variante, nil
}

func (r *VarianteRepo) SupprimereById(ctx context.Context, id string) (bool, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result := r.db.WithContext(opCtx).Where("id = ?", id).Delete(&models.Variante{})
	if result.Error != nil {
		return false, fmt.Errorf("failed to delete Variante: %w", result.Error)
	}
	return result.RowsAffected > 0, nil
}

func (r *VarianteRepo) CheckDuplicateCombination(ctx context.Context, produitID string, valeurOptionIDs []string) (bool, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var count int64
	err := r.db.WithContext(opCtx).
		Model(&models.Variante{}).
		Joins("JOIN variante_valeur_option vvo ON vvo.variante_id = variantes.id").
		Where("variantes.produit_id = ?", produitID).
		Where("vvo.valeur_option_id IN ?", valeurOptionIDs).
		Group("variantes.id").
		Having("COUNT(vvo.valeur_option_id) = ?", len(valeurOptionIDs)).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check duplicate combination: %w", err)
	}
	return count > 0, nil
}
