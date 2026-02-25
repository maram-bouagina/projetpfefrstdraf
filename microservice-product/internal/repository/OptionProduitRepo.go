package repository

import (
	"context"
	"fmt"
	"projet/internal/models"
	"time"

	"gorm.io/gorm"
)

type OptionProduitValeurRepo struct {
	db *gorm.DB
}

func NewOptionProduitValeurRepo(db *gorm.DB) *OptionProduitValeurRepo {
	return &OptionProduitValeurRepo{db: db}
}

// creation
func (r *OptionProduitValeurRepo) CreationOptProduit(ctx context.Context, optProduit *models.OptionProduit) (*models.OptionProduit, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := r.db.WithContext(opCtx).Create(optProduit).Error; err != nil {
		return nil, fmt.Errorf("failed to insert ProductOption: %w", err)
	}
	return optProduit, nil
}

func (r *OptionProduitValeurRepo) CreationValeurOption(ctx context.Context, optProduit *models.ValeurOption) (*models.ValeurOption, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := r.db.WithContext(opCtx).Create(optProduit).Error; err != nil {
		return nil, fmt.Errorf("failed to insert valeurOption: %w", err)
	}
	return optProduit, nil
}

// Liste
func (r *OptionProduitValeurRepo) ListeOptProduits(ctx context.Context, produitID string) ([]models.OptionProduit, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var produits []models.OptionProduit
	if err := r.db.WithContext(opCtx).
		Where("produit_id = ?", produitID).
		Preload("ValeurOpts", func(db *gorm.DB) *gorm.DB {
			return db.Order("position")
		}).
		Find(&produits).Error; err != nil {
		return nil, fmt.Errorf("find OptionProducts failed: %w", err)
	}
	return produits, nil
}

func (r *OptionProduitValeurRepo) ListeValeursOption(ctx context.Context, OptionID string) ([]models.ValeurOption, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var valeuropt []models.ValeurOption
	if err := r.db.WithContext(opCtx).Where("option_id = ?", OptionID).Find(&valeuropt).Error; err != nil {
		return nil, fmt.Errorf("find valeurOption failed: %w", err)
	}
	return valeuropt, nil
}

// GetById
func (r *OptionProduitValeurRepo) GetByIdOptionproduit(ctx context.Context, id string) (*models.OptionProduit, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var optProduit models.OptionProduit
	err := r.db.WithContext(opCtx).
		Where("id = ?", id).
		Preload("ValeurOpts", func(db *gorm.DB) *gorm.DB {
			return db.Order("position")
		}).
		First(&optProduit).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching ProductOption: %w", err)
	}
	return &optProduit, nil
}

func (r *OptionProduitValeurRepo) GetByIDValeurOption(ctx context.Context, id string) (*models.ValeurOption, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var valeurOpt models.ValeurOption
	err := r.db.WithContext(opCtx).Where("id = ?", id).First(&valeurOpt).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching valeurOption: %w", err)
	}
	return &valeurOpt, nil
}

// update
func (r *OptionProduitValeurRepo) UpdateProduitOption(ctx context.Context, id string, updates map[string]interface{}) (*models.OptionProduit, error) {
	/*yhdhr fil context mtaa bdd*/
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	/*9aad ybdati*/
	//milloul yimchi li table optProduit bModel ou baad bidhbt win bl id. ou baad yaaml l u^date
	result := r.db.WithContext(opCtx).Model(&models.OptionProduit{}).
		Where("id = ?", id).
		Updates(updates)

	/*ytesti l9aha walla mal9ahech w njhit wella*/
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update ProductOption: %w", result.Error)
	}
	//ml9a hatte ligne
	if result.RowsAffected == 0 {
		return nil, nil
	}

	/*ki nijhit 9aadin nlwjou bech nrja3ou lprod*/
	var optProduit models.OptionProduit
	if err := r.db.WithContext(opCtx).Where("id = ?", id).First(&optProduit).Error; err != nil {
		return nil, fmt.Errorf("ProductOption updated but failed to fetch: %w", err)
	}
	return &optProduit, nil
}
func (r *OptionProduitValeurRepo) UpdateValeurOpt(ctx context.Context, id string, updates map[string]interface{}) (*models.ValeurOption, error) {
	/*yhdhr fil context mtaa bdd*/
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	/*9aad ybdati*/
	//milloul yimchi li table optProduit bModel ou baad bidhbt win bl id. ou baad yaaml l u^date
	result := r.db.WithContext(opCtx).Model(&models.ValeurOption{}).
		Where("id = ?", id).
		Updates(updates)

	/*ytesti l9aha walla mal9ahech w njhit wella*/
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update ProductOption: %w", result.Error)
	}
	//ml9a hatte ligne
	if result.RowsAffected == 0 {
		return nil, nil
	}

	/*ki nijhit 9aadin nlwjou bech nrja3ou lprod*/
	var vOpt models.ValeurOption
	if err := r.db.WithContext(opCtx).Where("id = ?", id).First(&vOpt).Error; err != nil {
		return nil, fmt.Errorf("valeurOption updated but failed to fetch: %w", err)
	}
	return &vOpt, nil
}

// Suppression
func (r *OptionProduitValeurRepo) SupprimerOptPById(ctx context.Context, id string) (bool, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result := r.db.WithContext(opCtx).Where("id = ?", id).Delete(&models.OptionProduit{})
	if result.Error != nil {
		return false, fmt.Errorf("failed to delete ProductOption: %w", result.Error)
	}
	return result.RowsAffected > 0, nil
}

func (r *OptionProduitValeurRepo) SupprimerByIdVOpt(ctx context.Context, id string) (bool, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result := r.db.WithContext(opCtx).Where("id = ?", id).Delete(&models.ValeurOption{})
	if result.Error != nil {
		return false, fmt.Errorf("failed to delete valeurOption: %w", result.Error)
	}
	return result.RowsAffected > 0, nil
}

func (r *OptionProduitValeurRepo) CountOptionsByProduit(ctx context.Context, produitID string) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.OptionProduit{}).
		Where("produit_id = ?", produitID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// À ajouter dans OptionProduitValeurRepo
func (r *OptionProduitValeurRepo) CountValeursByOption(ctx context.Context, optionID string) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.ValeurOption{}).
		Where("option_id = ?", optionID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
