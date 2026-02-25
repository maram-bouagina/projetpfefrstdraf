package service

import (
	"context"
	"errors"
	"fmt"
	"projet/internal/dto"
	"projet/internal/models"
	"projet/internal/repository"
	"time"
)

type VarianteService struct {
	repo *repository.VarianteRepo
}

func NewVarianteService(repo *repository.VarianteRepo) *VarianteService {
	return &VarianteService{repo: repo}
}

// ------------------------------------------------------------
// Convertisseur
// ------------------------------------------------------------
func (s *VarianteService) toResponse(v models.Variante, prixDefautProduit float64) dto.VarianteResponse {
	valeursOpts := make([]dto.ValeurOptionResponse, len(v.ValeurOptions))
	for i, vo := range v.ValeurOptions {
		valeursOpts[i] = dto.ValeurOptionResponse{
			ID:       vo.ID,
			OptionID: vo.OptionID,
			Valeur:   vo.Valeur,
			Position: vo.Position,
		}
	}

	prixEffectif := prixDefautProduit
	if v.Prix != nil {
		prixEffectif = *v.Prix
	}

	return dto.VarianteResponse{
		ID:            v.ID,
		ProduitID:     v.ProduitID,
		SKU:           v.SKU,
		Prix:          v.Prix,
		QuantiteStock: v.QuantiteStock,
		CodeBarres:    v.CodeBarres,
		Poids:         v.Poids,
		Images:        v.Images,
		CreeLe:        v.CreeLe,
		MisAJourLe:    v.MisAJourLe,
		ValeurOptions: valeursOpts,
		PrixEffectif:  prixEffectif,
	}
}

// ------------------------------------------------------------
// Créer une variante
// ------------------------------------------------------------
func (s *VarianteService) Create(
	ctx context.Context,
	produitID string,
	req dto.RequeteCreationVariante,
	prixDefautProduit float64,
) (*dto.VarianteResponse, error) {

	// Vérifier les doublons
	if len(req.ValeurOptionIDs) > 0 {
		existe, err := s.repo.CheckDuplicateCombination(ctx, produitID, req.ValeurOptionIDs)
		if err != nil {
			return nil, fmt.Errorf("erreur vérification doublon: %v", err)
		}
		if existe {
			return nil, errors.New("cette combinaison existe déjà")
		}
	}

	// Créer la variante
	variante := &models.Variante{
		ProduitID:     produitID,
		SKU:           req.SKU,
		Prix:          req.Prix,
		QuantiteStock: req.QuantiteStock,
		CodeBarres:    req.CodeBarres,
		Poids:         req.Poids,
		Images:        req.Images,
		CreeLe:        time.Now(),
		MisAJourLe:    time.Now(),
	}

	creee, err := s.repo.CreationVariant(ctx, variante)
	if err != nil {
		return nil, fmt.Errorf("échec création: %v", err)
	}

	// Ajouter les relations (à implémenter dans le repo)
	if len(req.ValeurOptionIDs) > 0 {
		// Il faudra ajouter une méthode dans le repo pour attacher les valeurs
		// Pour l'instant, on suppose que c'est géré ailleurs
	}

	// Récupérer la variante complète
	finale, err := s.repo.GetByID(ctx, creee.ID)
	if err != nil {
		return nil, err
	}

	reponse := s.toResponse(*finale, prixDefautProduit)
	return &reponse, nil
}

// ------------------------------------------------------------
// Lister les variantes d'un produit
// ------------------------------------------------------------
func (s *VarianteService) ListByProduit(ctx context.Context, produitID string, prixDefautProduit float64) ([]dto.VarianteResponse, error) {
	if produitID == "" {
		return nil, errors.New("ID produit requis")
	}

	variantes, err := s.repo.ListeProduitID(ctx, produitID)
	if err != nil {
		return nil, err
	}

	resultats := make([]dto.VarianteResponse, len(variantes))
	for i, v := range variantes {
		resultats[i] = s.toResponse(v, prixDefautProduit)
	}
	return resultats, nil
}

// ------------------------------------------------------------
// Récupérer une variante par ID
// ------------------------------------------------------------
func (s *VarianteService) GetByID(ctx context.Context, id string, prixDefautProduit float64) (*dto.VarianteResponse, error) {
	variante, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if variante == nil {
		return nil, errors.New("variante non trouvée")
	}

	reponse := s.toResponse(*variante, prixDefautProduit)
	return &reponse, nil
}

// ------------------------------------------------------------
// Mettre à jour une variante
// ------------------------------------------------------------
func (s *VarianteService) Update(ctx context.Context, id string, req dto.RequeteUpdateVariante, prixDefautProduit float64) (*dto.VarianteResponse, error) {
	// Vérifier que la variante existe
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Préparer les modifs
	modifications := make(map[string]interface{})
	if req.SKU != nil {
		modifications["sku"] = *req.SKU
	}
	if req.Prix != nil {
		modifications["prix"] = *req.Prix
	}
	if req.QuantiteStock != nil {
		modifications["quantite_stock"] = *req.QuantiteStock
	}
	if req.CodeBarres != nil {
		modifications["code_barres"] = *req.CodeBarres
	}
	if req.Poids != nil {
		modifications["poids"] = *req.Poids
	}
	if req.Images != nil {
		modifications["images"] = req.Images
	}
	modifications["mis_a_jour_le"] = time.Now()

	// Mettre à jour
	modifiee, err := s.repo.Update(ctx, id, modifications)
	if err != nil {
		return nil, err
	}
	if modifiee == nil {
		return nil, errors.New("variante non trouvée après update")
	}

	reponse := s.toResponse(*modifiee, prixDefautProduit)
	return &reponse, nil
}

// ------------------------------------------------------------
// Supprimer une variante
// ------------------------------------------------------------
func (s *VarianteService) Delete(ctx context.Context, id string) error {
	supprimee, err := s.repo.SupprimereById(ctx, id)
	if err != nil {
		return err
	}
	if !supprimee {
		return errors.New("variante non trouvée")
	}
	return nil
}
